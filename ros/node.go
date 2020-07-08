package ros

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	modular "github.com/edwinhayes/logrus-modular"
	"github.com/edwinhayes/rosgo/xmlrpc"
	"github.com/sirupsen/logrus"
)

const (
	//APIStatusError is an API call which returned an Error
	APIStatusError = -1
	//APIStatusFailure is a failed API call
	APIStatusFailure = 0
	//APIStatusSuccess is a successful API call
	APIStatusSuccess = 1
	//Remap string constant for splitting components
	Remap = ":="
	//Default setting for OS interrupts
	defaultInterrupts = true
)

func processArguments(args []string) (NameMap, NameMap, NameMap, []string) {
	mapping := make(NameMap)
	params := make(NameMap)
	specials := make(NameMap)
	rest := make([]string, 0)
	for _, arg := range args {
		components := strings.Split(arg, Remap)
		if len(components) == 2 {
			key := components[0]
			value := components[1]
			if strings.HasPrefix(key, "__") {
				specials[key] = value
			} else if strings.HasPrefix(key, "_") {
				params[key[1:]] = value
			} else {
				mapping[key] = value
			}
		} else {
			rest = append(rest, arg)
		}
	}
	return mapping, params, specials, rest
}

// *defaultNode implements Node interface
// a defaultNode instance must be accessed in user goroutine.
type defaultNode struct {
	name             string
	namespace        string
	qualifiedName    string
	masterURI        string
	xmlrpcURI        string
	xmlrpcListener   net.Listener
	xmlrpcHandler    *xmlrpc.Handler
	subscribers      map[string]*defaultSubscriber
	publishers       sync.Map
	servers          map[string]*defaultServiceServer
	jobChan          chan func()
	interruptChan    chan os.Signal
	enableInterrupts bool
	logger           modular.ModuleLogger
	ok               bool
	okMutex          sync.RWMutex
	waitGroup        sync.WaitGroup
	logDir           string
	hostname         string
	listenIP         string
	homeDir          string
	nameResolver     *NameResolver
	nonRosArgs       []string
}

// serviceheader is the header returned from probing a ros service, containing all type information
type serviceHeader struct {
	callerid     string
	md5sum       string
	requestType  string
	responseType string
	serviceType  string
}

func listenRandomPort(address string, trialLimit int) (net.Listener, error) {
	var listener net.Listener
	var err error
	numTrial := 0
	for numTrial < trialLimit {
		port := 1024 + rand.Intn(65535-1024)
		addr := fmt.Sprintf("%s:%d", address, port)
		listener, err = net.Listen("tcp", addr)
		if err == nil {
			return listener, nil
		}
		numTrial++

	}
	return nil, fmt.Errorf("listenRandomPort exceeds trial limit")
}

func newDefaultNodeWithLogs(name string, logger *modular.ModuleLogger, args []string) (*defaultNode, error) {
	node, err := newDefaultNode(name, args)
	if err != nil {
		(*logger).Errorf("could not instantiate newDefaultNode : %v", err)
		return nil, err
	}

	node.logger = (*logger).(modular.ModuleLogger)
	return node, nil
}

func newDefaultNode(name string, args []string) (*defaultNode, error) {
	node := new(defaultNode)

	namespace, nodeName, err := qualifyNodeName(name)
	if err != nil {
		return nil, err
	}

	remapping, params, specials, rest := processArguments(args)

	node.homeDir = filepath.Join(os.Getenv("HOME"), ".ros")
	if homeDir := os.Getenv("ROS_HOME"); len(homeDir) > 0 {
		node.homeDir = homeDir
	}

	logger := logrus.New()
	rootLogger := modular.NewRootLogger(logger)
	if value, ok := specials["__ll"]; ok {
		val, err := strconv.ParseInt(value, 10, 32)
		if err == nil {
			logger.SetLevel(logrus.Level(val))
		}
	} else {
		logger.SetLevel(logrus.FatalLevel)
	}
	node.logger = rootLogger

	node.name = nodeName
	if value, ok := specials["__name"]; ok {
		node.name = value
	}

	node.namespace = namespace
	if ns := os.Getenv("ROS_NAMESPACE"); len(ns) > 0 {
		node.namespace = ns
	}
	if value, ok := specials["__ns"]; ok {
		node.namespace = value
	}
	if value, ok := specials["__si"]; ok {
		val, err := strconv.ParseBool(value)
		if err != nil {
			node.enableInterrupts = val
		}
	} else {
		node.enableInterrupts = defaultInterrupts
	}
	node.logDir = filepath.Join(node.homeDir, "log")
	if logDir := os.Getenv("ROS_LOG_DIR"); len(logDir) > 0 {
		node.logDir = logDir
	}
	if value, ok := specials["__log"]; ok {
		node.logDir = value
	}

	var onlyLocalhost bool
	node.hostname, onlyLocalhost = determineHost()
	if value, ok := specials["__hostname"]; ok {
		node.hostname = value
		onlyLocalhost = (value == "localhost")
	} else if value, ok := specials["__ip"]; ok {
		node.hostname = value
		onlyLocalhost = (value == "::1" || strings.HasPrefix(value, "127."))
	}
	if onlyLocalhost {
		node.listenIP = "127.0.0.1"
	} else {
		node.listenIP = "0.0.0.0"
	}

	node.masterURI = os.Getenv("ROS_MASTER_URI")
	if value, ok := specials["__master"]; ok {
		node.masterURI = value
	}

	node.nameResolver = newNameResolver(node.namespace, node.name, remapping)
	node.nonRosArgs = rest

	node.qualifiedName = node.namespace + node.name
	node.subscribers = make(map[string]*defaultSubscriber)
	node.servers = make(map[string]*defaultServiceServer)
	node.interruptChan = make(chan os.Signal)
	node.ok = true

	// Install signal handler
	if node.enableInterrupts == true {
		signal.Notify(node.interruptChan, os.Interrupt)
		go func() {
			<-node.interruptChan
			logger.Info("Interrupted")
			node.okMutex.Lock()
			node.ok = false
			node.okMutex.Unlock()
		}()
	}
	node.jobChan = make(chan func(), 100)

	logger.Debugf("Master URI = %s", node.masterURI)

	// Set parameters set by arguments
	for k, v := range params {
		_, err := callRosAPI(node.masterURI, "setParam", node.qualifiedName, k, v)
		if err != nil {
			return nil, err
		}
	}

	listener, err := listenRandomPort(node.listenIP, 10)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		// Not reached
		logger.Error(err)
		return nil, err
	}
	node.xmlrpcURI = fmt.Sprintf("http://%s:%s", node.hostname, port)
	logger.Debugf("listen on http://%s", listener.Addr().String())
	node.xmlrpcListener = listener
	m := map[string]xmlrpc.Method{
		"getBusStats":      func(callerId string) (interface{}, error) { return node.getBusStats(callerId) },
		"getBusInfo":       func(callerId string) (interface{}, error) { return node.getBusInfo(callerId) },
		"getMasterUri":     func(callerId string) (interface{}, error) { return node.getMasterURI(callerId) },
		"shutdown":         func(callerId string, msg string) (interface{}, error) { return node.shutdown(callerId, msg) },
		"getPid":           func(callerId string) (interface{}, error) { return node.getPid(callerId) },
		"getSubscriptions": func(callerId string) (interface{}, error) { return node.getSubscriptions(callerId) },
		"getPublications":  func(callerId string) (interface{}, error) { return node.getPublications(callerId) },
		"paramUpdate": func(callerId string, key string, value interface{}) (interface{}, error) {
			return node.paramUpdate(callerId, key, value)
		},
		"publisherUpdate": func(callerId string, topic string, publishers []interface{}) (interface{}, error) {
			return node.publisherUpdate(callerId, topic, publishers)
		},
		"requestTopic": func(callerId string, topic string, protocols []interface{}) (interface{}, error) {
			return node.requestTopic(callerId, topic, protocols)
		},
	}
	node.xmlrpcHandler = xmlrpc.NewHandler(m)
	go http.Serve(node.xmlrpcListener, node.xmlrpcHandler)
	logger.Debugf("Started %s", node.qualifiedName)
	return node, nil
}

func (node *defaultNode) OK() bool {
	node.okMutex.RLock()
	ok := node.ok
	node.okMutex.RUnlock()
	return ok
}

func (node *defaultNode) RemovePublisher(topic string) {
	name := node.nameResolver.remap(topic)
	if pub, ok := node.publishers.Load(name); ok {
		pub.(*defaultPublisher).Shutdown()
		node.publishers.Delete(name)
	}
}

func (node *defaultNode) Name() string {
	return node.name
}

func (node *defaultNode) getBusStats(callerID string) (interface{}, error) {
	return buildRosAPIResult(-1, "Not implemented", 0), nil
}

func (node *defaultNode) getBusInfo(callerID string) (interface{}, error) {
	return buildRosAPIResult(-1, "Not implemeted", 0), nil
}

func (node *defaultNode) getMasterURI(callerID string) (interface{}, error) {
	return buildRosAPIResult(0, "Success", node.masterURI), nil
}

func (node *defaultNode) shutdown(callerID string, msg string) (interface{}, error) {
	node.okMutex.Lock()
	node.ok = false
	node.okMutex.Unlock()
	return buildRosAPIResult(0, "Success", 0), nil
}

func (node *defaultNode) getPid(callerID string) (interface{}, error) {
	return buildRosAPIResult(0, "Success", os.Getpid()), nil
}

func (node *defaultNode) getSubscriptions(callerID string) (interface{}, error) {
	result := []interface{}{}
	for t, s := range node.subscribers {
		pair := []interface{}{t, s.msgType.Name()}
		result = append(result, pair)
	}
	return buildRosAPIResult(0, "Success", result), nil
}

func (node *defaultNode) getPublications(callerID string) (interface{}, error) {
	result := []interface{}{}
	node.publishers.Range(func(t interface{}, p interface{}) bool {
		pair := []interface{}{
			t.(string),
			p.(*defaultPublisher).msgType.Name(),
		}
		result = append(result, pair)
		return true
	})

	return buildRosAPIResult(0, "Success", result), nil
}

func (node *defaultNode) paramUpdate(callerID string, key string, value interface{}) (interface{}, error) {
	return buildRosAPIResult(-1, "Not implemented", 0), nil
}

func (node *defaultNode) publisherUpdate(callerID string, topic string, publishers []interface{}) (interface{}, error) {
	node.logger.Debug("Slave API publisherUpdate() called.")
	var code int32
	var message string
	if sub, ok := node.subscribers[topic]; !ok {
		node.logger.Debug("publisherUpdate() called without subscribing topic.")
		code = 0
		message = "No such topic"
	} else {
		pubUris := make([]string, len(publishers))
		for i, uri := range publishers {
			pubUris[i] = uri.(string)
		}
		sub.pubListChan <- pubUris
		code = 1
		message = "Success"
	}
	return buildRosAPIResult(code, message, 0), nil
}

func (node *defaultNode) requestTopic(callerID string, topic string, protocols []interface{}) (interface{}, error) {
	node.logger.Debugf("Slave API requestTopic(%s, %s, ...) called.", callerID, topic)
	var code int32
	var message string
	var value interface{}
	if pub, ok := node.publishers.Load(topic); !ok {
		node.logger.Debug("requestTopic() called with not publishing topic.")
		code = 0
		message = "No such topic"
		value = nil
	} else {
		selectedProtocol := make([]interface{}, 0)
		for _, v := range protocols {
			protocolParams := v.([]interface{})
			protocolName := protocolParams[0].(string)
			if protocolName == "TCPROS" {
				node.logger.Debug("TCPROS requested")
				selectedProtocol = append(selectedProtocol, "TCPROS")
				host, portStr, err := pub.(*defaultPublisher).hostAndPort()
				if err != nil {
					return nil, err
				}
				p, err := strconv.ParseInt(portStr, 10, 32)
				if err != nil {
					return nil, err
				}
				port := int(p)
				selectedProtocol = append(selectedProtocol, host)
				selectedProtocol = append(selectedProtocol, port)
				break
			}
		}
		node.logger.Debug(selectedProtocol)
		code = 1
		message = "Success"
		value = selectedProtocol
	}
	return buildRosAPIResult(code, message, value), nil
}

func (node *defaultNode) NewPublisher(topic string, msgType MessageType) (Publisher, error) {
	name := node.nameResolver.remap(topic)
	return node.NewPublisherWithCallbacks(name, msgType, nil, nil)
}

func (node *defaultNode) NewPublisherWithCallbacks(topic string, msgType MessageType, connectCallback, disconnectCallback func(SingleSubscriberPublisher)) (Publisher, error) {
	name := node.nameResolver.remap(topic)
	pub, ok := node.publishers.Load(topic)
	if !ok {
		_, err := callRosAPI(node.masterURI, "registerPublisher",
			node.qualifiedName,
			name, msgType.Name(),
			node.xmlrpcURI)
		if err != nil {
			node.logger.Errorf("Failed to call registerPublisher(): %s", err)
			return nil, err
		}

		pub = newDefaultPublisher(node, name, msgType, connectCallback, disconnectCallback)
		node.publishers.Store(name, pub)
		go pub.(*defaultPublisher).start(&node.waitGroup)
	}
	return pub.(*defaultPublisher), nil
}

// Master API for getSystemState
func (node *defaultNode) GetSystemState() ([]interface{}, error) {
	node.logger.Debug("Call Master API getSystemState")
	result, err := callRosAPI(node.masterURI, "getSystemState",
		node.qualifiedName)
	if err != nil {
		node.logger.Errorf("Failed to call getSystemState() for %s.", err)
		return nil, err
	}
	list, ok := result.([]interface{})
	if !ok {
		node.logger.Errorf("result is not []string but %s.", reflect.TypeOf(result).String())
	}
	node.logger.Trace("Result: ", list)
	return list, nil
}

// GetServiceTypes uses getSystemState, and probes the services to return service types
func (node *defaultNode) GetServiceTypes() (map[string]*serviceHeader, error) {
	// Get the system state
	sysState, err := node.GetSystemState()
	if err != nil {
		node.logger.Errorf("Failed to call getSystemState() for %s.", err)
		return nil, err
	}
	// result map
	serviceTypes := make(map[string]*serviceHeader)
	// Get the service list
	services := sysState[2].([]interface{})
	for _, s := range services {
		serviceItem := s.([]interface{})
		serviceName := serviceItem[0].(string)
		// Probe the service
		result, err := callRosAPI(node.masterURI, "lookupService", node.qualifiedName, serviceName)
		if err != nil {
			return nil, err
		}

		serviceRawURL, converted := result.(string)
		if !converted {
			return nil, fmt.Errorf("Result of 'lookupService' is not a string")
		}
		var serviceURL *url.URL
		if serviceURL, err = url.Parse(serviceRawURL); err != nil {
			return nil, err
		}
		var conn net.Conn
		if conn, err = net.Dial("tcp", serviceURL.Host); err != nil {
			return nil, err
		}

		// Write connection header
		var headers []header
		headers = append(headers, header{"probe", "1"})
		headers = append(headers, header{"md5sum", "*"})
		headers = append(headers, header{"callerid", node.qualifiedName})
		headers = append(headers, header{"service", serviceName})

		conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
		if err := writeConnectionHeader(headers, conn); err != nil {
			return nil, err
		}
		// Read reponse header
		conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
		resHeaders, err := readConnectionHeader(conn)
		if err != nil {
			return nil, err
		}
		// Convert headers to map
		resHeaderMap := make(map[string]string)
		for _, h := range resHeaders {
			resHeaderMap[h.key] = h.value
		}
		srvHeader := serviceHeader{
			callerid:     resHeaderMap["callerid"],
			md5sum:       resHeaderMap["md5sum"],
			requestType:  resHeaderMap["request_type"],
			responseType: resHeaderMap["response_type"],
			serviceType:  resHeaderMap["type"],
		}
		serviceTypes[serviceName] = &srvHeader
	}
	return serviceTypes, nil
}

// Master API call for getPublishedTopics
func (node *defaultNode) GetPublishedTopics(subgraph string) ([]interface{}, error) {
	node.logger.Debug("Call Master API getPublishedTopics")
	result, err := callRosAPI(node.masterURI, "getPublishedTopics",
		node.qualifiedName,
		subgraph)
	if err != nil {
		node.logger.Errorf("Failed to call getPublishedTopics() for %s.", err)
		return nil, err
	}
	list, ok := result.([]interface{})
	if !ok {
		node.logger.Errorf("result is not []string but %s.", reflect.TypeOf(result).String())
	}
	node.logger.Trace("Result: ", list)
	return list, nil
}

// Master API call for getTopicTypes
func (node *defaultNode) GetTopicTypes() []interface{} {
	node.logger.Debug("Call Master API getTopicTypes")
	result, err := callRosAPI(node.masterURI, "getTopicTypes",
		node.qualifiedName)
	if err != nil {
		node.logger.Errorf("Failed to call getTopicTypes() for %s.", err)
	}
	list, ok := result.([]interface{})
	if !ok {
		node.logger.Errorf("result is not []string but %s.", reflect.TypeOf(result).String())
	}
	node.logger.Debug("Result: ", list)
	return list
}

// RemoveSubscriber shuts down and deletes an existing topic subscriber.
func (node *defaultNode) RemoveSubscriber(topic string) {
	name := node.nameResolver.remap(topic)
	if sub, ok := node.subscribers[name]; ok {
		sub.Shutdown()
		delete(node.subscribers, name)
	}
}

func (node *defaultNode) NewSubscriber(topic string, msgType MessageType, callback interface{}) (Subscriber, error) {
	name := node.nameResolver.remap(topic)
	sub, ok := node.subscribers[name]
	if !ok {
		node.logger.Debug("Call Master API registerSubscriber")
		result, err := callRosAPI(node.masterURI, "registerSubscriber",
			node.qualifiedName,
			name,
			msgType.Name(),
			node.xmlrpcURI)
		if err != nil {
			node.logger.Errorf("Failed to call registerSubscriber() for %s.", err)
			return nil, err
		}
		list, ok := result.([]interface{})
		if !ok {
			node.logger.Errorf("result is not []string but %s.", reflect.TypeOf(result).String())
		}
		var publishers []string
		for _, item := range list {
			s, ok := item.(string)
			if !ok {
				node.logger.Error("Publisher list contains no string object")
			}
			publishers = append(publishers, s)
		}

		node.logger.Debugf("Publisher URI list: %v", publishers)

		sub = newDefaultSubscriber(name, msgType, callback)
		node.subscribers[name] = sub

		node.logger.Debugf("Start subscriber goroutine for topic '%s'", sub.topic)
		go sub.start(&node.waitGroup, node.qualifiedName, node.xmlrpcURI, node.masterURI, node.jobChan, &node.logger)
		node.logger.Debugf("Done")
		sub.pubListChan <- publishers
		node.logger.Debugf("Update publisher list for topic '%s'", sub.topic)
	} else {
		sub.callbacks = append(sub.callbacks, callback)
	}
	return sub, nil
}

func (node *defaultNode) NewServiceClient(service string, srvType ServiceType) ServiceClient {
	name := node.nameResolver.remap(service)
	client := newDefaultServiceClient(&node.logger, node.qualifiedName, node.masterURI, name, srvType)
	return client
}

func (node *defaultNode) NewServiceServer(service string, srvType ServiceType, handler interface{}) ServiceServer {
	name := node.nameResolver.remap(service)
	server, ok := node.servers[name]
	if ok {
		server.Shutdown()
	}
	server = newDefaultServiceServer(node, name, srvType, handler)
	if server == nil {
		return nil
	}
	node.servers[name] = server
	return server
}

func (node *defaultNode) SpinOnce() bool {
	timeoutChan := time.After(10 * time.Millisecond)
	select {
	case job := <-node.jobChan:
		job()
		return false
	case <-timeoutChan:
		break
	}
	return true
}

func (node *defaultNode) Spin() {
	for node.OK() {
		timeoutChan := time.After(1000 * time.Millisecond)
		select {
		case job := <-node.jobChan:
			node.logger.Debug("Execute job")
			job()
		case <-timeoutChan:
			break
		}
	}
}

func (node *defaultNode) Shutdown() {
	node.logger.Debug("Shutting node down")
	node.okMutex.Lock()
	node.ok = false
	node.okMutex.Unlock()
	node.logger.Debug("Shutdown subscribers")
	for _, s := range node.subscribers {
		s.Shutdown()
	}
	node.logger.Debug("Shutdown subscribers...done")
	node.logger.Debug("Shutdown publishers")
	node.publishers.Range(func(key interface{}, value interface{}) bool {
		value.(*defaultPublisher).Shutdown()
		return true
	})
	node.logger.Debug("Shutdown publishers...done")
	node.logger.Debug("Shutdown servers")
	for _, s := range node.servers {
		s.Shutdown()
	}
	node.logger.Debug("Shutdown servers...done")
	node.logger.Debug("Wait all goroutines")
	node.waitGroup.Wait()
	node.logger.Debug("Wait all goroutines...Done")
	node.logger.Debug("Close XMLRPC lisetner")
	node.xmlrpcListener.Close()
	node.logger.Debug("Close XMLRPC done")
	node.logger.Debug("Wait XMLRPC server shutdown")
	node.xmlrpcHandler.WaitForShutdown()
	node.logger.Debug("Wait XMLRPC server shutdown...Done")
	node.logger.Debug("Shutting node down completed")
	return
}

func (node *defaultNode) GetParam(key string) (interface{}, error) {
	name := node.nameResolver.remap(key)
	return callRosAPI(node.masterURI, "getParam", node.qualifiedName, name)
}

func (node *defaultNode) SetParam(key string, value interface{}) error {
	name := node.nameResolver.remap(key)
	_, e := callRosAPI(node.masterURI, "setParam", node.qualifiedName, name, value)
	return e
}

func (node *defaultNode) HasParam(key string) (bool, error) {
	name := node.nameResolver.remap(key)
	result, err := callRosAPI(node.masterURI, "hasParam", node.qualifiedName, name)
	if err != nil {
		return false, err
	}
	hasParam := result.(bool)
	return hasParam, nil
}

func (node *defaultNode) SearchParam(key string) (string, error) {
	result, err := callRosAPI(node.masterURI, "searchParam", node.qualifiedName, key)
	if err != nil {
		return "", err
	}
	foundKey := result.(string)
	return foundKey, nil
}

func (node *defaultNode) DeleteParam(key string) error {
	name := node.nameResolver.remap(key)
	_, err := callRosAPI(node.masterURI, "deleteParam", node.qualifiedName, name)
	return err
}

func (node *defaultNode) Logger() *modular.ModuleLogger {
	return &node.logger
}

func (node *defaultNode) NonRosArgs() []string {
	return node.nonRosArgs
}

func loadParamFromString(s string) (interface{}, error) {
	decoder := json.NewDecoder(strings.NewReader(s))
	var value interface{}
	err := decoder.Decode(&value)
	if err != nil {
		return nil, err
	}
	return value, err
}

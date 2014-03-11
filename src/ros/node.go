package ros

import (
    "fmt"
    "math/rand"
    "net"
    "net/http"
    "os"
    "os/signal"
    "reflect"
    "strconv"
    "sync"
    "time"
    "xmlrpc"
)

const (
    ApiStatusError   = -1
    ApiStatusFailure = 0
    ApiStatusSuccess = 1
)

// *defaultNode implements Node interface
// a defaultNode instance must be accessed in user goroutine.
type defaultNode struct {
    qualifiedName  string
    masterUri      string
    xmlrpcUri      string
    xmlrpcListener net.Listener
    xmlrpcHandler  *xmlrpc.Handler
    subscribers    map[string]*defaultSubscriber
    publishers     map[string]*defaultPublisher
    servers        map[string]*defaultServiceServer
    jobChan        chan func()
    interruptChan  chan os.Signal
    logger         Logger
    ok             bool
    okMutex        sync.RWMutex
    waitGroup      sync.WaitGroup
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
        } else {
            numTrial += 1
        }
    }
    return nil, fmt.Errorf("listenRandomPort exceeds trial limit.")
}

func newDefaultNode(name string) *defaultNode {
    node := new(defaultNode)
    node.qualifiedName = name
    node.subscribers = make(map[string]*defaultSubscriber)
    node.publishers = make(map[string]*defaultPublisher)
    node.servers = make(map[string]*defaultServiceServer)
    node.interruptChan = make(chan os.Signal)
    node.ok = true

    logger := NewDefaultLogger()
    node.logger = logger

    // Install signal handler
    signal.Notify(node.interruptChan, os.Interrupt)
    go func() {
        <-node.interruptChan
        logger.Info("Interrupted")
        node.okMutex.Lock()
        node.ok = false
        node.okMutex.Unlock()
    }()

    node.jobChan = make(chan func(), 100)

    node.masterUri = os.Getenv("ROS_MASTER_URI")
    logger.Debugf("Master URI = %s", node.masterUri)

    listener, err := listenRandomPort("127.0.0.1", 10)
    if err != nil {
        logger.Fatal(err)
    }
    node.xmlrpcUri = fmt.Sprintf("http://%s", listener.Addr().String())
    node.xmlrpcListener = listener
    m := map[string]xmlrpc.Method{
        "getBusStats":      func(callerId string) (interface{}, error) { return node.getBusStats(callerId) },
        "getBusInfo":       func(callerId string) (interface{}, error) { return node.getBusInfo(callerId) },
        "getMasterUri":     func(callerId string) (interface{}, error) { return node.getMasterUri(callerId) },
        "shutdown":         func(callerId string, msg string) (interface{}, error) { return node.shutdown(callerId, msg) },
        "getPid":           func(callerId string) (interface{}, error) { return node.getPid(callerId) },
        "getSubscriptions": func(callerId string) (interface{}, error) { return node.getSubscriptions(callerId) },
        "getPublications":  func(callerId string) (interface{}, error) { return node.getPublications(callerId) },
        "paramUpdate":      func(callerId string, key string, value interface{}) (interface{}, error) { return node.paramUpdate(callerId, key ,value) },
        "publisherUpdate":  func(callerId string, topic string, publishers []interface{}) (interface{}, error) { return node.publisherUpdate(callerId, topic, publishers) },
        "requestTopic":     func(callerId string, topic string, protocols []interface{}) (interface{}, error) { return node.requestTopic(callerId, topic, protocols) },
    }
    node.xmlrpcHandler = xmlrpc.NewHandler(m)
    go http.Serve(node.xmlrpcListener, node.xmlrpcHandler)
    logger.Debugf("Started %s", node.qualifiedName)
    return node
}


func (node *defaultNode) OK() bool {
    node.okMutex.RLock()
    ok := node.ok
    node.okMutex.RUnlock()
    return ok
}


func (node *defaultNode) getBusStats(callerId string) (interface{}, error) {
    return buildRosApiResult(-1, "Not implemented", 0), nil
}

func (node *defaultNode) getBusInfo(callerId string) (interface{}, error) {
    return buildRosApiResult(-1, "Not implemeted", 0), nil
}

func (node *defaultNode) getMasterUri(callerId string) (interface{}, error) {
    return buildRosApiResult(0, "Success", node.masterUri), nil
}

func (node *defaultNode) shutdown(callerId string, msg string) (interface{}, error) {
    node.okMutex.Lock()
    node.ok = false
    node.okMutex.Unlock()
    return buildRosApiResult(0, "Success", 0), nil
}

func (node *defaultNode) getPid(callerId string) (interface{}, error) {
    return buildRosApiResult(0, "Success", os.Getpid()), nil
}

func (node *defaultNode) getSubscriptions(callerId string) (interface{}, error) {
    result := []interface{}{}
    for t, s := range node.subscribers {
        pair := []interface{}{t, s.msgType.Name()}
        result = append(result, pair)
    }
    return buildRosApiResult(0, "Success", result), nil
}

func (node *defaultNode) getPublications(callerId string) (interface{}, error) {
    result := []interface{}{}
    for t, p := range node.publishers {
        pair := []interface{}{t, p.msgType.Name()}
        result = append(result, pair)
    }
    return buildRosApiResult(0, "Success", result), nil
}

func (node *defaultNode) paramUpdate(callerId string, key string, value interface{}) (interface{}, error) {
    return buildRosApiResult(-1, "Not implemented", 0), nil
}

func (node *defaultNode) publisherUpdate(callerId string, topic string, publishers []interface{}) (interface{}, error) {
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
    return buildRosApiResult(code, message, 0), nil
}


func (node *defaultNode) requestTopic(callerId string, topic string, protocols []interface{}) (interface{}, error) {
    node.logger.Debugf("Slave API requestTopic(%s, %s, ...) called.", callerId, topic)
    var code int32
    var message string
    var value interface{}
    if pub, ok := node.publishers[topic]; !ok {
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
                host, portStr := pub.hostAndPort()
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
    return buildRosApiResult(code, message, value), nil
}


func (node *defaultNode) NewPublisher(topic string, msgType MessageType) Publisher {
    pub, ok := node.publishers[topic]
    logger := node.logger
    if !ok {
        _, err := callRosApi(node.masterUri, "registerPublisher",
            node.qualifiedName,
            topic, msgType.Name(),
            node.xmlrpcUri)
        if err != nil {
            logger.Fatalf("Failed to call registerPublisher(): %s", err)
        }

        pub = newDefaultPublisher(logger, node.qualifiedName, node.xmlrpcUri, node.masterUri, topic, msgType)
        node.publishers[topic] = pub
        go pub.start(&node.waitGroup)
    }
    return pub
}

func (node *defaultNode) NewSubscriber(topic string, msgType MessageType, callback interface{}) Subscriber {
    sub, ok := node.subscribers[topic]
    logger := node.logger
    if !ok {
        node.logger.Debug("Call Master API registerSubscriber")
        result, err := callRosApi(node.masterUri, "registerSubscriber",
            node.qualifiedName,
            topic,
            msgType.Name(),
            node.xmlrpcUri)
        if err != nil {
            logger.Fatalf("Failed to call registerSubscriber() for %s.", err)
        }
        list, ok := result.([]interface{})
        if !ok {
            logger.Fatalf("result is not []string but %s.", reflect.TypeOf(result).String())
        }
        var publishers []string
        for _, item := range list {
            s, ok := item.(string)
            if !ok {
                logger.Fatal("Publisher list contains no string object")
            }
            publishers = append(publishers, s)
        }

        logger.Debugf("Publisher URI list: ", publishers)

        sub := newDefaultSubscriber(topic, msgType, callback)
        node.subscribers[topic] = sub

        logger.Debugf("Start subscriber goroutine for topic '%s'", sub.topic)
        go sub.start(&node.waitGroup, node.masterUri, node.qualifiedName, node.xmlrpcUri, node.jobChan, logger)
        logger.Debugf("Done")
        sub.pubListChan <- publishers
        logger.Debugf("Update publisher list for topic '%s'", topic)
    } else {
        sub.callbacks = append(sub.callbacks, callback)
    }
    return sub
}


func (node *defaultNode) NewServiceClient(service string, srvType ServiceType) ServiceClient {
    client := newDefaultServiceClient(node.logger, node.qualifiedName, node.masterUri, service, srvType)
    return client
}


func (node *defaultNode) NewServiceServer(service string, srvType ServiceType, handler interface{}) ServiceServer {
    server, ok := node.servers[service]
    if ok {
        server.Shutdown()
    }
    server = newDefaultServiceServer(node, service, srvType, handler)
    if server == nil {
        return nil
    }
    node.servers[service] = server
    return server
}

func (node *defaultNode) SpinOnce() {
    timeoutChan := time.After(10 * time.Millisecond)
    select {
    case job := <-node.jobChan:
        job()
    case <-timeoutChan:
       break
    }
}

func (node *defaultNode) Spin() {
    logger := node.logger
    for node.OK() {
        timeoutChan := time.After(10 * time.Millisecond)
        select {
        case job := <-node.jobChan:
            logger.Debug("Execute job")
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
    for _, p := range node.publishers {
        p.Shutdown()
    }
    node.logger.Debug("Shutdown publishers...done")
    node.logger.Debug("Shutdown servers")
    for _, s := range node.servers {
        s.Shutdown()
    }
    node.logger.Debug("Shutdown servers...done")
    node.logger.Debug("Close XMLRPC lisetner")
    node.xmlrpcListener.Close()
    node.logger.Debug("Close XMLRPC done")
    node.logger.Debug("Wait XMLRPC server shutdown")
    node.xmlrpcHandler.WaitForShutdown()
    node.logger.Debug("Wait XMLRPC server shutdown...Done")
    node.logger.Debug("Wait all goroutines")
    node.waitGroup.Wait()
    node.logger.Debug("Wait all goroutines...Done")
    node.logger.Debug("Shutting node down completed")
    return
}

func (node *defaultNode) GetParam(key string) (interface{}, error) {
    return callRosApi(node.masterUri, "getParam", node.qualifiedName, key)
}

func (node *defaultNode) SetParam(key string, value interface{}) error {
    _, e := callRosApi(node.masterUri, "setParam", node.qualifiedName, key, value)
    return e
}

func (node *defaultNode) HasParam(key string) (bool, error) {
    result, e := callRosApi(node.masterUri, "hasParam", node.qualifiedName, key)
    hasParam := result.(bool)
    return hasParam, e
}

func (node *defaultNode) SearchParam(key string) (string, error) {
    result, e := callRosApi(node.masterUri, "searchParam", node.qualifiedName, key)
    foundKey := result.(string)
    return foundKey, e
}

func (node *defaultNode) DeleteParam(key string) error {
    _, e := callRosApi(node.masterUri, "deleteParam", node.qualifiedName, key)
    return e
}

func (node *defaultNode) Logger() Logger {
    return node.logger
}


package ros

import (
	"time"
)

type Node interface {
	NewPublisher(topic string, msgType MessageType) Publisher
	// Create a publisher which gives you callbacks when subscribers
	// connect and disconnect.  The callbacks are called in their own
	// goroutines, so they don't need to return immediately to let the
	// connection proceed.
	NewPublisherWithCallbacks(topic string,
		msgType MessageType,
		connectCallback, disconnectCallback func(SingleSubscriberPublisher)) Publisher
	// callback should be a function which takes 0, 1, or 2 arguments.
	// If it takes 0 arguments, it will simply be called without the
	// message.  1-argument functions are the normal case, and the
	// argument should be of the generated message type.  If the
	// function takes 2 arguments, the first argument should be of the
	// generated message type and the second argument should be of
	// type MessageEvent.
	NewSubscriber(topic string, msgType MessageType, callback interface{}) Subscriber
	NewServiceClient(service string, srvType ServiceType) ServiceClient
	NewServiceServer(service string, srvType ServiceType, callback interface{}) ServiceServer

	OK() bool
	SpinOnce()
	Spin()
	Shutdown()

	GetParam(name string) (interface{}, error)
	SetParam(name string, value interface{}) error
	HasParam(name string) (bool, error)
	SearchParam(name string) (string, error)
	DeleteParam(name string) error

	Logger() Logger

	NonRosArgs() []string
}

func NewNode(name string, args []string) (Node, error) {
	return newDefaultNode(name, args)
}

type Publisher interface {
	Publish(msg Message)
	Shutdown()
}

// A publisher which only sends to one specific subscriber.  This is
// sent as an argument to the connect and disconnect callback
// functions passed to Node.NewPublisherWithCallbacks().
type SingleSubscriberPublisher interface {
	Publish(msg Message)
	GetSubscriberName() string
	GetTopic() string
}

type Subscriber interface {
	GetNumPublishers() int
	Shutdown()
}

// Optional second argument to a Subscriber callback.
type MessageEvent struct {
	PublisherName    string
	ReceiptTime      time.Time
	ConnectionHeader map[string]string
}

type ServiceHandler interface{}

type ServiceFactory interface {
	Name() string
	MD5Sum() string
}

type ServiceServer interface {
	Shutdown()
}

type ServiceClient interface {
	Call(srv Service) error
	Shutdown()
}

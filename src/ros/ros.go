package ros

type Node interface {
    NewPublisher(topic string, msgType MessageType) Publisher
    // Create a publisher which gives you callbacks when subscribers
    // connect and disconnect.  The callbacks are called in their own
    // goroutines, so they don't need to return immediately to let the
    // connection proceed.
    NewPublisherWithCallbacks(topic string,
        msgType MessageType,
        connectCallback, disconnectCallback func(SingleSubscriberPublisher)) Publisher
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
}

func NewNode(name string) Node {
    return newDefaultNode(name)
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
    Shutdown()
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

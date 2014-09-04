package ros

type Node interface {
	NewPublisher(topic string, msgType MessageType) Publisher
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

package ros

type ServiceType interface {
	MD5Sum() string
	Name() string
	RequestType() MessageType
	ResponseType() MessageType
	NewService() Service
}

type Service interface {
	ReqMessage() Message
	ResMessage() Message
}

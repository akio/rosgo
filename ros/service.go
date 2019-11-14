package ros

//ServiceType is the interface definition of a ROS Service
//This contains MD5sum and Name of service and the request and response message types
//NewService iinstantiates a new Service object
type ServiceType interface {
	MD5Sum() string
	Name() string
	RequestType() MessageType
	ResponseType() MessageType
	NewService() Service
}

//Service interface contains the Request and Response ROS messages
type Service interface {
	ReqMessage() Message
	ResMessage() Message
}

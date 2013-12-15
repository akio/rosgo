package ros

type ServiceType interface {
    MD5Sum() string
    Name() string
    RequestType() MessageType
    ResponseType() MessageType
}


type Service interface {
    Request() Message
    Response() Message
}

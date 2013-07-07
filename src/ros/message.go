package ros

type MessageType interface {
    Definition() string
    MD5Sum() string
    Name() string
    NewMessage() Message
}

type Message interface {
    Serialize() []byte
    Deserialize(buffer []byte) error
}

type ServiceType interface {
    MD5Sum() string
    Name() string
    RequestType() MessageType
    ResponseType() MessageType
}

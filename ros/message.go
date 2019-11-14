package ros

import (
	"bytes"
)

//MessageType struct which contains the interface functions for the important properties of a message
//Text, Name, and MD5 sum are data used to identify message types while NewMessage instantiates the message fields
type MessageType interface {
	Text() string
	MD5Sum() string
	Name() string
	NewMessage() Message
}

//Message struct which contains the message type, serialize and deserialize functions
type Message interface {
	Type() MessageType
	Serialize(buf *bytes.Buffer) error
	Deserialize(buf *bytes.Reader) error
}

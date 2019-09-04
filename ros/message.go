package ros

import (
	"bytes"
)

type MessageType interface {
	Text() string
	MD5Sum() string
	Name() string
	NewMessage() Message
}

type Message interface {
	Type() MessageType
	Serialize(buf *bytes.Buffer) error
	Deserialize(buf *bytes.Reader) error
}

type GenericMessageType struct {
	ros_type string
	fields   string
}

func (m GenericMessageType) Fields() string {
	return m.fields
}
func (m *GenericMessageType) SetFields(fields string) {
	m.fields = fields
}

func (m *GenericMessageType) SetMessageType(ros_type string) {
	m.ros_type = ros_type
}

func (m GenericMessageType) Name() string {
	return m.ros_type
}
func (m GenericMessageType) Text() string {
	return "Test!"
}
func (m GenericMessageType) MD5Sum() string {
	return "992ce8a1687cec8c8bd883ec73ca41d1"
}
func (m GenericMessageType) NewMessage() Message {
	return nil
}

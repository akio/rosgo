package ros

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/edwinhayes/rosgo/libgengo"
	"io"
	"os"
	"strings"
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
	checksum string
	fields   string
}

type GenericMessage struct {
	generic_type *GenericMessageType
	data         []interface{}
}

func (m GenericMessageType) Fields() string {
	return m.fields
}
func (m *GenericMessageType) SetFields(fields string) {
	m.fields = fields
}

func (m *GenericMessageType) SetMessageType(ros_type string) error {

	// Create context for our ROS install.
	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")
	context, err := libgengo.NewMsgContext(strings.Split(rosPkgPath, ":"))
	if err != nil {
		return err
	}

	// Load context for the target message.
	spec, err := context.LoadMsg(ros_type)
	if err != nil {
		return err
	}

	// Now we know the checksum!
	m.checksum = spec.MD5Sum
	fmt.Println(spec)
	fmt.Println(spec.MD5Sum)

	m.ros_type = ros_type
	return nil
}

func (m GenericMessageType) Name() string {
	return m.ros_type
}
func (m GenericMessageType) Text() string {
	return "Test!"
}
func (m GenericMessageType) MD5Sum() string {
	return m.checksum
}
func (m GenericMessageType) NewMessage() Message {
	instance := new(GenericMessage)
	instance.generic_type = &m
	return instance
}

func (m GenericMessage) Type() MessageType {
	return m.generic_type
}
func (m GenericMessage) Serialize(buf *bytes.Buffer) error {
	return errors.New("Not implemented.")
}
func (m *GenericMessage) Deserialize(buf *bytes.Reader) error {
	// Read from the buffer in chunks.
	chunk := make([]byte, 8)
	for {
		n, err := buf.Read(chunk)
		// Check whether we're at EOF or not.
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("Some mystery problem occurred while reading.")
		}
		// Else, handle the data byte by byte.
		for _, b := range chunk[:n] {
			// TODO - For now, we'll just store the data as individual bytes.
			m.data = append(m.data, b)
		}
	}

	// Not all done, since defer?
	return nil
}
func (m GenericMessage) String() string {
	// Just print out the data!
	return fmt.Sprint(m.data)
}

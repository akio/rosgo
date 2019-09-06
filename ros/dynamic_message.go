package ros

// IMPORT REQUIRED PACKAGES.

// TODO - Why is the syntax for import different to everywhere else?

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/edwinhayes/rosgo/libgengo"
	"os"
	"strings"
)

// DEFINE PUBLIC STRUCTURES.

type DynamicMessageType struct {
	spec *libgengo.MsgSpec
}

type DynamicMessage struct {
	dynamic_type *DynamicMessageType
	data         map[string]interface{}
}

// DEFINE PRIVATE STRUCTURES.

// DEFINE PUBLIC GLOBALS.

// DEFINE PRIVATE GLOBALS.

var context *libgengo.MsgContext // We'll try to preserve a single message context to avoid reloading each time.

// DEFINE PUBLIC STATIC FUNCTIONS.

func NewDynamicMessageType(ros_type string) (*DynamicMessageType, error) {
	// Create an empty message type.
	m := new(DynamicMessageType)

	// If we haven't created a message context yet, better do that.
	if context == nil {
		// Create context for our ROS install.
		rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")
		c, err := libgengo.NewMsgContext(strings.Split(rosPkgPath, ":"))
		if err != nil {
			return nil, err
		}
		context = c
	}

	// Load context for the target message.
	spec, err := context.LoadMsg(ros_type)
	if err != nil {
		return nil, err
	}

	// Now we know all about the message!
	m.spec = spec

	// We've successfully made a new message type matching the requested ROS type.
	return m, nil
}

// DEFINE PUBLIC RECEIVER FUNCTIONS.

//	DynamicMessageType

func (m DynamicMessageType) Name() string {
	return m.spec.FullName
}

func (m DynamicMessageType) Text() string {
	return m.spec.Text
}

func (m DynamicMessageType) MD5Sum() string {
	return m.spec.MD5Sum
}

func (m DynamicMessageType) NewMessage() Message {
	// Don't instantiate messages for incomplete types.
	if m.spec == nil {
		return nil
	}

	// But otherwise, make a new one.
	d := new(DynamicMessage)
	d.dynamic_type = &m
	return d
}

//	DynamicMessage

func (m DynamicMessage) Type() MessageType {
	return m.dynamic_type
}

func (m DynamicMessage) Serialize(buf *bytes.Buffer) error {
	return errors.New("Not implemented.")
}

func (m *DynamicMessage) Deserialize(buf *bytes.Reader) error {
	// THIS METHOD IS BASICALLY AN UNTEMPLATED COPY OF THE TEMPLATE IN LIBGENGO.

	var err error = nil
	m.data = make(map[string]interface{})
	// Iterate over each of the fields in the message.
	for _, field := range m.dynamic_type.spec.Fields {
		if field.IsArray {
			// The array starts with a declaration of the number of array elements.
			var size uint32
			if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
				return err
			}
			// Create an array of the target type.
			if field.ArrayLen < 0 {
				switch field.GoType {
				case "int8":
					m.data[field.GoName] = make([]int8, 0)
				case "string":
					m.data[field.GoName] = make([]string, 0)
				case "time":
					m.data[field.GoName] = make([]Time, 0)
				case "duration":
					m.data[field.GoName] = make([]Duration, 0)
				default:
					// In this case, it will probably be because the go_type is describing another ROS message, so we need to replace that with a nested DynamicMessage.
					m.data[field.GoName] = make([]Message, 0)
				}
			}
			for i := 0; i < int(size); i++ {
				if field.IsBuiltin {
					if field.Type == "string" {
						// The string will start with a declaration of the number of characters.
						var str_size uint32
						if err = binary.Read(buf, binary.LittleEndian, &str_size); err != nil {
							return err
						}
						data := make([]byte, int(str_size))
						if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
							return err
						}
						m.data[field.GoName] = append(m.data[field.GoName].([]string), string(data))
					} else if field.Type == "time" {
						var data Time
						// Time/duration types have two fields, so consume into these in two reads.
						if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
							return err
						}
						if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
							return err
						}
						m.data[field.GoName] = append(m.data[field.GoName].([]Time), data)
					} else if field.Type == "duration" {
						var data Duration
						// Time/duration types have two fields, so consume into these in two reads.
						if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
							return err
						}
						if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
							return err
						}
						m.data[field.GoName] = append(m.data[field.GoName].([]Duration), data)
					} else {
						// It's a regular primitive element.
						data := instantiate_scalar_type(field.GoType)
						if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
							return err
						}
						// Because not runtime expressions in type assertions in Go, we need to manually do this.
						switch field.GoType {
						case "int8":
							m.data[field.GoName] = append(m.data[field.GoName].([]int8), data.(int8))
						default:
							// Something went wrong.
							return errors.New("We haven't implemented this primitive yet!")
						}
					}
				} else {
					// Else it's not a builtin.
					msg_type, err := NewDynamicMessageType(field.Type)
					if err != nil {
						return err
					}
					msg := msg_type.NewMessage()
					if err = msg.Deserialize(buf); err != nil {
						return err
					}
					m.data[field.GoName] = append(m.data[field.GoName].([]Message), msg)
				}
			}
		}
		// Else it's not an array.  This is just the same as above, with the '[i]' bits removed.
		m.data[field.GoName] = instantiate_scalar_type(field.GoType)
		if field.IsBuiltin {
			if field.Type == "string" {
				// The string will start with a declaration of the number of characters.
				var str_size uint32
				if err = binary.Read(buf, binary.LittleEndian, &str_size); err != nil {
					return err
				}
				data := make([]byte, int(str_size))
				if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
					return err
				}
				m.data[field.GoName] = string(data)
			} else if field.Type == "time" {
				var data Time
				// Time/duration types have two fields, so consume into these in two reads.
				if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
					return err
				}
				if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
					return err
				}
				m.data[field.GoName] = data
			} else if field.Type == "duration" {
				var data Duration
				// Time/duration types have two fields, so consume into these in two reads.
				if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
					return err
				}
				if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
					return err
				}
				m.data[field.GoName] = data
			} else {
				// It's a regular primitive element.
				data := instantiate_scalar_type(field.GoType)
				if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
					return err
				}
				m.data[field.GoName] = data
			}
		} else {
			// Else it's not a builtin.
			msg_type, err := NewDynamicMessageType(field.Type)
			if err != nil {
				return err
			}
			m.data[field.GoName] = msg_type.NewMessage()
			if err = m.data[field.GoName].(Message).Deserialize(buf); err != nil {
				return err
			}
		}
	}

	// All done.
	return err
}

func (m DynamicMessage) String() string {
	// Just print out the data!
	return fmt.Sprint(m.data)
}

// DEFINE PRIVATE STATIC FUNCTIONS.

func instantiate_scalar_type(go_type string) interface{} {
	switch go_type {
	case "int8":
		var v int8
		return v
	case "string":
		var v string
		return v
	case "time":
		var v Time
		return v
	case "duration":
		var v Duration
		return v
	default:
		// In this case, it will probably be because the go_type is describing another ROS message, so we need to replace that with a nested DynamicMessage.
		var v Message
		return v
	}
}

// DEFINE PRIVATE RECEIVER FUNCTIONS.

// ALL DONE.

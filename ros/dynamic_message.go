package ros

// IMPORT REQUIRED PACKAGES.

// TODO - Why is the syntax for import different to everywhere else?

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/edwinhayes/rosgo/libgengo"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"reflect"
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

var known_messages map[string]string // Just for diagnostic purposes.

var context *libgengo.MsgContext // We'll try to preserve a single message context to avoid reloading each time.

// DEFINE PUBLIC STATIC FUNCTIONS.

func NewDynamicMessageType(ros_type string) (*DynamicMessageType, error) {
	return newDynamicMessageType_Nested(ros_type, "")
}

func newDynamicMessageType_Nested(ros_type string, parent_type string) (*DynamicMessageType, error) {
	// Create an empty message type.
	m := new(DynamicMessageType)

	// If we haven't created a message context yet, better do that.
	if context == nil {
		// Create context for our ROS install.
		rosPkgPath := os.Getenv("ROS_PACKAGE_PATH") + ":~/environment/goenv/src/github.com/edwinhayes/rosgo/test"
		c, err := libgengo.NewMsgContext(strings.Split(rosPkgPath, ":"))
		if err != nil {
			return nil, err
		}
		context = c
	}
	if known_messages == nil {
		known_messages = make(map[string]string)
	}

	// We need to try to look up the full name, in case we've just been given a short name.
	fullname := ros_type

	// The Header type has some special treatment.
	if ros_type == "Header" {
		fullname = "std_msgs/Header"
	} else {
		_, ok := context.GetMsgs()[fullname]
		if !ok {
			// Seems like the ros_type we were give wasn't the full name.

			// Message in the same package are allowed to use relative names, so try using the parent's full name.
			if parent_type != "" {
				pkgName := filepath.Base(parent_type)
				fullname = pkgName + "/" + fullname
			}
		}
	}

	// Load context for the target message.
	spec, err := context.LoadMsg(fullname)
	if err != nil {
		return nil, err
	}

	// Now we know all about the message!
	m.spec = spec
	//fmt.Println(spec)

	// We've successfully made a new message type matching the requested ROS type.
	known_messages[m.spec.FullName] = m.spec.MD5Sum
	return m, nil
}

func GetKnownMsgs() map[string]string {
	return known_messages
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
	// THIS METHOD IS BASICALLY AN UNTEMPLATED COPY OF THE TEMPLATE IN LIBGENGO.

	var err error = nil

	// Iterate over each of the fields in the message.
	for _, field := range m.dynamic_type.spec.Fields {
		if field.IsArray {
			// It's an array.

			// Look up the item.
			array, ok := m.data[field.GoName]
			if !ok {
				return errors.New("Field: " + field.Name + ": No data found.")
			}

			// If the array is not a fixed length, it begins with a declaration of the array size.
			if field.ArrayLen < 0 {
				var size uint32 = uint32(reflect.ValueOf(array).Len())
				if err := binary.Write(buf, binary.LittleEndian, size); err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
			}

			// Then we just write out all the elements one after another.
			for _, array_item := range array.([]interface{}) {
				// Need to handle each type appropriately.
				if field.IsBuiltin {
					if field.Type == "string" {
						// Make sure we've actually got a string.
						str, ok := array_item.(string)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected string.")
						}
						// The string should start with a declaration of the number of characters.
						var size_str uint32 = uint32(len(str))
						if err := binary.Write(buf, binary.LittleEndian, size_str); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						// Then write out the actual characters.
						data := []byte(str)
						if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}

					} else if field.Type == "time" {
						// Make sure we've actually got a time.
						t, ok := array_item.(Time)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected ros.Time.")
						}
						// Then write out the structure.
						if err := binary.Write(buf, binary.LittleEndian, t); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}

					} else if field.Type == "duration" {
						// Make sure we've actually got a duration.
						d, ok := array_item.(Duration)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected ros.Duration.")
						}
						// Then write out the structure.
						if err := binary.Write(buf, binary.LittleEndian, d); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}

					} else {
						// It's a primitive.

						// Because no runtime expressions in type assertions in Go, we need to manually do this.
						switch field.GoType {
						case "bool":
							// Make sure we've actually got a bool.
							v, ok := array_item.(bool)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected bool.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int8":
							// Make sure we've actually got a bool.
							v, ok := array_item.(int8)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected int8.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int16":
							// Make sure we've actually got a bool.
							v, ok := array_item.(int16)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected int16.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int32":
							// Make sure we've actually got a bool.
							v, ok := array_item.(int32)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected int32.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int64":
							// Make sure we've actually got a bool.
							v, ok := array_item.(int64)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected int64.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint8":
							// Make sure we've actually got a bool.
							v, ok := array_item.(uint8)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected uint8.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint16":
							// Make sure we've actually got a bool.
							v, ok := array_item.(uint16)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected uint16.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint32":
							// Make sure we've actually got a bool.
							v, ok := array_item.(uint32)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected uint32.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint64":
							// Make sure we've actually got a bool.
							v, ok := array_item.(uint64)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected uint64.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "float32":
							// Make sure we've actually got a bool.
							v, ok := array_item.(float32)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected float32.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "float64":
							// Make sure we've actually got a bool.
							v, ok := array_item.(float64)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected float64.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						default:
							// Something went wrong.
							return errors.New("We haven't implemented this primitive yet!")
						}
					}

				} else {
					// Else it's not a builtin.

					// Confirm the message we're holding is actually the correct type.
					msg, ok := array_item.(DynamicMessage)
					if !ok {
						return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(array_item).Name() + ", expected Message.")
					}
					if msg.dynamic_type.spec.FullName != field.Type {
						return errors.New("Field: " + field.Name + ": Found msg " + msg.dynamic_type.spec.FullName + ", expected " + field.Type + ".")
					}
					// Otherwise, we just recursively serialise it.
					if err = msg.Serialize(buf); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
				}
			}

		} else {
			// It's a scalar.

			// Look up the item.
			item, ok := m.data[field.GoName]
			if !ok {
				return errors.New("Field: " + field.Name + ": No data found.")
			}

			// Need to handle each type appropriately.
			if field.IsBuiltin {
				if field.Type == "string" {
					// Make sure we've actually got a string.
					str, ok := item.(string)
					if !ok {
						return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected string.")
					}
					// The string should start with a declaration of the number of characters.
					var size_str uint32 = uint32(len(str))
					if err := binary.Write(buf, binary.LittleEndian, size_str); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					// Then write out the actual characters.
					data := []byte(str)
					if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}

				} else if field.Type == "time" {
					// Make sure we've actually got a time.
					t, ok := item.(Time)
					if !ok {
						return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected ros.Time.")
					}
					// Then write out the structure.
					if err := binary.Write(buf, binary.LittleEndian, t); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}

				} else if field.Type == "duration" {
					// Make sure we've actually got a duration.
					d, ok := item.(Duration)
					if !ok {
						return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected ros.Duration.")
					}
					// Then write out the structure.
					if err := binary.Write(buf, binary.LittleEndian, d); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}

				} else {
					// It's a primitive.

					// Because no runtime expressions in type assertions in Go, we need to manually do this.
					switch field.GoType {
					case "bool":
						// Make sure we've actually got a bool.
						v, ok := item.(bool)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected bool.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "int8":
						// Make sure we've actually got a bool.
						v, ok := item.(int8)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected int8.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "int16":
						// Make sure we've actually got a bool.
						v, ok := item.(int16)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected int16.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "int32":
						// Make sure we've actually got a bool.
						v, ok := item.(int32)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected int32.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "int64":
						// Make sure we've actually got a bool.
						v, ok := item.(int64)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected int64.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "uint8":
						// Make sure we've actually got a bool.
						v, ok := item.(uint8)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected uint8.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "uint16":
						// Make sure we've actually got a bool.
						v, ok := item.(uint16)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected uint16.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "uint32":
						// Make sure we've actually got a bool.
						v, ok := item.(uint32)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected uint32.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "uint64":
						// Make sure we've actually got a bool.
						v, ok := item.(uint64)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected uint64.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "float32":
						// Make sure we've actually got a bool.
						v, ok := item.(float32)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected float32.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					case "float64":
						// Make sure we've actually got a bool.
						v, ok := item.(float64)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected float64.")
						}
						// Then write out the value.
						if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					default:
						// Something went wrong.
						return errors.New("We haven't implemented this primitive yet!")
					}
				}

			} else {
				// Else it's not a builtin.

				// Confirm the message we're holding is actually the correct type.
				msg, ok := item.(DynamicMessage)
				if !ok {
					return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected Message.")
				}
				if msg.dynamic_type.spec.FullName != field.Type {
					return errors.New("Field: " + field.Name + ": Found msg " + msg.dynamic_type.spec.FullName + ", expected " + field.Type + ".")
				}
				// Otherwise, we just recursively serialise it.
				if err = msg.Serialize(buf); err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
			}
		}
	}

	// All done.
	return err
}

func (m *DynamicMessage) Deserialize(buf *bytes.Reader) error {
	// THIS METHOD IS BASICALLY AN UNTEMPLATED COPY OF THE TEMPLATE IN LIBGENGO.

	// To give more sane results in the event of a decoding issue, we decode into a copy of the data field.
	var err error = nil
	tmp_data := make(map[string]interface{})
	m.data = nil

	// Iterate over each of the fields in the message.
	for _, field := range m.dynamic_type.spec.Fields {
		if field.IsArray {
			// It's an array.

			// The array may be a fixed length, or it may be dynamic.
			var size uint32 = uint32(field.ArrayLen)
			if field.ArrayLen < 0 {
				// The array is dynamic, so it starts with a declaration of the number of array elements.
				if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
			}
			// Create an array of the target type.
			switch field.GoType {
			case "bool":
				tmp_data[field.GoName] = make([]bool, 0)
			case "int8":
				tmp_data[field.GoName] = make([]int8, 0)
			case "int16":
				tmp_data[field.GoName] = make([]int16, 0)
			case "int32":
				tmp_data[field.GoName] = make([]int32, 0)
			case "int64":
				tmp_data[field.GoName] = make([]int64, 0)
			case "uint8":
				tmp_data[field.GoName] = make([]uint8, 0)
			case "uint16":
				tmp_data[field.GoName] = make([]uint16, 0)
			case "uint32":
				tmp_data[field.GoName] = make([]uint32, 0)
			case "uint64":
				tmp_data[field.GoName] = make([]uint64, 0)
			case "float32":
				tmp_data[field.GoName] = make([]float32, 0)
			case "float64":
				tmp_data[field.GoName] = make([]float64, 0)
			case "string":
				tmp_data[field.GoName] = make([]string, 0)
			case "time":
				tmp_data[field.GoName] = make([]Time, 0)
			case "duration":
				tmp_data[field.GoName] = make([]Duration, 0)
			default:
				// In this case, it will probably be because the go_type is describing another ROS message, so we need to replace that with a nested DynamicMessage.
				tmp_data[field.GoName] = make([]Message, 0)
			}
			// Iterate over each item in the array.
			for i := 0; i < int(size); i++ {
				if field.IsBuiltin {
					if field.Type == "string" {
						// The string will start with a declaration of the number of characters.
						var str_size uint32
						if err = binary.Read(buf, binary.LittleEndian, &str_size); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						data := make([]byte, int(str_size))
						if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						tmp_data[field.GoName] = append(tmp_data[field.GoName].([]string), string(data))
					} else if field.Type == "time" {
						var data Time
						// Time/duration types have two fields, so consume into these in two reads.
						if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						tmp_data[field.GoName] = append(tmp_data[field.GoName].([]Time), data)
					} else if field.Type == "duration" {
						var data Duration
						// Time/duration types have two fields, so consume into these in two reads.
						if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						tmp_data[field.GoName] = append(tmp_data[field.GoName].([]Duration), data)
					} else {
						// It's a regular primitive element.

						// Because not runtime expressions in type assertions in Go, we need to manually do this.
						switch field.GoType {
						case "bool":
							var data bool
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]bool), data)
						case "int8":
							var data int8
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]int8), data)
						case "int16":
							var data int16
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]int16), data)
						case "int32":
							var data int32
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]int32), data)
						case "int64":
							var data int64
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]int64), data)
						case "uint8":
							var data uint8
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]uint8), data)
						case "uint16":
							var data uint16
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]uint16), data)
						case "uint32":
							var data uint32
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]uint32), data)
						case "uint64":
							var data uint64
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]uint64), data)
						case "float32":
							var data float32
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]float32), data)
						case "float64":
							var data float64
							binary.Read(buf, binary.LittleEndian, &data)
							tmp_data[field.GoName] = append(tmp_data[field.GoName].([]float64), data)
						default:
							// Something went wrong.
							return errors.New("We haven't implemented this primitive yet!")
						}
						if err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					}
				} else {
					// Else it's not a builtin.
					msg_type, err := newDynamicMessageType_Nested(field.Type, m.dynamic_type.spec.FullName)
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					msg := msg_type.NewMessage()
					if err = msg.Deserialize(buf); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmp_data[field.GoName] = append(tmp_data[field.GoName].([]Message), msg)
				}
			}
		} else {
			// Else it's a scalar.  This is just the same as above, with the '[i]' bits removed.

			if field.IsBuiltin {
				if field.Type == "string" {
					// The string will start with a declaration of the number of characters.
					var str_size uint32
					if err = binary.Read(buf, binary.LittleEndian, &str_size); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					data := make([]byte, int(str_size))
					if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmp_data[field.GoName] = string(data)
				} else if field.Type == "time" {
					var data Time
					// Time/duration types have two fields, so consume into these in two reads.
					if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmp_data[field.GoName] = data
				} else if field.Type == "duration" {
					var data Duration
					// Time/duration types have two fields, so consume into these in two reads.
					if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmp_data[field.GoName] = data
				} else {
					// It's a regular primitive element.
					switch field.GoType {
					case "bool":
						var data bool
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "int8":
						var data int8
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "int16":
						var data int16
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "int32":
						var data int32
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "int64":
						var data int64
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "uint8":
						var data uint8
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "uint16":
						var data uint16
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "uint32":
						var data uint32
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "uint64":
						var data uint64
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "float32":
						var data float32
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					case "float64":
						var data float64
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmp_data[field.GoName] = data
					default:
						// Something went wrong.
						return errors.New("We haven't implemented this primitive yet!")
					}
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
				}
			} else {
				// Else it's not a builtin.
				msg_type, err := newDynamicMessageType_Nested(field.Type, m.dynamic_type.spec.FullName)
				if err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
				tmp_data[field.GoName] = msg_type.NewMessage()
				if err = tmp_data[field.GoName].(Message).Deserialize(buf); err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
			}
		}
	}

	// All done.
	m.data = tmp_data
	return err
}

func (m DynamicMessage) String() string {
	// Just print out the data!
	return fmt.Sprint(m.data)
}

// DEFINE PRIVATE STATIC FUNCTIONS.

// DEFINE PRIVATE RECEIVER FUNCTIONS.

// ALL DONE.

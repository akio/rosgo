package ros

// IMPORT REQUIRED PACKAGES.

// TODO - Why is the syntax for import different to everywhere else?

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/edwinhayes/rosgo/libgengo"
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strings"
)

// DEFINE PUBLIC STRUCTURES.

// DynamicMessageType abstracts the schema of a ROS Message whose schema is only known at runtime.  DynamicMessageTypes are created by looking up the relevant schema information from
// ROS Message definition files.  DynamicMessageType implements the rosgo MessageType interface, allowing it to be used throughout rosgo in the same manner as message schemas generated
// at compiletime by gengo.
type DynamicMessageType struct {
	spec *libgengo.MsgSpec
}

// DynamicMessage abstracts an instance of a ROS Message whose type is only known at runtime.  The schema of the message is denoted by the referenced DynamicMessageType, while the
// actual payload of the Message is stored in a map[string]interface{} which maps the name of each field to its value.  DynamicMessage implements the rosgo Message interface, allowing
// it to be used throughout rosgo in the same manner as message types generated at compiletime by gengo.
type DynamicMessage struct {
	dynamicType *DynamicMessageType
	data        map[string]interface{}
}

// DEFINE PRIVATE STRUCTURES.

// DEFINE PUBLIC GLOBALS.

// DEFINE PRIVATE GLOBALS.

var rosPkgPath string // Colon separated list of paths to search for message definitions on.

// TODO - Probably remove knownMessages to reduce memory consumption.

var knownMessages map[string]string // Just for diagnostic purposes.

var context *libgengo.MsgContext // We'll try to preserve a single message context to avoid reloading each time.

// DEFINE PUBLIC STATIC FUNCTIONS.

// SetRuntimePackagePath sets the ROS package search path which will be used by DynamicMessage to look up ROS message definitions at runtime.
func SetRuntimePackagePath(path string) {
	// We're not going to check that the result is valid, we'll just accept it blindly.
	rosPkgPath = path

	// All done.
	return
}

// GetRuntimePackagePath returns the ROS package search path which will be used by DynamicMessage to look up ROS message definitions at runtime.  By default, this will
// be equivalent to the ${ROS_PACKAGE_PATH} environment variable.
func GetRuntimePackagePath() string {
	// If a package path hasn't been set at the time of first use, by default we'll just use the ROS evironment default.
	if rosPkgPath == "" {
		rosPkgPath = os.Getenv("ROS_PACKAGE_PATH")
	}
	// All done.
	return rosPkgPath
}

// NewDynamicMessageType generates a DynamicMessageType corresponding to the specified typeName from the available ROS message definitions; typeName should be a fully-qualified
// ROS message type name.  The first time the function is run, a message 'context' is created by searching through the available ROS message definitions, then the ROS message to
// be used for the definition is looked up by name.  On subsequent calls, the ROS message type is looked up directly from the existing context.
func NewDynamicMessageType(typeName string) (*DynamicMessageType, error) {
	return newDynamicMessageTypeNested(typeName, "")
}

// newDynamicMessageTypeNested generates a DynamicMessageType from the available ROS message definitions.  The first time the function is run, a message 'context' is created by
// searching through the available ROS message definitions, then the ROS message type to use for the defintion is looked up by name.  On subsequent calls, the ROS message type
// is looked up directly from the existing context.  This 'nested' version of the function is able to be called recursively, where packageName should be the typeName of the
// parent ROS message; this is used internally for handling complex ROS messages.
func newDynamicMessageTypeNested(typeName string, packageName string) (*DynamicMessageType, error) {
	// Create an empty message type.
	m := new(DynamicMessageType)

	// If we haven't created a message context yet, better do that.
	if context == nil {
		// Create context for our ROS install.
		c, err := libgengo.NewMsgContext(strings.Split(GetRuntimePackagePath(), ":"))
		if err != nil {
			return nil, err
		}
		context = c
	}
	if knownMessages == nil {
		knownMessages = make(map[string]string)
	}

	// We need to try to look up the full name, in case we've just been given a short name.
	fullname := typeName

	// The Header type has some special treatment.
	if typeName == "Header" {
		fullname = "std_msgs/Header"
	} else {
		_, ok := context.GetMsgs()[fullname]
		if !ok {
			// Seems like the package_name we were give wasn't the full name.

			// Messages in the same package are allowed to use relative names, so try prefixing the package.
			if packageName != "" {
				fullname = packageName + "/" + fullname
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

	// We've successfully made a new message type matching the requested ROS type.
	knownMessages[m.spec.FullName] = m.spec.MD5Sum
	return m, nil
}

// GetKnownMsgs returns a map whose keys are all the ROS fullnames for which a DynamicMessageType has been successfully
// created (and thus are definitely able to be sent/received), and whose values are the corresponding MD5 sums for those
// messages.
func GetKnownMsgs() map[string]string {
	return knownMessages
}

// DEFINE PUBLIC RECEIVER FUNCTIONS.

//	DynamicMessageType

// Name returns the full ROS name of the message type; required for ros.MessageType.
func (t DynamicMessageType) Name() string {
	return t.spec.FullName
}

// Text returns the full ROS message specification for this message type; required for ros.MessageType.
func (t DynamicMessageType) Text() string {
	return t.spec.Text
}

// MD5Sum returns the ROS compatible MD5 sum of the message type; required for ros.MessageType.
func (t DynamicMessageType) MD5Sum() string {
	return t.spec.MD5Sum
}

// NewMessage creates a new DynamicMessage instantiating the message type; required for ros.MessageType.
func (t DynamicMessageType) NewMessage() Message {
	// Don't instantiate messages for incomplete types.
	if t.spec == nil {
		return nil
	}
	// But otherwise, make a new one.
	d := new(DynamicMessage)
	d.dynamicType = &t
	var err error
	d.data, err = zeroValueData(t.Name())
	if err != nil {
		return nil
	}
	return d
}

//zeroValueData creates the zeroValue (default) data map for a new dynamic message
func zeroValueData(s string) (map[string]interface{}, error) {
	//Create map
	d := make(map[string]interface{})

	//Instantiate new dynamic message type from string name parsed
	t, err := NewDynamicMessageType(s)
	if err != nil {
		return d, errors.Wrap(err, "Failed to create NewDynamicMessageType "+s)
	}
	//Range fields in the dynamic message type
	for _, field := range t.spec.Fields {
		if field.IsArray {
			// It's an array.
			//Create slices for GoType arrays
			switch field.GoType {
			case "bool":
				d[field.GoName] = make([]bool, 0)
			case "int8":
				d[field.GoName] = make([]int8, 0)
			case "int16":
				d[field.GoName] = make([]int16, 0)
			case "int32":
				d[field.GoName] = make([]int32, 0)
			case "int64":
				d[field.GoName] = make([]int64, 0)
			case "uint8":
				d[field.GoName] = make([]uint8, 0)
			case "uint16":
				d[field.GoName] = make([]uint16, 0)
			case "uint32":
				d[field.GoName] = make([]uint32, 0)
			case "uint64":
				d[field.GoName] = make([]uint64, 0)
			case "float32":
				d[field.GoName] = make([]float32, 0)
			case "float64":
				d[field.GoName] = make([]float64, 0)
			case "string":
				d[field.GoName] = make([]string, 0)
			case "time":
				d[field.GoName] = make([]Time, 0)
			case "duration":
				d[field.GoName] = make([]Duration, 0)
			default:
				// In this case, it will probably be because the go_type is describing another ROS message, so we need to replace that with a nested DynamicMessage.
				d[field.GoName] = make([]Message, 0)
			}
			var size uint32 = uint32(field.ArrayLen)
			//In the case the array length is static, we iterated through array items
			if field.ArrayLen != -1 {
				for i := 0; i < int(size); i++ {
					if field.IsBuiltin {
						//Append the goType zeroValues to their arrays
						switch field.GoType {
						case "bool":
							d[field.GoName] = append(d[field.GoName].([]bool), false)
						case "int8":
							d[field.GoName] = append(d[field.GoName].([]int8), 0)
						case "int16":
							d[field.GoName] = append(d[field.GoName].([]int16), 0)
						case "int32":
							d[field.GoName] = append(d[field.GoName].([]int32), 0)
						case "int64":
							d[field.GoName] = append(d[field.GoName].([]int64), 0)
						case "uint8":
							d[field.GoName] = append(d[field.GoName].([]uint8), 0)
						case "uint16":
							d[field.GoName] = append(d[field.GoName].([]uint16), 0)
						case "uint32":
							d[field.GoName] = append(d[field.GoName].([]uint32), 0)
						case "uint64":
							d[field.GoName] = append(d[field.GoName].([]uint64), 0)
						case "float32":
							d[field.GoName] = append(d[field.GoName].([]float32), 0.0)
						case "float64":
							d[field.GoName] = append(d[field.GoName].([]float64), 0.)
						case "string":
							d[field.GoName] = append(d[field.GoName].([]string), "")
						case "ros.Time":
							d[field.GoName] = append(d[field.GoName].([]Time), Time{})
						case "ros.Duration":
							d[field.GoName] = append(d[field.GoName].([]Duration), Duration{})
						default:
							// Something went wrong.
							return d, errors.Wrap(err, "Builtin field "+field.GoType+" not found")
						}
					} else {
						// Else it's not a builtin. Create a nested message type for values inside
						t2, err := newDynamicMessageTypeNested(field.Type, field.Package)
						if err != nil {
							return d, errors.Wrap(err, "Failed to create newDynamicMessageTypeNested "+field.Type)
						}
						msg := t2.NewMessage()
						//Append nested message map to message type array in main map
						d[field.GoName] = append(d[field.GoName].([]Message), msg)
					}
					//Else array is dynamic, by default we do not initialize any values in it
				}
			}
		} else if field.IsBuiltin {
			//If its a built in type
			switch field.GoType {
			case "string":
				d[field.GoName] = ""
			case "bool":
				d[field.GoName] = bool(false)
			case "int8":
				d[field.GoName] = int8(0)
			case "int16":
				d[field.GoName] = int16(0)
			case "int32":
				d[field.GoName] = int32(0)
			case "int64":
				d[field.GoName] = int64(0)
			case "uint8":
				d[field.GoName] = uint8(0)
			case "uint16":
				d[field.GoName] = uint16(0)
			case "uint32":
				d[field.GoName] = uint32(0)
			case "uint64":
				d[field.GoName] = uint64(0)
			case "float32":
				d[field.GoName] = float32(0.0)
			case "float64":
				d[field.GoName] = float64(0.0)
			case "ros.Time":
				d[field.GoName] = Time{}
			case "ros.Duration":
				d[field.GoName] = Duration{}
			default:
				return d, errors.Wrap(err, "Builtin field "+field.GoType+" not found")
			}
			//Else its a ros message type
		} else {
			//Create new dynamic message type nested
			t2, err := newDynamicMessageTypeNested(field.Type, field.Package)
			if err != nil {
				return d, errors.Wrap(err, "Failed to create dewDynamicMessageTypeNested "+field.Type)
			}
			//Append message as a map item
			d[field.GoName] = t2.NewMessage()
		}
	}
	return d, err
}

// Data returns the data map field of the DynamicMessage
func (m DynamicMessage) Data() map[string]interface{} {
	return m.data
}

// GenerateJSONSchema generates a (primitive) JSON schema for the associated DynamicMessageType; however note that since
// we are mostly interested in making schema's for particular _topics_, the function takes a string prefix, and string topic name, which are
// used to id the resulting schema.
func (t DynamicMessageType) GenerateJSONSchema(prefix string, topic string) ([]byte, error) {
	// The JSON schema for a message consist of the (recursive) properties names/types:
	schemaItems, err := t.generateJSONSchemaProperties(prefix + topic)
	if err != nil {
		return nil, err
	}

	// Plus some extra keywords:
	schemaItems["$schema"] = "https://json-schema.org/draft-07/schema#"
	schemaItems["$id"] = prefix + topic

	// The schema itself is created from the map of properties.
	schemaString, err := json.Marshal(schemaItems)
	if err != nil {
		return nil, err
	}

	// All done.
	return schemaString, nil
}

func (t DynamicMessageType) generateJSONSchemaProperties(topic string) (map[string]interface{}, error) {
	// Each message's schema indicates that it is an 'object' with some nested properties: those properties are the fields and their types.
	properties := make(map[string]interface{})
	schemaItems := make(map[string]interface{})
	schemaItems["type"] = "object"
	schemaItems["title"] = topic
	schemaItems["properties"] = properties

	// Iterate over each of the fields in the message.
	for _, field := range t.spec.Fields {
		if field.IsArray {
			// It's an array.

			// Arrays all have a type of 'array', regardless of that the hold, then the 'item' keyword determines what type goes in the array.
			propertyContent := make(map[string]interface{})
			properties[field.GoName] = propertyContent
			propertyContent["type"] = "array"
			propertyContent["title"] = topic + Sep + field.GoName
			arrayItems := make(map[string]interface{})
			propertyContent["items"] = arrayItems

			// Need to handle each type appropriately.
			if field.IsBuiltin {
				if field.Type == "string" {
					arrayItems["type"] = "string"
				} else if field.Type == "time" {
					timeItems := make(map[string]interface{})
					timeItems["Sec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "Sec"}
					timeItems["NSec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "NSec"}
					arrayItems["type"] = "object"
					arrayItems["properties"] = timeItems
				} else if field.Type == "duration" {
					timeItems := make(map[string]interface{})
					timeItems["Sec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "Sec"}
					timeItems["NSec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "NSec"}
					arrayItems["type"] = "object"
					arrayItems["properties"] = timeItems
				} else {
					// It's a primitive.
					var jsonType string
					if field.GoType == "int8" || field.GoType == "uint8" || field.GoType == "int16" || field.GoType == "uint16" ||
						field.GoType == "int32" || field.GoType == "uint32" || field.GoType == "int64" || field.GoType == "uint64" {
						jsonType = "integer"
					} else if field.GoType == "float32" || field.GoType == "float64" {
						jsonType = "number"
					} else if field.GoType == "bool" {
						jsonType = "bool"
					} else {
						// Something went wrong.
						return nil, errors.New("we haven't implemented this primitive yet")
					}
					arrayItems["type"] = jsonType
				}
			} else {
				// It's another nested message.

				// Generate the nested type.
				msgType, err := newDynamicMessageTypeNested(field.Type, field.Package)
				if err != nil {
					return nil, errors.Wrap(err, "Schema Field: "+field.Name)
				}

				// Recursively generate schema information for the nested type.
				schemaElement, err := msgType.generateJSONSchemaProperties(topic + Sep + field.GoName)
				if err != nil {
					return nil, errors.Wrap(err, "Schema Field:"+field.Name)
				}
				arrayItems["type"] = schemaElement
			}
		} else {
			// It's a scalar.
			if field.IsBuiltin {
				propertyContent := make(map[string]interface{})
				properties[field.GoName] = propertyContent
				propertyContent["title"] = topic + Sep + field.GoName

				if field.Type == "string" {
					propertyContent["type"] = "string"
				} else if field.Type == "time" {
					timeItems := make(map[string]interface{})
					timeItems["Sec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "Sec"}
					timeItems["NSec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "NSec"}
					propertyContent["type"] = "object"
					propertyContent["properties"] = timeItems
				} else if field.Type == "duration" {
					timeItems := make(map[string]interface{})
					timeItems["Sec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "Sec"}
					timeItems["NSec"] = map[string]string{"type": "integer", "title": topic + Sep + field.GoName + Sep + "NSec"}
					propertyContent["type"] = "object"
					propertyContent["properties"] = timeItems
				} else {
					// It's a primitive.
					var jsonType string
					if field.GoType == "int8" || field.GoType == "uint8" || field.GoType == "int16" || field.GoType == "uint16" ||
						field.GoType == "int32" || field.GoType == "uint32" || field.GoType == "int64" || field.GoType == "uint64" {
						jsonType = "integer"
						jsonType = "integer"
						jsonType = "integer"
					} else if field.GoType == "float32" || field.GoType == "float64" {
						jsonType = "number"
					} else if field.GoType == "bool" {
						jsonType = "bool"
					} else {
						// Something went wrong.
						return nil, errors.New("we haven't implemented this primitive yet")
					}
					propertyContent["type"] = jsonType
				}
			} else {
				// It's another nested message.

				// Generate the nested type.
				msgType, err := newDynamicMessageTypeNested(field.Type, field.Package)
				if err != nil {
					return nil, errors.Wrap(err, "Schema Field: "+field.Name)
				}

				// Recursively generate schema information for the nested type.
				schemaElement, err := msgType.generateJSONSchemaProperties(topic + Sep + field.GoName)
				if err != nil {
					return nil, errors.Wrap(err, "Schema Field:"+field.Name)
				}
				properties[field.GoName] = schemaElement
			}
		}
	}

	// All done.
	return schemaItems, nil
}

//	DynamicMessage

// Type returns the ROS type of a dynamic message; required for ros.Message.
func (m DynamicMessage) Type() MessageType {
	return m.dynamicType
}

// MarshalJSON provides a custom implementation of JSON marshalling, so that when the DynamicMessage is recursively
// marshalled using the standard Go json.marshal() mechanism, the resulting JSON representation is a compact representation
// of just the important message payload (and not the message definition).  It's important that this representation matches
// the schema generated by GenerateJSONSchema().
func (m DynamicMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

//UnmarshalJSON provides a custom implementation of JSON unmarshalling. Using the dynamicMessage provided, Msgspec is used to
//determine the individual parsing of each JSON encoded payload item into the correct Go type. It is important each type is
//correct so that the message serializes correctly and is understood by the ROS system
func (m DynamicMessage) UnmarshalJSON(buf []byte) error {

	var err error
	var tmpArray []interface{}
	//We range the message spec fields as this is more predictable than using JSON structure to loop
	for _, field := range m.dynamicType.spec.Fields {
		//If the field is an array
		if field.IsArray {
			//Here we grab the msgBytes and json decoded msgType for the field Name lookup
			msgBytes, msgType, _, _ := jsonparser.Get(buf, field.GoName)
			jsonType := fmt.Sprintf("%v", msgType)
			//For scenario where a []uint8 array is provided, ros byte array encoded as json string
			if jsonType == "string" {
				bytes, err := jsonparser.GetString(buf, field.GoName)
				if err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
				m.data[field.GoName] = []byte(bytes)
			} else if jsonType == "array" {
				//Continue if the payload provides an array
				//TODO: Utilize the array callback feature in the jsonparser package, structure will change significantly
				//Using a builtin json unmarshal is cheating and uses too many resources
				err := json.Unmarshal(msgBytes, &tmpArray)
				if err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
				//here we are checking the size of the array provided in json
				//If the array is in field spec is fixed length and it does not match, an error will occur
				size := len(tmpArray)
				if field.ArrayLen != -1 && field.ArrayLen != size {
					return errors.New("Fixed array length incorrect : " + field.GoName)
				}
				//Loop through array items
				for i := 0; i < size; i++ {
					//Get the json bytes of the array Item
					arrayBytes, err := json.Marshal(tmpArray[i])
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					//If it is a builtin goType
					if field.IsBuiltin {
						//Here we grab the msgBytes and json decoded msgType for the field Name lookup in array
						msgBytes, msgType, _, _ := jsonparser.Get(arrayBytes, field.GoName)
						jsonType := fmt.Sprintf("%v", msgType)
						if jsonType == "bool" {
							//Lookup bool field inside array
							data, err := jsonparser.GetBoolean(msgBytes, field.GoName)
							if err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
							m.data[field.GoName] = append(m.data[field.GoName].([]bool), bool(data))
						} else if jsonType == "number" {
							//Lookup number inside array
							data, err := jsonparser.GetFloat(msgBytes, field.GoName)
							if err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
							switch field.GoType {
							case "int8":
								m.data[field.GoName] = append(m.data[field.GoName].([]int8), int8(data))
							case "int16":
								m.data[field.GoName] = append(m.data[field.GoName].([]int16), int16(data))
							case "int32":
								m.data[field.GoName] = append(m.data[field.GoName].([]int32), int32(data))
							case "int64":
								m.data[field.GoName] = append(m.data[field.GoName].([]int64), int64(data))
							case "uint8":
								m.data[field.GoName] = append(m.data[field.GoName].([]uint8), uint8(data))
							case "uint16":
								m.data[field.GoName] = append(m.data[field.GoName].([]uint16), uint16(data))
							case "uint32":
								m.data[field.GoName] = append(m.data[field.GoName].([]uint32), uint32(data))
							case "uint64":
								m.data[field.GoName] = append(m.data[field.GoName].([]uint64), uint64(data))
							case "float32":
								m.data[field.GoName] = append(m.data[field.GoName].([]float32), float32(data))
							case "float64":
								m.data[field.GoName] = append(m.data[field.GoName].([]float64), float64(data))
							}
						} else if jsonType == "string" {
							//Lookup string inside array
							data, err := jsonparser.GetString(msgBytes, field.GoName)
							if err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
							if field.GoType == "string" {
								m.data[field.GoName] = append(m.data[field.GoName].([]string), string(data))
							} else {
								return errors.New("Expected type []string, found : " + field.GoType)
							}
						} else if jsonType == "object" {
							//Time and Duration are provided as JSON objects
							switch field.GoType {
							case "ros.Time":
								sec, err := jsonparser.GetFloat(msgBytes, "Sec")
								nSec, err := jsonparser.GetFloat(msgBytes, "NSec")
								if err != nil {
									return errors.Wrap(err, "Field: "+field.Name)
								}
								tmpTime := Time{}
								tmpTime.Sec = uint32(sec)
								tmpTime.NSec = uint32(nSec)
								m.data[field.GoName] = append(m.data[field.GoName].([]Time), tmpTime)
							case "ros.Duration":
								sec, err := jsonparser.GetFloat(msgBytes, "Sec")
								nSec, err := jsonparser.GetFloat(msgBytes, "NSec")
								if err != nil {
									return errors.Wrap(err, "Field: "+field.Name)
								}
								tmpDuration := Duration{}
								tmpDuration.Sec = uint32(sec)
								tmpDuration.NSec = uint32(nSec)
								m.data[field.GoName] = append(m.data[field.GoName].([]Duration), tmpDuration)
							}
						}
					} else {
						//We have a nested message array
						t2, err := newDynamicMessageTypeNested(field.Type, field.Package)
						if err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						nestedMsg := t2.NewMessage().(*DynamicMessage)
						if err = nestedMsg.UnmarshalJSON(arrayBytes); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						m.data[field.GoName] = append(m.data[field.GoName].([]Message), nestedMsg)
					}
				}
			}
		} else {
			//Here we grab the msgBytes and json decoded msgType for the field Name lookup
			msgBytes, msgType, _, _ := jsonparser.Get(buf, field.GoName)
			jsonType := fmt.Sprintf("%v", msgType)
			if jsonType == "bool" {
				//Lookup bool field inside json message
				data, err := jsonparser.GetBoolean(msgBytes, field.GoName)
				if err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
				m.data[field.GoName] = bool(data)
			} else if jsonType == "number" {
				//Lookup number field inside json message
				data, err := jsonparser.GetFloat(buf, field.GoName)
				if err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
				switch field.GoType {
				case "int8":
					m.data[field.GoName] = int8(data)
				case "int16":
					m.data[field.GoName] = int16(data)
				case "int32":
					m.data[field.GoName] = int32(data)
				case "int64":
					m.data[field.GoName] = int64(data)
				case "uint8":
					m.data[field.GoName] = uint8(data)
				case "uint16":
					m.data[field.GoName] = uint16(data)
				case "uint32":
					m.data[field.GoName] = uint32(data)
				case "uint64":
					m.data[field.GoName] = uint64(data)
				case "float32":
					m.data[field.GoName] = float32(data)
				case "float64":
					m.data[field.GoName] = float64(data)
				default:
					return errors.New("Built-in type not found : " + field.GoType)
				}
			} else if jsonType == "string" {
				//Lookup string field inside json message
				data, err := jsonparser.GetString(buf, field.GoName)
				if err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
				if field.GoType == "string" {
					m.data[field.GoName] = data
				} else {
					return errors.New("Expected type []string, found : " + field.GoType)
				}
			} else if jsonType == "object" {
				//Time and Duration are provided as JSON objects
				switch field.GoType {
				case "ros.Time":
					sec, err := jsonparser.GetFloat(msgBytes, "Sec")
					nSec, err := jsonparser.GetFloat(msgBytes, "NSec")
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmpTime := Time{}
					tmpTime.Sec = uint32(sec)
					tmpTime.NSec = uint32(nSec)
					m.data[field.GoName] = tmpTime
				case "ros.Duration":
					sec, err := jsonparser.GetFloat(msgBytes, "Sec")
					nSec, err := jsonparser.GetFloat(msgBytes, "NSec")
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmpDuration := Duration{}
					tmpDuration.Sec = uint32(sec)
					tmpDuration.NSec = uint32(nSec)
					m.data[field.GoName] = tmpDuration
				default:
					//We've got a nested ros message
					t2, err := newDynamicMessageTypeNested(field.Type, field.Package)
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					nestedMsg := t2.NewMessage().(*DynamicMessage)
					if err = nestedMsg.UnmarshalJSON(msgBytes); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					m.data[field.GoName] = nestedMsg
				}
			}
		}
	}
	return err
}

// Serialize converts a DynamicMessage into a TCPROS bytestream allowing it to be published to other nodes; required for
// ros.Message.
func (m DynamicMessage) Serialize(buf *bytes.Buffer) error {
	// THIS METHOD IS BASICALLY AN UNTEMPLATED COPY OF THE TEMPLATE IN LIBGENGO.

	var err error = nil

	// Iterate over each of the fields in the message.
	for _, field := range m.dynamicType.spec.Fields {
		if field.IsArray {
			// It's an array.

			// Look up the item.
			array, ok := m.data[field.GoName]
			if !ok {
				return errors.New("Field: " + field.Name + ": No data found.")
			}

			// If the array is not a fixed length, it begins with a declaration of the array size.
			var size uint32
			if field.ArrayLen < 0 {
				size = uint32(reflect.ValueOf(array).Len())
				if err := binary.Write(buf, binary.LittleEndian, size); err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
			} else {
				size = uint32(field.ArrayLen)
			}

			// Then we just write out all the elements one after another.
			arrayValue := reflect.ValueOf(array)
			for i := uint32(0); i < size; i++ {
				//Casting the array item to interface type
				var arrayItem interface{} = arrayValue.Index(int(i)).Interface()
				// Need to handle each type appropriately.
				if field.IsBuiltin {
					if field.Type == "string" {
						// Make sure we've actually got a string.
						str, ok := arrayItem.(string)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected string.")
						}
						// The string should start with a declaration of the number of characters.
						var sizeStr uint32 = uint32(len(str))
						if err := binary.Write(buf, binary.LittleEndian, sizeStr); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						// Then write out the actual characters.
						data := []byte(str)
						if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}

					} else if field.Type == "time" {
						// Make sure we've actually got a time.
						t, ok := arrayItem.(Time)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected ros.Time.")
						}
						// Then write out the structure.
						if err := binary.Write(buf, binary.LittleEndian, t); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}

					} else if field.Type == "duration" {
						// Make sure we've actually got a duration.
						d, ok := arrayItem.(Duration)
						if !ok {
							return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected ros.Duration.")
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
							v, ok := arrayItem.(bool)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected bool.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int8":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(int8)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected int8.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int16":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(int16)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected int16.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int32":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(int32)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected int32.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "int64":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(int64)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected int64.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint8":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(uint8)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected uint8.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint16":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(uint16)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected uint16.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint32":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(uint32)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected uint32.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "uint64":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(uint64)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected uint64.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "float32":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(float32)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected float32.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						case "float64":
							// Make sure we've actually got a bool.
							v, ok := arrayItem.(float64)
							if !ok {
								return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected float64.")
							}
							// Then write out the value.
							if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
								return errors.Wrap(err, "Field: "+field.Name)
							}
						default:
							// Something went wrong.
							return errors.New("we haven't implemented this primitive yet")
						}
					}

				} else {
					// Else it's not a builtin.

					// Confirm the message we're holding is actually the correct type.
					msg, ok := arrayItem.(*DynamicMessage)
					if !ok {
						return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(arrayItem).Name() + ", expected Message.")
					}
					if msg.dynamicType.spec.ShortName != field.Type {
						return errors.New("Field: " + field.Name + ": Found msg " + msg.dynamicType.spec.ShortName + ", expected " + field.Type + ".")
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
					var sizeStr uint32 = uint32(len(str))
					if err := binary.Write(buf, binary.LittleEndian, sizeStr); err != nil {
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
						return errors.New("we haven't implemented this primitive yet")
					}
				}

			} else {
				// Else it's not a builtin.

				// Confirm the message we're holding is actually the correct type.
				msg, ok := item.(*DynamicMessage)
				if !ok {
					return errors.New("Field: " + field.Name + ": Found " + reflect.TypeOf(item).Name() + ", expected Message.")
				}
				if msg.dynamicType.spec.ShortName != field.Type {
					return errors.New("Field: " + field.Name + ": Found msg " + msg.dynamicType.spec.ShortName + ", expected " + field.Type + ".")
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

// Deserialize parses a byte stream into a DynamicMessage, thus reconstructing the fields of a received ROS message; required
// for ros.Message.
func (m *DynamicMessage) Deserialize(buf *bytes.Reader) error {
	// THIS METHOD IS BASICALLY AN UNTEMPLATED COPY OF THE TEMPLATE IN LIBGENGO.

	// To give more sane results in the event of a decoding issue, we decode into a copy of the data field.
	var err error = nil
	tmpData := make(map[string]interface{})
	m.data = nil

	// Iterate over each of the fields in the message.
	for _, field := range m.dynamicType.spec.Fields {
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
				tmpData[field.GoName] = make([]bool, 0)
			case "int8":
				tmpData[field.GoName] = make([]int8, 0)
			case "int16":
				tmpData[field.GoName] = make([]int16, 0)
			case "int32":
				tmpData[field.GoName] = make([]int32, 0)
			case "int64":
				tmpData[field.GoName] = make([]int64, 0)
			case "uint8":
				tmpData[field.GoName] = make([]uint8, 0)
			case "uint16":
				tmpData[field.GoName] = make([]uint16, 0)
			case "uint32":
				tmpData[field.GoName] = make([]uint32, 0)
			case "uint64":
				tmpData[field.GoName] = make([]uint64, 0)
			case "float32":
				tmpData[field.GoName] = make([]float32, 0)
			case "float64":
				tmpData[field.GoName] = make([]float64, 0)
			case "string":
				tmpData[field.GoName] = make([]string, 0)
			case "time":
				tmpData[field.GoName] = make([]Time, 0)
			case "duration":
				tmpData[field.GoName] = make([]Duration, 0)
			default:
				// In this case, it will probably be because the go_type is describing another ROS message, so we need to replace that with a nested DynamicMessage.
				tmpData[field.GoName] = make([]Message, 0)
			}
			// Iterate over each item in the array.
			for i := 0; i < int(size); i++ {
				if field.IsBuiltin {
					if field.Type == "string" {
						// The string will start with a declaration of the number of characters.
						var strSize uint32
						if err = binary.Read(buf, binary.LittleEndian, &strSize); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						data := make([]byte, int(strSize))
						if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						tmpData[field.GoName] = append(tmpData[field.GoName].([]string), string(data))
					} else if field.Type == "time" {
						var data Time
						// Time/duration types have two fields, so consume into these in two reads.
						if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						tmpData[field.GoName] = append(tmpData[field.GoName].([]Time), data)
					} else if field.Type == "duration" {
						var data Duration
						// Time/duration types have two fields, so consume into these in two reads.
						if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
						tmpData[field.GoName] = append(tmpData[field.GoName].([]Duration), data)
					} else {
						// It's a regular primitive element.

						// Because not runtime expressions in type assertions in Go, we need to manually do this.
						switch field.GoType {
						case "bool":
							var data bool
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]bool), data)
						case "int8":
							var data int8
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]int8), data)
						case "int16":
							var data int16
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]int16), data)
						case "int32":
							var data int32
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]int32), data)
						case "int64":
							var data int64
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]int64), data)
						case "uint8":
							var data uint8
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]uint8), data)
						case "uint16":
							var data uint16
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]uint16), data)
						case "uint32":
							var data uint32
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]uint32), data)
						case "uint64":
							var data uint64
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]uint64), data)
						case "float32":
							var data float32
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]float32), data)
						case "float64":
							var data float64
							binary.Read(buf, binary.LittleEndian, &data)
							tmpData[field.GoName] = append(tmpData[field.GoName].([]float64), data)
						default:
							// Something went wrong.
							return errors.New("we haven't implemented this primitive yet")
						}
						if err != nil {
							return errors.Wrap(err, "Field: "+field.Name)
						}
					}
				} else {
					// Else it's not a builtin.
					msgType, err := newDynamicMessageTypeNested(field.Type, field.Package)
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					msg := msgType.NewMessage()
					if err = msg.Deserialize(buf); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmpData[field.GoName] = append(tmpData[field.GoName].([]Message), msg)
				}
			}
		} else {
			// Else it's a scalar.  This is just the same as above, with the '[i]' bits removed.

			if field.IsBuiltin {
				if field.Type == "string" {
					// The string will start with a declaration of the number of characters.
					var strSize uint32
					if err = binary.Read(buf, binary.LittleEndian, &strSize); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					data := make([]byte, int(strSize))
					if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmpData[field.GoName] = string(data)
				} else if field.Type == "time" {
					var data Time
					// Time/duration types have two fields, so consume into these in two reads.
					if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmpData[field.GoName] = data
				} else if field.Type == "duration" {
					var data Duration
					// Time/duration types have two fields, so consume into these in two reads.
					if err = binary.Read(buf, binary.LittleEndian, &data.Sec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					if err = binary.Read(buf, binary.LittleEndian, &data.NSec); err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
					tmpData[field.GoName] = data
				} else {
					// It's a regular primitive element.
					switch field.GoType {
					case "bool":
						var data bool
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "int8":
						var data int8
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "int16":
						var data int16
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "int32":
						var data int32
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "int64":
						var data int64
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "uint8":
						var data uint8
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "uint16":
						var data uint16
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "uint32":
						var data uint32
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "uint64":
						var data uint64
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "float32":
						var data float32
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					case "float64":
						var data float64
						err = binary.Read(buf, binary.LittleEndian, &data)
						tmpData[field.GoName] = data
					default:
						// Something went wrong.
						return errors.New("we haven't implemented this primitive yet")
					}
					if err != nil {
						return errors.Wrap(err, "Field: "+field.Name)
					}
				}
			} else {
				// Else it's not a builtin.
				msgType, err := newDynamicMessageTypeNested(field.Type, field.Package)
				if err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
				tmpData[field.GoName] = msgType.NewMessage()
				if err = tmpData[field.GoName].(Message).Deserialize(buf); err != nil {
					return errors.Wrap(err, "Field: "+field.Name)
				}
			}
		}
	}

	// All done.
	m.data = tmpData
	return err
}

func (m DynamicMessage) String() string {
	// Just print out the data!
	return fmt.Sprint(m.dynamicType.Name(), "::", m.data)
}

// DEFINE PRIVATE STATIC FUNCTIONS.

// DEFINE PRIVATE RECEIVER FUNCTIONS.

// ALL DONE.

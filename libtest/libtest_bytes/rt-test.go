package libtest_bytes

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/edwinhayes/rosgo/ros"
	"testing"
)

//RTTest tests the serialization, deserialization, and JSON functions of dynamic_message.
//Each function is compared with correct hard coded values to verify functionality
//A map of all message definitions is created and cycled through to test serialization of each message
func RTTest(t *testing.T) {
	//var msgType *ros.DynamicMessageType
	var err error

	//Instantiate a new dynamic message type
	msgType, err := ros.NewDynamicMessageType("geometry_msgs/Twist")
	if err != nil {
		t.Error("failed to get message definition; ", err)
		return
	}
	//Instantiate the sub message. This is not necessary for operation,
	// but for us to manually setup data to test the main message type
	nestedMsgType, err := ros.NewDynamicMessageType("geometry_msgs/Vector3")
	if err != nil {
		t.Error("failed to get message definition; ", err)
		return
	}
	//Example JSON payload, Marshaled JSON
	examplePayload := `{"angular":{"x":1,"y":2,"z":3},"linear":{"x":1,"y":2,"z":3}}`

	//Declaring example bytes taken from external ROS source
	rawmsg := "000000000000f03f00000000000000400000000000000840000000000000f03f00000000000000400000000000000840"
	exampleBytes, err := hex.DecodeString(rawmsg)

	//Example message data
	exampleMsg := "geometry_msgs/Twist::map[angular:geometry_msgs/Vector3::map[x:1.00000 y:2.00000 z:3.00000] linear:geometry_msgs/Vector3::map[x:1.00000 y:2.00000 z:3.00000]]"

	//Example schema
	exampleSchema := `{"$id":"/ros/testy","$schema":"https://json-schema.org/draft-07/schema#","properties":{"x":{"title":"/ros/testy/x","type":"number"},"y":{"title":"/ros/testy/y","type":"number"},"z":{"title":"/ros/testy/z","type":"number"}},"title":"/ros/testy","type":"object"}`
	//Generating a schema for geometry_msgs/Vector3 on topic chatty
	schema, err := nestedMsgType.GenerateJSONSchema("/ros/", "testy")
	if err != nil {
		t.Error("failed to get generate JSON schema; ", err)
		return
	}
	//Converting json schema into string for comparison to input schema
	rosgoSchema := string(schema)

	//Creating new message instances of the message types to be used for serialization/deserialization tests
	dynamicMsg := msgType.NewMessage().(*ros.DynamicMessage)
	dynamicBlankMsg := msgType.NewMessage().(*ros.DynamicMessage)
	returnMsg := msgType.NewMessage().(*ros.DynamicMessage)
	nestedDynamicMsg := nestedMsgType.NewMessage().(*ros.DynamicMessage)

	//Declaring some sample data for serialization
	d := dynamicMsg.Data()
	d2 := nestedDynamicMsg.Data()
	d2["x"] = ros.JsonFloat64{F: 1.0}
	d2["y"] = ros.JsonFloat64{F: 2.0}
	d2["z"] = ros.JsonFloat64{F: 3.0}
	d["angular"] = nestedDynamicMsg
	d["linear"] = nestedDynamicMsg

	//Using UnmarshalJSON method on a set of example bytes to compare with example Message
	err = dynamicBlankMsg.UnmarshalJSON([]byte(examplePayload))

	//Serializing message into bytes buffer
	var buf bytes.Buffer
	err = dynamicMsg.Serialize(&buf)
	if err != nil {
		t.Error("failed to serialize message; ", err)
		return
	}
	rosgoBytes := buf.Bytes()

	var buf2 bytes.Buffer
	err = dynamicBlankMsg.Serialize(&buf2)
	if err != nil {
		t.Error("failed to serialize message; ", err)
		return
	}
	jsonBytes := buf2.Bytes()

	//Deserializing message into bytes reader
	reader := bytes.NewReader(buf.Bytes())
	err = returnMsg.Deserialize(reader)
	if err != nil {
		t.Error("failed to deserialize message; ", err)
	}
	rosgoMsg := fmt.Sprintf("%v", returnMsg)

	//Using MarshalJSON method on dynamic message to create JSON payload
	payloadMsg, err := dynamicMsg.MarshalJSON()
	if err != nil {
		t.Error("failed to marshal JSON; ", err)
		return
	}
	//Convert to string and compare to example JSON payload
	rosgoPayload := fmt.Sprintf("%s", payloadMsg)
	if rosgoPayload != examplePayload {
		t.Error("marshalled JSON incorrect: " + rosgoPayload + " vs " + examplePayload)
		return
	}

	//Comparing byte slice arrays to check Serialization worked
	res := bytes.Compare(exampleBytes, rosgoBytes)
	if res != 0 {
		t.Error("Serialized Message incorrect: " + string(rosgoBytes) + " vs " + string(exampleBytes))
		return
	}
	//Comparing deserialized ros messages to check Deserialization worked
	if rosgoMsg != exampleMsg {
		t.Error("Deserialized message incorrect: " + rosgoMsg + " vs " + exampleMsg)
		return
	}
	//Comparing unmarshalled payload to check unmarshalJSON worked
	res = bytes.Compare(jsonBytes, rosgoBytes)
	if res != 0 {
		t.Error("Unmarshalled message incorrect; ", err)
		return
	}
	//Comparing json schema to example schema to check GenerateJSONSchema worked
	if rosgoSchema != exampleSchema {
		t.Error("JSON Schema information incorrect; ", err)
		return
	}
	return
}

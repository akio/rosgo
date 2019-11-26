package libtest_bytes

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/edwinhayes/rosgo/libgengo"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"strings"
	"testing"
)

//RTTest tests the serialization, deserialization, and JSON functions of dynamic_message.
//Each function is compared with correct hard coded values to verify functionality
//A map of all message definitions is created and cycled through to test serialization of each message
func RTTest(t *testing.T) {
	//var msgType *ros.DynamicMessageType
	var err error

	//Make a node
	node, err := ros.NewNode("/rosgo", os.Args)
	if err != nil {
		t.Error("Failed to make node;", err)
		return
	}
	node.Logger().SetSeverity(ros.LogLevelWarn)
	defer node.Shutdown()

	//Generate a map of all message types
	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")
	allMessages, err := libgengo.FindAllMessages(strings.Split(rosPkgPath, ":"))

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
	examplePayload := "{\"Angular\":{\"X\":1,\"Y\":2,\"Z\":3},\"Linear\":{\"X\":0,\"Y\":0,\"Z\":0}}"

	//Declaring example bytes taken from external ROS source
	rawmsg := "000000000000000000000000000000000000000000000000000000000000f03f00000000000000400000000000000840"
	exampleBytes, err := hex.DecodeString(rawmsg)

	//Example message data
	exampleMsg := "geometry_msgs/Twist::map[Angular:geometry_msgs/Vector3::map[X:1 Y:2 Z:3] Linear:geometry_msgs/Vector3::map[X:0 Y:0 Z:0]]"

	//Example schema
	exampleSchema := "{\"$id\":\"/ros/chatty\",\"$schema\":\"https://json-schema.org/draft-07/schema#\",\"properties\":{\"X\":{\"title\":\"/ros/chatty/X\",\"type\":\"number\"},\"Y\":{\"title\":\"/ros/chatty/Y\",\"type\":\"number\"},\"Z\":{\"title\":\"/ros/chatty/Z\",\"type\":\"number\"}},\"title\":\"/ros/chatty\",\"type\":\"object\"}"
	//Generating a schema for geometry_msgs/Vector3 on topic chatty
	schema, err := nestedMsgType.GenerateJSONSchema("/ros/", "chatty")
	if err != nil {
		t.Error("failed to get generate JSON schema; ", err)
		return
	}
	//Converting json schema into string for comparison to input schema
	rosgoSchema := string(schema)

	//Creating new message instances of the message types to be used for serialization/deserialization tests
	dynamicMsg := msgType.NewMessage().(*ros.DynamicMessage)
	returnMsg := msgType.NewMessage().(*ros.DynamicMessage)
	nestedDynamicMsg := nestedMsgType.NewMessage().(*ros.DynamicMessage)

	//Declaring some sample data for serialization
	d := dynamicMsg.Data()
	d2 := nestedDynamicMsg.Data()
	d2["X"] = float64(1)
	d2["Y"] = float64(2)
	d2["Z"] = float64(3)
	d["Angular"] = nestedDynamicMsg

	//Serializing message into bytes buffer
	var buf bytes.Buffer
	err = dynamicMsg.Serialize(&buf)
	if err != nil {
		t.Error("failed to serialize message; ", err)
		return
	}
	rosgoBytes := buf.Bytes()

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
		t.Error("Marshalled JSON incorrect; ", err)
		return
	}

	//Range through all messages
	for node.OK() {
		for message := range allMessages {

			//Create new dynamicMessageType with message from map
			msgType, err := ros.NewDynamicMessageType(message)
			if err != nil {
				t.Error("failed to get message definition; ", err)
				return
			}
			//Instantiate new message type with zero values
			dynamicMsg := msgType.NewMessage().(*ros.DynamicMessage)
			//Create a new publisher based on message name
			pubName := fmt.Sprintf("/shakedown/%s", message)
			pub := node.NewPublisher(pubName, dynamicMsg.Type())
			if pub == nil {
				t.Error("failed to create publisher; ", err)
				return
			}

			//Publish message
			msg := ros.Message(dynamicMsg)
			pub.Publish(msg)

			//Serializing message into new bytes buffer
			var buf bytes.Buffer
			err = dynamicMsg.Serialize(&buf)
			if err != nil {
				t.Error("failed to serialize message; ", err)
				return
			}
			//Remove publisher after a millisecond
			//	time.Sleep(time.Millisecond)
			node.RemovePublisher(pubName)

		}
		return
	}

	//Comparing byte slice arrays to check Serialization worked
	res := bytes.Compare(exampleBytes, rosgoBytes)
	if res != 0 {
		t.Error("Serialized Message incorrect; ", err)
		return
	}
	//Comparing deserialized ros messages to check Deserialization worked
	if rosgoMsg != exampleMsg {
		t.Error("Deserialized message incorrect; ", err)
		return
	}
	//Comparing json schema to example schema to check GenerateJSONSchema worked
	if rosgoSchema != exampleSchema {
		t.Error("JSON Schema information incorrect; ", err)
		return
	}
	return
}

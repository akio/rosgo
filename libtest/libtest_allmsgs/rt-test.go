package libtest_allmsgs

import (
	"bytes"
	"fmt"
	"github.com/edwinhayes/rosgo/libgengo"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"strings"
	"testing"
)

//RTTest searches all ros message files on the ros environment (opt/ros + catkin)
//For each message, a new dynamicMessageType is create and a default newMessage,
//and published to a new topic to test publishing and serialization, followed by publisher deletion
func RTTest(t *testing.T) {
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
			pub, err := node.NewPublisher(pubName, dynamicMsg.Type())
			if err != nil {
				t.Error("failed to create publisher; ", err)
				return
			}

			//Publish message
			msg := ros.Message(dynamicMsg)
			pub.Publish(msg)

			// //Serializing message into new bytes buffer
			var buf bytes.Buffer
			err = dynamicMsg.Serialize(&buf)
			if err != nil {
				t.Error("failed to serialize message; ", err)
				return
			}
			node.RemovePublisher(pubName)

		}
		return
	}
}

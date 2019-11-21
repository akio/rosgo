package libtest_dynamic_message

import (
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
)

var message string

//callback retrieves mapped data from ros.DynamicMessage.data
func callback(msg *ros.DynamicMessage) {
	message = (msg.Data()["Data"].(string))
}

//RTTest creates a new DynamicMessageType from a known type and instantiates a DynamicMessage with data
//Data is published to a topic and subscribed to with the dynamic types, to test rosgo integration of dynamic types
//This also checks message is recieved correctly, to test dynamic Serialize and Deserialize
func RTTest(t *testing.T) {

	//Make a node
	node, err := ros.NewNode("/rosgo", os.Args)
	if err != nil {
		t.Error("Error instantiating node; ", err)
		return
	}
	node.Logger().SetSeverity(ros.LogLevelWarn)
	defer node.Shutdown()

	//Make a dynamicMessageType
	msgType, err := ros.NewDynamicMessageType("std_msgs/String")
	if err != nil {
		t.Error("Error creating dynamic message type: ", err)
		return
	}

	//Make a dynamicMessage with data
	dynamicMsg := msgType.NewMessage().(*ros.DynamicMessage)
	d := dynamicMsg.Data()
	d["Data"] = "hello"

	//Make a publisher and subscriber
	pub := node.NewPublisher("/chatter", dynamicMsg.Type())
	if pub == nil {
		t.Error("Failed to make publisher with dynamic message type")
	}
	node.NewSubscriber("/chatter", msgType, callback)

	for node.OK() {
		node.SpinOnce()

		//Publish a message
		msg := ros.Message(dynamicMsg)
		pub.Publish(msg)

		//When message recieved, check it is correct
		if message != "" {
			if message == "hello" {
				//Test has passed
				return
			}
			t.Error("Recieved wrong message")
			return
		}
	}
	//Should not get here
	t.Error("Node test failed")

}

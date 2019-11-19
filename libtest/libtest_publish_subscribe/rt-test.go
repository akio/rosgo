package libtest_publish_subscribe

import (
	"fmt"
	"github.com/edwinhayes/rosgo/libtest/msgs/std_msgs"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
)

var message string
var subscription int

func callback(msg *std_msgs.String) {
	message = string(msg.Data)
}

// RTTest performs a run-time test of using rosgo to create publisher and subscriber nodes
// The test must remove publishers and subscribers and reinitialize them with same name to pass
func RTTest(t *testing.T) {
	// Instantiate a ROS node.
	node, err := ros.NewNode("/rosgo", os.Args)
	if err != nil {
		t.Error(err)
		return
	}
	node.Logger().SetSeverity(ros.LogLevelWarn)

	defer node.Shutdown()

	subscription = 1

	pub := node.NewPublisher("rosgomessage", std_msgs.MsgString)
	if pub == nil {
		t.Error("NewPublisher failed; ", pub)
		return
	}

	// Create a publisher on the node.

	for node.OK() {
		node.SpinOnce()

		// Try to publish a message.
		var m std_msgs.String
		if subscription == 1 {
			m.Data = "First Subscriber"
			pub.Publish(&m)
		} else {
			m.Data = "Second Subscriber"
			pub := node.NewPublisher("rosgomessage", std_msgs.MsgString)
			if pub == nil {
				t.Error("NewPublisher failed; ", pub)

				fmt.Println("got here")
				return
			}
			pub.Publish(&m)
		}

		//Try to subscribe to the message
		node.NewSubscriber("rosgomessage", std_msgs.MsgString, callback)
		//Check the message is the same as what we published
		if message != "" {
			if message == "First Subscriber" {
				//Correct message recieved
				//Shutdown publisher and subscriber and initiate second test
				node.RemoveSubscriber("rosgomessage")
				node.RemovePublisher("rosgomessage")
				subscription = 2
				message = ""
			} else if message == "Second Subscriber" {
				//Second subscription worked
				return
			} else {
				//An incorrect message has been recieved
				t.Error("Wrong message recieved", message)
				return
			}
		}
	}
	t.Error("Node test failed")
	return
}

// ALL DONE.

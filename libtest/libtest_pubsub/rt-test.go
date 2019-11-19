package libtest_pubsub

import (
	"github.com/edwinhayes/rosgo/libtest/msgs/std_msgs"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
)

var message string

func callback(msg *std_msgs.String) {
	message = string(msg.Data)
}

// RTTest performs a run-time test of using rosgo to create publisher and subscriber nodes
// The test pubslihes a message to a new topic and subscribes to it on a new node
// The test will fail is the message is not recieved on the other end correctly
func RTTest(t *testing.T) {
	// Instantiate a ROS node.
	node, err := ros.NewNode("/rosgo", os.Args)
	if err != nil {
		t.Error(err)
		return
	}
	node.Logger().SetSeverity(ros.LogLevelWarn)

	defer node.Shutdown()

	// Create a publisher on the node.
	pub := node.NewPublisher("rosgomessage", std_msgs.MsgString)

	for node.OK() {
		node.SpinOnce()
		// Try to publish a message.
		var m std_msgs.String
		m.Data = "Hello"
		pub.Publish(&m)

		//Try to subscribe to the message
		node.NewSubscriber("rosgomessage", std_msgs.MsgString, callback)

		//Check the message is the same as what we published
		if message != "" {
			if message == "Hello" {
				//Correct message recieved
				return
			}
			//An incorrect message has been recieved
			t.Error("Wrong message recieved", message)
			return
		}
	}
	t.Error("Node test failed")
	return
}

// ALL DONE.

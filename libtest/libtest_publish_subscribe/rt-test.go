package libtest_publish_subscribe

import (
	"github.com/edwinhayes/rosgo/libtest/msgs/std_msgs"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
)

var message string
var subscription int
var eventname string
var m std_msgs.String

//Subscriber callback with event check
func callback(msg *std_msgs.String, event ros.MessageEvent) {
	message = string(msg.Data)
	eventname = event.PublisherName
}

//onConnect callback to run on publisher startup
func onConnect(pub ros.SingleSubscriberPublisher) {
	if pub.GetTopic() != "/rosgomessage" {
		message = "onConnect"
	}
	m.Data = "First Subscriber"
	pub.Publish(&m)
	subscription = 2
}

//onDisconnect callback to run on publisher shutdown
func onDisconnect(pub ros.SingleSubscriberPublisher) {
	if pub.GetTopic() != "/rosgomessage" {
		message = "onDisconnect"
	} else {
		message = ""
	}
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
	defer node.Shutdown()

	subscription = 1

	for node.OK() {
		node.SpinOnce()

		// Try to publish a message.
		if subscription == 1 {
			pub, err := node.NewPublisherWithCallbacks("rosgomessage", std_msgs.MsgString, onConnect, onDisconnect)
			if err != nil {
				t.Error("NewPublisher failed; ", pub, err)
				return
			}
		} else {
			m.Data = "Second Subscriber"
			pub, err := node.NewPublisher("rosgomessage", std_msgs.MsgString)
			if err != nil {
				t.Error("NewPublisher failed; ", pub, err)
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

			} else if message == "Second Subscriber" {
				//Second subscription worked
				if eventname == "/rosgo" {
					return
				}
				t.Error("Wrong message event", eventname)
				return
			} else {
				//An incorrect message has been recieved
				t.Error("Failed callback", message)
				return
			}
		}
	}
	t.Error("Node test failed")
	return
}

// ALL DONE.

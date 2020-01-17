package libtest_dynamic_message

import (
	"github.com/edwinhayes/rosgo/ros"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var message string

const targetMessage string = "hello world"

var chanRx chan string

// callback retrieves data from ros.DynamicMessage.data
func callback(msg *ros.DynamicMessage) {
	message = (msg.Data()["data"].(string))
	chanRx <- message
}

// RTTest creates a new DynamicMessageType from a known type and instantiates a DynamicMessage with data
// Data is published to a topic and subscribed to with the dynamic types, to test rosgo integration of dynamic types
// This also checks message is recieved correctly, to test dynamic Serialize and Deserialize
func RTTest(t *testing.T) {
	// Make a node.
	node, err := ros.NewNode("/rosgo", os.Args)
	if err != nil {
		t.Error("error instantiating node; ", err)
		return
	}
	node.SetLogLevel(logrus.ErrorLevel)
	defer node.Shutdown()

	// Make a dynamicMessageType.
	msgType, err := ros.NewDynamicMessageType("std_msgs/String")
	if err != nil {
		t.Error("error creating dynamic message type: ", err)
		return
	}

	// Make a dynamicMessage with data.
	dynamicMsg := msgType.NewMessage().(*ros.DynamicMessage)
	dynamicMsg.Data()["data"] = targetMessage

	// Make a publisher and subscriber.
	pub, err := node.NewPublisher("/chatter", dynamicMsg.Type())
	if err != nil {
		t.Error("failed to make publisher with dynamic message type")
		return
	}
	node.NewSubscriber("/chatter", msgType, callback)

	chanRx = make(chan string, 1)
	chanNodeStop := make(chan struct{})
	defer func() { chanNodeStop <- struct{}{}; <-chanNodeStop }()

	go func() {
		for node.OK() {
			select {
			case <-chanNodeStop:
				chanNodeStop <- struct{}{}
				return
			default:
				//Publish a message
				msg := ros.Message(dynamicMsg)
				pub.Publish(msg)

				node.SpinOnce()
			}
		}
	}()

	// Look for the rx callback to report something.
	select {
	case rx := <-chanRx:
		if rx == targetMessage {
			// Test passes.
			return
		}
		t.Error("rx mismatch:", rx, "vs", targetMessage)
		return
	case <-time.After(500 * time.Millisecond):
		t.Error("timeout waiting for rx")
		return
	}
}

package libtest_talker_with_callbacks

//go:generate gengo msg std_msgs/String
import (
	"fmt"
	"github.com/edwinhayes/rosgo/libtest/msgs/std_msgs"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
	"time"
)

//RTTest creates a new publisher node with callbacks
//This tests the newpublisherwithcallbacks functionality which invokes
//callback methods on connection and disconnection
func RTTest(t *testing.T) {
	node, err := ros.NewNode("/rosgo", os.Args)
	if err != nil {
		t.Error(err)
		return
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelWarn)
	pub := node.NewPublisherWithCallbacks("selftest/string", std_msgs.MsgString, onConnect, onDisconnect)
	if pub == nil {
		t.Error("NewPublisherWithCallbacks failed; ", pub)
		return
	}

	for node.OK() {
		node.SpinOnce()
		var msg std_msgs.String
		msg.Data = fmt.Sprintf("Hello World! The time is %s.", time.Now().String())
		pub.Publish(&msg)
		return
	}
	return
}

func onConnect(pub ros.SingleSubscriberPublisher) {
	fmt.Printf("-------Connect callback: node %s topic %s\n", pub.GetSubscriberName(), pub.GetTopic())
	var msg std_msgs.String
	msg.Data = fmt.Sprintf("hello %s", pub.GetSubscriberName())
	pub.Publish(&msg)
}

func onDisconnect(pub ros.SingleSubscriberPublisher) {
	fmt.Printf("-------Disconnect callback: node %s topic %s\n", pub.GetSubscriberName(), pub.GetTopic())
}

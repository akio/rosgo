package libtest_listener_with_event

//go:generate gengo msg std_msgs/String
import (
	"fmt"
	"github.com/edwinhayes/rosgo/libtest/msgs/std_msgs"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
)

func callback(msg *std_msgs.String, event ros.MessageEvent) {
	fmt.Printf("Received: %s from %s, header = %v, time = %v\n",
		msg.Data, event.PublisherName, event.ConnectionHeader, event.ReceiptTime)
}

//RTTest creates a listener with ros.messageevent
func RTTest(t *testing.T) {
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		t.Error(err)
		return
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	sub := node.NewSubscriber("/chatter", std_msgs.MsgString, callback)
	if sub == nil {
		t.Error("NewSubscriber failed; ", sub)
		return
	}
	node.Spin()
	return
}

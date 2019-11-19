package libtest_listener

import (
	"github.com/edwinhayes/rosgo/libtest/msgs/std_msgs"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
)

func callback(msg *std_msgs.String) {
	return
}

//RTTest creates a node which subscribes to the standard rosout topic on the ros system
//Does not test whether message is deserialized, or even recieved
func RTTest(t *testing.T) {
	// Instantiate a ROS node.
	node, err := ros.NewNode("/rosgo", os.Args)
	if err != nil {
		t.Error(err)
		return
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelWarn)
	// Create a subscriber on the node.
	node.NewSubscriber("selftest/string", std_msgs.MsgString, callback)
	return
}

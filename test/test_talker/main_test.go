package main

import (
	"fmt"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"std_msgs"
	"testing"
	"time"
)

//Run with env ROS_MASTER_URI="Your rosmaster uri" go test -timeout 30s github.com/edwinhayes/rosgo/test/test_talker
func Test(t *testing.T) {
	node, err := ros.NewNode("/talker", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	pub := node.NewPublisher("/chatter", std_msgs.MsgString)

	var msg std_msgs.String
	msg.Data = fmt.Sprintf("hello %s", time.Now().String())
	fmt.Println(msg.Data)
	pub.Publish(&msg)
	if err != nil {
		t.Error(err)
	}
}

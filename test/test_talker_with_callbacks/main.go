package main

//go:generate gengo msg std_msgs/String
import (
	"fmt"
	"os"
	"std_msgs"
	"time"

	"github.com/fetchrobotics/rosgo/ros"
)

func main() {
	node, err := ros.NewNode("/talker", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	pub := node.NewPublisherWithCallbacks("/chatter", std_msgs.MsgString, onConnect, onDisconnect)

	for node.OK() {
		node.SpinOnce()
		var msg std_msgs.String
		msg.Data = fmt.Sprintf("hello %s", time.Now().String())
		fmt.Println(msg.Data)
		pub.Publish(&msg)
		time.Sleep(time.Second)
	}
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

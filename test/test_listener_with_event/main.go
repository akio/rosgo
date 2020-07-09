package main

//go:generate gengo msg std_msgs/String
import (
	"fmt"
	"os"
	"std_msgs"

	"github.com/fetchrobotics/rosgo/ros"
)

func callback(msg *std_msgs.String, event ros.MessageEvent) {
	fmt.Printf("Received: %s from %s, header = %v, time = %v\n",
		msg.Data, event.PublisherName, event.ConnectionHeader, event.ReceiptTime)
}

func main() {
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	node.NewSubscriber("/chatter", std_msgs.MsgString, callback)
	node.Spin()
}

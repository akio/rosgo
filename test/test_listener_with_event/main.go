package main

import (
	"fmt"
	"ros"
	"rosgo_tests"
)

func callback(msg *rosgo_tests.Hello, event ros.MessageEvent) {
	fmt.Printf("Received: %s from %s, header = %v, time = %v\n",
		msg.Data, event.PublisherName, event.ConnectionHeader, event.ReceiptTime)
}

func main() {
	node := ros.NewNode("/listener")
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	node.NewSubscriber("/chatter", rosgo_tests.MsgHello, callback)
	node.Spin()
}

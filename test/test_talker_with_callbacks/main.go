package main

import (
	"fmt"
	"ros"
	"rosgo_tests"
	"time"
)

func main() {
	node := ros.NewNode("/talker")
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	pub := node.NewPublisherWithCallbacks("/chatter", rosgo_tests.MsgHello, onConnect, onDisconnect)

	for node.OK() {
		node.SpinOnce()
		var msg rosgo_tests.Hello
		msg.Data = fmt.Sprintf("hello %s", time.Now().String())
		fmt.Println(msg.Data)
		pub.Publish(&msg)
		time.Sleep(time.Second)
	}
}

func onConnect(pub ros.SingleSubscriberPublisher) {
	fmt.Printf("-------Connect callback: node %s topic %s\n", pub.GetSubscriberName(), pub.GetTopic())
	var msg rosgo_tests.Hello
	msg.Data = fmt.Sprintf("hello %s", pub.GetSubscriberName())
	pub.Publish(&msg)
}

func onDisconnect(pub ros.SingleSubscriberPublisher) {
	fmt.Printf("-------Disconnect callback: node %s topic %s\n", pub.GetSubscriberName(), pub.GetTopic())
}

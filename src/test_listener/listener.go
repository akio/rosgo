package main

import (
    "fmt"
    "ros"
    "rosgo_test"
)

func callback(msg *rosgo_test.Hello) {
    fmt.Printf("Received: %s\n", msg.Data)
}

func main() {
    node := ros.NewNode("/listener")
    defer node.Shutdown()
    node.Logger().SetSeverity(ros.LogLevelDebug)
    node.NewSubscriber("/chatter", rosgo_test.TypeOfHello(), callback)
    node.Spin()
}

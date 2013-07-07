package main

import (
    "fmt"
    "ros"
    "std_msgs"
)

func callback(msg *std_msgs.String) {
    fmt.Printf("Received: %s\n", msg.Data)
}

func main() {
    node := ros.NewNode("/listener")
    defer node.Shutdown()
    sub := node.NewSubscriber("/chatter", std_msgs.TypeOfString(), callback)
    defer sub.Shutdown()
    node.Spin()
}

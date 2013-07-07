package main

import (
    "fmt"
    "ros"
    "std_msgs"
    "time"
)

func main() {
    node := ros.NewNode("/talker")
    defer node.Shutdown()
    pub := node.NewPublisher("/chatter", std_msgs.TypeOfString())
    defer pub.Shutdown()

    for node.OK() {
        node.SpinOnce()
        var msg std_msgs.String
        msg.Data = fmt.Sprintf("hello %s", time.Now().String())
        fmt.Println(msg.Data)
        pub.Publish(&msg)
        time.Sleep(time.Second)
    }
}

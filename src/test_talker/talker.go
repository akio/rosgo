package main

import (
    "fmt"
    "time"
    "ros"
    "rosgo_test"
)

func main() {
    node := ros.NewNode("/talker")
    defer node.Shutdown()
    node.Logger().SetSeverity(ros.LogLevelDebug)
    pub := node.NewPublisher("/chatter", rosgo_test.TypeOfHello)

    for node.OK() {
        node.SpinOnce()
        var msg rosgo_test.Hello
        msg.Data = fmt.Sprintf("hello %s", time.Now().String())
        fmt.Println(msg.Data)
        pub.Publish(&msg)
        time.Sleep(time.Second)
    }
}

package main

import (
    "os"
    "fmt"
    "ros"
    "rosgo_test"
)


func callback(srv *rosgo_test.AddTwoInts) error {
    srv.Response.Sum = srv.Request.A + srv.Request.B
    fmt.Printf("%d + %d = %d\n", srv.Request.A, srv.Request.B, srv.Response.Sum)
    return nil
}

func main() {
    node := ros.NewNode("server")
    defer node.Shutdown()
    logger := node.Logger()
    logger.SetSeverity(ros.LogLevelDebug)
    server := node.NewServiceServer("/add_two_ints", rosgo_test.TypeOfAddTwoInts, callback)
    if server == nil {
        fmt.Println("Failed to initialize '/add_two_ints' service server")
        os.Exit(1)
    }
    defer server.Shutdown()
    node.Spin()
}

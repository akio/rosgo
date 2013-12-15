package main

import (
    "fmt"
    "ros"
    "test_srvs"
)


func main() {
    node := ros.NewNode("client")
    defer node.Shutdown()
    cli := node.NewServiceClient("add_two_ints", test_srvs.TypeOfAddTwoInts)
    defer cli.Shutdown()
    var srv test_srvs.AddTwoInts 
    if err := cli.Call(&srv); err != nil {
        fmt.Print(err)
    } else {
        fmt.Printf("%d + %d = %d", srv.Req.A, srv.Req.B, srv.Res.Result) }
}

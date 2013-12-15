package main

import (
    "os"
    "fmt"
    "strconv"
    "ros"
    "rosgo_test"
)


func main() {
    if len(os.Args) != 3 {
        fmt.Print("USAGE: test_client <int> <int>")
        os.Exit(1)
    }

    node := ros.NewNode("client")
    defer node.Shutdown()
    logger := node.Logger()
    logger.SetSeverity(ros.LogLevelDebug)
    cli := node.NewServiceClient("/add_two_ints", rosgo_test.TypeOfAddTwoInts)
    defer cli.Shutdown()
    var srv rosgo_test.AddTwoInts 
    var err error
    var a, b int64
    a, err = strconv.ParseInt(os.Args[1], 10, 32)
    if err != nil {
        fmt.Print(err)
        fmt.Println()
        os.Exit(1)
    }
    b, err = strconv.ParseInt(os.Args[2], 10, 32)
    if err != nil {
        fmt.Print(err)
        fmt.Println()
        os.Exit(1)
    }
    srv.Req.A = int32(a)
    srv.Req.B = int32(b)
    if err := cli.Call(&srv); err != nil {
        fmt.Print(err)
        fmt.Println()
    } else {
        fmt.Printf("%d + %d = %d", srv.Req.A, srv.Req.B, srv.Res.Result)
        fmt.Println()
    }
}

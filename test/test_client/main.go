package main

//go:generate gengo srv rospy_tutorials/AddTwoInts
import (
	"fmt"
	"os"
	"rospy_tutorials"
	"strconv"

	"github.com/fetchrobotics/rosgo/ros"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Print("USAGE: test_client <int> <int>")
		os.Exit(-1)
	}

	node, err := ros.NewNode("client", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	logger := node.Logger()
	logger.SetSeverity(ros.LogLevelDebug)
	cli := node.NewServiceClient("/add_two_ints", rospy_tutorials.SrvAddTwoInts)
	defer cli.Shutdown()
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
	var srv rospy_tutorials.AddTwoInts
	srv.Request.A = a
	srv.Request.B = b
	if err = cli.Call(&srv); err != nil {
		fmt.Print(err)
		fmt.Println()
	} else {
		fmt.Printf("%d + %d = %d",
			srv.Request.A, srv.Request.B, srv.Response.Sum)
		fmt.Println()
	}
}

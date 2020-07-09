package main

//go:generate gengo srv rospy_tutorials/AddTwoInts
import (
	"fmt"
	"os"
	"rospy_tutorials"

	"github.com/fetchrobotics/rosgo/ros"
)

func callback(srv *rospy_tutorials.AddTwoInts) error {
	srv.Response.Sum = srv.Request.A + srv.Request.B
	fmt.Printf("%d + %d = %d\n", srv.Request.A, srv.Request.B, srv.Response.Sum)
	return nil
}

func main() {
	node, err := ros.NewNode("server", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	logger := node.Logger()
	logger.SetSeverity(ros.LogLevelDebug)
	server := node.NewServiceServer("/add_two_ints", rospy_tutorials.SrvAddTwoInts, callback)
	if server == nil {
		fmt.Println("Failed to initialize '/add_two_ints' service server")
		os.Exit(-1)
	}
	defer server.Shutdown()
	node.Spin()
}

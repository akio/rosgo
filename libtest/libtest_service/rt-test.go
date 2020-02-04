package libtest_service

//go:generate gengo srv rospy_tutorials/AddTwoInts
import (
	"github.com/edwinhayes/rosgo/libtest/msgs/rospy_tutorials"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"testing"
)

var service rospy_tutorials.AddTwoInts

//Callback function for ros service
func callback(srv *rospy_tutorials.AddTwoInts) error {
	srv.Response.Sum = srv.Request.A + srv.Request.B
	return nil
}

//Go routine function to spin server node to be run in separate thread
func spinServer(node ros.Node, quit <-chan bool) {
	//Initialize server
	server := node.NewServiceServer("/add_two_ints", rospy_tutorials.SrvAddTwoInts, callback)
	defer server.Shutdown()
	for {
		select {
		case <-quit:
			server.Shutdown()
			return
		default:
			node.SpinOnce()
		}
	}
}

//RTTest creates two separate nodes for a client and server. It makes these in separate threads so they can spin simultaneously
//A service for add_two_ints is called and response is checked
func RTTest(t *testing.T) {
	//func main() {
	//Initialize nodes ; skip error tests
	node, err := ros.NewNode("client", os.Args)
	node2, err := ros.NewNode("server", os.Args)
	//Defer node shutdown
	defer node.Shutdown()
	defer node2.Shutdown()

	//Initialize Client
	cli := node.NewServiceClient("/add_two_ints", rospy_tutorials.SrvAddTwoInts)
	if cli == nil {
		t.Error("Failed to initialize client")
	}
	defer cli.Shutdown()

	//Initialize server thread
	quitThread := make(chan bool)
	go spinServer(node2, quitThread)

	for node.OK() {
		//Create  and call service request
		service.Request.A = 10
		service.Request.B = 10
		if err = cli.Call(&service); err != nil {
		}

		//When a response is recieved
		if service.Response.Sum != 0 {
			//Check if response is correct
			if service.Response.Sum == 20 {
				cli.Shutdown()
				defer close(quitThread)
				return
			}
			//Response incorrect
			cli.Shutdown()
			defer close(quitThread)
			t.Error("Incorrect response recieved from server")
		}
		//Spin client node
		node.SpinOnce()
	}
}

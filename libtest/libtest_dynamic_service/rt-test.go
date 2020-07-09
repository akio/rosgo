package libtest_dynamic_service

import (
	"os"
	"testing"

	"github.com/edwinhayes/rosgo/ros"
)

var service *ros.DynamicService

//Callback function for ros service
func callback(srv ros.Service) error {
	req := srv.ReqMessage().(*ros.DynamicMessage)
	if req.Data()["data"] == true {
		res := srv.ResMessage().(*ros.DynamicMessage)
		res.Data()["success"] = true
		res.Data()["message"] = "dynamic message worked"
	}
	return nil
}

//Go routine function to spin server node to be run in separate thread
func spinServer(node ros.Node, quit <-chan bool) {

	//Initialize server - Server can keep using static service for now
	server := node.NewServiceServer("/add_two_ints", service.Type(), callback)
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
	var err error
	node, err := ros.NewNode("client", os.Args)
	node2, err := ros.NewNode("server", os.Args)
	//Defer node shutdown
	defer node.Shutdown()
	defer node2.Shutdown()

	// Create dynamic service
	srv, err := ros.NewDynamicServiceType("std_srvs/SetBool")
	if err != nil {
		t.Errorf("failed to create service: %s", err)
	}
	service = srv.NewService().(*ros.DynamicService)
	req := service.ReqMessage().(*ros.DynamicMessage)
	req.Data()["data"] = true

	//Initialize Client
	cli := node.NewServiceClient("/add_two_ints", service.Type())
	if cli == nil {
		t.Error("Failed to initialize client")
	}
	defer cli.Shutdown()

	//Initialize server thread
	quitThread := make(chan bool)
	go spinServer(node2, quitThread)

	for node.OK() {
		//Create  and call service request
		if err = cli.Call(service); err != nil {
		}
		res := service.ResMessage().(*ros.DynamicMessage)
		if res.Data()["success"] == true {
			// Succeeded
			cli.Shutdown()
			defer close(quitThread)
			return
		}
		//Spin client node
		_ = node.SpinOnce()
	}
}

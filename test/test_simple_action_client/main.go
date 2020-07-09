package main

import (
	"actionlib_tutorials"
	"fmt"
	"os"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

func main() {
	node, err := ros.NewNode("test_fibonacci_client", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	logger := node.Logger()
	defer node.Shutdown()
	go node.Spin()

	ac := actionlib.NewSimpleActionClient(node, "fibonacci", actionlib_tutorials.ActionFibonacci)
	logger.Info("Waiting for server to start")

	started := ac.WaitForServer(ros.NewDuration(0, 0))
	if !started {
		logger.Info("Action server failed to start within timeout period.")
		return
	}

	logger.Info("Action server started, sending goal.")
	goal := &actionlib_tutorials.FibonacciGoal{Order: 20}
	ac.SendGoal(goal, nil, nil, nil)

	finished := ac.WaitForResult(ros.NewDuration(60, 0))
	if finished {
		state, err := ac.GetState()
		if err != nil {
			logger.Errorf("Error getting state: %v", err)
			return
		}
		logger.Infof("Action finished: %v", state)
	} else {
		logger.Errorf("Action did not finish before the timeout")
	}
}

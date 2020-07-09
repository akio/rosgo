package main

//go:generate gengo action actionlib_tutorials/Fibonacci

import (
	"actionlib_tutorials"
	"fmt"
	"os"

	"github.com/fetchrobotics/rosgo/actionlib"

	"github.com/fetchrobotics/rosgo/ros"
)

type fibonacciClient struct {
	node   ros.Node
	name   string
	ac     actionlib.SimpleActionClient
	logger ros.Logger
}

func newfibonacciClient(node ros.Node, name string) *fibonacciClient {
	fc := &fibonacciClient{
		node:   node,
		ac:     actionlib.NewSimpleActionClient(node, name, actionlib_tutorials.ActionFibonacci),
		logger: node.Logger(),
	}

	fc.logger.Info("Waiting for server to start")
	fc.ac.WaitForServer(ros.NewDuration(0, 0))
	fc.logger.Info("Server started")
	return fc
}

func (fc *fibonacciClient) activeCb() {
	fc.logger.Info("Goal just went active")
}

func (fc *fibonacciClient) feedbackCb(fb *actionlib_tutorials.FibonacciFeedback) {
	fc.logger.Infof("Got feedback of from server: %v", fb.Sequence)
}

func (fc *fibonacciClient) doneCb(state uint8, result *actionlib_tutorials.FibonacciResult) {
	fc.logger.Infof("Finished in state %v", state)
	fc.logger.Infof("Sequence: %v", result.Sequence)
	fc.node.Shutdown()
}

func (fc *fibonacciClient) sendGoal(order int32) {
	goal := &actionlib_tutorials.FibonacciGoal{Order: order}
	fc.ac.SendGoal(goal, fc.doneCb, fc.activeCb, fc.feedbackCb)
}

func main() {
	node, err := ros.NewNode("test_fibonacci_client", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fc := newfibonacciClient(node, "fibonacci")
	fc.sendGoal(10)

	node.Spin()
}

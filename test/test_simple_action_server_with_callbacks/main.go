package main

//go:generate gengo action actionlib_tutorials/Fibonacci
import (
	"actionlib_tutorials"
	"fmt"
	"os"
	"time"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

// fibonacciServer implements a fibonacci simple action server
// using the execute callback
type fibonacciServer struct {
	node   ros.Node
	logger ros.Logger
	as     actionlib.SimpleActionServer
	name   string
	fb     ros.Message
	result ros.Message
}

// newfibonacciServer creates a new fibonacci action server and starts it
func newFibonacciServer(node ros.Node, name string) {
	s := &fibonacciServer{
		node:   node,
		name:   name,
		logger: node.Logger(),
	}

	s.as = actionlib.NewSimpleActionServer(node, name,
		actionlib_tutorials.ActionFibonacci, s.executeCallback, false)
	s.as.Start()
}

func (s *fibonacciServer) executeCallback(goal *actionlib_tutorials.FibonacciGoal) {
	feed := &actionlib_tutorials.FibonacciFeedback{}
	feed.Sequence = append(feed.Sequence, 0)
	feed.Sequence = append(feed.Sequence, 1)
	success := true

	for i := 1; i < int(goal.Order); i++ {
		if s.as.IsPreemptRequested() {
			success = false
			if err := s.as.SetPreempted(nil, ""); err != nil {
				s.logger.Error(err)
			}
			break
		}

		val := feed.Sequence[i] + feed.Sequence[i-1]
		feed.Sequence = append(feed.Sequence, val)

		s.as.PublishFeedback(feed)
		time.Sleep(1000 * time.Millisecond)
	}

	if success {
		result := &actionlib_tutorials.FibonacciResult{Sequence: feed.Sequence}
		if err := s.as.SetSucceeded(result, "goal"); err != nil {
			s.logger.Error(err)
		}
	}
}

func main() {
	node, err := ros.NewNode("test_fibonacci_server", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()

	newFibonacciServer(node, "fibonacci")
	node.Spin()
}

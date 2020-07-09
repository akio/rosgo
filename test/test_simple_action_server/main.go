package main

//go:generate gengo action actionlib_tutorials/Averaging
import (
	"actionlib_tutorials"
	"fmt"
	"math"
	"os"
	"std_msgs"

	"github.com/fetchrobotics/rosgo/actionlib"
	"github.com/fetchrobotics/rosgo/ros"
)

type averagingServer struct {
	node      ros.Node
	as        actionlib.SimpleActionServer
	dataCount int32
	goal      int32
	sum       float32
	sumSq     float64
	feedback  *actionlib_tutorials.AveragingFeedback
	result    *actionlib_tutorials.AveragingResult
	sub       ros.Subscriber
	logger    ros.Logger
}

func newAveragingServer(node ros.Node, name string) {
	avg := new(averagingServer)
	avg.node = node
	avg.logger = node.Logger()
	avg.as = actionlib.NewSimpleActionServer(node, name, actionlib_tutorials.ActionAveraging, nil, false)
	avg.sub = node.NewSubscriber("/random_number", std_msgs.MsgFloat32, avg.analysisCallback)

	avg.as.RegisterGoalCallback(avg.goalCallback)
	avg.as.RegisterPreemptCallback(avg.preemptCallback)
	avg.as.Start()
}

func (avg *averagingServer) goalCallback() {
	avg.dataCount = 0
	avg.sum = 0
	avg.sumSq = 0

	goal, err := avg.as.AcceptNewGoal()
	if err != nil {
		avg.logger.Errorf("Error accepting new goal: %v", err)
		return
	}

	avgGoal, ok := goal.(*actionlib_tutorials.AveragingGoal)
	if !ok {
		avg.logger.Errorf("Error accepting new goal: expected averaging action goal")
		return
	}

	avg.goal = avgGoal.Samples
}

func (avg *averagingServer) preemptCallback() {
	avg.as.SetPreempted(nil, "")
}

func (avg *averagingServer) analysisCallback(msg *std_msgs.Float32) {
	if !avg.as.IsActive() {
		return
	}

	avg.dataCount++
	avg.sum += msg.Data
	avg.feedback.Sample = avg.dataCount
	avg.feedback.Data = msg.Data
	avg.feedback.Mean = avg.sum / float32(avg.dataCount)
	avg.sumSq = math.Pow(float64(msg.Data), 2)
	avg.feedback.StdDev = float32(math.Sqrt(math.Abs((avg.sumSq/float64(msg.Data) - math.Pow(float64(avg.feedback.Mean), 2)))))

	if avg.dataCount > avg.goal {
		avg.result.Mean = avg.feedback.Mean
		avg.result.StdDev = avg.feedback.StdDev

		if avg.result.Mean < 5.0 {
			avg.logger.Info("Averaging action aborted")
			avg.as.SetAborted(avg.result, "")
		} else {
			avg.logger.Info("Averaging action succeeded")
			avg.as.SetSucceeded(avg.result, "")
		}
	}
}

func main() {
	node, err := ros.NewNode("test_averaging_server", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()

	newAveragingServer(node, "averaging")
	node.Spin()
}

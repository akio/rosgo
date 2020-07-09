package actionlib

import (
	"actionlib_msgs"

	"github.com/fetchrobotics/rosgo/ros"
)

func NewActionClient(node ros.Node, action string, actionType ActionType) ActionClient {
	return newDefaultActionClient(node, action, actionType)
}

func NewActionServer(node ros.Node, action string, actionType ActionType, goalCb, cancelCb interface{}, autoStart bool) ActionServer {
	return newDefaultActionServer(node, action, actionType, goalCb, cancelCb, autoStart)
}

func NewSimpleActionClient(node ros.Node, action string, actionType ActionType) SimpleActionClient {
	return newSimpleActionClient(node, action, actionType)
}

func NewSimpleActionServer(node ros.Node, action string, actionType ActionType, executeCb interface{}, autoStart bool) SimpleActionServer {
	return newSimpleActionServer(node, action, actionType, executeCb, autoStart)
}

func NewServerGoalHandlerWithGoal(as ActionServer, goal ActionGoal) ServerGoalHandler {
	return newServerGoalHandlerWithGoal(as, goal)
}

func NewServerGoalHandlerWithGoalId(as ActionServer, goalID *actionlib_msgs.GoalID) ServerGoalHandler {
	return newServerGoalHandlerWithGoalId(as, goalID)
}

type ActionClient interface {
	WaitForServer(timeout ros.Duration) bool
	SendGoal(goal ros.Message, transitionCallback interface{}, feedbackCallback interface{}) ClientGoalHandler
	CancelAllGoals()
	CancelAllGoalsBeforeTime(stamp ros.Time)
}

type ActionServer interface {
	Start()
	Shutdown()
	PublishResult(status actionlib_msgs.GoalStatus, result ros.Message)
	PublishFeedback(status actionlib_msgs.GoalStatus, feedback ros.Message)
	PublishStatus()
	RegisterGoalCallback(interface{})
	RegisterCancelCallback(interface{})
}

type SimpleActionClient interface {
	SendGoal(goal ros.Message, doneCb, activeCb, feedbackCb interface{})
	SendGoalAndWait(goal ros.Message, executeTimeout, preeptTimeout ros.Duration) (uint8, error)
	WaitForServer(timeout ros.Duration) bool
	WaitForResult(timeout ros.Duration) bool
	GetResult() (ros.Message, error)
	GetState() (uint8, error)
	GetGoalStatusText() (string, error)
	CancelAllGoals()
	CancelAllGoalsBeforeTime(stamp ros.Time)
	CancelGoal() error
	StopTrackingGoal()
}

type SimpleActionServer interface {
	Start()
	IsNewGoalAvailable() bool
	IsPreemptRequested() bool
	IsActive() bool
	SetSucceeded(result ros.Message, text string) error
	SetAborted(result ros.Message, text string) error
	SetPreempted(result ros.Message, text string) error
	AcceptNewGoal() (ros.Message, error)
	PublishFeedback(feedback ros.Message)
	GetDefaultResult() ros.Message
	RegisterGoalCallback(callback interface{}) error
	RegisterPreemptCallback(callback interface{})
}

type ClientGoalHandler interface {
	IsExpired() bool
	GetCommState() (CommState, error)
	GetGoalStatus() (uint8, error)
	GetGoalStatusText() (string, error)
	GetTerminalState() (uint8, error)
	GetResult() (ros.Message, error)
	Resend() error
	Cancel() error
}

type ServerGoalHandler interface {
	SetAccepted(string) error
	SetCancelled(ros.Message, string) error
	SetRejected(ros.Message, string) error
	SetAborted(ros.Message, string) error
	SetSucceeded(ros.Message, string) error
	SetCancelRequested() bool
	PublishFeedback(ros.Message)
	GetGoal() ros.Message
	GetGoalId() actionlib_msgs.GoalID
	GetGoalStatus() actionlib_msgs.GoalStatus
	Equal(ServerGoalHandler) bool
	NotEqual(ServerGoalHandler) bool
	Hash() uint32
	GetHandlerDestructionTime() ros.Time
	SetHandlerDestructionTime(ros.Time)
}

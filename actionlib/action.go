package actionlib

import (
	"actionlib_msgs"
	"std_msgs"

	"github.com/fetchrobotics/rosgo/ros"
)

type ActionType interface {
	MD5Sum() string
	Name() string
	GoalType() ros.MessageType
	FeedbackType() ros.MessageType
	ResultType() ros.MessageType
	NewAction() Action
}

type Action interface {
	GetActionGoal() ActionGoal
	GetActionFeedback() ActionFeedback
	GetActionResult() ActionResult
}

type ActionGoal interface {
	ros.Message
	GetHeader() std_msgs.Header
	GetGoalId() actionlib_msgs.GoalID
	GetGoal() ros.Message
	SetHeader(std_msgs.Header)
	SetGoalId(actionlib_msgs.GoalID)
	SetGoal(ros.Message)
}

type ActionFeedback interface {
	ros.Message
	GetHeader() std_msgs.Header
	GetStatus() actionlib_msgs.GoalStatus
	GetFeedback() ros.Message
	SetHeader(std_msgs.Header)
	SetStatus(actionlib_msgs.GoalStatus)
	SetFeedback(ros.Message)
}

type ActionResult interface {
	ros.Message
	GetHeader() std_msgs.Header
	GetStatus() actionlib_msgs.GoalStatus
	GetResult() ros.Message
	SetHeader(std_msgs.Header)
	SetStatus(actionlib_msgs.GoalStatus)
	SetResult(ros.Message)
}

package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"reflect"
	"std_msgs"
	"sync"
	"time"

	"github.com/fetchrobotics/rosgo/ros"
)

type defaultActionServer struct {
	node             ros.Node
	autoStart        bool
	started          bool
	action           string
	actionType       ActionType
	actionResult     ros.MessageType
	actionResultType ros.MessageType
	actionFeedback   ros.MessageType
	actionGoal       ros.MessageType
	statusMutex      sync.RWMutex
	statusFrequency  ros.Rate
	statusTimer      *time.Ticker
	handlers         map[string]*serverGoalHandler
	handlersTimeout  ros.Duration
	handlersMutex    sync.Mutex
	goalCallback     interface{}
	cancelCallback   interface{}
	lastCancel       ros.Time
	pubQueueSize     int
	subQueueSize     int
	goalSub          ros.Subscriber
	cancelSub        ros.Subscriber
	resultPub        ros.Publisher
	feedbackPub      ros.Publisher
	statusPub        ros.Publisher
	statusPubChan    chan struct{}
	goalIDGen        *goalIDGenerator
	shutdownChan     chan struct{}
}

func newDefaultActionServer(node ros.Node, action string, actType ActionType, goalCb interface{}, cancelCb interface{}, start bool) *defaultActionServer {
	return &defaultActionServer{
		node:            node,
		autoStart:       start,
		started:         false,
		action:          action,
		actionType:      actType,
		actionResult:    actType.ResultType(),
		actionFeedback:  actType.FeedbackType(),
		actionGoal:      actType.GoalType(),
		handlersTimeout: ros.NewDuration(60, 0),
		goalCallback:    goalCb,
		cancelCallback:  cancelCb,
		lastCancel:      ros.Now(),
	}
}

func (as *defaultActionServer) init() {
	as.statusPubChan = make(chan struct{}, 10)
	as.shutdownChan = make(chan struct{}, 10)

	// setup goal id generator and goal handlers
	as.goalIDGen = newGoalIDGenerator(as.node.Name())
	as.handlers = map[string]*serverGoalHandler{}

	// setup action result type so that we can create default result messages
	res := as.actionResult.NewMessage().(ActionResult).GetResult()
	as.actionResultType = res.Type()

	// get frequency from ros params
	as.statusFrequency = ros.NewRate(5.0)

	// get queue sizes from ros params
	// queue sizes not implemented by ros.Node yet
	as.pubQueueSize = 50
	as.subQueueSize = 50

	as.goalSub = as.node.NewSubscriber(fmt.Sprintf("%s/goal", as.action), as.actionType.GoalType(), as.internalGoalCallback)
	as.cancelSub = as.node.NewSubscriber(fmt.Sprintf("%s/cancel", as.action), actionlib_msgs.MsgGoalID, as.internalCancelCallback)
	as.resultPub = as.node.NewPublisher(fmt.Sprintf("%s/result", as.action), as.actionType.ResultType())
	as.feedbackPub = as.node.NewPublisher(fmt.Sprintf("%s/feedback", as.action), as.actionType.FeedbackType())
	as.statusPub = as.node.NewPublisher(fmt.Sprintf("%s/status", as.action), actionlib_msgs.MsgGoalStatusArray)
}

func (as *defaultActionServer) Start() {
	logger := as.node.Logger()
	defer func() {
		logger.Debug("defaultActionServer.start exit")
		as.started = false
	}()

	// initialize subscribers and publishers
	as.init()

	// start status publish ticker that notifies at 5hz
	as.statusTimer = time.NewTicker(time.Second / 5.0)
	defer as.statusTimer.Stop()

	as.started = true

	for {
		select {
		case <-as.shutdownChan:
			return

		case <-as.statusTimer.C:
			as.PublishStatus()

		case <-as.statusPubChan:
			arr := as.getStatus()
			as.statusPub.Publish(arr)
		}
	}
}

// PublishResult publishes action result message
func (as *defaultActionServer) PublishResult(status actionlib_msgs.GoalStatus, result ros.Message) {
	msg := as.actionResult.NewMessage().(ActionResult)
	msg.SetHeader(std_msgs.Header{Stamp: ros.Now()})
	msg.SetStatus(status)
	msg.SetResult(result)
	as.resultPub.Publish(msg)
}

// PublishFeedback publishes action feedback messages
func (as *defaultActionServer) PublishFeedback(status actionlib_msgs.GoalStatus, feedback ros.Message) {
	msg := as.actionFeedback.NewMessage().(ActionFeedback)
	msg.SetHeader(std_msgs.Header{Stamp: ros.Now()})
	msg.SetStatus(status)
	msg.SetFeedback(feedback)
	as.feedbackPub.Publish(msg)
}

func (as *defaultActionServer) getStatus() *actionlib_msgs.GoalStatusArray {
	as.handlersMutex.Lock()
	defer as.handlersMutex.Unlock()
	var statusList []actionlib_msgs.GoalStatus

	if as.node.OK() {
		for id, gh := range as.handlers {
			handlerTime := gh.GetHandlerDestructionTime()
			destroyTime := handlerTime.Add(as.handlersTimeout)

			if !handlerTime.IsZero() && destroyTime.Cmp(ros.Now()) <= 0 {
				delete(as.handlers, id)
				continue
			}

			statusList = append(statusList, gh.GetGoalStatus())
		}
	}

	goalStatus := &actionlib_msgs.GoalStatusArray{}
	goalStatus.Header.Stamp = ros.Now()
	goalStatus.StatusList = statusList
	return goalStatus
}

func (as *defaultActionServer) PublishStatus() {
	as.statusPubChan <- struct{}{}
}

// internalCancelCallback recieves cancel message from client
func (as *defaultActionServer) internalCancelCallback(goalID *actionlib_msgs.GoalID, event ros.MessageEvent) {
	as.handlersMutex.Lock()
	defer as.handlersMutex.Unlock()

	goalFound := false
	logger := as.node.Logger()
	logger.Debug("Action server has received a new cancel request")

	for id, gh := range as.handlers {
		cancelAll := (goalID.Id == "" && goalID.Stamp.IsZero())
		cancelCurrent := (goalID.Id == id)

		st := gh.GetGoalStatus()
		cancelBeforeStamp := (!goalID.Stamp.IsZero() && st.GoalId.Stamp.Cmp(goalID.Stamp) <= 0)

		if cancelAll || cancelCurrent || cancelBeforeStamp {
			if goalID.Id == st.GoalId.Id {
				goalFound = true
			}

			if gh.SetCancelRequested() {
				args := []reflect.Value{reflect.ValueOf(goalID)}
				fun := reflect.ValueOf(as.cancelCallback)
				numArgsNeeded := fun.Type().NumIn()

				if numArgsNeeded <= 1 {
					fun.Call(args[0:numArgsNeeded])
				}
			}
		}
	}

	if goalID.Id != "" && !goalFound {
		gh := newServerGoalHandlerWithGoalId(as, goalID)
		as.handlers[goalID.Id] = gh
		gh.SetHandlerDestructionTime(ros.Now())
	}

	if goalID.Stamp.Cmp(as.lastCancel) > 0 {
		as.lastCancel = goalID.Stamp
	}
}

// internalGoalCallback recieves the goals from client and checks if
// the goalID already exists in the status list. If not, it will call
// server's goalCallback with goal that was recieved from the client.
func (as *defaultActionServer) internalGoalCallback(goal ActionGoal, event ros.MessageEvent) {
	as.handlersMutex.Lock()
	defer as.handlersMutex.Unlock()

	logger := as.node.Logger()
	goalID := goal.GetGoalId()

	for id, gh := range as.handlers {
		if goalID.Id == id {
			st := gh.GetGoalStatus()
			logger.Debugf("Goal %s was already in the status list with status %+v", goalID.Id, st.Status)
			if st.Status == actionlib_msgs.RECALLING {
				st.Status = actionlib_msgs.RECALLED
				result := as.actionResultType.NewMessage()
				as.PublishResult(st, result)
			}

			gh.SetHandlerDestructionTime(ros.Now())
			return
		}
	}

	id := goalID.Id
	if len(id) == 0 {
		id = as.goalIDGen.generateID()
		goal.SetGoalId(actionlib_msgs.GoalID{
			Id:    id,
			Stamp: goalID.Stamp,
		})
	}

	gh := newServerGoalHandlerWithGoal(as, goal)
	as.handlers[id] = gh
	if !goalID.Stamp.IsZero() && goalID.Stamp.Cmp(as.lastCancel) <= 0 {
		gh.SetCancelled(nil, "timestamp older than last goal cancel")
		return
	}

	args := []reflect.Value{reflect.ValueOf(goal), reflect.ValueOf(event)}
	fun := reflect.ValueOf(as.goalCallback)
	numArgsNeeded := fun.Type().NumIn()

	if numArgsNeeded <= 1 {
		fun.Call(args[0:numArgsNeeded])
	}
}

func (as *defaultActionServer) getHandler(id string) *serverGoalHandler {
	handler := as.handlers[id]
	return handler
}

// RegisterGoalCallback replaces existing goal callback function with newly
// provided goal callback function.
func (as *defaultActionServer) RegisterGoalCallback(goalCb interface{}) {
	as.goalCallback = goalCb
}

func (as *defaultActionServer) RegisterCancelCallback(cancelCb interface{}) {
	as.cancelCallback = cancelCb
}

func (as *defaultActionServer) Shutdown() {
	as.shutdownChan <- struct{}{}
}

package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"std_msgs"
	"sync"

	"github.com/fetchrobotics/rosgo/ros"
)

type defaultActionClient struct {
	started          bool
	node             ros.Node
	action           string
	actionType       ActionType
	actionResult     ros.MessageType
	actionResultType ros.MessageType
	actionFeedback   ros.MessageType
	actionGoal       ros.MessageType
	goalPub          ros.Publisher
	cancelPub        ros.Publisher
	resultSub        ros.Subscriber
	feedbackSub      ros.Subscriber
	statusSub        ros.Subscriber
	logger           ros.Logger
	handlers         []*clientGoalHandler
	handlersMutex    sync.RWMutex
	goalIDGen        *goalIDGenerator
	statusReceived   bool
	callerID         string
}

func newDefaultActionClient(node ros.Node, action string, actType ActionType) *defaultActionClient {
	ac := &defaultActionClient{
		node:           node,
		action:         action,
		actionType:     actType,
		actionResult:   actType.ResultType(),
		actionFeedback: actType.FeedbackType(),
		actionGoal:     actType.GoalType(),
		logger:         node.Logger(),
		statusReceived: false,
		goalIDGen:      newGoalIDGenerator(node.Name()),
	}

	ac.goalPub = node.NewPublisher(fmt.Sprintf("%s/goal", action), actType.GoalType())
	ac.cancelPub = node.NewPublisher(fmt.Sprintf("%s/cancel", action), actionlib_msgs.MsgGoalID)
	ac.resultSub = node.NewSubscriber(fmt.Sprintf("%s/result", action), actType.ResultType(), ac.internalResultCallback)
	ac.feedbackSub = node.NewSubscriber(fmt.Sprintf("%s/feedback", action), actType.FeedbackType(), ac.internalFeedbackCallback)
	ac.statusSub = node.NewSubscriber(fmt.Sprintf("%s/status", action), actionlib_msgs.MsgGoalStatusArray, ac.internalStatusCallback)

	return ac
}

func (ac *defaultActionClient) SendGoal(goal ros.Message, transitionCb, feedbackCb interface{}) ClientGoalHandler {
	if !ac.started {
		ac.logger.Error("[ActionClient] Trying to send a goal on an inactive ActionClient")
	}

	ag := ac.actionType.GoalType().NewMessage().(ActionGoal)
	goalID := actionlib_msgs.GoalID{Id: ac.goalIDGen.generateID(), Stamp: ros.Now()}
	header := std_msgs.Header{Stamp: ros.Now()}

	ag.SetGoal(goal)
	ag.SetGoalId(goalID)
	ag.SetHeader(header)
	ac.PublishActionGoal(ag)

	handler := newClientGoalHandler(ac, ag, transitionCb, feedbackCb)

	ac.handlersMutex.Lock()
	ac.handlers = append(ac.handlers, handler)
	ac.handlersMutex.Unlock()

	return handler
}

func (ac *defaultActionClient) CancelAllGoals() {
	if !ac.started {
		ac.logger.Error("[ActionClient] Trying to cancel goals on an inactive ActionClient")
		return
	}

	ac.cancelPub.Publish(&actionlib_msgs.GoalID{})
}

func (ac *defaultActionClient) CancelAllGoalsBeforeTime(stamp ros.Time) {
	if !ac.started {
		ac.logger.Error("[ActionClient] Trying to cancel goals on an inactive ActionClient")
		return
	}

	cancelMsg := &actionlib_msgs.GoalID{Stamp: stamp}
	ac.cancelPub.Publish(cancelMsg)
}

func (ac *defaultActionClient) Shutdown() {
	ac.handlersMutex.Lock()
	defer ac.handlersMutex.Unlock()

	ac.started = false
	for _, h := range ac.handlers {
		h.Shutdown(false)
	}

	ac.handlers = nil
	ac.node.Shutdown()
}

func (ac *defaultActionClient) PublishActionGoal(ag ActionGoal) {
	if ac.started {
		ac.goalPub.Publish(ag)
	}
}

func (ac *defaultActionClient) PublishCancel(cancel *actionlib_msgs.GoalID) {
	if ac.started {
		ac.cancelPub.Publish(cancel)
	}
}

func (ac *defaultActionClient) WaitForServer(timeout ros.Duration) bool {
	started := false
	ac.logger.Info("[ActionClient] Waiting action server to start")
	rate := ros.CycleTime(ros.NewDuration(0, 10000000))
	waitStart := ros.Now()

LOOP:
	for !started {
		gSubs := ac.goalPub.GetNumSubscribers()
		cSubs := ac.cancelPub.GetNumSubscribers()
		fPubs := ac.feedbackSub.GetNumPublishers()
		rPubs := ac.resultSub.GetNumPublishers()
		sPubs := ac.statusSub.GetNumPublishers()
		started = (gSubs > 0 && cSubs > 0 && fPubs > 0 && rPubs > 0 && sPubs > 0)

		now := ros.Now()
		diff := now.Diff(waitStart)
		if !timeout.IsZero() && diff.Cmp(timeout) >= 0 {
			break LOOP
		}

		rate.Sleep()
	}

	if started {
		ac.started = started
	}

	return started
}

func (ac *defaultActionClient) DeleteGoalHandler(gh *clientGoalHandler) {
	ac.handlersMutex.Lock()
	defer ac.handlersMutex.Unlock()

	for i, h := range ac.handlers {
		if h == gh {
			ac.handlers[i] = ac.handlers[len(ac.handlers)-1]
			ac.handlers[len(ac.handlers)-1] = nil
			ac.handlers = ac.handlers[:len(ac.handlers)-1]
		}
	}
}

func (ac *defaultActionClient) internalResultCallback(result ActionResult, event ros.MessageEvent) {
	ac.handlersMutex.RLock()
	defer ac.handlersMutex.RUnlock()

	for _, h := range ac.handlers {
		if err := h.updateResult(result); err != nil {
			ac.logger.Error(err)
		}
	}
}

func (ac *defaultActionClient) internalFeedbackCallback(feedback ActionFeedback, event ros.MessageEvent) {
	ac.handlersMutex.RLock()
	defer ac.handlersMutex.RUnlock()

	for _, h := range ac.handlers {
		h.updateFeedback(feedback)
	}
}

func (ac *defaultActionClient) internalStatusCallback(statusArr *actionlib_msgs.GoalStatusArray, event ros.MessageEvent) {
	ac.handlersMutex.RLock()
	defer ac.handlersMutex.RUnlock()

	if !ac.statusReceived {
		ac.statusReceived = true
		ac.logger.Debug("Recieved first status message from action server ")
	} else if ac.callerID != event.PublisherName {
		ac.logger.Debug("Previously received status from %s, now from %s. Did the action server change", ac.callerID, event.PublisherName)
	}

	ac.callerID = event.PublisherName
	for _, h := range ac.handlers {
		if err := h.updateStatus(statusArr); err != nil {
			ac.logger.Error(err)
		}
	}
}

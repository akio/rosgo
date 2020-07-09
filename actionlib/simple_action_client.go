package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"reflect"
	"time"

	modular "github.com/edwinhayes/logrus-modular"
	"github.com/edwinhayes/rosgo/ros"
)

const (
	SimpleStatePending uint8 = 0
	SimpleStateActive  uint8 = 1
	SimpleStateDone    uint8 = 2
)

type simpleActionClient struct {
	ac          *defaultActionClient
	simpleState uint8
	gh          ClientGoalHandler
	doneCb      interface{}
	activeCb    interface{}
	feedbackCb  interface{}
	doneChan    chan struct{}
	logger      *modular.ModuleLogger
}

func newSimpleActionClient(node ros.Node, action string, actionType ActionType) *simpleActionClient {
	return &simpleActionClient{
		ac:          newDefaultActionClient(node, action, actionType),
		simpleState: SimpleStateDone,
		doneChan:    make(chan struct{}, 10),
		logger:      node.Logger(),
	}
}

func (sc *simpleActionClient) SendGoal(goal ros.Message, doneCb, activeCb, feedbackCb interface{}) {
	sc.StopTrackingGoal()
	sc.doneCb = doneCb
	sc.activeCb = activeCb
	sc.feedbackCb = feedbackCb

	sc.setSimpleState(SimpleStatePending)
	sc.gh = sc.ac.SendGoal(goal, sc.transitionHandler, sc.feedbackHandler)
}

func (sc *simpleActionClient) SendGoalAndWait(goal ros.Message, executeTimeout, preeptTimeout ros.Duration) (uint8, error) {
	logger := *sc.logger
	sc.SendGoal(goal, nil, nil, nil)
	if !sc.WaitForResult(executeTimeout) {
		logger.Debug("Cancelling goal")
		sc.CancelGoal()
		if sc.WaitForResult(preeptTimeout) {
			logger.Debug("Preempt finished within specified timeout")
		} else {
			logger.Debug("Preempt did not finish within specified timeout")
		}
	}

	return sc.GetState()
}

func (sc *simpleActionClient) WaitForServer(timeout ros.Duration) bool {
	return sc.ac.WaitForServer(timeout)
}

func (sc *simpleActionClient) WaitForResult(timeout ros.Duration) bool {
	logger := *sc.logger
	if sc.gh == nil {
		logger.Errorf("[SimpleActionClient] Called WaitForResult when no goal exists")
		return false
	}

	waitStart := ros.Now()
	waitStart = waitStart.Add(timeout)

LOOP:
	for {
		select {
		case <-sc.doneChan:
			break LOOP
		case <-time.After(100 * time.Millisecond):
		}

		if !timeout.IsZero() && waitStart.Cmp(ros.Now()) <= 0 {
			break LOOP
		}
	}

	return sc.simpleState == SimpleStateDone
}

func (sc *simpleActionClient) GetResult() (ros.Message, error) {
	if sc.gh == nil {
		return nil, fmt.Errorf("called get result when no goal running")
	}

	return sc.gh.GetResult()
}

func (sc *simpleActionClient) GetState() (uint8, error) {
	if sc.gh == nil {
		return actionlib_msgs.LOST, fmt.Errorf("called get state when no goal running")
	}

	status, err := sc.gh.GetGoalStatus()
	if err != nil {
		return actionlib_msgs.LOST, err
	}

	if status == actionlib_msgs.RECALLING {
		status = actionlib_msgs.PENDING
	} else if status == actionlib_msgs.PREEMPTING {
		status = actionlib_msgs.ACTIVE
	}

	return status, nil
}

func (sc *simpleActionClient) GetGoalStatusText() (string, error) {
	if sc.gh == nil {
		return "", fmt.Errorf("called GetGoalStatusText when no goal is running")
	}

	return sc.gh.GetGoalStatusText()
}

func (sc *simpleActionClient) CancelAllGoals() {
	sc.ac.CancelAllGoals()
}

func (sc *simpleActionClient) CancelAllGoalsBeforeTime(stamp ros.Time) {
	sc.ac.CancelAllGoalsBeforeTime(stamp)
}

func (sc *simpleActionClient) CancelGoal() error {
	if sc.gh == nil {
		return nil
	}

	return sc.gh.Cancel()
}

func (sc *simpleActionClient) StopTrackingGoal() {
	sc.gh = nil
}

func (sc *simpleActionClient) transitionHandler(gh ClientGoalHandler) {
	logger := *sc.logger
	commState, err := gh.GetCommState()
	if err != nil {
		logger.Errorf("Error getting CommState: %v", err)
		return
	}

	errMsg := fmt.Errorf("received comm state %s when in simple state %d with SimpleActionClient in NS %s",
		commState, sc.simpleState, sc.ac.node.Name())

	var callbackType string
	var args []reflect.Value
	switch commState {
	case Active:
		switch sc.simpleState {
		case SimpleStatePending:
			sc.setSimpleState(SimpleStateActive)
			callbackType = "active"

		case SimpleStateDone:
			logger.Errorf("[SimpleActionClient] %v", errMsg)
		}

	case Recalling:
		switch sc.simpleState {
		case SimpleStateActive, SimpleStateDone:
			logger.Errorf("[SimpleActionClient] %v", errMsg)
		}

	case Preempting:
		switch sc.simpleState {
		case SimpleStatePending:
			sc.setSimpleState(SimpleStateActive)
			callbackType = "active"

		case SimpleStateDone:
			logger.Errorf("[SimpleActionClient] %v", errMsg)
		}

	case Done:
		switch sc.simpleState {
		case SimpleStatePending, SimpleStateActive:
			sc.setSimpleState(SimpleStateDone)
			sc.sendDone()

			if sc.doneCb == nil {
				break
			}

			status, err := gh.GetGoalStatus()
			if err != nil {
				logger.Errorf("[SimpleActionClient] Error getting status: %v", err)
				break
			}

			result, err := gh.GetResult()
			if err != nil {
				logger.Errorf("[SimpleActionClient] Error getting result: %v", err)
				break
			}

			callbackType = "done"
			args = append(args, reflect.ValueOf(status))
			args = append(args, reflect.ValueOf(result))

		case SimpleStateDone:
			logger.Errorf("[SimpleActionClient] received DONE twice")
		}
	}

	if len(callbackType) > 0 {
		sc.runCallback(callbackType, args)
	}
}

func (sc *simpleActionClient) sendDone() {
	logger := *sc.logger
	select {
	case sc.doneChan <- struct{}{}:
	default:
		logger.Errorf("[SimpleActionClient] Error sending done notification. Channel full.")
	}
}

func (sc *simpleActionClient) feedbackHandler(gh ClientGoalHandler, msg ros.Message) {
	if sc.gh == nil || sc.gh != gh {
		return
	}

	sc.runCallback("feedback", []reflect.Value{reflect.ValueOf(msg)})
}

func (sc *simpleActionClient) setSimpleState(state uint8) {
	logger := *sc.logger
	logger.Debugf("[SimpleActionClient] Transitioning from %d to %d", sc.simpleState, state)
	sc.simpleState = state
}

func (sc *simpleActionClient) runCallback(cbType string, args []reflect.Value) {
	var callback interface{}
	logger := *sc.logger
	switch cbType {
	case "active":
		callback = sc.activeCb
	case "feedback":
		callback = sc.feedbackCb
	case "done":
		callback = sc.doneCb
	default:
		logger.Errorf("[SimpleActionClient] Unknown callback %s", cbType)
	}

	if callback == nil {
		return
	}

	fun := reflect.ValueOf(callback)
	numArgsNeeded := fun.Type().NumIn()

	if numArgsNeeded > len(args) {
		logger.Errorf("[SimpleActionClient] Unexpected arguments:"+
			"callback %s expects %d arguments but %d arguments provided", cbType, numArgsNeeded, len(args))
		return
	}

	logger.Debugf("[SimpleActionClient] Calling %s callback with %d arguments", cbType, len(args))

	fun.Call(args[0:numArgsNeeded])
}

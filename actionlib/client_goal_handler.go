package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"reflect"

	"github.com/fetchrobotics/rosgo/ros"
)

type clientGoalHandler struct {
	actionClient *defaultActionClient
	stateMachine *clientStateMachine
	actionGoal   ActionGoal
	actionGoalID string
	transitionCb interface{}
	feedbackCb   interface{}
	logger       ros.Logger
}

func newClientGoalHandler(ac *defaultActionClient, ag ActionGoal, transitionCb, feedbackCb interface{}) *clientGoalHandler {
	return &clientGoalHandler{
		actionClient: ac,
		stateMachine: newClientStateMachine(),
		actionGoal:   ag,
		actionGoalID: ag.GetGoalId().Id,
		transitionCb: transitionCb,
		feedbackCb:   feedbackCb,
		logger:       ac.logger,
	}
}

func findGoalStatus(statusArr *actionlib_msgs.GoalStatusArray, id string) *actionlib_msgs.GoalStatus {
	var status actionlib_msgs.GoalStatus
	for _, st := range statusArr.StatusList {
		if st.GoalId.Id == id {
			status = st
			break
		}
	}

	return &status
}

func (gh *clientGoalHandler) GetCommState() (CommState, error) {
	if gh.stateMachine == nil {
		return Lost, fmt.Errorf("trying to get state on an inactive ClientGoalHandler")
	}

	return gh.stateMachine.getState(), nil
}

func (gh *clientGoalHandler) GetGoalStatus() (uint8, error) {
	if gh.stateMachine == nil {
		return actionlib_msgs.LOST, fmt.Errorf("trying to get goal status on an inactive ClientGoalHandler")
	}

	return gh.stateMachine.getGoalStatus().Status, nil
}

func (gh *clientGoalHandler) GetGoalStatusText() (string, error) {
	if gh.stateMachine == nil {
		return "", fmt.Errorf("trying to get goal status text on an inactive ClientGoalHandler")
	}

	return gh.stateMachine.getGoalStatus().Text, nil
}

func (gh *clientGoalHandler) GetTerminalState() (uint8, error) {
	if gh.stateMachine == nil {
		return 0, fmt.Errorf("trying to get goal status on inactive clientGoalHandler")
	}

	if gh.stateMachine.state != Done {
		gh.actionClient.logger.Warnf("Asking for terminal state when we are in state %v", gh.stateMachine.state)
	}

	// implement get status
	goalStatus := gh.stateMachine.getGoalStatus().Status
	if goalStatus == actionlib_msgs.PREEMPTED ||
		goalStatus == actionlib_msgs.SUCCEEDED ||
		goalStatus == actionlib_msgs.ABORTED ||
		goalStatus == actionlib_msgs.REJECTED ||
		goalStatus == actionlib_msgs.RECALLED ||
		goalStatus == actionlib_msgs.LOST {

		return goalStatus, nil
	}

	gh.actionClient.logger.Warnf("Asking for terminal state when latest goal is in %v", goalStatus)
	return actionlib_msgs.LOST, nil
}

func (gh *clientGoalHandler) GetResult() (ros.Message, error) {
	if gh.stateMachine == nil {
		return nil, fmt.Errorf("trying to get goal status on inactive clientGoalHandler")
	}

	result := gh.stateMachine.getGoalResult()

	if result == nil {
		return nil, fmt.Errorf("trying to get result when no result has been recieved")
	}

	return result.GetResult(), nil
}

func (gh *clientGoalHandler) Resend() error {
	if gh.stateMachine == nil {
		return fmt.Errorf("trying to call resend on inactive client goal hanlder")
	}

	gh.actionClient.goalPub.Publish(gh.actionGoal)
	return nil
}

func (gh *clientGoalHandler) IsExpired() bool {
	return gh.stateMachine == nil
}

func (gh *clientGoalHandler) Cancel() error {
	if gh.stateMachine == nil {
		return fmt.Errorf("trying to call cancel on inactive client goal hanlder")
	}

	cancelMsg := &actionlib_msgs.GoalID{
		Stamp: ros.Now(),
		Id:    gh.actionGoalID}

	gh.actionClient.cancelPub.Publish(cancelMsg)
	gh.stateMachine.transitionTo(WaitingForCancelAck, gh, gh.transitionCb)
	return nil
}

func (gh *clientGoalHandler) Shutdown(deleteFromManager bool) {
	gh.stateMachine = nil
	if deleteFromManager {
		gh.actionClient.DeleteGoalHandler(gh)
	}
}

func (gh *clientGoalHandler) updateFeedback(af ActionFeedback) {
	if gh.actionGoalID != af.GetStatus().GoalId.Id {
		return
	}

	if gh.feedbackCb != nil && gh.stateMachine.getState() != Done {
		fun := reflect.ValueOf(gh.feedbackCb)
		args := []reflect.Value{reflect.ValueOf(gh), reflect.ValueOf(af.GetFeedback())}
		numArgsNeeded := fun.Type().NumIn()

		if numArgsNeeded == 2 {
			fun.Call(args)
		}
	}
}

func (gh *clientGoalHandler) updateResult(result ActionResult) error {
	if gh.actionGoalID != result.GetStatus().GoalId.Id {
		return nil
	}

	status := result.GetStatus()
	state := gh.stateMachine.getState()

	gh.stateMachine.setGoalStatus(status.GoalId, status.Status, status.Text)
	gh.stateMachine.setGoalResult(result)

	if state == WaitingForGoalAck ||
		state == WaitingForCancelAck ||
		state == Pending ||
		state == Active ||
		state == WaitingForResult ||
		state == Recalling ||
		state == Preempting {

		statusArr := new(actionlib_msgs.GoalStatusArray)
		statusArr.StatusList = append(statusArr.StatusList, result.GetStatus())
		if err := gh.updateStatus(statusArr); err != nil {
			return err
		}

		gh.stateMachine.transitionTo(Done, gh, gh.transitionCb)
		return nil
	} else if state == Done {
		return fmt.Errorf("got a result when we are in the `DONE` state")
	} else {
		return fmt.Errorf("unknown state %v", state)
	}
}

func (gh *clientGoalHandler) updateStatus(statusArr *actionlib_msgs.GoalStatusArray) error {
	state := gh.stateMachine.getState()
	if state == Done {
		return nil
	}

	status := findGoalStatus(statusArr, gh.actionGoalID)
	if status == nil {
		if state != WaitingForGoalAck &&
			state != WaitingForResult &&
			state != Done {

			gh.logger.Warn("Transitioning goal to `Lost`")
			gh.stateMachine.setAsLost()
			gh.stateMachine.transitionTo(Done, gh, gh.transitionCb)
		}
		return nil
	}

	gh.stateMachine.setGoalStatus(status.GoalId, status.Status, status.Text)
	nextStates, err := gh.stateMachine.getTransitions(*status)
	if err != nil {
		return fmt.Errorf("error getting transitions: %v", err)
	}

	for e := nextStates.Front(); e != nil; e = e.Next() {
		gh.stateMachine.transitionTo(e.Value.(CommState), gh, gh.transitionCb)
	}

	return nil
}

package actionlib

import (
	"actionlib_msgs"
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

type CommState uint8

const (
	WaitingForGoalAck CommState = iota
	Pending
	Active
	WaitingForResult
	WaitingForCancelAck
	Recalling
	Preempting
	Done
	Lost
)

func (cs CommState) String() string {
	switch cs {
	case WaitingForGoalAck:
		return "WAITING_FOR_GOAL_ACK"
	case Pending:
		return "PENDING"
	case Active:
		return "ACTIVE"
	case WaitingForResult:
		return "WAITING_FOR_RESULT"
	case WaitingForCancelAck:
		return "WAITING_FOR_CANCEL_ACK"
	case Recalling:
		return "RECALLING"
	case Preempting:
		return "PREEMPTING"
	case Done:
		return "DONE"
	case Lost:
		return "LOST"
	default:
		return "UNKNOWN"
	}
}

type clientStateMachine struct {
	state      CommState
	goalStatus actionlib_msgs.GoalStatus
	goalResult ActionResult
	mutex      sync.RWMutex
}

func newClientStateMachine() *clientStateMachine {
	return &clientStateMachine{
		state:      WaitingForGoalAck,
		goalStatus: actionlib_msgs.GoalStatus{Status: actionlib_msgs.PENDING},
	}
}

func (sm *clientStateMachine) getState() CommState {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return sm.state
}

func (sm *clientStateMachine) getGoalStatus() actionlib_msgs.GoalStatus {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return sm.goalStatus
}

func (sm *clientStateMachine) getGoalResult() ActionResult {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return sm.goalResult
}

func (sm *clientStateMachine) setState(state CommState) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.state = state
}

func (sm *clientStateMachine) setGoalStatus(id actionlib_msgs.GoalID, status uint8, text string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.goalStatus.GoalId = id
	sm.goalStatus.Status = status
	sm.goalStatus.Text = text
}

func (sm *clientStateMachine) setGoalResult(result ActionResult) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.goalResult = result
}

func (sm *clientStateMachine) setAsLost() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.goalStatus.Status = uint8(Lost)
}

func (sm *clientStateMachine) transitionTo(state CommState, gh ClientGoalHandler, callback interface{}) {
	sm.setState(state)
	if callback != nil {
		fun := reflect.ValueOf(callback)
		args := []reflect.Value{reflect.ValueOf(gh)}
		numArgsNeeded := fun.Type().NumIn()

		if numArgsNeeded <= 1 {
			fun.Call(args[:numArgsNeeded])
		}
	}
}

func (sm *clientStateMachine) getTransitions(goalStatus actionlib_msgs.GoalStatus) (stateList list.List, err error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	status := goalStatus.Status

	switch sm.state {
	case WaitingForGoalAck:
		switch status {
		case actionlib_msgs.PENDING:
			stateList.PushBack(Pending)
			break
		case actionlib_msgs.ACTIVE:
			stateList.PushBack(Active)
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(Pending)
			stateList.PushBack(WaitingForCancelAck)
			break
		case actionlib_msgs.RECALLING:
			stateList.PushBack(Pending)
			stateList.PushBack(Recalling)
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(Pending)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			break
		}
		break

	case Pending:
		switch status {
		case actionlib_msgs.PENDING:
			break
		case actionlib_msgs.ACTIVE:
			stateList.PushBack(Active)
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.RECALLING:
			stateList.PushBack(Recalling)
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Active)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Active)
			stateList.PushBack(Preempting)
			break
		}
		break
	case Active:
		switch status {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("invalid transition from Active to Pending")
			break
		case actionlib_msgs.ACTIVE:
			break
		case actionlib_msgs.REJECTED:
			err = fmt.Errorf("invalid transition from Active to Rejected")
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("invalid transition from Active to Recalling")
			break
		case actionlib_msgs.RECALLED:
			err = fmt.Errorf("invalid transition from Active to Recalled")
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Preempting)
			break
		}
		break
	case WaitingForResult:
		switch status {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("invalid transition from WaitingForResult to Pending")
			break
		case actionlib_msgs.ACTIVE:
			break
		case actionlib_msgs.REJECTED:
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("invalid transition from WaitingForResult to Recalling")
			break
		case actionlib_msgs.RECALLED:
			break
		case actionlib_msgs.PREEMPTED:
			break
		case actionlib_msgs.SUCCEEDED:
			break
		case actionlib_msgs.ABORTED:
			break
		case actionlib_msgs.PREEMPTING:
			err = fmt.Errorf("invalid transition from WaitingForResult to Preempting")
			break
		}
		break
	case WaitingForCancelAck:
		switch status {
		case actionlib_msgs.PENDING:
			break
		case actionlib_msgs.ACTIVE:
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.RECALLING:
			stateList.PushBack(Recalling)
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Recalling)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Preempting)
			break
		}
		break
	case Recalling:
		switch status {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("invalid transition from Recalling to Pending")
			break
		case actionlib_msgs.ACTIVE:
			err = fmt.Errorf("invalid transition from Recalling to Active")
			break
		case actionlib_msgs.REJECTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.RECALLING:
			break
		case actionlib_msgs.RECALLED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(Preempting)
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			stateList.PushBack(Preempting)
			break
		}
		break
	case Preempting:
		switch status {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("invalid transition from Preempting to Pending")
			break
		case actionlib_msgs.ACTIVE:
			err = fmt.Errorf("invalid transition from Preempting to Active")
			break
		case actionlib_msgs.REJECTED:
			err = fmt.Errorf("invalid transition from Preempting to Rejected")
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("invalid transition from Preempting to Recalling")
			break
		case actionlib_msgs.RECALLED:
			err = fmt.Errorf("invalid transition from Preempting to Recalled")
			break
		case actionlib_msgs.PREEMPTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.SUCCEEDED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.ABORTED:
			stateList.PushBack(WaitingForResult)
			break
		case actionlib_msgs.PREEMPTING:
			break
		}
		break
	case Done:
		switch status {
		case actionlib_msgs.PENDING:
			err = fmt.Errorf("invalid transition from Done to Pending")
			break
		case actionlib_msgs.ACTIVE:
			err = fmt.Errorf("invalid transition from Done to Active")
			break
		case actionlib_msgs.REJECTED:
			break
		case actionlib_msgs.RECALLING:
			err = fmt.Errorf("invalid transition from Done to Recalling")
			break
		case actionlib_msgs.RECALLED:
			break
		case actionlib_msgs.PREEMPTED:
			break
		case actionlib_msgs.SUCCEEDED:
			break
		case actionlib_msgs.ABORTED:
			break
		case actionlib_msgs.PREEMPTING:
			err = fmt.Errorf("invalid transition from Done to Preempting")
			break
		}
		break
	}

	return
}

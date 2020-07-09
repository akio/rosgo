package actionlib

import (
	"actionlib_msgs"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/fetchrobotics/rosgo/ros"
)

type simpleActionServer struct {
	actionServer          *defaultActionServer
	currentGoal           *serverGoalHandler
	nextGoal              *serverGoalHandler
	goalMutex             sync.RWMutex
	newGoal               bool
	preemptRequest        bool
	newGoalPreemptRequest bool
	logger                ros.Logger
	goalCallback          interface{}
	preemptCallback       interface{}
	executeCb             interface{}
	executorCh            chan struct{}
}

func newSimpleActionServer(node ros.Node, action string, actType ActionType, executeCb interface{}, start bool) *simpleActionServer {
	s := new(simpleActionServer)
	s.actionServer = newDefaultActionServer(node, action, actType, s.internalGoalCallback, s.internalPreemptCallback, start)
	s.newGoal = false
	s.preemptRequest = false
	s.newGoalPreemptRequest = false
	s.executeCb = executeCb
	s.logger = node.Logger()
	s.executorCh = make(chan struct{}, 100)
	return s
}

func (s *simpleActionServer) Start() {
	if s.executeCb != nil {
		go s.goalExecutor()
	}

	go s.actionServer.Start()
}

func (s *simpleActionServer) IsNewGoalAvailable() bool {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	return s.newGoal
}

func (s *simpleActionServer) IsPreemptRequested() bool {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	return s.preemptRequest
}

func (s *simpleActionServer) AcceptNewGoal() (ros.Message, error) {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if !s.newGoal || s.nextGoal == nil {
		return nil, fmt.Errorf("attempting to accept the next goal when a new goal is not available")
	}

	// check if we need to send a preempted message for the goal that we're currently pursuing
	if s.IsActive() && s.currentGoal != nil && s.currentGoal.NotEqual(s.nextGoal) {
		s.currentGoal.SetCancelled(s.GetDefaultResult(),
			"This goal was canceled because another goal was received by the simple action server")
	}

	s.logger.Debug("Accepting a new goal")

	// accept the next goal
	s.currentGoal = s.nextGoal
	s.newGoal = false

	// set preempt to request to equal the preempt state of the new goal
	s.preemptRequest = s.newGoalPreemptRequest
	s.newGoalPreemptRequest = false

	// set the status of the current goal to be active
	s.currentGoal.SetAccepted("This goal has been accepted by the simple action server")

	return s.currentGoal.GetGoal(), nil
}

func (s *simpleActionServer) IsActive() bool {
	if s.currentGoal == nil || s.currentGoal.GetGoalId().Id == "" {
		return false
	}

	status := s.currentGoal.GetGoalStatus().Status
	if status == actionlib_msgs.ACTIVE || status == actionlib_msgs.PREEMPTING {
		return true
	}

	return false
}

func (s *simpleActionServer) SetSucceeded(result ros.Message, text string) error {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.SetSucceeded(result, text)
}

func (s *simpleActionServer) SetAborted(result ros.Message, text string) error {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.SetAborted(result, text)
}

func (s *simpleActionServer) SetPreempted(result ros.Message, text string) error {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if result == nil {
		result = s.GetDefaultResult()
	}

	return s.currentGoal.SetCancelled(result, text)
}

func (s *simpleActionServer) PublishFeedback(feedback ros.Message) {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	s.currentGoal.PublishFeedback(feedback)
}

func (s *simpleActionServer) GetDefaultResult() ros.Message {
	return s.actionServer.actionResultType.NewMessage()
}

func (s *simpleActionServer) RegisterGoalCallback(cb interface{}) error {
	if s.executeCb != nil {
		return fmt.Errorf("execute callback if present: not registering goal callback")
	}

	s.goalCallback = cb

	return nil
}

func (s *simpleActionServer) RegisterPreemptCallback(cb interface{}) {
	s.preemptCallback = cb
}

func (s *simpleActionServer) internalGoalCallback(ag ActionGoal) {
	goalHandler := s.actionServer.getHandler(ag.GetGoalId().Id)
	s.logger.Infof("[SimpleActionServer] Server received new goal with id %s", goalHandler.GetGoalId().Id)

	var goalStamp, nextGoalStamp ros.Time
	goalStamp = goalHandler.GetGoalId().Stamp
	if s.nextGoal != nil {
		nextGoalStamp = s.nextGoal.GetGoalId().Stamp
	}

	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	if (s.currentGoal == nil || goalStamp.Cmp(s.currentGoal.GetGoalId().Stamp) >= 0) &&
		(s.nextGoal == nil || nextGoalStamp.Cmp(s.currentGoal.GetGoalId().Stamp) >= 0) {

		if (s.nextGoal != nil) &&
			(s.currentGoal == nil || s.nextGoal.NotEqual(s.currentGoal)) {
			s.nextGoal.SetCancelled(s.GetDefaultResult(),
				"This goal was canceled because another goal was received by the simple action server")
		}

		s.nextGoal = goalHandler
		s.newGoal = true
		s.newGoalPreemptRequest = false
		goal := goalHandler.GetGoal()
		args := []reflect.Value{reflect.ValueOf(goal)}

		if s.IsActive() {
			s.preemptRequest = true
			if err := s.runCallback("preempt", args); err != nil {
				s.logger.Error(err)
			}
		}

		if err := s.runCallback("goal", args); err != nil {
			s.logger.Error(err)
		}

		// notify executor that a new goal is available
		select {
		case s.executorCh <- struct{}{}:
		default:
			s.logger.Error("[SimpleActionServer] Exectuor new goal notification error: Channel full.")
		}
	} else {
		goalHandler.SetCancelled(s.GetDefaultResult(),
			"This goal was canceled because another goal was received by the simple action server")
	}
}

func (s *simpleActionServer) internalPreemptCallback(gID *actionlib_msgs.GoalID) {
	s.goalMutex.Lock()
	defer s.goalMutex.Unlock()

	goalHandler := s.actionServer.getHandler(gID.Id)
	s.logger.Infof("[SimpleActionServer] Server received preempt call for goal with id %s",
		goalHandler.GetGoalId().Id)

	if goalHandler.GetGoalId().Id == s.currentGoal.GetGoalId().Id {
		s.preemptRequest = true
		goal := goalHandler.GetGoal()
		args := []reflect.Value{reflect.ValueOf(goal)}
		if err := s.runCallback("preempt", args); err != nil {
			s.logger.Error(err)
		}
	} else {
		s.newGoalPreemptRequest = true
	}
}

func (s *simpleActionServer) goalExecutor() {
	intervalCh := time.NewTicker(1 * time.Second)
	defer intervalCh.Stop()

	for s.actionServer.node.OK() {
		select {
		case <-s.executorCh:
			if err := s.execute(); err != nil {
				s.logger.Error(err)
				return
			}

		case <-intervalCh.C:
			if err := s.execute(); err != nil {
				s.logger.Error(err)
				return
			}
		}
	}
}

func (s *simpleActionServer) execute() error {
	if s.IsActive() {
		return fmt.Errorf("should never reach this code with an active goal")
	}

	if s.IsNewGoalAvailable() {
		goal, err := s.AcceptNewGoal()
		if err != nil {
			return err
		}

		if s.executeCb == nil {
			return fmt.Errorf("execute callback must exist. This is a bug in SimpleActionServer")
		}

		args := []reflect.Value{reflect.ValueOf(goal)}
		if err := s.runCallback("execute", args); err != nil {
			return err
		}

		if s.IsActive() {
			s.logger.Warn("Your executeCallback did not set the goal to a terminal status." +
				"This is a bug in your ActionServer implementation. Fix your code!" +
				"For now, the ActionServer will set this goal to aborted")
			if err := s.SetAborted(nil, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *simpleActionServer) runCallback(cbType string, args []reflect.Value) error {
	var callback interface{}
	switch cbType {
	case "goal":
		callback = s.goalCallback
	case "preempt":
		callback = s.preemptCallback
	case "execute":
		callback = s.executeCb
	default:
		return fmt.Errorf("unknown callback type called")
	}

	if callback == nil {
		return nil
	}

	fun := reflect.ValueOf(callback)
	numArgsNeeded := fun.Type().NumIn()

	if numArgsNeeded <= 1 {
		fun.Call(args[0:numArgsNeeded])
	} else {
		return fmt.Errorf("unexepcted number of arguments for callback")
	}

	return nil
}

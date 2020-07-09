package actionlib

import (
	"fmt"
	"sync"

	"github.com/fetchrobotics/rosgo/ros"
)

type goalIDGenerator struct {
	goals      int
	goalsMutex sync.RWMutex
	nodeName   string
}

func newGoalIDGenerator(nodeName string) *goalIDGenerator {
	return &goalIDGenerator{
		nodeName: nodeName,
	}
}

func (g *goalIDGenerator) generateID() string {
	g.goalsMutex.Lock()
	defer g.goalsMutex.Unlock()

	g.goals++

	timeNow := ros.Now()
	return fmt.Sprintf("%s-%d-%d-%d", g.nodeName, g.goals, timeNow.Sec, timeNow.NSec)
}

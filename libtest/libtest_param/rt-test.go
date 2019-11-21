package libtest_param

import (
	"github.com/edwinhayes/rosgo/ros"
	"log"
	"os"
	"testing"
)

//RTTest creates a node and makes ros api calls to the ros parameter server
//Various parameter API are called are checked for error return, and bad return values
func RTTest(t *testing.T) {
	node, err := ros.NewNode("/test_param", os.Args)
	if err != nil {
		t.Error("Failed to initialize node")
		return
	}
	defer node.Shutdown()

	//Calls HasParam API for paramater /rosdistro
	if hasParam, err := node.HasParam("/rosdistro"); err != nil {
		t.Error("HasParam api call failed", err)
	} else {
		if !hasParam {
			t.Error("No parameter set")
		}
	}

	//Calls SearchParam API for paramater /rosdistro
	if foundKey, err := node.SearchParam("rosdistro"); err != nil {
		t.Error("SearchParam api call failed", err)
	} else {
		if foundKey != "/rosdistro" {
			t.Error("No parameter found")
		}
	}

	//Calls GetParam API for paramater /rosdistro and checks return is string
	if param, err := node.GetParam("/rosdistro"); err != nil {
		t.Error("GetParam api call failed", err)
	} else {
		if value, ok := param.(string); !ok {
			t.Error("Bad parameter Recieved", value, ok)
		}
	}

	//Calls SetParam API to set a new parameter /test_param
	if err := node.SetParam("/test_param", 42); err != nil {
		t.Error("SetParam api call failed", err)
	}

	//Calls GetParam API on /test_param we just set
	if param, err := node.GetParam("/test_param"); err != nil {
		t.Error("GetParam api call failed", err)
	} else {
		if value, ok := param.(int32); ok {
			if value != 42 {
				t.Error("Test Param value is wrong", value)
			}
		} else {
			t.Error("Could not retrieve test param value")
		}
	}

	if err := node.DeleteParam("/test_param"); err != nil {
		log.Fatalf("DeleteParam failed: %v", err)
	}
	return
}

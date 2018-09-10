package ros

import (
	"fmt"
	"github.com/akio/rosgo/xmlrpc"
)

func callRosApi(calleeUri string, method string, args ...interface{}) (interface{}, error) {
	result, err := xmlrpc.Call(calleeUri, method, args...)
	if err != nil {
		return nil, err
	}

	var ok bool
	var xs []interface{}
	var code int32
	var message string
	var value interface{}
	if xs, ok = result.([]interface{}); !ok {
		return nil, fmt.Errorf("Malformed ROS API result.")
	}
	if len(xs) != 3 {
		err := fmt.Errorf("Malformed ROS API result. Length must be 3 but %d", len(xs))
		return nil, err
	}
	if code, ok = xs[0].(int32); !ok {
		return nil, fmt.Errorf("Status code is not int.")
	}
	if message, ok = xs[1].(string); !ok {
		return nil, fmt.Errorf("Message is not string.")
	}
	value = xs[2]

	if code != ApiStatusSuccess {
		err := fmt.Errorf("ROS Master API call failed with code %d: %s", code, message)
		return nil, err
	}
	return value, nil
}

// Build XMLRPC ready array from ROS API result triplet.
func buildRosApiResult(code int32, message string, value interface{}) interface{} {
	result := make([]interface{}, 3)
	result[0] = code
	result[1] = message
	result[2] = value
	return result
}

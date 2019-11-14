package ros

import (
	"fmt"
	"github.com/edwinhayes/rosgo/xmlrpc"
)

//callRosApi performs an XML-RPC call to the ROS system. CalleeUri is the address to send the request
//Method is the method to be called in the request. Args is an interface of values that are required
//by the method call. Returns interface of the XML response from callee.
func callRosAPI(calleeURI string, method string, args ...interface{}) (interface{}, error) {
	result, err := xmlrpc.Call(calleeURI, method, args...)
	if err != nil {
		return nil, err
	}

	var ok bool
	var xs []interface{}
	var code int32
	var message string
	var value interface{}
	if xs, ok = result.([]interface{}); !ok {
		return nil, fmt.Errorf("malformed ROS API result")
	}
	if len(xs) != 3 {
		err := fmt.Errorf("Malformed ROS API result. Length must be 3 but %d", len(xs))
		return nil, err
	}
	if code, ok = xs[0].(int32); !ok {
		return nil, fmt.Errorf("status code is not int")
	}
	if message, ok = xs[1].(string); !ok {
		return nil, fmt.Errorf("message is not string")
	}
	value = xs[2]

	if code != APIStatusSuccess {
		err := fmt.Errorf("ROS Master API call failed with code %d: %s", code, message)
		return nil, err
	}
	return value, nil
}

// Build XMLRPC ready array from ROS API result triplet.
func buildRosAPIResult(code int32, message string, value interface{}) interface{} {
	result := make([]interface{}, 3)
	result[0] = code
	result[1] = message
	result[2] = value
	return result
}

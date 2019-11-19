package main

//go:generate gengo msg std_msgs/String
import (
	"github.com/edwinhayes/rosgo/libtest/libtest_listener_with_event"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_listener_with_event.RTTest(t)
}

package main

//go:generate gengo msg rosgraph_msgs/Log
import (
	"github.com/edwinhayes/rosgo/libtest/libtest_listener"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_listener.RTTest(t)
}

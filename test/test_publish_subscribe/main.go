package main

//go:generate gengo msg std_msgs/String
import (
	"github.com/edwinhayes/rosgo/libtest/libtest_publish_subscribe"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_publish_subscribe.RTTest(t)
}

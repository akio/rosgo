package main

//go:generate gengo msg std_msgs/String
import (
	"github.com/edwinhayes/rosgo/libtest/libtest_service"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_service.RTTest(t)
}

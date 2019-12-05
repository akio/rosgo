package main

import (
	"github.com/edwinhayes/rosgo/libtest/libtest_dynamic_message"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_dynamic_message.RTTest(t)
}

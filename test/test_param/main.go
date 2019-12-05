package main

import (
	"github.com/edwinhayes/rosgo/libtest/libtest_param"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_param.RTTest(t)
}

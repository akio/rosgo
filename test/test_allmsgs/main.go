package main

import (
	"github.com/edwinhayes/rosgo/libtest/libtest_allmsgs"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_allmsgs.RTTest(t)
}

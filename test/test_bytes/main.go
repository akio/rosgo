package main

import (
	"github.com/edwinhayes/rosgo/libtest/libtest_bytes"
	"testing"
)

func main() {
	t := new(testing.T)
	libtest_bytes.RTTest(t)
}

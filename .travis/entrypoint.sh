#!/bin/bash
source /opt/ros/kinetic/setup.bash
export PATH=$PWD/bin:/usr/local/go/bin:$PATH
export GOPATH=$PWD:/usr/local/go

roscore &
go install github.com/akio/rosgo/gengo
go generate github.com/akio/rosgo/test/test_message
go test github.com/akio/rosgo/xmlrpc
go test github.com/akio/rosgo/xmlrpc
go test github.com/akio/rosgo/ros
go install github.com/akio/rosgo/test/test_message
go test github.com/akio/rosgo/test/test_message


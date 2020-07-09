#!/bin/bash
source /opt/ros/melodic/setup.bash
export PATH=$PWD/bin:/usr/local/go/bin:$PATH
export GOPATH=$PWD:/usr/local/go

roscore &
go install github.com/edwinhayes/rosgo/gengo
go generate github.com/edwinhayes/rosgo/test/test_message
go test github.com/edwinhayes/rosgo/xmlrpc
go test github.com/edwinhayes/rosgo/ros
go test github.com/edwinhayes/rosgo/test/test_message


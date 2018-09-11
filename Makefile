all: test

.PHONY: test

test:
	go test github.com/akio/rosgo/test/test_message
	go install github.com/akio/rosgo/test/test_listener
	go install github.com/akio/rosgo/test/test_listener_with_event
	go install github.com/akio/rosgo/test/test_param
	go install github.com/akio/rosgo/test/test_server
	go install github.com/akio/rosgo/test/test_talker
	go install github.com/akio/rosgo/test/test_talker_with_callbacks


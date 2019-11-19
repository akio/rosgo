package main

// Disable code generation, since we're trying to void doing this:   //go:generate gengo msg std_msgs/String

// IMPORT REQUIRED PACKAGES.

// TODO - Why is the syntax for import different to everywhere else?

import (
	"fmt"
	"github.com/edwinhayes/rosgo/ros"
	"github.com/edwinhayes/rosgo/libtest/libtest_talker"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// DEFINE PUBLIC STRUCTURES.

// DEFINE PRIVATE STRUCTURES.

// DEFINE PUBLIC GLOBALS.

// DEFINE PRIVATE GLOBALS.

var g_node ros.Node

var subscribers map[string]ros.Subscriber
var publishers map[string]ros.Publisher

// DEFINE PUBLIC STATIC FUNCTIONS.

// DEFINE PRIVATE STATIC FUNCTIONS.

func callback(msg *ros.DynamicMessage, event ros.MessageEvent) {
	pub_name := strings.Trim(event.PublisherName, "/")
	topic_name := strings.Trim(event.ConnectionHeader["topic"], "/")
	topic_type := msg.Type().Name()
	g_node.Logger().Info("Received from ", pub_name, ": ", topic_name, " : ", msg)

	// Try not to loopback into oblivion...
	if pub_name == g_node.Name() {
		return
	}

	// Check whether we already created a publisher for this topic.
	var pub ros.Publisher
	var ok bool
	if pub, ok = publishers[topic_name]; !ok {
		// Create a publisher for rebroadcasting the messages we recieve.
		var out_msg *ros.DynamicMessageType
		var err error
		if out_msg, err = ros.NewDynamicMessageType(topic_type); err != nil {
			g_node.Logger().Error("Oh noes!")
			return
		}
		pub = g_node.NewPublisher("echo_"+topic_name, out_msg)
		publishers[topic_name] = pub
		g_node.Logger().Info(g_node.Name(), " now has ", len(publishers), " publishers echoing topics.")
	}

	// Republish the message, just to show that we know how to send it.
	out_msg := ros.Message(msg)
	pub.Publish(out_msg)
}

func poll_for_topics(node ros.Node, quit <-chan bool) {
	// Create a ticker to tell us to do stuff periodically.
	ticker := time.NewTicker(1 * time.Second)

	// Loop forever, or until rx on chan quit.
	node.Logger().Info("Starting goroutine to poll for topics...")
	for {
		select {
		case <-quit:
			node.Logger().Info("Stopping polling for topics...")
			return
		case <-ticker.C:
			// Fetch list of available topics (i.e. those with publishers) from the master.
			topic_list := node.GetPublishedTopics("")

			// Try to iterate over each of the topics in the list.
			for _, v := range topic_list {
				topic := v.([]interface{})
				topic_name := topic[0].(string)
				topic_type := topic[1].(string)

				// Check if we have a subscriber for this topic already.
				if _, ok := subscribers[topic_name]; !ok {
					// Apparently not, so we try to subscribe.
					node.Logger().Info("Attempting to subscribe to topic: ", topic_name)

					// Create a generic message, which tries to look up the important checksum via gengo.
					var m *ros.DynamicMessageType
					var err error
					if m, err = ros.NewDynamicMessageType(topic_type); err != nil {
						node.Logger().Error("Couldn't set message type: ", topic_type, " : ", err)
						continue
					}

					// Generate schema for the topic, and print it out.
					var schema []byte
					if schema, err = m.GenerateJSONSchema("/ros", topic_name); err != nil {
						node.Logger().Error("Couldn't generate scheme: ", topic_name, " : ", err)
						continue
					}
					node.Logger().Info("Schema for ", topic_name)
					node.Logger().Info(string(schema))

					// Then subscribe to the topic, and if we're successful, keep a note so we don't try to subscribe again.
					s := node.NewSubscriber(topic_name, m, callback)
					if s != nil {
						subscribers[topic_name] = s
					}
				}
			}
		}
	}
	// Not all done, since defer?
}

func main() {
	// Run diagnostic tests.
	t := new (testing.T)
	libtest_talker.RTTest(t)
	if t.Failed() {
		fmt.Println("rosgo self-test failed.")
		os.Exit(-2)
	}

	// Create our node.
	node_name := "thinger_" + strconv.Itoa(os.Getpid())
	node, err := ros.NewNode(node_name, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	node.Logger().Info("Created new node: ", node_name)
	g_node = node

	// Configure node logging.
	node.Logger().SetSeverity(ros.LogLevelInfo)

	// Change the package search path so we can use custom messages.
	ros.SetRuntimePackagePath(ros.GetRuntimePackagePath() + ":/home/ubuntu/environment/goenv/src/github.com/edwinhayes/rosgo/test/test_talker/vendor")

	// We'll keep lists of ROS subscribers and publishers, so we can identify topics which we still need to subscribe to or publish.
	subscribers = make(map[string]ros.Subscriber)
	publishers = make(map[string]ros.Publisher)

	// Spawn a routine to look for new topics which get published.
	quit_poll_for_topics := make(chan bool)
	go poll_for_topics(node, quit_poll_for_topics)
	defer close(quit_poll_for_topics)

	// Setup a signal handler to catch the keyboard interrupt.
	quit_mainloop := make(chan os.Signal, 2)
	signal.Notify(quit_mainloop, os.Interrupt, syscall.SIGTERM)

	// Wait forever.
	node.Logger().Info("Spinning...")
	defer node.Logger().Info("Shutting down...")
	for {
		select {
		case <-quit_mainloop:
			node.Logger().Debug("Received SIGTERM.")
			return
		default:
			node.SpinOnce()
		}
	}

	// Not all done, since defer?
}

// ALL DONE.

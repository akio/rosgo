package main

// Disable code generation, since we're trying to void doing this:   //go:generate gengo msg std_msgs/String

// IMPORT REQUIRED PACKAGES.

// TODO - Why is the syntax for import different to everywhere else?

import (
	"fmt"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DEFINE PUBLIC STRUCTURES.

// DEFINE PRIVATE STRUCTURES.

// DEFINE PUBLIC GLOBALS.

// DEFINE PRIVATE GLOBALS.

var subscribers map[string]ros.Subscriber

// DEFINE PUBLIC STATIC FUNCTIONS.

// DEFINE PRIVATE STATIC FUNCTIONS.

func callback(msg *ros.GenericMessageType) {
	fmt.Printf("Received: %s\n", msg.Fields)
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
					m := new(ros.GenericMessageType)
					m.SetMessageType(topic_type)
					s := node.NewSubscriber("/chatter", m, callback)
					// If we subscribe successfully, keep note so we don't have to subscribe again.
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
	// Create our node.
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()

	// Configure node logging.
	node.Logger().SetSeverity(ros.LogLevelDebug)

	// We'll keep a list of ROS subscribers, so we can identify topics which we still need to subscribe to.
	subscribers = make(map[string]ros.Subscriber)

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

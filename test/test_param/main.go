package main

import (
	"fmt"
	"github.com/fetchrobotics/rosgo/ros"
	"log"
	"os"
)

func main() {
	node, err := ros.NewNode("/test_param", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()

	if hasParam, err := node.HasParam("/rosdistro"); err != nil {
		log.Fatalf("HasParam failed: %v", err)
	} else {
		if !hasParam {
			log.Fatal("HasParam() failed.")
		}
	}

	if foundKey, err := node.SearchParam("rosdistro"); err != nil {
		log.Fatalf("SearchParam failed: %v", err)
	} else {
		if foundKey != "/rosdistro" {
			log.Fatal("SearchParam() failed.")
		}
	}

	if param, err := node.GetParam("/rosdistro"); err != nil {
		log.Fatalf("GetParam: %v", err)
	} else {
		if value, ok := param.(string); !ok {
			log.Fatal("GetParam() failed.")
		} else {
			if value != "kinetic\n" {
				log.Fatalf("Expected 'kinetic\\n' but '%s'", value)
			}
		}
	}

	if err := node.SetParam("/test_param", 42); err != nil {
		log.Fatalf("SetParam failed: %v", err)
	}

	if param, err := node.GetParam("/test_param"); err != nil {
		log.Fatalf("GetParam failed: %v", err)
	} else {
		if value, ok := param.(int32); ok {
			if value != 42 {
				log.Fatalf("Expected 42 but %d", value)
			}
		} else {
			log.Fatal("GetParam('/test_param') failed.")
		}
	}

	if err := node.DeleteParam("/test_param"); err != nil {
		log.Fatalf("DeleteParam failed: %v", err)
	}

	log.Print("Success")
}

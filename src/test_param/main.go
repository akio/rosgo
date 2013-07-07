package main

import (
    "ros"
    "log"
)

func main() {
    node := ros.NewNode("/test_param")
    defer node.Shutdown()

    if hasParam, err := node.HasParam("/rosdistro"); err != nil {
        log.Fatal(err)
    } else {
        if !hasParam {
            log.Fatal("HasParam() failed.")
        }
    }


    if foundKey, err := node.SearchParam("/rosdistro"); err != nil {
        log.Fatal(err)
    } else {
        if foundKey != "/rosdistro" {
            log.Fatal("SearchParam() failed.")
        }
    }

    if param, err := node.GetParam("/rosdistro"); err != nil {
        log.Fatal(err)
    } else {
        if value, ok := param.(string); !ok {
            log.Fatal("GetParam() failed.")
        } else {
            if value != "groovy\n" {
                log.Fatalf("Expected 'groovy\\n' but '%s'", value)
            }
        }
    }

    if err := node.SetParam("/test_param", 42); err != nil {
        log.Fatal(err)
    }

    if param, err := node.GetParam("/test_param"); err != nil {
        log.Fatal(err)
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
        log.Fatal(err)
    }
    
    log.Print("Success")
}


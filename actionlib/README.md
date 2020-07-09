# actionlib [WIP]

## Package Summary

A pure go implementation for ROS action library built on top of ROSGO. This package is unstable and the API can change in future. 

## Prerequisites

This library uses messages `GoalID`, `GoalStatus` and `GoalStatusArray` from `actionlib_msgs` package. Please generate Go code for the messages in `actionlib_msgs` package and place them in your `$GOPATH/src`.

Use the following commands after install `gengo`.

```cmd
gengo -out=$GOPATH/src msg actionlib_msgs/GoalID
gengo -out=$GOPATH/src msg actionlib_msgs/GoalStatus
gengo -out=$GOPATH/src msg actionlib_msgs/GoalStatusArray
```

## Status

This package implements all the features of actionlib library but is still very unstable and is still a work in progress to fix known issues and make this packge more robust. Following are the features that are implemented and what's to be added in the future.

### Implemented

- Action Client
- Action Server
- Simple Action Client
- Simple Action Server
- Client Goal Handler
- Server Goal Handler
- Go code generation from action definitons

### To Be Added

- Tests
- Documentation
- Fix for golint issues
- Go mod

## How To Use

Examples of client and server usage can be found in `rosgo/test` folder.

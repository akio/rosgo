## rosgo

[![GoDoc](https://godoc.org/github.com/fetchrobotics/rosgo?status.svg)](https://godoc.org/github.com/fetchrobotics/rosgo) 
[![Build Status](https://travis-ci.org/fetchrobotics/rosgo.svg?branch=master)](https://travis-ci.org/fetchrobotics/rosgo)

## Package Summary

**rosgo** is pure Go implementation of [ROS](http://www.ros.org/) client library.

- Author: Akio Ochiai
- Maintainer: Fetch Robotics
- License: Apache License 2.0
- Source: git [https://github.com/fetchrobotics/rosgo](https://github.com/fetchrobotics/rosgo)
- ROS Version Support: [Indigo] [Jade] [Melodic]

## Prerequisites

To use this library you should have installed ROS: [Install](wiki.ros.org/melodic/Installation/Ubuntu).
To run the tests please install all sensor msgs: `sudo apt install ros-melodic-desktop-full` for Ubuntu

## Status

**rosgo** is under development to implement all features of [ROS Client Library Requiements](http://www.ros.org/wiki/Implementing%20Client%20Libraries).

At present, following basic functions are provided.

- Parameter API (get/set/search....)
- ROS Slave API (with some exceptions)
- Publisher/Subscriber API (with TCPROS)
- Remapping
- Message Generation

Work to do:

- Action Servers
- Go Module Support
- Tutorials
- Bus Statistics
- ROS 2 Support

## How to use

Please look in the [test](test) folder for how to use rosgo in your projects.

## See also

- [rosgo in ROS Wiki](http://www.ros.org/wiki/rosgo)

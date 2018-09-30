package ros

import (
	"os"
	"testing"
)

func TestDetermineHost(t *testing.T) {
	os.Unsetenv("ROS_HOSTNAME")
	os.Unsetenv("ROS_IP")

	var host string
	var localOnly bool

	// ROS_HOSTNAME: localhost, ROS_IP: nil
	os.Setenv("ROS_HOSTNAME", "localhost")
	host, localOnly = determineHost()
	if host != "localhost" {
		t.Error("ROS_HOSTNAME is not addressed")
	}
	if localOnly != true {
		t.Errorf("localOnly flag is wrong for %s", host)
	}

	// ROS_HOSTNAME: hostname.in.env.var, ROS_IP: nil
	os.Setenv("ROS_HOSTNAME", "hostname.in.env.var")
	host, localOnly = determineHost()
	if host != "hostname.in.env.var" {
		t.Error("ROS_HOSTNAME is not addressed")
	}
	if localOnly != false {
		t.Errorf("localOnly flag is wrong for %s", host)
	}

	// ROS_HOSTNAME: hostname.in.env.var, ROS_IP: 1.2.3.4
	os.Setenv("ROS_IP", "1.2.3.4")
	host, localOnly = determineHost()
	if host != "hostname.in.env.var" {
		t.Error("ROS_HOSTNAME is not addressed when ROS_IP is set")
	}
	if localOnly != false {
		t.Errorf("localOnly flag is wrong for %s", host)
	}

	// ROS_HOSTNAME: nil, ROS_IP: 1.2.3.4
	os.Unsetenv("ROS_HOSTNAME")
	host, localOnly = determineHost()
	if host != "1.2.3.4" {
		t.Error("ROS_IP is not addressed")
	}
	if localOnly != false {
		t.Errorf("localOnly flag is wrong for %s", host)
	}

	// ROS_HOSTNAME: nil, ROS_IP: 127.0.0.1
	os.Setenv("ROS_IP", "127.0.0.1")
	host, localOnly = determineHost()
	if host != "127.0.0.1" {
		t.Error("ROS_HOSTNAME is not addressed when ROS_IP is set")
	}
	if localOnly != true {
		t.Errorf("localOnly flag is wrong for %s", host)
	}
}

package ros

import (
	"net"
	"os"
	"strings"
)

func determineHost() (string, bool) {
	// If the user set ROS_HOSTNAME, use it as is
	if rosHostname, ok := os.LookupEnv("ROS_HOSTNAME"); ok {
		return rosHostname, (rosHostname == "localhost")
	}

	// If the user set ROS_IP, use it as is
	if rosIP, ok := os.LookupEnv("ROS_IP"); ok {
		return rosIP, (rosIP == "::1" || strings.HasPrefix(rosIP, "127."))
	}

	// Try using the hostname
	if osHostname, err := os.Hostname(); err == nil && osHostname != "localhost" {
		return osHostname, false
	}

	// Fall back on the interface IP
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				return ipnet.IP.String(), false
			}
		}
	}
	// Fall back to the loopback UP
	return "127.0.0.1", true
}

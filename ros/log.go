package ros

import (
	"github.com/sirupsen/logrus"
)

// DefaultLogger represents a default logging instance
var logger *logrus.Logger

// DefaultLogger returns an instance of the default logger
func DefaultLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	return logrus.StandardLogger()
}

// NewLogger returns a new instance of a logger
func NewLogger() *logrus.Logger {
	return logrus.New()
}

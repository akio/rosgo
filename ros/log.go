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

<<<<<<< HEAD
// NewLogger returns a new instance of a logger
func NewLogger() *logrus.Logger {
	return logrus.New()
=======
func (logger *defaultLogger) Info(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelInfo) {
		msg := fmt.Sprintf("[INFO] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Infof(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelInfo) {
		log.Printf("[INFO] "+format, v...)
	}
}

func (logger *defaultLogger) Warn(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelWarn) {
		msg := fmt.Sprintf("[WARN] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Warnf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelWarn) {
		log.Printf("[WARN] "+format, v...)
	}
}

func (logger *defaultLogger) Error(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelError) {
		msg := fmt.Sprintf("[ERROR] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Errorf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelError) {
		log.Printf("[ERROR]"+format, v...)
	}
}

func (logger *defaultLogger) Fatal(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelFatal) {
		msg := fmt.Sprintf("[FATAL] %s", fmt.Sprint(v...))
		log.Println(msg)
		os.Exit(1)
	}
}

func (logger *defaultLogger) Fatalf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelFatal) {
		log.Printf("[FATAL] "+format, v...)
		os.Exit(1)
	}
>>>>>>> 24a6463ff109d57010e214746b042cd6742395da
}

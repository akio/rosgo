package ros

import (
    "log"
    "os"
)

type Logger interface {
    Debug(v ...interface{})
    Debugf(format string, v ...interface{})
    Info(v ...interface{})
    Infof(format string, v ...interface{})
    Warn(v ...interface{})
    Warnf(format string, v ...interface{})
    Error(v ...interface{})
    Errorf(format string, v ...interface{})
    Fatal(v ...interface{})
    Fatalf(format string, v ...interface{})
}

type defaultLogger struct {
    severity LogLevel
}

func NewDefaultLogger(level LogLevel) *defaultLogger {
    logger := new(defaultLogger)
    logger.severity = LogLevelInfo
    return logger
}

func (logger *defaultLogger) Debug(v ...interface{}) {
    if int(logger.severity) <= int(LogLevelDebug) {
        log.Println(v...)
    }
}

func (logger *defaultLogger) Debugf(format string, v ...interface{}) {
    if int(logger.severity) <= int(LogLevelDebug) {
        log.Printf(format, v...)
        log.Println()
    }
}

func (logger *defaultLogger) Info(v ...interface{}) {
    if int(logger.severity) <= int(LogLevelInfo) {
        log.Println(v...)
    }
}

func (logger *defaultLogger) Infof(format string, v ...interface{}) {
    if int(logger.severity) <= int(LogLevelInfo) {
        log.Printf(format, v...)
        log.Println()
    }
}

func (logger *defaultLogger) Warn(v ...interface{}) {
    if int(logger.severity) <= int(LogLevelWarn) {
        log.Println(v...)
    }
}

func (logger *defaultLogger) Warnf(format string, v ...interface{}) {
    if int(logger.severity) <= int(LogLevelWarn) {
        log.Printf(format, v...)
        log.Println()
    }
}

func (logger *defaultLogger) Error(v ...interface{}) {
    if int(logger.severity) <= int(LogLevelError) {
        log.Println(v...)
    }
}

func (logger *defaultLogger) Errorf(format string, v ...interface{}) {
    if int(logger.severity) <= int(LogLevelError) {
        log.Printf(format, v...)
        log.Println()
    }
}

func (logger *defaultLogger) Fatal(v ...interface{}) {
    if int(logger.severity) <= int(LogLevelFatal) {
        log.Fatalln(v...)
    }
}

func (logger *defaultLogger) Fatalf(format string, v ...interface{}) {
    if int(logger.severity) <= int(LogLevelFatal) {
        log.Printf(format, v...)
        log.Println()
        os.Exit(1)
    }
}

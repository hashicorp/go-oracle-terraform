package opc

import (
	"log"
	"os"
)

const (
	LogOff LogLevelType = iota * 0x1000
	LogDebug
)

// Needs to be fleshed out further
type LogLevelType uint

// Logger interface. Should be satisfied by Terraform's logger as well as the Default logger
type Logger interface {
	Log(...interface{})
}

type LoggerFunc func(...interface{})

func (f LoggerFunc) Log(args ...interface{}) {
	f(args...)
}

// Returns a default logger if one isn't specified during configuration
func NewDefaultLogger() Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Default logger to satisfy the logger interface
type defaultLogger struct {
	logger *log.Logger
}

func (l defaultLogger) Log(args ...interface{}) {
	l.logger.Println(args...)
}

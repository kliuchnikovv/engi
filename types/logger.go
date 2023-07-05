package types

import (
	"fmt"
	basicLog "log"
)

type Method string

const (
	Trace   Method = "TRACE"
	Info    Method = "INFO"
	Warning Method = "WARN"
	Error   Method = "ERROR"
)

type (
	Logger interface {
		Write(string, string, ...interface{})
		Trace(string, ...interface{})
		Info(string, ...interface{})
		Error(string, ...interface{})
		Warning(string, ...interface{})
	}

	Log struct {
		Logger

		channel chan error
	}
)

func NewLog(logger Logger) *Log {
	return &Log{
		Logger: logger,
	}
}

func (log *Log) Channel() chan error {
	return log.channel
}

// SendError - sends error to channel and writes it in log.
func (log *Log) SendError(format string, args ...interface{}) {
	if log.channel != nil {
		log.channel <- fmt.Errorf(format, args...)
	}

	log.Error(format, args...)
}

func (log *Log) Write(
	method Method,
	format string,
	args ...interface{},
) {
	if log.Logger == nil {
		basicLog.Printf("%s: %s", method, fmt.Sprintf(format, args...))
		return
	}

	var logFunction func(Logger, string, ...interface{})

	switch method {
	case Trace:
		logFunction = Logger.Trace
	case Info:
		logFunction = Logger.Info
	case Error:
		logFunction = Logger.Error
	case Warning:
		logFunction = Logger.Warning
	default:
		logFunction = func(logger Logger, format string, args ...interface{}) {
			logger.Write(string(method), format, args...)
		}
	}

	logFunction(log.Logger, format, args...)
}

func (log *Log) Trace(format string, args ...interface{}) {
	log.Write(Trace, format, args...)
}

func (log *Log) Info(format string, args ...interface{}) {
	log.Write(Info, format, args...)
}

func (log *Log) Warning(format string, args ...interface{}) {
	log.Write(Warning, format, args...)
}

func (log *Log) Error(format string, args ...interface{}) {
	log.Write(Error, format, args...)
}

func (log *Log) SetChannelCapacity(capacity int) {
	log.channel = make(chan error, capacity)
}

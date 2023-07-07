package logger

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

func New(logger Logger) *Log {
	return &Log{
		Logger: logger,
	}
}

func (log *Log) Channel() chan error {
	return log.channel
}

// SendErrorf - sends error to channel and writes it in log.
func (log *Log) SendErrorf(format string, args ...interface{}) {
	if log.channel != nil {
		log.channel <- fmt.Errorf(format, args...)
	}

	log.Errorf(format, args...)
}

func (log *Log) Writef(
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

func (log *Log) Tracef(format string, args ...interface{}) {
	log.Writef(Trace, format, args...)
}

func (log *Log) Infof(format string, args ...interface{}) {
	log.Writef(Info, format, args...)
}

func (log *Log) Warningf(format string, args ...interface{}) {
	log.Writef(Warning, format, args...)
}

func (log *Log) Errorf(format string, args ...interface{}) {
	log.Writef(Error, format, args...)
}

func (log *Log) SetChannelCapacity(capacity int) {
	log.channel = make(chan error, capacity)
}

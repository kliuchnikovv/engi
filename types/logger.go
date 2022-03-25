package types

import (
	"fmt"
	basicLog "log"
)

type (
	Logger interface {
		Infof(string, ...interface{})
		Errorf(string, ...interface{})
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

// SendErrorf - sends error to channel and writes it in log.
func (log *Log) SendErrorf(format string, args ...interface{}) {
	if log.channel != nil {
		log.channel <- fmt.Errorf(format, args...)
	}

	log.Errorf(format, args...)
}

func (log *Log) Infof(format string, args ...interface{}) {
	if log.Logger == nil {
		basicLog.Printf(format, args...)
	} else {
		log.Logger.Infof(format, args...)
	}
}

func (log *Log) Errorf(format string, args ...interface{}) {
	if log.Logger == nil {
		basicLog.Printf("ERROR: %s", fmt.Sprintf(format, args...))
	} else {
		log.Logger.Errorf(format, args...)
	}
}

func (log *Log) SetChannelCapacity(capacity int) {
	log.channel = make(chan error, capacity)
}

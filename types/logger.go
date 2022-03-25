package types

import (
	"fmt"
	"log"
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

func (e *Log) SendErrorf(format string, args ...interface{}) {
	if e.channel != nil {
		e.channel <- fmt.Errorf(format, args...)
	}

	e.Errorf(format, args...)
}

func (e *Log) Infof(format string, args ...interface{}) {
	if e.Logger == nil {
		log.Printf(format, args...)
	} else {
		e.Logger.Infof(format, args...)
	}
}

func (e *Log) Errorf(format string, args ...interface{}) {
	if e.Logger == nil {
		log.Printf("ERROR: %s", fmt.Sprintf(format, args...))
	} else {
		e.Logger.Errorf(format, args...)
	}
}

func (e *Log) SetChannelCapacity(capacity int) {
	e.channel = make(chan error, capacity)
}

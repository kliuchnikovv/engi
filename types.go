package webapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

type (
	ResultFunc      func(*Context)
	RouterFunc      func(string)
	HandlerFunc     func(*Context) error
	MarshalerFunc   func(interface{}) ([]byte, error)
	UnmarshalerFunc func([]byte, interface{}) error
	parameter       struct {
		raw          []string
		parsed       interface{}
		wasRequested bool
	}
)

type Responser interface {
	SetPayload(interface{})
	SetError(error)
}

type AsIsResponse struct {
	XMLName  xml.Name `xml:"response" json:"-"`
	Response string   `xml:",chardata"`
}

func (obj *AsIsResponse) SetPayload(object interface{}) {
	obj.Response = fmt.Sprint(object)
}

func (obj *AsIsResponse) SetError(err error) {
	obj.Response = err.Error()
}

func (obj *AsIsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.Response)
}

type ResponseObject struct {
	XMLName     xml.Name    `xml:"response" json:"-"`
	Result      interface{} `json:"result,omitempty" xml:"result,omitempty"`
	ErrorString string      `json:"error,omitempty" xml:"error,omitempty"`

	code int
}

func NewErrorResponse(code int, format string, args ...interface{}) ResponseObject {
	return ResponseObject{
		code:        code,
		ErrorString: fmt.Sprintf(format, args...),
	}
}

func (a *ResponseObject) SetPayload(object interface{}) {
	a.Result = object
}

func (a *ResponseObject) SetError(err error) {
	a.ErrorString = err.Error()
}

func (a ResponseObject) Error() string {
	return a.ErrorString
}

type Logger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

type Log struct {
	Logger

	channel chan error
}

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

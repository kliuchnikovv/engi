package webapi

import (
	"fmt"
)

// TODO: rename
type ResponseObject struct {
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

func (a ResponseObject) Error() string {
	return a.ErrorString
}

func (a *ResponseObject) SetPayload(object interface{}) {
	a.Result = object
}

func (a *ResponseObject) SetError(err error) {
	a.ErrorString = err.Error()
}

type Responser interface {
	SetPayload(interface{})
	SetError(error)
}

// TODO: add xml tags
type AsIsObject string

func (obj *AsIsObject) SetPayload(object interface{}) {
	*obj = AsIsObject(fmt.Sprint(object))
}

func (obj *AsIsObject) SetError(err error) {
	*obj = AsIsObject(err.Error())
}

package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type (
	Marshaler   func(interface{}) ([]byte, error)
	Unmarshaler func([]byte, interface{}) error
	Responser   interface {
		// SetPayload - sets response payload into object.
		SetPayload(interface{})
		// SetError - sets error response into object.
		SetError(error)
	}
)

// AsIsResponse - returns payload without any wrapping (even errors).
type AsIsResponse struct {
	XMLName  xml.Name    `xml:"response" json:"-"`
	Code     int         `xml:"-" json:"-"`
	Response interface{} `xml:",chardata" json:"response,omitempty"`
}

// SetPayload - sets response payload into object.
func (obj *AsIsResponse) SetPayload(object interface{}) {
	obj.Response = object
}

// SetError - sets error response into object.
func (obj *AsIsResponse) SetError(err error) {
	obj.Response = err.Error()
}

func (obj *AsIsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.Response)
}

type ResponseObject struct {
	XMLName     xml.Name    `json:"-" xml:"response"`
	Code        int         `json:"-" xml:"-"`
	Result      interface{} `json:"result,omitempty" xml:"result,omitempty"`
	ErrorString string      `json:"error,omitempty" xml:"error,omitempty"`
}

// SetPayload - sets response payload into object.
func (a *ResponseObject) SetPayload(object interface{}) {
	a.Result = object
}

// SetError - sets error response into object.
func (a *ResponseObject) SetError(err error) {
	a.ErrorString = err.Error()
}

func (a ResponseObject) Error() string {
	return a.ErrorString
}

func NewErrorResponse(code int, format string, args ...interface{}) ResponseObject {
	return ResponseObject{
		Code:        code,
		ErrorString: fmt.Sprintf(format, args...),
	}
}

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
		SetPayload(interface{})
		SetError(error)
	}
)

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
	XMLName     xml.Name    `json:"-" xml:"response"`
	Result      interface{} `json:"result,omitempty" xml:"result,omitempty"`
	ErrorString string      `json:"error,omitempty" xml:"error,omitempty"`

	Code int `json:"-" xml:"-"`
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

func NewErrorResponse(code int, format string, args ...interface{}) ResponseObject {
	return ResponseObject{
		Code:        code,
		ErrorString: fmt.Sprintf(format, args...),
	}
}

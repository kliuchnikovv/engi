package response

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// AsIs - returns payload without any wrapping (even errors).
type AsIs struct {
	XMLName  xml.Name    `json:"-"                  xml:"response"`
	Code     int         `json:"-"                  xml:"-"`
	Response interface{} `json:"response,omitempty" xml:",chardata"`
}

// SetPayload - sets response payload into object.
func (obj *AsIs) SetPayload(object interface{}) {
	obj.Response = object
}

// SetError - sets error response into object.
func (obj *AsIs) SetError(err error) {
	obj.Response = err.Error()
}

func (obj *AsIs) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.Response)
}

type AsObject struct {
	XMLName     xml.Name    `json:"-"                xml:"response"`
	Code        int         `json:"-"                xml:"-"`
	Result      interface{} `json:"result,omitempty" xml:"result,omitempty"`
	ErrorString string      `json:"error,omitempty"  xml:"error,omitempty"`
}

// SetPayload - sets response payload into object.
func (a *AsObject) SetPayload(object interface{}) {
	a.Result = object
}

// SetError - sets error response into object.
func (a *AsObject) SetError(err error) {
	a.ErrorString = err.Error()
}

func (a AsObject) Error() string {
	return a.ErrorString
}

func AsError(code int, format string, args ...interface{}) *AsObject {
	return &AsObject{
		Code:        code,
		ErrorString: fmt.Sprintf(format, args...),
	}
}

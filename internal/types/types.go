package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type (
	Unmarshaler func([]byte, interface{}) error
	Marshaler   struct {
		ContentType func() string
		Marshal     func(interface{}) ([]byte, error)
	}

	Responser interface {
		// SetPayload - sets response payload into object.
		SetPayload(payload interface{})
		// SetError - sets error response into object.
		SetError(err error)
	}

	Logger interface {
		Info()
	}
)

func NewJSONMarshaler() Marshaler {
	return Marshaler{
		ContentType: func() string { return "application/json" },
		Marshal:     json.Marshal,
	}
}

func NewXMLMarshaler() Marshaler {
	return Marshaler{
		ContentType: func() string { return "application/xml" },
		Marshal: func(i interface{}) ([]byte, error) {
			bytes, err := xml.Marshal(i)
			if err != nil {
				return nil, err
			}

			// Should append header for proper visualization.
			return append([]byte(xml.Header), bytes...), nil
		},
	}
}

// ResponseAsIs - returns payload without any wrapping (even errors).
type ResponseAsIs struct {
	XMLName  xml.Name    `json:"-"                  xml:"response"`
	Code     int         `json:"-"                  xml:"-"`
	Response interface{} `json:"response,omitempty" xml:",chardata"`
}

// SetPayload - sets response payload into object.
func (obj *ResponseAsIs) SetPayload(object interface{}) {
	obj.Response = object
}

// SetError - sets error response into object.
func (obj *ResponseAsIs) SetError(err error) {
	obj.Response = err.Error()
}

func (obj *ResponseAsIs) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.Response)
}

type ResponseAsObject struct {
	XMLName     xml.Name    `json:"-"                xml:"response"`
	Code        int         `json:"-"                xml:"-"`
	Result      interface{} `json:"result,omitempty" xml:"result,omitempty"`
	ErrorString string      `json:"error,omitempty"  xml:"error,omitempty"`
}

// SetPayload - sets response payload into object.
func (a *ResponseAsObject) SetPayload(object interface{}) {
	a.Result = object
}

// SetError - sets error response into object.
func (a *ResponseAsObject) SetError(err error) {
	a.ErrorString = err.Error()
}

func (a ResponseAsObject) Error() string {
	return a.ErrorString
}

func AsError(code int, format string, args ...interface{}) *ResponseAsObject {
	return &ResponseAsObject{
		Code:        code,
		ErrorString: fmt.Sprintf(format, args...),
	}
}

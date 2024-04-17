package types

import (
	"encoding/json"
	"encoding/xml"
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
)

func NewJSONMarshaler() *Marshaler {
	return &Marshaler{
		ContentType: func() string { return "application/json" },
		Marshal:     json.Marshal,
	}
}

func NewXMLMarshaler() *Marshaler {
	return &Marshaler{
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

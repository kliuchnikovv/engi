package webapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)

// Context - provides methods for extracting data from query and response back.
type Context struct {
	*Request
	*Response
}

// NewContext - returns new Context instance from 'echo.Context'.
func NewContext(
	response http.ResponseWriter,
	request *http.Request,
	responseMarshaler MarshalerFunc,
	responseObject Responser,
) *Context {
	return &Context{
		Request:  NewRequest(request),
		Response: NewResponse(response, responseMarshaler, responseObject),
	}
}

func (ctx *Context) readBody() error {
	defer ctx.request.Body.Close()

	bytes, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil && !errors.Is(err, http.ErrBodyReadAfterClose) {
		return NewErrorResponse(http.StatusInternalServerError, "reading body failed: %s", err.Error())
	}

	if len(bytes) != 0 {
		ctx.body.raw = []string{string(bytes)}
	}

	if len(ctx.body.raw) == 0 {
		return NewErrorResponse(http.StatusBadRequest, "no required body provided")
	}

	return err
}

func (ctx *Context) getUnmarshaler() (UnmarshalerFunc, error) {
	var (
		contentType = ctx.request.Header.Get("Content-type")
		unmarshal   UnmarshalerFunc
	)

	switch contentType {
	case "application/json":
		unmarshal = json.Unmarshal
	case "application/xml":
		unmarshal = xml.Unmarshal
	case "text/plain":
		unmarshal = func(b []byte, i interface{}) error {
			typed, ok := i.(*string)
			if !ok {
				return NewErrorResponse(http.StatusInternalServerError, "pointer must be of type '*string'")
			}

			*typed = string(b)

			return nil
		}
	default:
		return nil, NewErrorResponse(http.StatusBadRequest, "content-type not supported: %s", contentType)
	}

	return func(bytes []byte, pointer interface{}) error {
		if err := unmarshal(bytes, pointer); err != nil {
			return NewErrorResponse(http.StatusInternalServerError, "unmarshaling body failed: %s", err.Error())
		}

		ctx.body.wasRequested = true
		ctx.body.parsed = pointer

		return nil
	}, nil
}

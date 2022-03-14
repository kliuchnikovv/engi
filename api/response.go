package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	*Context
}

func NewResponse(ctx *Context) Response {
	return Response{
		Context: ctx,
	}
}

func (response *Response) WithJSON(code int, payload interface{}) error {
	var (
		result json.RawMessage
		err    error
	)

	switch typed := payload.(type) {
	case []byte:
		result = typed
	default:
		result, err = json.Marshal(payload)
		if err != nil {
			return response.InternalServerError(err.Error())
		}
	}

	return response.Context.Context.JSON(code, map[string]interface{}{
		"result": result,
	})
}

func (response *Response) WithoutContent(code int) error {
	return response.Context.Context.NoContent(http.StatusNoContent)
}

/// OK's

func (response *Response) OK(payload interface{}) error {
	return response.WithJSON(http.StatusOK, payload)
}

func (response *Response) Created() error {
	return response.WithoutContent(http.StatusCreated)
}

func (response *Response) NoContent() error {
	return response.WithoutContent(http.StatusNoContent)
}

/// Errors

func (response *Response) Error(code int, err error) error {
	return response.Context.Context.JSON(code, map[string]string{
		"error": err.Error(),
	})
}

func (response *Response) BadRequest(format string, args ...interface{}) error {
	return response.Error(http.StatusBadRequest, fmt.Errorf(format, args...))
}

func (response *Response) Forbidden(format string, args ...interface{}) error {
	return response.Error(http.StatusForbidden, fmt.Errorf(format, args...))
}

func (response *Response) NotFound(format string, args ...interface{}) error {
	return response.Error(http.StatusNotFound, fmt.Errorf(format, args...))
}

func (response *Response) MethodNotAllowed(format string, args ...interface{}) error {
	return response.Error(http.StatusMethodNotAllowed, fmt.Errorf(format, args...))
}

func (response *Response) InternalServerError(format string, args ...interface{}) error {
	return response.Error(http.StatusInternalServerError, fmt.Errorf(format, args...))
}

// func (api *Response) responseResult(code int, payload interface{}) error {
// 	// var (
// 	// 	result json.RawMessage
// 	// 	err    error
// 	// )
// 	// switch typed := payload.(type) {
// 	// case []byte:
// 	// 	result = typed
// 	// default:
// 	// 	result, err = json.Marshal(payload)
// 	// 	if err != nil {
// 	// 		return InternalError(ctx, err.Error())
// 	// 	}
// 	// }
// 	// return ctx.JSON(code, map[string]interface{}{
// 	// 	"result": result,
// 	// })
// }

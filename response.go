package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response - provide methods for creating responses.
type Response struct {
	*Context
}

func NewResponse(ctx *Context) Response {
	return Response{
		Context: ctx,
	}
}

// WithJSON - responses with provided custom code and body.
func (response *Response) WithJSON(code int, payload interface{}) error {
	var (
		result json.RawMessage
		err    error
	)

	switch typed := payload.(type) {
	case string:
		result = []byte(typed)
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

// WithourContent - responses with provided custom code and no body.
func (response *Response) WithoutContent(code int) error {
	return response.Context.Context.NoContent(http.StatusNoContent)
}

// Error - responses custom error with provided code and message.
func (response *Response) Error(code int, err error) error {
	return response.Context.Context.JSON(code, map[string]string{
		"error": err.Error(),
	})
}

// OK - writes payload into json's 'result' field with 200 http code.
func (response *Response) OK(payload interface{}) error {
	return response.WithJSON(http.StatusOK, payload)
}

// Created - responses with 201 http code and no content.
func (response *Response) Created() error {
	return response.WithoutContent(http.StatusCreated)
}

// NoContent - responses with 204 http code and no content.
func (response *Response) NoContent() error {
	return response.WithoutContent(http.StatusNoContent)
}

// BadRequest - responses with 400 code and provided message.
func (response *Response) BadRequest(format string, args ...interface{}) error {
	return response.Error(http.StatusBadRequest, fmt.Errorf(format, args...))
}

// Forbidden - responses with 403 error code and provided message.
func (response *Response) Forbidden(format string, args ...interface{}) error {
	return response.Error(http.StatusForbidden, fmt.Errorf(format, args...))
}

// NotFound - responses with 404 error code and provided message.
func (response *Response) NotFound(format string, args ...interface{}) error {
	return response.Error(http.StatusNotFound, fmt.Errorf(format, args...))
}

// MethodNotAllowed - responses with 405 error code and provided message.
func (response *Response) MethodNotAllowed(format string, args ...interface{}) error {
	return response.Error(http.StatusMethodNotAllowed, fmt.Errorf(format, args...))
}

// InternalServerError - responses with 500 error code and provided message.
func (response *Response) InternalServerError(format string, args ...interface{}) error {
	return response.Error(http.StatusInternalServerError, fmt.Errorf(format, args...))
}

package param

import (
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/webapi/types"
)

// Response - provide methods for creating responses.
type Response struct {
	writer    http.ResponseWriter
	marshaler types.Marshaler
	object    types.Responser
}

func NewResponse(
	writer http.ResponseWriter,
	marshaler types.Marshaler,
	object types.Responser,
) *Response {
	return &Response{
		writer:    writer,
		marshaler: marshaler,
		object:    object,
	}
}

// WithJSON - responses with provided custom code and body.
func (resp *Response) WithJSON(code int, payload interface{}) error {
	resp.object.SetPayload(payload)

	bytes, err := resp.marshaler(resp.object)
	if err != nil {
		return err
	}

	if _, err := resp.writer.Write(bytes); err != nil {
		return err
	}

	return nil
}

// Error - responses custom error with provided code and message.
func (resp *Response) Error(code int, format string, args ...interface{}) error {
	resp.object.SetError(fmt.Errorf(format, args...))

	bytes, err := resp.marshaler(resp.object)
	if err != nil {
		return err
	}

	resp.writer.WriteHeader(code)
	_, err = resp.writer.Write(bytes)

	return err
}

// WithourContent - responses with provided custom code and no body.
func (resp *Response) WithoutContent(code int) {
	resp.writer.WriteHeader(code)
}

// OK - writes payload into json's 'result' field with 200 http code.
func (resp *Response) OK(payload interface{}) error {
	return resp.WithJSON(http.StatusOK, payload)
}

// Created - responses with 201 http code and no content.
func (resp *Response) Created() error {
	resp.WithoutContent(http.StatusCreated)
	return nil
}

// NoContent - responses with 204 http code and no content.
func (resp *Response) NoContent() error {
	resp.WithoutContent(http.StatusNoContent)
	return nil
}

// BadRequest - responses with 400 code and provided message.
func (resp *Response) BadRequest(format string, args ...interface{}) error {
	return resp.Error(http.StatusBadRequest, format, args...)
}

// Forbidden - responses with 403 error code and provided message.
func (resp *Response) Forbidden(format string, args ...interface{}) error {
	return resp.Error(http.StatusForbidden, format, args...)
}

// NotFound - responses with 404 error code and provided message.
func (resp *Response) NotFound(format string, args ...interface{}) error {
	return resp.Error(http.StatusNotFound, format, args...)
}

// MethodNotAllowed - responses with 405 error code and provided message.
func (resp *Response) MethodNotAllowed(format string, args ...interface{}) error {
	return resp.Error(http.StatusMethodNotAllowed, format, args...)
}

// InternalServerError - responses with 500 error code and provided message.
func (resp *Response) InternalServerError(format string, args ...interface{}) error {
	return resp.Error(http.StatusInternalServerError, format, args...)
}

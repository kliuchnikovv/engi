package response

import (
	"fmt"
	"net/http"

	"github.com/kliuchnikovv/engi/internal/types"
)

// TODO: add gRPC and RPC support

type Responser interface {
	// ResponseWriter - returns http.ResponseWriter associated with request.
	ResponseWriter() http.ResponseWriter
	// Object - responses with provided custom code and body.
	// Body will be marshaled using service-defined object and marshaler.
	Object(code int, payload interface{}) error
	// WithourContent - responses with provided custom code and no body.
	WithoutContent(code int) error
	// Error - responses custom error with provided code and error.
	Error(code int, err error) error
	// Errorf - responses custom error with provided code and formatted string message.
	Errorf(code int, format string, args ...interface{}) error
	// OK - writes payload into json's 'result' field with 200 http code.
	OK(payload interface{}) error
	// Created - responses with 201 http code and no content.
	Created() error
	// NoContent - responses with 204 http code and no content.
	NoContent() error
	// BadRequest - responses with 400 code and provided formatted string message.
	BadRequest(format string, args ...interface{}) error
	// Forbidden - responses with 403 error code and provided formatted string message.
	Forbidden(format string, args ...interface{}) error
	// NotFound - responses with 404 error code and provided formatted string message.
	NotFound(format string, args ...interface{}) error
	// MethodNotAllowed - responses with 405 error code and provided formatted string message.
	MethodNotAllowed(format string, args ...interface{}) error
	// InternalServerError - responses with 500 error code and provided formatted string message.
	InternalServerError(format string, args ...interface{}) error
}

// Response - provide methods for creating responses.
type Response struct {
	writer    http.ResponseWriter
	marshaler types.Marshaler
	object    types.Responser
}

func New(
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

func (resp *Response) Object(code int, payload interface{}) error {
	resp.object.SetPayload(payload)

	bytes, err := resp.marshaler.Marshal(resp.object)
	if err != nil {
		return err
	}

	var contentType = resp.marshaler.ContentType()
	if contentType != "" {
		resp.writer.Header().Add("Content-Type", contentType)
	}

	resp.writer.WriteHeader(code)

	if _, err := resp.writer.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (resp *Response) Error(code int, err error) error {
	resp.object.SetError(err)

	bytes, err := resp.marshaler.Marshal(resp.object)
	if err != nil {
		return err
	}

	var contentType = resp.marshaler.ContentType()
	if contentType != "" {
		resp.writer.Header().Add("Content-Type", contentType)
	}

	resp.writer.WriteHeader(code)
	_, err = resp.writer.Write(bytes)

	return err
}

func (resp *Response) Errorf(code int, format string, args ...any) error {
	return resp.Error(code,
		fmt.Errorf(format, args...),
	)
}

func (resp *Response) WithoutContent(code int) error {
	resp.writer.WriteHeader(code)
	return nil // in purpose of unification
}

func (resp *Response) OK(payload interface{}) error {
	return resp.Object(http.StatusOK, payload)
}

func (resp *Response) Created() error {
	return resp.WithoutContent(http.StatusCreated)
}

func (resp *Response) NoContent() error {
	return resp.WithoutContent(http.StatusNoContent)
}

func (resp *Response) BadRequest(format string, args ...interface{}) error {
	return resp.Errorf(http.StatusBadRequest, format, args...)
}

func (resp *Response) Forbidden(format string, args ...interface{}) error {
	return resp.Errorf(http.StatusForbidden, format, args...)
}

func (resp *Response) NotFound(format string, args ...interface{}) error {
	return resp.Errorf(http.StatusNotFound, format, args...)
}

func (resp *Response) MethodNotAllowed(format string, args ...interface{}) error {
	return resp.Errorf(http.StatusMethodNotAllowed, format, args...)
}

func (resp *Response) InternalServerError(format string, args ...interface{}) error {
	return resp.Errorf(http.StatusInternalServerError, format, args...)
}

func (resp *Response) ResponseWriter() http.ResponseWriter {
	return resp.writer
}


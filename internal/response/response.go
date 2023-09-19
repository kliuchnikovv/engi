package response

import (
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/webapi/internal/types"
)

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

	if _, err := resp.writer.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (resp *Response) Error(code int, format string, args ...interface{}) error {
	resp.object.SetError(fmt.Errorf(format, args...))

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
	return resp.Error(http.StatusBadRequest, format, args...)
}

func (resp *Response) Forbidden(format string, args ...interface{}) error {
	return resp.Error(http.StatusForbidden, format, args...)
}

func (resp *Response) NotFound(format string, args ...interface{}) error {
	return resp.Error(http.StatusNotFound, format, args...)
}

func (resp *Response) MethodNotAllowed(format string, args ...interface{}) error {
	return resp.Error(http.StatusMethodNotAllowed, format, args...)
}

func (resp *Response) InternalServerError(format string, args ...interface{}) error {
	return resp.Error(http.StatusInternalServerError, format, args...)
}

func (resp *Response) ResponseWriter() http.ResponseWriter {
	return resp.writer
}

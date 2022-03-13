package types

import "time"

// Describes echo or gin
// type Contexter interface {
// 	QueryParameter(string) string
// 	Bind(interface{}) error
// 	JSON(int, interface{}) error
// 	NoContent(int) error
// }

type Query interface {
	String(key string, isObligatory bool) (string, error)
	Bool(key string, isObligatory bool) (bool, error)
	Integer(key string, isObligatory bool) (int, error)
	Time(key string, isObligatory bool, layout string) (*time.Time, error)
	Body(pointer interface{}) error
}

type Response interface {
	WithJSON(code int, payload interface{}) error
	WithoutContent(code int) error
	Error(code int, err error) error

	OK(payload interface{}) error
	Created() error
	NoContent() error
	BadRequest(format string, args ...interface{}) error
	Forbidden(format string, args ...interface{}) error
	NotFound(format string, args ...interface{}) error
	MethodNotAllowed(format string, args ...interface{}) error
	InternalServerError(format string, args ...interface{}) error
}

// TODO: naming
type APIContext interface {
	QueryParameter(string) string
	Bind(interface{}) error
	JSON(int, interface{}) error
	Status(int) error
}

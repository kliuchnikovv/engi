package errors

import (
	"fmt"
	"strings"
)

type Field struct {
	key   string
	value interface{}
}

func (f Field) String() string {
	return fmt.Sprintf("%s: %v", f.key, f.value)
}

type Error struct {
	msg error

	fields []Field
}

func New(msg error, fields ...Field) Error {
	return Error{
		msg:    msg,
		fields: fields,
	}
}

func Newf() {

}

func (e Error) Error() string {
	var builder = strings.Builder{}

	if _, err := builder.WriteString(e.msg.Error()); err != nil {
		panic(err)
	}

	for _, field := range e.fields {
		if _, err := builder.WriteString(" " + field.String()); err != nil {
			panic(err)
		}
	}

	return builder.String()
}

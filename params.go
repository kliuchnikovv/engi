package webapi

import (
	"net/http"
	"strconv"
	"time"
)

// WithBool - queries mandatory boolean parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Bool(key)'.
func WithBool(key string, configs ...ParameterConfig) HandlerFunc {
	return func(ctx *Context) error {
		return ctx.extractParam(key, configs, func(param string) (interface{}, error) {
			return strconv.ParseBool(param)
		})
	}
}

// WithInteger - queries mandatory integer parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Integer(key)'.
func WithInteger(key string, configs ...ParameterConfig) HandlerFunc {
	return func(ctx *Context) error {
		return ctx.extractParam(key, configs, func(param string) (interface{}, error) {
			var (
				intBase = 10
				bitSize = 64
			)

			result, err := strconv.ParseInt(param, intBase, bitSize)
			if err != nil {
				return nil, NewErrorResponse(http.StatusBadRequest, "parameter '%s' not of type int", key)
			}

			return result, err
		})
	}
}

// WithFloat - queries mandatory floating point number parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Float(key)'.
func WithFloat(key string, configs ...ParameterConfig) HandlerFunc {
	return func(ctx *Context) error {
		return ctx.extractParam(key, configs, func(param string) (interface{}, error) {
			var bitSize = 64

			result, err := strconv.ParseFloat(param, bitSize)
			if err != nil {
				return nil, NewErrorResponse(http.StatusBadRequest, "parameter '%s' not of type float", key)
			}

			return result, err
		})
	}
}

// WithString - queries mandatory string parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.String(key)'.
func WithString(key string, configs ...ParameterConfig) HandlerFunc {
	return func(ctx *Context) error {
		return ctx.extractParam(key, configs, func(param string) (interface{}, error) {
			return param, nil
		})
	}
}

// WithTime - queries mandatory time parameter from request by 'key' using 'layout'.
// Result can be retrieved from context using 'context.QueryParams.Time(key, layout)'.
func WithTime(key, layout string, configs ...ParameterConfig) HandlerFunc {
	return func(ctx *Context) error {
		return ctx.extractParam(key, configs, func(param string) (interface{}, error) {
			result, err := time.Parse(layout, param)
			if err != nil {
				return nil, NewErrorResponse(http.StatusBadRequest,
					"could not parse '%s' param to datetime using '%s' layout", key, layout,
				)
			}

			return result, err
		})
	}
}

// WithBody - takes pointer to structure and saves casted request body into context.
// Result can be retrieved from context using 'context.QueryParams.Body()'.
func WithBody(pointer interface{}) HandlerFunc {
	return func(ctx *Context) error {
		unmarshal, err := ctx.getUnmarshaler()
		if err != nil {
			return err
		}

		return ctx.extractBody(unmarshal, pointer)
	}
}

// WithCustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
// Result can be retrieved from context using 'context.QueryParams.Body()'.
func WithCustomBody(unmarshal UnmarshalerFunc, pointer interface{}) HandlerFunc {
	return func(ctx *Context) error {
		// TODO: context.QueryParams.CustomBody()
		return ctx.extractBody(unmarshal, pointer)
	}
}

func Description(s string) ParameterConfig {
	return func(p *parameter) error {
		p.description = s

		return nil
	}
}

func NotEmpty(p *parameter) error {
	var isNotEmpty func() bool

	switch typed := p.parsed.(type) {
	case bool:
		// Bool can't be empty
		return nil
	case int64:
		isNotEmpty = func() bool { return typed != 0 }
	case float64:
		isNotEmpty = func() bool { return typed != 0 }
	case string:
		isNotEmpty = func() bool { return len(typed) != 0 }
	case time.Time:
		isNotEmpty = func() bool { return typed.UnixNano() != 0 }
	}

	if isNotEmpty() {
		return nil
	}

	return NewErrorResponse(http.StatusBadRequest,
		"'%s' shouldn't be empty", p.name,
	)
}

func Greater(than float64) ParameterConfig {
	return func(p *parameter) error {
		var greater func() bool

		switch typed := p.parsed.(type) {
		case int64:
			greater = func() bool { return typed > int64(than) }
		case float64:
			greater = func() bool { return typed > than }
		case string:
			greater = func() bool { return len(typed) > int(than) }
		case time.Time:
			greater = func() bool { return typed.UnixNano() > int64(than) }
		}

		if greater() {
			return nil
		}

		return NewErrorResponse(http.StatusBadRequest,
			"'%s' should be greater than %f", p.name, than,
		)
	}
}

func Less(than float64) ParameterConfig {
	return func(p *parameter) error {
		var greater func() bool

		switch typed := p.parsed.(type) {
		case int64:
			greater = func() bool { return typed < int64(than) }
		case float64:
			greater = func() bool { return typed < than }
		case string:
			greater = func() bool { return len(typed) < int(than) }
		case time.Time:
			greater = func() bool { return typed.UnixNano() < int64(than) }
		}

		if greater() {
			return nil
		}

		return NewErrorResponse(http.StatusBadRequest,
			"'%s' should be less than %f", p.name, than,
		)
	}
}

func OR(first, second ParameterConfig) ParameterConfig {
	return func(p *parameter) error {
		var (
			err1 = first(p)
			err2 = second(p)
		)

		if err1 == nil || err2 == nil {
			return nil
		}

		return NewErrorResponse(http.StatusBadRequest,
			"'%s' failed checks: %s and %s", p.name, err1, err2,
		)
	}
}

func AND(first, second ParameterConfig) ParameterConfig {
	return func(p *parameter) error {
		var (
			err1 = first(p)
			err2 = second(p)
		)

		if err1 == nil && err2 == nil {
			return nil
		}

		if err1 != nil && err2 != nil {
			return NewErrorResponse(http.StatusBadRequest,
				"'%s' failed checks: %s and %s", p.name, err1, err2,
			)
		}

		if err1 != nil {
			return NewErrorResponse(http.StatusBadRequest,
				"'%s' failed check: %s", p.name, err1,
			)
		}

		return NewErrorResponse(http.StatusBadRequest,
			"'%s' failed check: %s", p.name, err2,
		)
	}
}

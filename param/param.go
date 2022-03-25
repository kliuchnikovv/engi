package param

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/KlyuchnikovV/webapi/types"
)

type (
	ParametersOption func(*Parameter) error
	HandlersOption   func(*Request) error

	Parameter struct {
		raw          []string
		parsed       interface{}
		wasRequested bool

		name        string
		description string
	}
)

// WithBool - queries mandatory boolean Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Bool(key)'.
func WithBool(key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(request, key, options, func(param string) (interface{}, error) {
			return strconv.ParseBool(param)
		})
	}
}

// WithInteger - queries mandatory integer Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Integer(key)'.
func WithInteger(key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(request, key, options, func(param string) (interface{}, error) {
			var (
				intBase = 10
				bitSize = 64
			)

			result, err := strconv.ParseInt(param, intBase, bitSize)
			if err != nil {
				return nil, types.NewErrorResponse(http.StatusBadRequest, "Parameter '%s' not of type int", key)
			}

			return result, err
		})
	}
}

// WithFloat - queries mandatory floating point number Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Float(key)'.
func WithFloat(key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(request, key, options, func(param string) (interface{}, error) {
			var bitSize = 64

			result, err := strconv.ParseFloat(param, bitSize)
			if err != nil {
				return nil, types.NewErrorResponse(http.StatusBadRequest, "Parameter '%s' not of type float", key)
			}

			return result, err
		})
	}
}

// WithString - queries mandatory string Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.String(key)'.
func WithString(key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(request, key, options, func(param string) (interface{}, error) {
			return param, nil
		})
	}
}

// WithTime - queries mandatory time Parameter from request by 'key' using 'layout'.
// Result can be retrieved from context using 'context.QueryParams.Time(key, layout)'.
func WithTime(key, layout string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(request, key, options, func(param string) (interface{}, error) {
			result, err := time.Parse(layout, param)
			if err != nil {
				return nil, types.NewErrorResponse(http.StatusBadRequest,
					"could not parse '%s' param to datetime using '%s' layout", key, layout,
				)
			}

			return result, err
		})
	}
}

// WithBody - takes pointer to structure and saves casted request body into context.
// Result can be retrieved from context using 'context.QueryParams.Body()'.
func WithBody(pointer interface{}) HandlersOption {
	return func(request *Request) error {
		unmarshal, err := getUnmarshaler(request)
		if err != nil {
			return err
		}

		return extractBody(request, unmarshal, pointer)
	}
}

// WithCustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
// Result can be retrieved from context using 'context.QueryParams.Body()'.
func WithCustomBody(unmarshal types.Unmarshaler, pointer interface{}) HandlersOption {
	return func(request *Request) error {
		return extractBody(request, unmarshal, pointer)
	}
}

func Description(s string) ParametersOption {
	return func(p *Parameter) error {
		p.description = s

		return nil
	}
}

func NotEmpty(p *Parameter) error {
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

	return types.NewErrorResponse(http.StatusBadRequest,
		"'%s' shouldn't be empty", p.name,
	)
}

func Greater(than float64) ParametersOption {
	return func(p *Parameter) error {
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

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' should be greater than %f", p.name, than,
		)
	}
}

func Less(than float64) ParametersOption {
	return func(p *Parameter) error {
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

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' should be less than %f", p.name, than,
		)
	}
}

func OR(first, second ParametersOption) ParametersOption {
	return func(p *Parameter) error {
		var (
			err1 = first(p)
			err2 = second(p)
		)

		if err1 == nil || err2 == nil {
			return nil
		}

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' failed checks: %s and %s", p.name, err1, err2,
		)
	}
}

func AND(first, second ParametersOption) ParametersOption {
	return func(p *Parameter) error {
		var (
			err1 = first(p)
			err2 = second(p)
		)

		if err1 == nil && err2 == nil {
			return nil
		}

		if err1 != nil && err2 != nil {
			return types.NewErrorResponse(http.StatusBadRequest,
				"'%s' failed checks: %s and %s", p.name, err1, err2,
			)
		}

		if err1 != nil {
			return types.NewErrorResponse(http.StatusBadRequest,
				"'%s' failed check: %s", p.name, err1,
			)
		}

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' failed check: %s", p.name, err2,
		)
	}
}

// extractParam - extracting parameter from context, calls middleware and saves to 'context.queryParameters[key]'.
// After this parameter can be retrieved from context using 'context.Query' methods.
func extractParam(
	request *Request,
	key string,
	configs []ParametersOption,
	convert func(string) (interface{}, error),
) error {
	var param = request.getParameter(key)
	if len(param) == 0 {
		return types.NewErrorResponse(http.StatusBadRequest, "parameter '%s' not found", key)
	}

	result, err := convert(param)
	if err != nil {
		return err
	}

	if result != nil {
		var parameter = request.parameters[key]

		request.parameters[key] = Parameter{
			name:         key,
			parsed:       result,
			raw:          parameter.raw,
			description:  parameter.description,
			wasRequested: true,
		}
	}

	var parameter = request.parameters[key]
	for _, config := range configs {
		if err := config(&parameter); err != nil {
			return err
		}
	}

	parameter.name = key
	request.parameters[key] = parameter

	return nil
}

func extractBody(request *Request, unmarshaler types.Unmarshaler, pointer interface{}) error {
	if request.body.parsed != nil {
		// If already parsed - skip parsing
		return nil
	}

	if err := readBody(request); err != nil {
		return err
	}

	if len(request.body.raw) == 0 {
		return types.NewErrorResponse(http.StatusInternalServerError, "no body found after reading")
	}

	return unmarshaler([]byte(request.body.raw[0]), pointer)
}

func getUnmarshaler(request *Request) (types.Unmarshaler, error) {
	var (
		contentType = request.request.Header.Get("Content-type")
		unmarshal   types.Unmarshaler
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
				return types.NewErrorResponse(http.StatusInternalServerError, "pointer must be of type '*string'")
			}

			*typed = string(b)

			return nil
		}
	default:
		return nil, types.NewErrorResponse(http.StatusBadRequest, "content-type not supported: %s", contentType)
	}

	return func(bytes []byte, pointer interface{}) error {
		if err := unmarshal(bytes, pointer); err != nil {
			return types.NewErrorResponse(http.StatusInternalServerError, "unmarshaling body failed: %s", err.Error())
		}

		request.body.wasRequested = true
		request.body.parsed = pointer

		return nil
	}, nil
}

func readBody(request *Request) error {
	defer request.request.Body.Close()

	bytes, err := ioutil.ReadAll(request.request.Body)
	if err != nil && !errors.Is(err, http.ErrBodyReadAfterClose) {
		return types.NewErrorResponse(http.StatusInternalServerError, "reading body failed: %s", err.Error())
	}

	if len(bytes) != 0 {
		request.body.raw = []string{string(bytes)}
	}

	if len(request.body.raw) == 0 {
		return types.NewErrorResponse(http.StatusBadRequest, "no required body provided")
	}

	return err
}

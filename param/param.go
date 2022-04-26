package param

import (
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

// Bool - queries mandatory boolean Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Bool(key)'.
func InPathBool(key string, options ...ParametersOption) HandlersOption {
	return boolParam(inPath, key, options...)
}

// QueryBool - queries mandatory boolean Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.QueryBool(key)'.
func QueryBool(key string, options ...ParametersOption) HandlersOption {
	return boolParam(inPath, key, options...)
}

// Integer - queries mandatory integer Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Integer(key)'.
func InPathInteger(key string, options ...ParametersOption) HandlersOption {
	return integerParam(inPath, key, options...)
}

// QueryInteger - queries mandatory integer Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.QueryInteger(key)'.
func QueryInteger(key string, options ...ParametersOption) HandlersOption {
	return integerParam(inPath, key, options...)
}

// Float - queries mandatory floating point number Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.Float(key)'.
func InPathFloat(key string, options ...ParametersOption) HandlersOption {
	return floatParam(inPath, key, options...)
}

// QueryFloat - queries mandatory floating point number Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.QueryFloat(key)'.
func QueryFloat(key string, options ...ParametersOption) HandlersOption {
	return floatParam(inPath, key, options...)
}

// String - queries mandatory string Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.String(key)'.
func InPathString(key string, options ...ParametersOption) HandlersOption {
	return stringParam(inPath, key, options...)
}

// QueryString - queries mandatory string Parameter from request by 'key'.
// Result can be retrieved from context using 'context.QueryParams.QueryString(key)'.
func QueryString(key string, options ...ParametersOption) HandlersOption {
	return stringParam(inPath, key, options...)
}

// Time - queries mandatory time Parameter from request by 'key' using 'layout'.
// Result can be retrieved from context using 'context.QueryParams.Time(key, layout)'.
func InPathTime(key, layout string, options ...ParametersOption) HandlersOption {
	return timeParam(inPath, key, layout, options...)
}

// QueryTime - queries mandatory time Parameter from request by 'key' using 'layout'.
// Result can be retrieved from context using 'context.QueryParams.QueryTime(key, layout)'.
func QueryTime(key, layout string, options ...ParametersOption) HandlersOption {
	return timeParam(inPath, key, layout, options...)
}

// Body - takes pointer to structure and saves casted request body into context.
// Result can be retrieved from context using 'context.QueryParams.Body()'.
func Body(pointer interface{}) HandlersOption {
	return func(request *Request) error {
		unmarshal, err := getUnmarshaler(request)
		if err != nil {
			return err
		}

		return extractBody(request, unmarshal, pointer)
	}
}

// CustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
// Result can be retrieved from context using 'context.QueryParams.Body()'.
func CustomBody(unmarshal types.Unmarshaler, pointer interface{}) HandlersOption {
	return func(request *Request) error {
		return extractBody(request, unmarshal, pointer)
	}
}

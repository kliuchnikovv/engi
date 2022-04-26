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

type paramPlacing string

const (
	inPath paramPlacing = "in-path"
	query  paramPlacing = "query"
)

func boolParam(place paramPlacing, key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(place, request, key, options, func(param string) (interface{}, error) {
			return strconv.ParseBool(param)
		})
	}
}

func integerParam(place paramPlacing, key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(place, request, key, options, func(param string) (interface{}, error) {
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

func floatParam(place paramPlacing, key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(place, request, key, options, func(param string) (interface{}, error) {
			var bitSize = 64

			result, err := strconv.ParseFloat(param, bitSize)
			if err != nil {
				return nil, types.NewErrorResponse(http.StatusBadRequest, "Parameter '%s' not of type float", key)
			}

			return result, err
		})
	}
}

func stringParam(place paramPlacing, key string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(place, request, key, options, func(param string) (interface{}, error) {
			return param, nil
		})
	}
}

func timeParam(place paramPlacing, key, layout string, options ...ParametersOption) HandlersOption {
	return func(request *Request) error {
		return extractParam(place, request, key, options, func(param string) (interface{}, error) {
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

// extractParam - extracting parameter from context, calls middleware and saves to 'context.parameters[from][key]'.
// After this parameter can be retrieved from context using 'context.Query' methods.
func extractParam(
	from paramPlacing,
	request *Request,
	key string,
	configs []ParametersOption,
	convert func(string) (interface{}, error),
) error {
	var param = request.getParameter(from, key)
	if len(param) == 0 {
		return types.NewErrorResponse(http.StatusBadRequest, "parameter '%s' not found", key)
	}

	result, err := convert(param)
	if err != nil {
		return err
	}

	if result != nil {
		var parameter = request.parameters[from][key]

		request.parameters[from][key] = Parameter{
			name:         key,
			parsed:       result,
			raw:          parameter.raw,
			description:  parameter.description,
			wasRequested: true,
		}
	}

	var parameter = request.parameters[from][key]
	for _, config := range configs {
		if err := config(&parameter); err != nil {
			return err
		}
	}

	parameter.name = key
	request.parameters[from][key] = parameter

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

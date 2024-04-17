package request

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/KlyuchnikovV/engi/api/parameter/placing"
	"github.com/KlyuchnikovV/engi/api/response"
	"github.com/KlyuchnikovV/engi/internal/types"
)

// ExtractParam - extracting parameter from context, calls middleware and saves to 'context.parameters[from][key]'.
// After this parameter can be retrieved from context using 'context.Query' methods.
func ExtractParam(
	key string,
	paramPlacing placing.Placing,
	request *Request,
	configs []Option,
	convert func(string) (interface{}, error),
) error {
	var param = request.GetParameter(key, paramPlacing)
	if len(param) == 0 {
		return fmt.Errorf("parameter '%s' not found", key)
	}

	result, err := convert(param)
	if err != nil {
		return err
	}

	if result != nil {
		var parameter = request.Parameters[paramPlacing][key]

		request.Parameters[paramPlacing][key] = Parameter{
			Name:         key,
			Parsed:       result,
			raw:          parameter.raw,
			Description:  parameter.Description,
			wasRequested: true,
		}
	}

	var parameter = request.Parameters[paramPlacing][key]
	for _, config := range configs {
		if err := config(&parameter); err != nil {
			return err
		}
	}

	parameter.Name = key
	request.Parameters[paramPlacing][key] = parameter

	return nil
}

func ExtractBody(
	request *Request,
	unmarshaler types.Unmarshaler,
	pointer interface{},
	configs []Option,
) error {
	if request.body.Parsed == nil {
		if err := readBody(request); err != nil {
			return err
		}

		if len(request.body.raw) == 0 {
			return fmt.Errorf("no body found after reading")
		}
	}

	if err := unmarshaler([]byte(request.body.raw[0]), pointer); err != nil {
		return err
	}

	for _, config := range configs {
		if err := config(&request.body); err != nil {
			return err
		}
	}

	return nil
}

func GetUnmarshaler(request *Request) (types.Unmarshaler, error) {
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
				return response.AsError(http.StatusInternalServerError, "pointer must be of type '*string'")
			}

			*typed = string(b)

			return nil
		}
	default:
		return nil, response.AsError(http.StatusBadRequest, "content-type not supported: %s", contentType)
	}

	return func(bytes []byte, pointer interface{}) error {
		if err := unmarshal(bytes, pointer); err != nil {
			return response.AsError(http.StatusInternalServerError, "unmarshaling body failed: %s", err.Error())
		}

		request.body.wasRequested = true
		request.body.Parsed = pointer

		return nil
	}, nil
}

func readBody(request *Request) error {
	defer request.request.Body.Close()

	bytes, err := io.ReadAll(request.request.Body)
	if err != nil && !errors.Is(err, http.ErrBodyReadAfterClose) {
		return response.AsError(http.StatusInternalServerError, "reading body failed: %s", err.Error())
	}

	if len(bytes) != 0 {
		request.body.raw = []string{string(bytes)}
	}

	if len(request.body.raw) == 0 {
		return response.AsError(http.StatusBadRequest, "no required body provided")
	}

	return err
}

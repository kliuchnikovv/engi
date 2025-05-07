package request

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KlyuchnikovV/engi/parameter/placing"
	"github.com/KlyuchnikovV/engi/response"
)

const (
	IntBase int = 10
	BitSize int = 64
)

type (
	Option          func(*Parameter) error
	Middleware      func(r *Request, w http.ResponseWriter) *response.AsObject
	ParamsValidator interface {
		Validate(param string) error
	}

	Requester interface {
		// Headers - returns request headers.
		Headers() map[string][]string
		// Parameters - returns all parsed parameters.
		Parameters() map[placing.Placing]map[string]string
		// GetParameter - returns parameter value from defined place.
		GetParameter(value string, place placing.Placing) string
		// GetRequest - return http.Request object associated with request.
		GetRequest() *http.Request
		// Body - returns request body.
		// Body must be requested by 'api.Body(pointer)' or 'api.CustomBody(unmarshaler, pointer)'.
		Body() interface{}
		// Bool - returns boolean parameter.
		// Mandatory parameter should be requested by 'api.Bool'.
		// Otherwise, parameter will be obtained by key and its value will be checked for truth.
		Bool(value string, place placing.Placing) bool
		// Integer - returns integer parameter.
		// Mandatory parameter should be requested by 'api.Integer'.
		// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
		Integer(value string, place placing.Placing) int64
		// Float - returns floating point number parameter.
		// Mandatory parameter should be requested by 'api.Float'.
		// Otherwise, parameter will be obtained by key and its value will be converted to float64.
		Float(value string, place placing.Placing) float64
		// String - returns String parameter.
		// Mandatory parameter should be requested by 'api.String'.
		// Otherwise, parameter will be obtained by key.
		String(value string, place placing.Placing) string
		// Time - returns date-time parameter.
		// Mandatory parameter should be requested by 'api.Time'.
		// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
		Time(key string, layout string, paramPlacing placing.Placing) time.Time
	}
)

type Parameter struct {
	raw          []string
	Parsed       interface{}
	wasRequested bool

	Name        string
	Description string
}

// Request - provide methods for extracting r parameters from context.
type Request struct {
	request *http.Request

	body       Parameter
	parameters map[placing.Placing]map[string]Parameter

	Description string
}

func New(request *http.Request) *Request {
	var (
		headers    = request.Header
		cookies    = request.Cookies()
		parameters = request.URL.Query()
		r          = Request{
			request:    request,
			parameters: make(map[placing.Placing]map[string]Parameter),
		}
	)

	r.parameters[placing.InQuery] = make(map[string]Parameter, len(parameters))

	for key, param := range parameters {
		r.parameters[placing.InQuery][key] = Parameter{
			Name: key,
			raw:  param,
		}
	}

	r.parameters[placing.InCookie] = make(map[string]Parameter, len(cookies))

	for _, cookie := range cookies {
		r.parameters[placing.InCookie][cookie.Name] = Parameter{
			Name: cookie.Name,
			raw:  []string{cookie.Value},
		}
	}

	r.parameters[placing.InHeader] = make(map[string]Parameter, len(headers))

	for key, value := range headers {
		r.parameters[placing.InHeader][key] = Parameter{
			Name: key,
			raw:  value,
		}
	}

	return &r
}

func (r *Request) Bool(key string, paramPlacing placing.Placing) bool {
	if r.isMandatoryParam(key, paramPlacing) {
		if result, ok := r.parameters[paramPlacing][key].Parsed.(bool); ok {
			return result
		}

		panic(fmt.Errorf("conversion parameter to bool failed (key: %s)", key))
	}

	result, _ := strconv.ParseBool(r.String(key, paramPlacing))

	return result
}

func (r *Request) Integer(key string, paramPlacing placing.Placing) int64 {
	if r.isMandatoryParam(key, paramPlacing) {
		if result, ok := r.parameters[paramPlacing][key].Parsed.(int64); ok {
			return result
		}

		panic(fmt.Errorf("conversion parameter to int64 failed (key: %s)", key))
	}

	result, _ := strconv.ParseInt(r.String(key, paramPlacing), IntBase, BitSize)

	return result
}

func (r *Request) Float(key string, paramPlacing placing.Placing) float64 {
	if r.isMandatoryParam(key, paramPlacing) {
		if result, ok := r.parameters[paramPlacing][key].Parsed.(float64); ok {
			return result
		}

		panic(fmt.Errorf("conversion parameter to float64 failed (key: %s)", key))
	}

	result, _ := strconv.ParseFloat(r.String(key, paramPlacing), BitSize)

	return result
}

func (r *Request) String(key string, paramPlacing placing.Placing) string {
	if r.isMandatoryParam(key, paramPlacing) {
		if result, ok := r.parameters[paramPlacing][key].Parsed.(string); ok {
			return result
		}

		panic(fmt.Errorf("conversion parameter to string failed (key: %s)", key))
	}

	return r.GetParameter(key, paramPlacing)
}

func (r *Request) Time(key, layout string, paramPlacing placing.Placing) time.Time {
	if r.isMandatoryParam(key, paramPlacing) {
		if result, ok := r.parameters[paramPlacing][key].Parsed.(time.Time); ok {
			return result
		}

		panic(fmt.Errorf("conversion parameter to time failed (key: %s)", key))
	}

	result, _ := time.Parse(layout, r.String(key, paramPlacing))

	return result
}

func (r *Request) Parameters() map[placing.Placing]map[string]string {
	var parameters = make(map[placing.Placing]map[string]string)

	for place, params := range r.parameters {
		parameters[place] = make(map[string]string)

		for name := range params {
			parameters[place][name] = r.GetParameter(name, place)
		}
	}

	return parameters
}

func (r *Request) Body() interface{} {
	return r.body.Parsed
}

func (r *Request) isMandatoryParam(key string, paramPlacing placing.Placing) bool {
	switch paramPlacing {
	case placing.InPath:
		return true
	case placing.InQuery:
		param, ok := r.parameters[paramPlacing][key]
		return ok && param.wasRequested
	default:
		return false
	}
}

func (r *Request) GetParameter(key string, paramPlacing placing.Placing) string {
	if _, ok := r.parameters[paramPlacing][key]; !ok {
		return ""
	}

	if len(r.parameters[paramPlacing][key].raw) > 1 {
		return strings.Join(r.parameters[paramPlacing][key].raw, ", ")
	}

	return r.parameters[paramPlacing][key].raw[0]
}

func (r *Request) GetRequest() *http.Request {
	return r.request
}

func (r *Request) AddInPathParameter(key string, value string) {
	if r.parameters[placing.InPath] == nil {
		r.parameters[placing.InPath] = make(map[string]Parameter)
	}

	r.parameters[placing.InPath][key] = Parameter{
		raw:          []string{value},
		wasRequested: true,
		Name:         key,
	}
}

func (r *Request) Headers() map[string][]string {
	return r.request.Header
}

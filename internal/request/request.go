package request

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KlyuchnikovV/webapi/placing"
)

const (
	IntBase int = 10
	BitSize int = 64
)

type (
	HandlerParams func(*Request, http.ResponseWriter) error
	Option        func(*Parameter) error
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

// Bool - returns boolean parameter.
// Mandatory parameter should be requested by 'api.Bool'.
// Otherwise, parameter will be obtained by key and its value will be checked for truth.
func (r *Request) Bool(key string, paramPlacing placing.Placing) bool {
	if r.isMandatoryParam(key, paramPlacing) {
		return r.parameters[paramPlacing][key].Parsed.(bool)
	}

	result, _ := strconv.ParseBool(r.String(key, paramPlacing))

	return result
}

// QueryInteger - returns integer parameter.
// Mandatory parameter should be requested by 'api.Integer'.
// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
func (r *Request) Integer(key string, paramPlacing placing.Placing) int64 {
	if r.isMandatoryParam(key, paramPlacing) {
		return r.parameters[paramPlacing][key].Parsed.(int64)
	}

	result, _ := strconv.ParseInt(r.String(key, paramPlacing), IntBase, BitSize)

	return result
}

// QueryFloat - returns floating point number parameter.
// Mandatory parameter should be requested by 'api.Float'.
// Otherwise, parameter will be obtained by key and its value will be converted to float64.
func (r *Request) Float(key string, paramPlacing placing.Placing) float64 {
	if r.isMandatoryParam(key, paramPlacing) {
		return r.parameters[paramPlacing][key].Parsed.(float64)
	}

	result, _ := strconv.ParseFloat(r.String(key, paramPlacing), BitSize)

	return result
}

// QueryString - returns String parameter.
// Mandatory parameter should be requested by 'api.String'.
// Otherwise, parameter will be obtained by key.
func (r *Request) String(key string, paramPlacing placing.Placing) string {
	if r.isMandatoryParam(key, paramPlacing) {
		return r.parameters[paramPlacing][key].Parsed.(string)
	}

	return r.GetParameter(key, paramPlacing)
}

// QueryTime - returns date-time parameter.
// Mandatory parameter should be requested by 'api.Time'.
// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
func (r *Request) Time(key, layout string, paramPlacing placing.Placing) time.Time {
	if r.isMandatoryParam(key, paramPlacing) {
		return r.parameters[paramPlacing][key].Parsed.(time.Time)
	}

	result, _ := time.Parse(layout, r.String(key, paramPlacing))

	return result
}

// All - returns all parameters.
func (r *Request) All() map[string]string {
	var parameters = make(map[string]string)

	for place, params := range r.parameters {
		for name := range params {
			parameters[name] = r.GetParameter(name, place)
		}
	}

	return parameters
}

// Body - returns request body.
// Body must be requested by 'api.Body(pointer)' or 'api.CustomBody(unmarshaler, pointer)'.
func (r *Request) Body() interface{} {
	return r.body.Parsed
}

// isMandatoryParam - checks if parameter was requested.
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

package param

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	intBase = 10
	bitSize = 64
)

// Request - provide methods for extracting r parameters from context.
type Request struct {
	request *http.Request

	body       Parameter
	parameters map[paramPlacing]map[string]Parameter
}

func NewRequest(request *http.Request) *Request {
	var (
		parameters = request.URL.Query()
		r          = Request{
			request:    request,
			parameters: make(map[paramPlacing]map[string]Parameter),
		}
	)

	r.parameters[query] = make(map[string]Parameter, len(parameters))

	for key, param := range parameters {
		r.parameters[query][key] = Parameter{
			name: key,
			raw:  param,
		}
	}

	return &r
}

// InPathBool - returns boolean in path parameter.
// Mandatory parameter should be requested by 'api.InPathBool(key)'.
// Otherwise, parameter will be obtained by key and its value will be checked for truth.
func (r *Request) InPathBool(key string) bool {
	return r.bool(inPath, key)
}

// QueryBool - returns boolean query parameter.
// Mandatory parameter should be requested by 'api.QueryBool(key)'.
// Otherwise, parameter will be obtained by key and its value will be checked for truth.
func (r *Request) QueryBool(key string) bool {
	return r.bool(query, key)
}

func (r *Request) bool(from paramPlacing, key string) bool {
	if r.isMandatoryParam(from, key) {
		return r.parameters[from][key].parsed.(bool)
	}

	result, _ := strconv.ParseBool(r.string(from, key))

	return result
}

// InPathInteger - returns integer in path parameter.
// Mandatory parameter should be requested by 'api.InPathInteger(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
func (r *Request) InPathInteger(key string) int64 {
	return r.integer(inPath, key)
}

// QueryInteger - returns integer query parameter.
// Mandatory parameter should be requested by 'api.QueryInteger(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
func (r *Request) QueryInteger(key string) int64 {
	return r.integer(query, key)
}

func (r *Request) integer(from paramPlacing, key string) int64 {
	if r.isMandatoryParam(from, key) {
		return r.parameters[from][key].parsed.(int64)
	}

	result, _ := strconv.ParseInt(r.string(from, key), intBase, bitSize)

	return result
}

// InPathFloat - returns floating point number in path parameter.
// Mandatory parameter should be requested by 'api.InPathFloat(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted to float64.
func (r *Request) InPathFloat(key string) float64 {
	return r.float(inPath, key)
}

// QueryFloat - returns floating point number query parameter.
// Mandatory parameter should be requested by 'api.QueryFloat(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted to float64.
func (r *Request) QueryFloat(key string) float64 {
	return r.float(query, key)
}

func (r *Request) float(from paramPlacing, key string) float64 {
	if r.isMandatoryParam(from, key) {
		return r.parameters[from][key].parsed.(float64)
	}

	result, _ := strconv.ParseFloat(r.string(from, key), bitSize)

	return result
}

// InPathString - returns string in path parameter.
// Mandatory parameter should be requested by 'api.InPathString(key)'.
// Otherwise, parameter will be obtained by key.
func (r *Request) InPathString(key string) string {
	return r.string(inPath, key)
}

// QueryString - returns string query parameter.
// Mandatory parameter should be requested by 'api.QueryString(key)'.
// Otherwise, parameter will be obtained by key.
func (r *Request) QueryString(key string) string {
	return r.string(query, key)
}

func (r *Request) string(from paramPlacing, key string) string {
	if r.isMandatoryParam(from, key) {
		return r.parameters[from][key].parsed.(string)
	}

	return r.getParameter(from, key)
}

// InPathTime - returns date-time in path parameter.
// Mandatory parameter should be requested by 'api.InPathTime(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
func (r *Request) InPathTime(key, layout string) time.Time {
	return r.time(inPath, key, layout)
}

// QueryTime - returns date-time  query parameter.
// Mandatory parameter should be requested by 'api.QueryTime(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
func (r *Request) QueryTime(key, layout string) time.Time {
	return r.time(query, key, layout)
}

func (r *Request) time(from paramPlacing, key, layout string) time.Time {
	if r.isMandatoryParam(from, key) {
		return r.parameters[from][key].parsed.(time.Time)
	}

	result, _ := time.Parse(layout, r.string(from, key))

	return result
}

// All - returns all parameters.
func (r *Request) All() map[string]string {
	var parameters = make(map[string]string)

	for place, params := range r.parameters {
		for name := range params {
			parameters[name] = r.getParameter(place, name)
		}
	}

	return parameters
}

// Body - returns request body.
// Body must be requested by 'api.Body(pointer)' or 'api.CustomBody(unmarshaler, pointer)'.
func (r *Request) Body() interface{} {
	return r.body.parsed
}

// isMandatoryParam - checks if parameter was requested.
func (r *Request) isMandatoryParam(from paramPlacing, key string) bool {
	switch from {
	case inPath:
		return true
	case query:
		param, ok := r.parameters[from][key]
		return ok && param.wasRequested
	default:
		return false
	}
}

func (r *Request) getParameter(from paramPlacing, key string) string {
	if _, ok := r.parameters[from][key]; !ok {
		return ""
	}

	if len(r.parameters[from][key].raw) > 1 {
		return strings.Join(r.parameters[from][key].raw, ", ")
	}

	return r.parameters[from][key].raw[0]
}

func (r *Request) Request() *http.Request {
	return r.request
}

func (r *Request) AddInPathParameter(key string, value string) {
	if r.parameters[inPath] == nil {
		r.parameters[inPath] = make(map[string]Parameter)
	}

	r.parameters[inPath][key] = Parameter{
		raw:          []string{value},
		wasRequested: true,
		name:         key,
	}
}

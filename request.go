package webapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type parameter struct {
	raw string

	wasRequested bool
	parsed       interface{}
}

// TODO: rename to request?

// Request - provide methods for extracting r parameters from context.
type Request struct {
	request *http.Request

	body       parameter
	parameters map[string]parameter
}

func NewRequest(request *http.Request) *Request {
	var (
		parameters = request.URL.Query()
		r          = Request{
			request:    request,
			parameters: make(map[string]parameter, len(parameters)),
		}
	)

	for key, param := range parameters {
		r.parameters[key] = parameter{
			raw: param[0], // TODO: save slice
		}
	}

	return &r
}

// Bool - returns boolean parameter.
// Mandatory parameter should be requested by 'api.WithBool(key)'.
// Otherwise, parameter will be obtained by key and its value will be checked for truth.
func (r *Request) Bool(key string) bool {
	if r.IsMandatoryParam(key) {
		return r.parameters[key].parsed.(bool)
	}

	return strings.ToLower(r.String(key)) == "true"
}

// Integer - returns integer parameter.
// Mandatory parameter should be requested by 'api.WithInt(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
func (r *Request) Integer(key string) int64 {
	if r.IsMandatoryParam(key) {
		return r.parameters[key].parsed.(int64)
	}

	var (
		intBase = 10
		bitSize = 64
	)

	result, _ := strconv.ParseInt(r.String(key), intBase, bitSize)

	return result
}

// IsMandatoryParam - checks if parameter was requested.
func (r *Request) IsMandatoryParam(key string) bool {
	param, ok := r.parameters[key]
	return ok && param.wasRequested
}

// Float - returns floating point number parameter.
// Mandatory parameter should be requested by 'api.WithFloat(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted to float64.
func (r *Request) Float(key string) float64 {
	if r.IsMandatoryParam(key) {
		return r.parameters[key].parsed.(float64)
	}

	var bitSize = 64

	result, _ := strconv.ParseFloat(r.String(key), bitSize)

	return result
}

// String - returns string parameter.
// Mandatory parameter should be requested by 'api.WithString(key)'.
// Otherwise, parameter will be obtained by key.
func (r *Request) String(key string) string {
	if r.IsMandatoryParam(key) {
		return r.parameters[key].parsed.(string)
	}

	param, ok := r.parameters[key]
	if !ok {
		return ""
	}

	return param.raw
}

// Time - returns boolean parameter.
// Mandatory parameter should be requested by 'api.WithTime(key, layout)'.
// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
func (r *Request) Time(key, layout string) time.Time {
	if r.IsMandatoryParam(key) {
		return r.parameters[key].parsed.(time.Time)
	}

	result, _ := time.Parse(layout, r.String(key))

	return result
}

// All - returns all parameters.
func (r *Request) All() map[string]string {
	var parameters = make(map[string]string)

	for name, param := range r.parameters {
		parameters[name] = param.raw
	}

	return parameters
}

func (r *Request) getParam(key string) string {
	return r.parameters[key].raw
}

func (r *Request) updateParam(key string, value interface{}) {
	param := r.parameters[key]

	param.wasRequested = true
	param.parsed = value

	r.parameters[key] = param
}

// TODO: refactor
func (r *Request) Bind(pointer interface{}) error {
	defer r.request.Body.Close()

	bytes, err := ioutil.ReadAll(r.request.Body)
	if err != nil && !errors.Is(err, http.ErrBodyReadAfterClose) {
		return err
	}

	if len(bytes) != 0 {
		r.body.raw = string(bytes)
	}

	if len(r.body.raw) == 0 {
		return fmt.Errorf("no body provided")
	}

	if r.body.parsed != nil {
		return nil
	}

	var (
		ct          = r.request.Header.Get("Content-type")
		unmarshaler func([]byte, interface{}) error
	)

	switch ct {
	case "application/json":
		unmarshaler = json.Unmarshal
	case "application/xml":
		unmarshaler = xml.Unmarshal
	default:
		return fmt.Errorf("content-type not supported: %s", ct)
	}

	if err := unmarshaler([]byte(r.body.raw), pointer); err != nil {
		return err
	}

	r.body.wasRequested = true
	r.body.parsed = pointer

	return nil
}

func (r *Request) Body() interface{} {
	return r.body.parsed
}

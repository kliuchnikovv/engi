package webapi

import (
	"strconv"
	"strings"
	"time"
)

// QueryParams - provide methods for extracting query parameters from context.
type QueryParams struct {
	*Context
}

func NewQueryParams(ctx *Context) QueryParams {
	return QueryParams{
		Context: ctx,
	}
}

// Bool - returns boolean parameter.
// Mandatory parameter should be requested by 'api.WithBool(key)'.
// Otherwise, parameter will be obtained by key and its value will be checked for truth.
func (query *QueryParams) Bool(key string) bool {
	if query.IsMandatoryParam(key) {
		return query.queryParameters[key].(bool)
	}

	return strings.ToLower(query.String(key)) == "true"
}

// Integer - returns integer parameter.
// Mandatory parameter should be requested by 'api.WithInt(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
func (query *QueryParams) Integer(key string) int64 {
	if query.IsMandatoryParam(key) {
		return query.queryParameters[key].(int64)
	}

	var (
		intBase = 10
		bitSize = 64
	)

	result, _ := strconv.ParseInt(query.String(key), intBase, bitSize)

	return result
}

// IsMandatoryParam - checks if parameter was requested.
func (query *QueryParams) IsMandatoryParam(key string) bool {
	_, ok := query.requestedParams[key]
	return ok
}

// Float - returns floating point number parameter.
// Mandatory parameter should be requested by 'api.WithFloat(key)'.
// Otherwise, parameter will be obtained by key and its value will be converted to float64.
func (query *QueryParams) Float(key string) float64 {
	if query.IsMandatoryParam(key) {
		return query.queryParameters[key].(float64)
	}

	var bitSize = 64

	result, _ := strconv.ParseFloat(query.String(key), bitSize)

	return result
}

// String - returns string parameter.
// Mandatory parameter should be requested by 'api.WithString(key)'.
// Otherwise, parameter will be obtained by key.
func (query *QueryParams) String(key string) string {
	if query.IsMandatoryParam(key) {
		return query.queryParameters[key].(string)
	}

	return query.Context.Context.QueryParam(key)
}

// Time - returns boolean parameter.
// Mandatory parameter should be requested by 'api.WithTime(key, layout)'.
// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
func (query *QueryParams) Time(key, layout string) time.Time {
	if query.IsMandatoryParam(key) {
		return query.queryParameters[key].(time.Time)
	}

	result, _ := time.Parse(layout, query.String(key))

	return result
}

// All - returns all parameters.
func (query *QueryParams) All() map[string]string {
	var parameters = make(map[string]string)

	for _, name := range query.Context.Context.ParamNames() {
		parameters[name] = query.Context.Context.Param(name)
	}

	return parameters
}

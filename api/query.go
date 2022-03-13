package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TODO: errors

type QueryParams struct {
	*Context
}

func NewQueryParams(ctx *Context) QueryParams {
	return QueryParams{
		Context: ctx,
	}
}

func (query *QueryParams) Body() interface{} {
	if query.bodyRequested {
		return query.body
	}
	return nil
}

func (query *QueryParams) Bool(key string) bool {
	if query.IsObligatoryParam(key) {
		return query.queryParameters[key].(bool)
	}

	return strings.ToLower(query.String(key)) == "true"
}

func (query *QueryParams) Integer(key string) int64 {
	if query.IsObligatoryParam(key) {
		return query.queryParameters[key].(int64)
	}

	result, _ := strconv.ParseInt(query.String(key), 10, 64)
	return result
}

func (query *QueryParams) IsObligatoryParam(key string) bool {
	_, ok := query.requestedParams[key]
	return ok
}

func (query *QueryParams) Float(key string) float64 {
	if query.IsObligatoryParam(key) {
		return query.queryParameters[key].(float64)
	}

	result, _ := strconv.ParseFloat(query.String(key), 64)
	return result
}

func (query *QueryParams) String(key string) string {
	if query.IsObligatoryParam(key) {
		return query.queryParameters[key].(string)
	}
	return query.Context.Context.QueryParam(key)
}

func (query *QueryParams) Time(key, layout string) time.Time {
	if query.IsObligatoryParam(key) {
		return query.queryParameters[key].(time.Time)
	}

	result, _ := time.Parse(layout, query.String(key))
	return result
}

func (query *QueryParams) Parameters() map[string]string {
	var parameters = make(map[string]string, len(query.queryParameters))
	for i, param := range query.queryParameters {
		switch typed := param.(type) {
		case string:
			parameters[i] = typed
		case []byte:
			parameters[i] = string(typed)
		default:
			parameters[i] = fmt.Sprintf("%v", typed)
		}
	}
	return parameters
}

func (query *QueryParams) StringParamURL(key string) string {
	return query.Context.Context.Param(key)
}

func (query *QueryParams) URLParameters() map[string]string {
	var (
		names  = query.Context.Context.ParamNames()
		result = make(map[string]string, len(names))
	)
	for _, name := range names {
		result[name] = query.Context.Context.Param(name)
	}
	return result
}

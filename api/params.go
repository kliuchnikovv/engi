package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func (api *BasicAPI) getParam(
	ctx *Context,
	key string,
	convert func(string) (interface{}, error),
) error {
	ctx.requestedParams[key] = struct{}{}

	var param = ctx.Context.QueryParam(key)
	if len(param) == 0 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Errorf("parameter '%s' not found", key),
		}
	}

	result, err := convert(param)
	if err != nil {
		return err
	}

	ctx.queryParameters[key] = result
	return nil
}

func (api *BasicAPI) WithBody(pointer interface{}) HandlerFunc {
	return func(ctx *Context) error {
		ctx.bodyRequested = true

		if err := ctx.Context.Bind(pointer); err != nil {
			return err
		}

		ctx.body = pointer

		return nil
	}
}

func (api *BasicAPI) WithBool(key string) HandlerFunc {
	return func(ctx *Context) error {
		return api.getParam(ctx, key, func(param string) (interface{}, error) {
			return strings.ToLower(param) == "true", nil
		})
	}
}

func (api *BasicAPI) WithInt(key string) HandlerFunc {
	return func(ctx *Context) error {
		return api.getParam(ctx, key, func(param string) (interface{}, error) {
			result, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return nil, &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  fmt.Errorf("parameter '%s' not of type int", key),
					Internal: err,
				}
			}
			return result, err
		})
	}
}

func (api *BasicAPI) WithFloat(key string) HandlerFunc {
	return func(ctx *Context) error {
		return api.getParam(ctx, key, func(param string) (interface{}, error) {
			result, err := strconv.ParseFloat(param, 64)
			if err != nil {
				return nil, &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  fmt.Errorf("parameter '%s' not of type int", key),
					Internal: err,
				}
			}
			return result, err
		})
	}
}

func (api *BasicAPI) WithString(key string) HandlerFunc {
	return func(ctx *Context) error {
		return api.getParam(ctx, key, func(param string) (interface{}, error) {
			return param, nil
		})
	}
}

func (api *BasicAPI) WithTime(key, layout string) HandlerFunc {
	return func(ctx *Context) error {
		return api.getParam(ctx, key, func(param string) (interface{}, error) {
			result, err := time.Parse(layout, param)
			if err != nil {
				return nil, &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  fmt.Errorf("could not parse '%s' param to datetime using '%s' layout", key, layout),
					Internal: err,
				}
			}
			return result, err
		})
	}
}

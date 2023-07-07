package validate

import (
	"net/http"
	"strings"
	"time"

	"github.com/KlyuchnikovV/webapi/internal/request"
	"github.com/KlyuchnikovV/webapi/response"
)

// NotEmpty - checks if parameter is not empty by it's type.
// NOTE: boolean parameter will be ignored.
func NotEmpty(p *request.Parameter) error {
	var isNotEmpty func() bool

	switch typed := p.Parsed.(type) {
	case bool:
		// Bool can't be empty
		return nil
	case int64:
		isNotEmpty = func() bool { return typed != 0 }
	case float64:
		isNotEmpty = func() bool { return typed != 0 }
	case string:
		isNotEmpty = func() bool { return len(typed) != 0 }
	case time.Time:
		isNotEmpty = func() bool { return typed.UnixNano() != 0 }
	default:
		isNotEmpty = func() bool { return typed != nil }
	}

	if isNotEmpty() {
		return nil
	}

	return response.NewError(http.StatusBadRequest,
		"'%s' shouldn't be empty", p.Name,
	)
}

// Greater - checks if parameter greater than a number.
// NOTES:
//   - for 'int' and 'float' parameters - simple values comparison;
//   - for 'string' - comparing with it's length;
//   - for 'time' - comparing with time.Unix() value in seconds;
func Greater(than float64) request.Option {
	return func(p *request.Parameter) error {
		var greater func() bool

		switch typed := p.Parsed.(type) {
		case int64:
			greater = func() bool { return typed > int64(than) }
		case float64:
			greater = func() bool { return typed > than }
		case string:
			greater = func() bool { return len(typed) > int(than) }
		case time.Time:
			greater = func() bool { return typed.Unix() > int64(than) }
		}

		if greater() {
			return nil
		}

		return response.NewError(http.StatusBadRequest,
			"'%s' should be greater than %f", p.Name, than,
		)
	}
}

// Less - checks if parameter less than a number.
// NOTES:
//   - for 'int' and 'float' parameters - simple values comparison;
//   - for 'string' - comparing with it's length;
//   - for 'time' - comparing with time.Unix() value in seconds;
func Less(than float64) request.Option {
	return func(p *request.Parameter) error {
		var greater func() bool

		switch typed := p.Parsed.(type) {
		case int64:
			greater = func() bool { return typed < int64(than) }
		case float64:
			greater = func() bool { return typed < than }
		case string:
			greater = func() bool { return len(typed) < int(than) }
		case time.Time:
			greater = func() bool { return typed.UnixNano() < int64(than) }
		}

		if greater() {
			return nil
		}

		return response.NewError(http.StatusBadRequest,
			"'%s' should be less than %f", p.Name, than,
		)
	}
}

// OR - combines several parameter checks and passes if one of them successful.
func OR(opts ...request.Option) request.Option {
	return func(p *request.Parameter) error {
		var (
			passed bool
			errs   []string
		)

		for _, option := range opts {
			if err := option(p); err == nil {
				passed = true
				break
			} else {
				errs = append(errs, err.Error())
			}
		}

		if passed {
			return nil
		}

		return response.NewError(http.StatusBadRequest,
			"'%s' failed check: %s", p.Name, strings.Join(errs, " and "),
		)
	}
}

// AND - combines several parameter checks and failing if one of them failed.
func AND(opts ...request.Option) request.Option {
	return func(p *request.Parameter) error {
		var err error

		for _, option := range opts {
			if err = option(p); err != nil {
				break
			}
		}

		if err == nil {
			return nil
		}

		return response.NewError(http.StatusBadRequest,
			"'%s' failed check: %s", p.Name, err,
		)
	}
}

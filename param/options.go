package param

import (
	"net/http"
	"strings"
	"time"

	"github.com/KlyuchnikovV/webapi/types"
)

// Description - adds string description to parameter.
// Can be used in errors description or in documentation.
func Description(s string) ParametersOption {
	return func(p *Parameter) error {
		p.description = s

		return nil
	}
}

// NotEmpty - checks if parameter is not empty by it's type.
// NOTE: boolean parameter will be ignored.
func NotEmpty(p *Parameter) error {
	var isNotEmpty func() bool

	switch typed := p.parsed.(type) {
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
	}

	if isNotEmpty() {
		return nil
	}

	return types.NewErrorResponse(http.StatusBadRequest,
		"'%s' shouldn't be empty", p.name,
	)
}

// Greater - checks if parameter greater than a number.
// NOTES:
//    - for 'int' and 'float' parameters - simple values comparison;
//    - for 'string' - comparing with it's length;
//    - for 'time' - comparing with time.Unix() value in seconds;
func Greater(than float64) ParametersOption {
	return func(p *Parameter) error {
		var greater func() bool

		switch typed := p.parsed.(type) {
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

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' should be greater than %f", p.name, than,
		)
	}
}

// Less - checks if parameter less than a number.
// NOTES:
//    - for 'int' and 'float' parameters - simple values comparison;
//    - for 'string' - comparing with it's length;
//    - for 'time' - comparing with time.Unix() value in seconds;
func Less(than float64) ParametersOption {
	return func(p *Parameter) error {
		var greater func() bool

		switch typed := p.parsed.(type) {
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

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' should be less than %f", p.name, than,
		)
	}
}

// OR - combines several parameter checks and passes if one of them successful.
func OR(options ...ParametersOption) ParametersOption {
	return func(p *Parameter) error {
		var (
			passed bool
			errs   []string
		)

		for _, option := range options {
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

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' failed check: %s", p.name, strings.Join(errs, " and "),
		)
	}
}

// AND - combines several parameter checks and failing if one of them failed.
func AND(options ...ParametersOption) ParametersOption {
	return func(p *Parameter) error {
		var err error

		for _, option := range options {
			if err = option(p); err != nil {
				break
			}
		}

		if err == nil {
			return nil
		}

		return types.NewErrorResponse(http.StatusBadRequest,
			"'%s' failed check: %s", p.name, err,
		)
	}
}

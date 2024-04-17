package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	corsOriginMatchAll       string = "*"
	corsOriginHeader         string = "Origin"
	corsAllowOriginHeader    string = "Access-Control-Allow-Origin"
	corsAllowHeadersHeader   string = "Access-Control-Allow-Headers"
	corsRequestMethodHeader  string = "Access-Control-Request-Method"
	corsRequestHeadersHeader string = "Access-Control-Request-Headers"
	corsOptionMethod         string = http.MethodOptions
)

var (
	ErrOriginNotAllowed           = errors.New("origin is not allowed")
	ErrMissingConanicalHeader     = errors.New("missing canonical header")
	ErrCORSMethodHeaderNotFound   = errors.New("CORS-Method header not found")
	ErrCORSMethodHeaderNotAllowed = errors.New("CORS-Method header not allowed")

	defaultCorsHeaders = []string{
		"Accept", "Accept-Language", "Content-Language", "Origin",
	}
)

func (route *Route) cors(
	r *http.Request,
	w http.ResponseWriter,
) (int, error) {
	var origin = r.Header.Get(corsOriginHeader)

	if !contains(route.allowedOrigins, origin) && !contains(route.allowedOrigins, corsOriginMatchAll) {
		return http.StatusForbidden, fmt.Errorf("%w: '%s'", ErrOriginNotAllowed, origin)
	}

	w.Header().Set(corsAllowOriginHeader, origin)

	if r.Method != corsOptionMethod {
		return 0, nil
	}

	var (
		requestHeaders = strings.Split(r.Header.Get(corsRequestHeadersHeader), ",")
		allowedHeaders = make([]string, 0, len(requestHeaders))
	)

	for _, v := range requestHeaders {
		canonicalHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
		if canonicalHeader == "" || contains(defaultCorsHeaders, canonicalHeader) {
			continue
		}

		if !contains(route.allowedHeaders, canonicalHeader) {
			return http.StatusForbidden, ErrMissingConanicalHeader
		}

		allowedHeaders = append(allowedHeaders, canonicalHeader)
	}

	if len(allowedHeaders) > 0 {
		w.Header().Set(corsAllowHeadersHeader, strings.Join(allowedHeaders, ","))
	}

	if _, ok := r.Header[corsRequestMethodHeader]; !ok {
		return http.StatusBadRequest, ErrCORSMethodHeaderNotFound
	}

	method := r.Header.Get(corsRequestMethodHeader)
	if !contains(route.allowedMethods, method) {
		return http.StatusMethodNotAllowed, ErrCORSMethodHeaderNotAllowed
	}

	return 0, nil
}

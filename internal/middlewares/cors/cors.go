package cors

import (
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/response"
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

// OriginValidator takes an origin string and returns whether or not that origin is allowed.
type OriginValidator func(string) bool

type CORS struct {
	allowedHeaders []string
	allowedMethods []string
	allowedOrigins []string
}

func (c *CORS) Handle(request *request.Request, writer http.ResponseWriter) *response.AsObject {
	var (
		r                  = request.GetRequest()
		origin             = r.Header.Get(corsOriginHeader)
		defaultCorsHeaders = []string{"Accept", "Accept-Language", "Content-Language", "Origin"}
	)

	if !contains(c.allowedOrigins, origin) && !contains(c.allowedOrigins, corsOriginMatchAll) {
		return response.AsError(http.StatusForbidden, "origin '%s' is not allowed", origin)
	}

	writer.Header().Set(corsAllowOriginHeader, origin)

	if r.Method != corsOptionMethod {
		return nil
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

		if !contains(c.allowedHeaders, canonicalHeader) {
			return response.AsError(http.StatusForbidden, "")
		}

		allowedHeaders = append(allowedHeaders, canonicalHeader)
	}

	if len(allowedHeaders) > 0 {
		writer.Header().Set(corsAllowHeadersHeader, strings.Join(allowedHeaders, ","))
	}

	if _, ok := r.Header[corsRequestMethodHeader]; !ok {
		return response.AsError(http.StatusBadRequest, "CORS-Method header not found")
	}

	method := r.Header.Get(corsRequestMethodHeader)
	if !contains(c.allowedMethods, method) {
		return response.AsError(http.StatusMethodNotAllowed, "CORS-Method header not found")
	}

	return &response.AsObject{Code: 200}
}

func contains(slice []string, item string) bool {
	if len(slice) == 0 {
		return true
	}

	for _, i := range slice {
		if i == item {
			return true
		}
	}

	return false
}

type CORSOption func(*CORS)

func AllowedHeaders(headers ...string) CORSOption {
	return func(ch *CORS) {
		for _, v := range headers {
			normalizedHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
			if normalizedHeader == "" {
				continue
			}

			if !contains(ch.allowedHeaders, normalizedHeader) {
				ch.allowedHeaders = append(ch.allowedHeaders, normalizedHeader)
			}
		}
	}
}

func AllowedMethods(methods ...string) CORSOption {
	return func(ch *CORS) {
		ch.allowedMethods = make([]string, 0, len(methods))

		for _, v := range methods {
			var method = strings.ToUpper(strings.TrimSpace(v))
			if method == "" {
				continue
			}

			if !contains(ch.allowedMethods, method) {
				ch.allowedMethods = append(ch.allowedMethods, method)
			}
		}
	}
}

func AllowedOrigins(origins ...string) CORSOption {
	return func(ch *CORS) {
		for _, v := range origins {
			if v == corsOriginMatchAll {
				ch.allowedOrigins = []string{corsOriginMatchAll}
				return
			}
		}

		ch.allowedOrigins = origins
	}
}

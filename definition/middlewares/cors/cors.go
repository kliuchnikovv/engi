package cors

import (
	"context"
	"net/http"
	"strings"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/routes"
)

type corsAllowedOrigins []string

func (origins corsAllowedOrigins) Handle(ctx context.Context, req *request.Request, resp *response.Response) error {
	var (
		r      = req.GetRequest()
		w      = resp.ResponseWriter()
		origin = r.Header.Get(corsOriginHeader)
	)

	if !contains(origins, origin) && !contains(origins, corsOriginMatchAll) {
		return resp.Forbidden("%s: '%s'", ErrOriginNotAllowed, origin)
	}

	w.Header().Set(corsAllowOriginHeader, origin)

	return nil
}

func (origins corsAllowedOrigins) Docs(*routes.Route) {
	panic("unimplemented")
}

func (origins corsAllowedOrigins) Priority() int {
	return 10
}

type corsAllowedHeaders []string

func (headers corsAllowedHeaders) Handle(ctx context.Context, req *request.Request, resp *response.Response) error {
	if req.GetRequest().Method != corsOptionMethod {
		return nil
	}

	var (
		r = req.GetRequest()
		w = resp.ResponseWriter()

		requestHeaders = strings.Split(r.Header.Get(corsRequestHeadersHeader), ",")
		allowedHeaders = make([]string, 0, len(requestHeaders))
	)

	for _, v := range requestHeaders {
		canonicalHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
		if canonicalHeader == "" || contains(defaultCorsHeaders, canonicalHeader) {
			continue
		}

		if !contains(headers, canonicalHeader) {
			return resp.Forbidden(ErrMissingConanicalHeader.Error())
		}

		allowedHeaders = append(allowedHeaders, canonicalHeader)
	}

	if len(allowedHeaders) > 0 {
		w.Header().Set(corsAllowHeadersHeader, strings.Join(allowedHeaders, ","))
	}

	return nil
}

func (headers corsAllowedHeaders) Docs(*routes.Route) {
	panic("unimplemented")
}

func (headers corsAllowedHeaders) Priority() int {
	return 11
}

type corsAllowedMethods []string

func (methods corsAllowedMethods) Handle(ctx context.Context, req *request.Request, resp *response.Response) error {
	var r = req.GetRequest()

	if _, ok := r.Header[corsRequestMethodHeader]; !ok {
		return resp.BadRequest(ErrCORSMethodHeaderNotFound.Error())
	}

	method := r.Header.Get(corsRequestMethodHeader)
	if !contains(methods, method) {
		return resp.MethodNotAllowed(ErrCORSMethodHeaderNotAllowed.Error())
	}

	return nil
}

func (methods corsAllowedMethods) Docs(*routes.Route) {
	panic("unimplemented")
}

func (methods corsAllowedMethods) Priority() int {
	return 12 // TODO: make external priority map
}

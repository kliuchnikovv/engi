package routes

import (
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
)

type Option interface {
	Bind(*Route) error
	Handle(*request.Request, *response.Response) error
	Docs(*Route)
}

// func AllowedHeaders(headers ...string) Option {
// 	return func(route *Route) {
// 		for _, v := range headers {
// 			normalizedHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
// 			if normalizedHeader == "" {
// 				continue
// 			}

// 			if !contains(route.allowedHeaders, normalizedHeader) {
// 				route.allowedHeaders = append(route.allowedHeaders, normalizedHeader)
// 			}
// 		}
// 	}
// }

// func AllowedMethods(methods ...string) Option {
// 	return func(route *Route) {
// 		route.allowedMethods = make([]string, 0, len(methods))

// 		for _, v := range methods {
// 			var method = strings.ToUpper(strings.TrimSpace(v))
// 			if method == "" {
// 				continue
// 			}

// 			if !contains(route.allowedMethods, method) {
// 				route.allowedMethods = append(route.allowedMethods, method)
// 			}
// 		}
// 	}
// }

// func AllowedOrigins(origins ...string) Option {
// 	return func(route *Route) {
// 		for _, v := range origins {
// 			if v == corsOriginMatchAll {
// 				route.allowedOrigins = []string{corsOriginMatchAll}
// 				return
// 			}
// 		}

// 		route.allowedOrigins = origins
// 	}
// }

package cors

import "net/http"

const (
	corsOriginMatchAll       string = "*"
	corsOriginHeader         string = "Origin"
	corsAllowOriginHeader    string = "Access-Control-Allow-Origin"
	corsAllowHeadersHeader   string = "Access-Control-Allow-Headers"
	corsRequestMethodHeader  string = "Access-Control-Request-Method"
	corsRequestHeadersHeader string = "Access-Control-Request-Headers"
	corsOptionMethod         string = http.MethodOptions
)

var defaultCorsHeaders = []string{
	"Accept", "Accept-Language", "Content-Language", "Origin",
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

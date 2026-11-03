package engi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/middlewares/auth"
	"github.com/stretchr/testify/assert"
)

type bearerService struct {
	middleware engi.Middleware
	handled    bool
}

func (*bearerService) Prefix() string { return "private" }
func (s *bearerService) Routers() engi.Routes {
	return engi.Routes{engi.GET("resource"): engi.Handle(s.handleGet, s.middleware)}
}

func (s *bearerService) handleGet(_ context.Context, _ engi.Request, resp engi.Response) error {
	s.handled = true
	return resp.OK("protected")
}

func TestBearerAuthentication(t *testing.T) {
	t.Parallel()

	const secret = "secret"
	cases := []struct {
		name, header, token string
		status              int
	}{
		{"valid", "Bearer secret", secret, http.StatusOK},
		{"case insensitive scheme", "bEaReR secret", secret, http.StatusOK},
		{"multiple spaces", "Bearer   secret", secret, http.StatusOK},
		{"invalid", "Bearer wrong", "wrong", http.StatusUnauthorized},
		{"case sensitive token", "Bearer SECRET", "SECRET", http.StatusUnauthorized},
		{"missing", "", "", http.StatusUnauthorized},
		{"wrong scheme", "Token secret", "", http.StatusUnauthorized},
		{"missing separator", "Bearersecret", "", http.StatusUnauthorized},
		{"empty token", "Bearer ", "", http.StatusUnauthorized},
		{"whitespace token", "Bearer   ", "", http.StatusUnauthorized},
		{"extra field", "Bearer secret other", "", http.StatusUnauthorized},
		{"tab separator", "Bearer\tsecret", "", http.StatusUnauthorized},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, callback := range []bool{false, true} {
				var tokens []string

				middleware := auth.BearerToken(secret)
				if callback {
					middleware = auth.BearerFunc(func(token string) bool {
						tokens = append(tokens, token)
						return token == secret
					})
				}
				service := &bearerService{middleware: middleware}
				engine := engi.New("")
				if err := engine.RegisterServices(service); err != nil {
					t.Fatal(err)
				}
				req := httptest.NewRequest(http.MethodGet, "/private/resource", nil)
				req.Header.Set("Authorization", tc.header)
				recorder := httptest.NewRecorder()
				engine.Server().Handler.ServeHTTP(recorder, req)
				assert.Equal(t, tc.status, recorder.Code, "callback=%v", callback)
				assert.Equal(t, tc.status == http.StatusOK, service.handled, "callback=%v", callback)

				if callback {
					var expected []string
					if tc.token != "" {
						expected = []string{tc.token}
					}

					assert.Equal(t, expected, tokens, "validator inputs")
				}
			}
		})
	}
}

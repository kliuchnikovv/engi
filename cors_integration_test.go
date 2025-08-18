package engi_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/middlewares/cors"
	"github.com/stretchr/testify/assert"
)

// corsTestService implements ServiceAPI with CORS middleware
type corsTestService struct{}

func (s *corsTestService) Prefix() string {
	return "test"
}

func (s *corsTestService) Middlewares() []engi.Middleware {
	return []engi.Middleware{
		cors.AllowedOrigins("https://example.com"),
		cors.AllowedMethods("GET", "POST", "PUT"),
		cors.AllowedHeaders("Content-Type", "Authorization"),
	}
}

func (s *corsTestService) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("hello"): engi.Handle(s.handleGet),
		engi.PST("hello"): engi.Handle(s.handlePost),
	}
}

func (s *corsTestService) handleGet(ctx context.Context, req engi.Request, resp engi.Response) error {
	return resp.OK("hello get")
}

func (s *corsTestService) handlePost(ctx context.Context, req engi.Request, resp engi.Response) error {
	return resp.OK("hello post")
}

func TestCORS_Integration_OPTIONS(t *testing.T) {
	// Initialize Engine
	eng := engi.New("")

	// Register service with CORS middleware
	err := eng.RegisterServices(&corsTestService{})
	assert.NoError(t, err)

	// Start test HTTP server
	server := httptest.NewServer(eng.Server().Handler)
	defer server.Close()

	// Create OPTIONS request (preflight CORS request)
	req, err := http.NewRequest(http.MethodOptions, server.URL+"/test/hello", nil)
	assert.NoError(t, err)

	// Set CORS headers
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Validate response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `null`, string(body)) // OPTIONS should return null from response.OK(nil)

	// Validate CORS headers
	assert.Equal(t, "https://example.com", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET,POST,PUT", resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}

func TestCORS_Integration_RegularRequest(t *testing.T) {
	// Initialize Engine
	eng := engi.New("")

	// Register service with CORS middleware
	err := eng.RegisterServices(&corsTestService{})
	assert.NoError(t, err)

	// Start test HTTP server
	server := httptest.NewServer(eng.Server().Handler)
	defer server.Close()

	// Create regular GET request
	req, err := http.NewRequest(http.MethodGet, server.URL+"/test/hello", nil)
	assert.NoError(t, err)

	// Set Origin header
	req.Header.Set("Origin", "https://example.com")

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Validate response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `"hello get"`, string(body))

	// Validate CORS origin header is set
	assert.Equal(t, "https://example.com", resp.Header.Get("Access-Control-Allow-Origin"))
	// Methods and Headers should NOT be set for regular requests
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
}

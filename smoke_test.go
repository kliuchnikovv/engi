package engi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// pingService is a simple ServiceAPI with one GET /ping/ endpoint
type pingService struct{}

// Prefix returns the route prefix for the service
func (s *pingService) Prefix() string {
	return "ping"
}

// Routers defines the mapping of paths to handlers
func (s *pingService) Routers() Routes {
	return Routes{
		"": GET(s.handlePing),
	}
}

// handlePing writes "pong" as plain text
func (s *pingService) handlePing(ctx context.Context, req Request, resp Response) error {
	return resp.OK("pong")
}

func TestPingService_Smoke(t *testing.T) {
	// Initialize Engine without additional prefix
	eng := New("")

	// Register our pingService
	err := eng.RegisterServices(&pingService{})
	assert.NoError(t, err)

	// Start test HTTP server using Engine's handler
	server := httptest.NewServer(eng.server.Handler)
	defer server.Close()

	// Perform GET request to /ping/
	resp, err := http.Get(server.URL + "/ping/")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Validate status and body
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `"pong"`, string(body))
}

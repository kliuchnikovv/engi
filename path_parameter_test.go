package engi_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/parameter/path"
	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/definition/validate"
	"github.com/stretchr/testify/assert"
)

// notesService serves a single path parameter, declared with the public
// "{name}" path syntax.
type notesService struct{}

func (s *notesService) Prefix() string {
	return "notes"
}

func (s *notesService) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("{id}"): engi.Handle(s.handleGet, path.Integer("id", validate.Greater(0))),
	}
}

func (s *notesService) handleGet(ctx context.Context, req engi.Request, resp engi.Response) error {
	return resp.OK(req.Integer("id", placing.InPath))
}

// assetsService serves a trailing wildcard route.
type assetsService struct{}

func (s *assetsService) Prefix() string {
	return "assets"
}

func (s *assetsService) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("*filepath"): engi.Handle(s.handleGet, path.String("filepath")),
	}
}

func (s *assetsService) handleGet(ctx context.Context, req engi.Request, resp engi.Response) error {
	return resp.OK(req.String("filepath", placing.InPath))
}

// TestPathParameterRouting checks that a route declared with a path parameter is
// reachable and that the parameter reaches the handler parsed into its type.
func TestPathParameterRouting(t *testing.T) {
	eng := engi.New("")

	err := eng.RegisterServices(&notesService{})
	assert.NoError(t, err)

	server := httptest.NewServer(eng.Server().Handler)
	defer server.Close()

	t.Run("parameter is extracted", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/notes/42")
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "42", string(body))
	})

	t.Run("parameter is validated", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/notes/abc")
		assert.NoError(t, err)
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("parameter outside of the allowed range is rejected", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/notes/0")
		assert.NoError(t, err)
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// TestWildcardRouting checks that a trailing wildcard route is reachable and that
// the wildcard reaches the handler, including the slashes it swallowed.
func TestWildcardRouting(t *testing.T) {
	eng := engi.New("")

	err := eng.RegisterServices(&assetsService{})
	assert.NoError(t, err)

	server := httptest.NewServer(eng.Server().Handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/assets/css/main.css")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `"css/main.css"`, string(body))
}

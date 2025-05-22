package routes_test

import (
	"testing"

	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/routes"
	"github.com/stretchr/testify/assert"
)

func TestTree_Get(t *testing.T) {
	pathfinder := routes.NewTrie[int]()

	// Dummy handlers
	handlerRoot := 1
	handlerUsers := 2
	handlerProfile := 3
	handlerAsset := 4

	// Register routes (now require method)
	pathfinder.Add("GET", "/", handlerRoot)
	pathfinder.Add("GET", "/users", handlerUsers)
	pathfinder.Add("GET", "/users/:id/profile", handlerProfile)
	pathfinder.Add("GET", "/assets/*filepath", handlerAsset)

	tests := []struct {
		path         string
		expectFound  bool
		expectValue  interface{}
		expectParams map[placing.Placing]map[string]string
	}{
		{"/", true, &handlerRoot, map[placing.Placing]map[string]string{}},
		{"/users", true, &handlerUsers, map[placing.Placing]map[string]string{}},
		{"/users/123", false, nil, nil},
		{"/unknown", false, nil, nil},
		{"/users/123/profile", true, &handlerProfile, map[placing.Placing]map[string]string{
			placing.InPath: {"id": "123"}},
		},
		{"/assets/css/main.css", true, &handlerAsset, map[placing.Placing]map[string]string{
			placing.InPath: {"filepath": "css/main.css"}},
		},
	}

	for _, tc := range tests {
		req := &request.Request{}
		got, err := pathfinder.Get(req, "GET", tc.path)

		if tc.expectFound {
			assert.NoError(t, err, "path %s: unexpected error", tc.path)
			assert.Equal(t, tc.expectValue, got, "path %s: handler mismatch", tc.path)
			assert.Equal(t, tc.expectParams, req.Parameters(), "path %s: params mismatch", tc.path)
		} else {
			assert.Equal(t, routes.ErrNotHandled, err, "path %s: expected ErrNotHandled", tc.path)
		}
	}
}

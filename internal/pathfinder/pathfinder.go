package pathfinder

import (
	"context"
	"regexp"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
)

var parameterRegexp = regexp.MustCompile("{[a-zA-Z]*}")

type Handler func(ctx context.Context, request *request.Request, response *response.Response) error

type PathFinder struct {
	exactHandlers  map[string]Handler
	regexpHandlers []HanlderNode
}

func NewPathFinder() *PathFinder {
	return &PathFinder{
		exactHandlers:  make(map[string]Handler),
		regexpHandlers: make([]HanlderNode, 0),
	}
}

func (finder *PathFinder) Add(path string, handler Handler) {
	if !parameterRegexp.MatchString(path) {
		finder.exactHandlers[path] = handler
		return
	}

	var (
		parts = strings.Split(path, "/")
		tree  = NewHandlerNode(parts[0], nil)
	)

	tree.Add(handler, parts[1:]...)

	finder.regexpHandlers = append(finder.regexpHandlers, tree)
}

func (finder *PathFinder) Handle(
	request *request.Request,
	uri string,
) Handler {
	handler, ok := finder.exactHandlers[uri]
	if ok {
		return handler
	}

	for _, regexpHandler := range finder.regexpHandlers {
		if handler := regexpHandler.GetHandler(request, uri); handler != nil {
			return handler
		}
	}

	return nil
}

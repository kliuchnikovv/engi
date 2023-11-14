package pathfinder

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
)

var (
	ErrNotHandled = errors.New("not handled")
	valueRegexp   = regexp.MustCompile("^[^/]+")
)

type HandlerNodes []HanlderNode

type HanlderNode interface {
	Add(handler Handler, parts ...string)
	Handle(ctx context.Context, request *request.Request, response *response.Response, path string) error
	Equal(node HanlderNode) bool
}

func NewHandlerNode(parameter string, handler Handler) HanlderNode {
	if parameterRegexp.MatchString(parameter) {
		return NewRegexpHandler(
			strings.Trim(parameter, "{}"),
			handler,
		)
	}

	return NewStringHandler(parameter, handler)
}

type StringHandler struct {
	pattern string
	nodes   []HanlderNode
	Handler
}

func NewStringHandler(pattern string, handler Handler) *StringHandler {
	return &StringHandler{
		pattern: pattern,
		Handler: handler,
		nodes:   make([]HanlderNode, 0),
	}
}

func (s *StringHandler) Add(handler Handler, parts ...string) {
	if len(parts) == 0 {
		return
	}

	var newNode = NewHandlerNode(parts[0], handler)
	for _, node := range s.nodes {
		if node.Equal(newNode) {
			node.Add(handler, parts[1:]...)
			return
		}
	}

	newNode.Add(handler, parts[1:]...)
	s.nodes = append(s.nodes, newNode)
}

func (s *StringHandler) Handle(
	ctx context.Context,
	request *request.Request,
	response *response.Response,
	path string,
) error {
	path = strings.TrimLeft(path, "/")

	if !strings.HasPrefix(path, s.pattern) {
		return ErrNotHandled
	}

	var subPath, _ = strings.CutPrefix(path, s.pattern)

	if len(subPath) == 0 {
		return s.Handler(ctx, request, response)
	}

	for _, node := range s.nodes {
		err := node.Handle(ctx, request, response, subPath)
		if err == nil {
			return nil
		}

		if !errors.Is(err, ErrNotHandled) {
			return err
		}
	}

	return ErrNotHandled
}

func (s *StringHandler) Equal(n HanlderNode) bool {
	m2, ok := n.(*StringHandler)
	if !ok {
		return false
	}

	return s.pattern == m2.pattern
}

type RegexpHandler struct {
	name    string
	pattern regexp.Regexp
	nodes   []HanlderNode
	Handler
}

func NewRegexpHandler(name string, handler Handler) *RegexpHandler {
	return &RegexpHandler{
		name:    name,
		pattern: *valueRegexp,
		Handler: handler,
		nodes:   make([]HanlderNode, 0),
	}
}

func (r *RegexpHandler) Add(handler Handler, parts ...string) {
	if len(parts) == 0 {
		return
	}

	var newNode = NewHandlerNode(parts[0], handler)
	for _, node := range r.nodes {
		if node.Equal(newNode) {
			node.Add(handler, parts[1:]...)
			return
		}
	}

	newNode.Add(handler, parts[1:]...)
	r.nodes = append(r.nodes, newNode)
}

func (r *RegexpHandler) Handle(
	ctx context.Context,
	request *request.Request,
	response *response.Response,
	path string,
) error {
	path = strings.TrimLeft(path, "/")

	var parts = strings.Split(path, "/")
	if len(parts) == 0 {
		return ErrNotHandled
	}

	if !r.pattern.MatchString(parts[0]) {
		return ErrNotHandled
	}

	request.AddInPathParameter(r.name, parts[0])

	if len(parts) == 1 {
		return r.Handler(ctx, request, response)
	}

	for _, node := range r.nodes {
		err := node.Handle(ctx, request, response, parts[1])
		if err == nil {
			return nil
		}

		if !errors.Is(err, ErrNotHandled) {
			return err
		}
	}

	return ErrNotHandled
}

func (r *RegexpHandler) Equal(n HanlderNode) bool {
	p2, ok := n.(*RegexpHandler)
	if !ok {
		return false
	}

	return r.name == p2.name
}

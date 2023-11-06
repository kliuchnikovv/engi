package internal

import (
	"regexp"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/context"
)

var (
	parameterRegexp = regexp.MustCompile("{[a-zA-Z]*}")
	valueRegexp     = regexp.MustCompile("^[a-zA-Z0-9]+") // TODO: extend it
)

type HandlerNodes []HanlderNode

type HanlderNode interface {
	Add(context.Handler, ...string)
	Handle(string, *context.Context) bool
	Equal(HanlderNode) bool
}

func NewHandlerNode(parameter string, handler context.Handler) HanlderNode {
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
	context.Handler
}

func NewStringHandler(pattern string, handler context.Handler) *StringHandler {
	return &StringHandler{
		pattern: pattern,
		Handler: handler,
		nodes:   make([]HanlderNode, 0),
	}
}

func (s *StringHandler) Add(handler context.Handler, parts ...string) {
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

func (s *StringHandler) Handle(path string, ctx *context.Context) bool {
	path = strings.TrimLeft(path, "/")

	if !strings.HasPrefix(path, s.pattern) {
		return false
	}

	var subPath, _ = strings.CutPrefix(path, s.pattern)

	if len(subPath) == 0 {
		if err := s.Handler(ctx); err != nil {
			return false
		}
		return true
	}

	for _, node := range s.nodes {
		if node.Handle(subPath, ctx) {
			return true
		}
	}

	return false
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
	context.Handler
}

func NewRegexpHandler(name string, handler context.Handler) *RegexpHandler {
	return &RegexpHandler{
		name:    name,
		pattern: *valueRegexp,
		Handler: handler,
		nodes:   make([]HanlderNode, 0),
	}
}

func (r *RegexpHandler) Add(handler context.Handler, parts ...string) {
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

func (r *RegexpHandler) Handle(path string, ctx *context.Context) bool {
	path = strings.TrimLeft(path, "/")

	var parts = strings.Split(path, "/")
	if len(parts) == 0 {
		return false
	}

	if !r.pattern.MatchString(parts[0]) {
		return false
	}

	ctx.Request.AddInPathParameter(r.name, parts[0])

	if len(parts) == 1 {
		if err := r.Handler(ctx); err != nil {
			return false
		}
		return true
	}

	for _, node := range r.nodes {
		if node.Handle(parts[1], ctx) {
			return true
		}
	}

	return false
}

func (r *RegexpHandler) Equal(n HanlderNode) bool {
	p2, ok := n.(*RegexpHandler)
	if !ok {
		return false
	}

	return r.name == p2.name
}

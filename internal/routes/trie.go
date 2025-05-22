package routes

import (
	"errors"
	"strings"

	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/internal/request"
)

// ErrNotHandled indicates no route matched
var ErrNotHandled = errors.New("not handled")

// Trie represents the routing trie supporting static, parameter, and wildcard segments
type Trie[T any] struct {
	root *node[T]
}

// node is a trie node
type node[T any] struct {
	segment    string
	children   []*node[T]
	isParam    bool         // :param segment
	isCatchAll bool         // *param segment
	handlers   map[string]T // method -> handler
}

// NewTrie initializes and returns an empty Tree
func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{root: &node[T]{handlers: make(map[string]T)}}
}

// Add registers a handler at the given pattern and method
// Pattern supports ":param" and final "*wildcard"
func (t *Trie[T]) Add(method, pattern string, handler T) {
	segments := split(pattern)
	cur := t.root
	for _, seg := range segments {
		var child *node[T]
		for _, c := range cur.children {
			if c.segment == seg {
				child = c
				break
			}
		}

		if child == nil {
			child = &node[T]{
				segment:    seg,
				isParam:    strings.HasPrefix(seg, ":"),
				isCatchAll: strings.HasPrefix(seg, "*"),
				handlers:   make(map[string]T),
			}
			cur.children = append(cur.children, child)
		}

		cur = child
		if cur.isCatchAll {
			// wildcard must be last
			break
		}
	}
	cur.handlers[method] = handler
}

// Get finds a handler for path and method, populates req.Params, or returns ErrNotHandled
func (t *Trie[T]) Get(req *request.Request, method, path string) (*T, error) {
	segments := split(path)
	params := make(map[string]string)
	h := t.root.search(segments, params, method)
	if h == nil {
		return nil, ErrNotHandled
	}

	request.SetParameters(req, placing.InPath, params)
	return h, nil
}

func (n *node[T]) search(segments []string, params map[string]string, method string) *T {
	if len(segments) == 0 {
		if handler, ok := n.handlers[method]; ok {
			return &handler
		}
		return nil
	}
	seg := segments[0]

	// 1. try exact match
	for _, c := range n.children {
		if c.segment == seg {
			if h := c.search(segments[1:], params, method); h != nil {
				return h
			}
		}
	}

	// 2. parameter match
	for _, c := range n.children {
		if !c.isParam {
			continue
		}

		key := strings.TrimPrefix(c.segment, ":")
		params[key] = seg
		if h := c.search(segments[1:], params, method); h != nil {
			return h
		}
	}

	// 3. wildcard match
	for _, c := range n.children {
		if !c.isCatchAll {
			continue
		}

		key := strings.TrimPrefix(c.segment, "*")
		params[key] = strings.Join(segments, "/")
		if handler, ok := c.handlers[method]; ok {
			return &handler
		}
	}

	return nil
}

// split trims and splits the path into segments
func split(path string) []string {
	clean := strings.Trim(path, "/")
	if clean == "" {
		return nil
	}

	return strings.Split(clean, "/")
}

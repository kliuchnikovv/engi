package tree

import (
	"errors"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/parameter/placing"
)

// ErrNotHandled indicates no route matched
var ErrNotHandled = errors.New("not handled")

// Tree represents the routing trie supporting static, parameter, and wildcard segments
type Tree[T any] struct {
	root *node[T]
}

// node is a trie node
type node[T any] struct {
	segment    string
	children   []*node[T]
	isParam    bool // :param segment
	isCatchAll bool // *param segment
	handler    *T
}

// NewTree initializes and returns an empty Tree
func NewTree[T any]() *Tree[T] {
	return &Tree[T]{root: &node[T]{}}
}

// Add registers a handler at the given pattern
// Pattern supports ":param" and final "*wildcard"
func (t *Tree[T]) Add(pattern string, handler T) {
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
			}
			cur.children = append(cur.children, child)
		}

		cur = child
		if cur.isCatchAll {
			// wildcard must be last
			break
		}
	}
	cur.handler = &handler
}

// Get finds a handler for path, populates req.Params, or returns ErrNotHandled
func (t *Tree[T]) Get(req *request.Request, path string) (*T, error) {
	segments := split(path)
	params := make(map[string]string)
	h := t.root.search(segments, params)
	if h == nil {
		return nil, ErrNotHandled
	}

	request.SetParameters(req, placing.InPath, params)
	return h, nil
}

func (n *node[T]) search(segments []string, params map[string]string) *T {
	if len(segments) == 0 {
		return n.handler
	}
	seg := segments[0]

	// 1. try exact match
	for _, c := range n.children {
		if c.segment == seg {
			return c.search(segments[1:], params)
		}
	}

	// 2. parameter match
	for _, c := range n.children {
		if !c.isParam {
			continue
		}

		key := strings.TrimPrefix(c.segment, ":")
		params[key] = seg
		return c.search(segments[1:], params)
	}

	// 3. wildcard match
	for _, c := range n.children {
		if !c.isCatchAll {
			continue
		}

		key := strings.TrimPrefix(c.segment, "*")
		params[key] = strings.Join(segments, "/")
		return c.handler
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

package tree

import "fmt"

type (
	Node struct {
		value    string
		children []*Node
	}
)

func NewTree(values ...string) Node {
	var tree = Node{
		value:    "",
		children: make([]*Node, 0),
	}

	if len(values) == 0 {
		return tree
	}

	tree.value = values[0]

	var currentNode = &tree

	for _, value := range values[1:] {
		var newNode = Node{
			value:    value,
			children: make([]*Node, 0),
		}

		currentNode.children = append(currentNode.children, &newNode)
		currentNode = &newNode
	}

	return tree
}

func (n *Node) Add(values ...string) *Node {
	if len(values) == 0 {
		return nil
	}

	var currentNode = n

	if n.value == values[0] {
		values = values[1:]
	}

	for _, value := range values {
		if i := currentNode.childWithValue(value); i != -1 {
			currentNode = currentNode.children[i]
			continue
		}

		var newNode = &Node{
			value:    value,
			children: make([]*Node, 0),
		}

		currentNode.children = append(currentNode.children, newNode)
		currentNode = newNode
	}

	return currentNode
}

func (n Node) childWithValue(value string) int {
	for i, child := range n.children {
		if child.value == value {
			return i
		}
	}

	return -1
}

func (n *Node) Optimize() {
	for _, child := range n.children {
		child.Optimize()
	}

	if len(n.children) == 1 {
		n.value = fmt.Sprintf("%s/%s", n.value, n.children[0].value)
		n.children = n.children[0].children
	}
}

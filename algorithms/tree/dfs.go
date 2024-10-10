package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/stacks/arraystack"
)

func DFS[V any](root *Node) []V {
	output := make([]V, 0)
	stack := arraystack.New()
	visited := hashset.New() 

	if root == nil {
		return output
	}

	stack.Push(root)
	for !stack.Empty() {
		// Get the current node
		elem, _ := stack.Pop()
		output = append(output, elem.(*Node).Value.(V))

		// Check if the node is visited already
		oldLen := visited.Size()
		visited.Add(elem)
		if (oldLen == visited.Size()) {
			continue
		}

		// Go through all children of the current node in reverse order
		for i := len(elem.(*Node).Children) - 1; i >= 0; i-- {
		    stack.Push(elem.(*Node).Children[i])
		}
	}

	return output 
}

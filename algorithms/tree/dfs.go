package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/stacks/arraystack"
)

func DFS[V any](root *Node) []V {
	output := make([]V, 0)

	if root == nil {
		return output
	}

	stack := arraystack.New()
	visited := hashset.New() 

	visited.Add(root)
	stack.Push(root)
	for !stack.Empty() {
		// Get the current node
		elem, ok := stack.Pop()
		if !ok {
			fmt.Println("failed to remove element form the stack")
		}
		output = append(output, elem.(*Node).Value.(V))

		// Check if the node is visited already
		oldLen := visited.Size()
		visited.Add(elem)
		if (oldLen == visited.Size()) {
			continue
		}

		// Go through all children of the current node
		for _, node := range elem.(*Node).Children {
			stack.Push(node)
		}
	}

	return output 
}

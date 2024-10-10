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
		node, _ := stack.Pop()
		output = append(output, node.(*Node).Value.(V))

		if visited.Contains(node) {
			continue
		}
		visited.Add(node)

		for i := len(node.(*Node).Children) - 1; i >= 0; i-- {
		    stack.Push(node.(*Node).Children[i])
		}
	}

	return output 
}

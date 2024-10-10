package main

import (
	"github.com/emirpasic/gods/queues/arrayqueue"
	"github.com/emirpasic/gods/sets/hashset"
)

func BFS[T any](root *Node) []T {
	output := make([]T, 0)
	queue := arrayqueue.New()
	visited := hashset.New()

	if root == nil {
		return output
	}

	queue.Enqueue(root)
	for !queue.Empty() {
		tmp, _ := queue.Dequeue()
		node := tmp.(*Node)

		if visited.Contains(node) {
			continue
		}
		visited.Add(node)

		output = append(output, node.Value.(T))

		for _, child := range node.Children {
			if !visited.Contains(child) {
				queue.Enqueue(child)
			}
		}
	}

	return output
}
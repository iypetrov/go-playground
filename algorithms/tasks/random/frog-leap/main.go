package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/stacks/arraystack"
)

func StartState(n int) string {
	leftFrogs := strings.Repeat(">", n)
	emptyField := "_"
	rightFrogs := strings.Repeat("<", n)

	return leftFrogs + emptyField + rightFrogs
}

func FinalState(n int) string {
	leftFrogs := strings.Repeat("<", n)
	emptyField := "_"
	rightFrogs := strings.Repeat(">", n)

	return leftFrogs + emptyField + rightFrogs
}

func GeneratePossiblePositions(current string) []string {
	var positions []string
	length := len(current)

	for i := 0; i < length; i++ {
		if current[i] == '_' {
			continue
		}

		if current[i] == '>' {
			if i < len(current)-1 && current[i+1] == '_' {
				tmp := []rune(current)
				tmp[i] = '_'
				tmp[i+1] = '>'
				positions = append(positions, string(tmp))
			}
			if i < len(current)-2 && current[i+2] == '_' {
				tmp := []rune(current)
				tmp[i] = '_'
				tmp[i+2] = '>'
				positions = append(positions, string(tmp))
			}
		}

		if current[i] == '<' {
			if i > 0 && current[i-1] == '_' {
				tmp := []rune(current)
				tmp[i] = '_'
				tmp[i-1] = '<'
				positions = append(positions, string(tmp))
			}
			if i > 1 && current[i-2] == '_' {
				tmp := []rune(current)
				tmp[i] = '_'
				tmp[i-2] = '<'
				positions = append(positions, string(tmp))
			}
		}
	}

	return positions
}

type Node struct {
	Value    string
	Parent   *Node
	Children []*Node
}

func GenerateChildren(node *Node) {
	if node == nil {
		return
	}

	vals := GeneratePossiblePositions(node.Value)
	children := make([]*Node, 0)
	for _, val := range vals {
		tmp := Node{
			Value:    val,
			Parent:   node,
			Children: []*Node{},
		}
		children = append(children, &tmp)
	}

	node.Children = children
}

func FrogLeap(n int) []string {
	root := Node{
		Value:    StartState(n),
		Parent:   nil,
		Children: []*Node{},
	}

	return DFS(&root, n)
}

func DFS(root *Node, n int) []string {
	stack := arraystack.New()
	visited := hashset.New()
	var fNode *Node
	result := []string{}

	if root == nil {
		return result
	}

	stack.Push(root)
	for !stack.Empty() {
		tmp, _ := stack.Pop()
		node := tmp.(*Node)
		if node.Value == FinalState(n) {
			fNode = node
			break
		}

		if visited.Contains(node) {
			continue
		}
		visited.Add(node)

		GenerateChildren(node)
		for i := len(node.Children) - 1; i >= 0; i-- {
			child := node.Children[i]
			if !visited.Contains(child) {
				stack.Push(child)
			}
		}
	}

	if fNode == nil {
		return result
	}

	for fNode != nil {
		result = append(result, fNode.Value)
		fNode = fNode.Parent
	}

	slices.Reverse(result)
	return result
}

func main() {
	fmt.Println("hello frog-leap")
}

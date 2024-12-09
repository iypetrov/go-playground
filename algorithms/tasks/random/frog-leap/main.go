package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/emirpasic/gods/stacks/arraystack"
)

func swap(frogs string, i, j int) string {
	runes := []rune(frogs)
	runes[i], runes[j] = runes[j], runes[i]
	return string(runes)
}

// sometimes tests might fail because it can start with both left and right side
func dfs(frogs, goal string, visited map[string]bool, stack *arraystack.Stack) bool {
	if frogs == goal {
		return true
	}

	blank := strings.Index(frogs, "_")

	children := make(map[string]bool)

	//  1 left
	if blank >= 1 && frogs[blank-1] == '>' {
		children[swap(frogs, blank, blank-1)] = true
	}
	// 1 right
	if blank <= len(frogs)-2 && frogs[blank+1] == '<' {
		children[swap(frogs, blank, blank+1)] = true
	}
	// 2 left
	if blank >= 2 && frogs[blank-2] == '>' {
		children[swap(frogs, blank, blank-2)] = true
	}
	// 2 right 
	if blank <= len(frogs)-3 && frogs[blank+2] == '<' {
		children[swap(frogs, blank, blank+2)] = true
	}

	for child := range children {
		if !visited[child] {
			visited[child] = true
			if dfs(child, goal, visited, stack) {
				stack.Push(child)
				return true
			}
		}
	}
	return false
}

func FrogLeap(n int) []string {
	init := strings.Repeat(">", n) + "_" + strings.Repeat("<", n)
	goal := strings.Repeat("<", n) + "_" + strings.Repeat(">", n)

	visited := make(map[string]bool)
	visited[init] = true

	stack := arraystack.New()

	start := time.Now()
	if !dfs(init, goal, visited, stack) {
		fmt.Println("no solution found")
		return []string{}
	}
	stack.Push(init)

	duration := time.Since(start).Milliseconds()

	steps := []string{}
	for !stack.Empty() {
		step, _ := stack.Pop()
		steps = append(steps, step.(string))
	}

	fmt.Printf("solved in %dms\n", duration)

	return steps
}

func main() {
	fmt.Println("hello frog-leap")
}

package main

import "fmt"

type Node struct {
	val      any
	parent   *Node
	children []*Node
}

func main() {
	fmt.Println("hello tree")
}

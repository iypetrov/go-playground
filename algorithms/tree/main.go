package main

import "fmt"

type Node struct {
	Value    any
	Children []*Node
}

func main() {
	fmt.Println("hello tree")
}

package main

import (
	"fmt"
	"os"
	"strconv"
)

func divide(a, b int) (int, error) {
	defer func() {
        if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
        }
    }()
	result := a / b
	return result, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <int1> <int2>")
		return
	}

	a, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("First argument must be an integer")
		return
	}

	b, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Second argument must be an integer")
		return
	}

	result, err := divide(a, b)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Result: %d\n", result)
}

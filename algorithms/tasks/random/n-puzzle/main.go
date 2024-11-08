package main

import (
	"fmt"
	"math"
)

func NPuzzle(n int, arr []int) []string {
	if n < 3 {
		panic("n is not valid")
	}
	k := int(math.Sqrt(float64(n + 1)))
	if !(k > 1 && k*k-1 == n) {
		panic("n is not valid")
	}

	

	return []string{}
}

func main() {
	fmt.Println("hello n-puzzle")
}

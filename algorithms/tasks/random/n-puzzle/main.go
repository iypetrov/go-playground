package main

import (
	"fmt"
)

func FinalState(numberSlots int, indexZero int) []int {
	index := 0
	result := make([]int, numberSlots+1)
	for i := 0; i <= numberSlots; i++ {
		if i == indexZero {
			result[i] = 0
			continue
		}
		index++
		result[i] = index
	}

	return result
}

func NPuzzle(numberSlots int, indexZero int, positions []int) (int, []string) {
	if indexZero == -1 {
		indexZero = numberSlots
	}

	return 0, []string{}
}

func main() {
	fmt.Println("hello n-puzzle")
}

package main

import (
	"fmt"
)

func FinalPosition(numberSlots int, indexZero int) []int {
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

func EvaluatePosition(currentPosition []int, desiredPosition []int) int {
	if len(currentPosition) != len(desiredPosition) {
		return -1
	}

	score := 0
	for i, item := range currentPosition {
		if item == desiredPosition[i] {
			score++
		}
	}

	return score
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

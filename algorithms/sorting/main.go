package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

// worst: O(N^2)
// avg: O(N^2)
// best: O(N^2)
// space: O(1)
func BubbleSort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

// worst: O(N^2)
// avg: O(N^2)
// best: O(N^2)
// space: O(1)
func SelectionSort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		indexMin := i
		for j := i; j < n-1; j++ {
			if arr[j] < arr[i] {
				indexMin = j
			}
		}

		arr[i], arr[indexMin] = arr[indexMin], arr[i]
	}
	return arr
}

func main() {
	fmt.Println("hello sorting")
}

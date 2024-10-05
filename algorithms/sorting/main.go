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
	sortedArr := make([]int, len(arr))
	copy(sortedArr, arr)

	n := len(sortedArr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if sortedArr[j] > sortedArr[j+1] {
				sortedArr[j], sortedArr[j+1] = sortedArr[j+1], sortedArr[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
	copy(arr, sortedArr)
	return sortedArr
}

// QuickSort sorts a slice of integers using the QuickSort algorithm.
func QuickSort(arr []int) []int {
	sortedArr := make([]int, len(arr))
	copy(sortedArr, arr)
	quickSortHelper(sortedArr, 0, len(sortedArr)-1)
	copy(arr, sortedArr)
	return sortedArr
}

func quickSortHelper(arr []int, low, high int) {
	if low < high {
		p := partition(arr, low, high)
		quickSortHelper(arr, low, p-1)
		quickSortHelper(arr, p+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return i
}

func main() {
	fmt.Println("hello sorting")
}

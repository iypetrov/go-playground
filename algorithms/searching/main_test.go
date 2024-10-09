package main

import (
	"testing"
)

func TestSearchAlgorithms(t *testing.T) {
	tests := []struct {
		input    []int
		target   int
		expected int // expected index of the target, -1 if not found
	}{
		{input: []int{}, target: 5, expected: -1},
		{input: []int{1}, target: 1, expected: 0},
		{input: []int{1}, target: 2, expected: -1},
		{input: []int{1, 2}, target: 2, expected: 1},
		{input: []int{1, 2}, target: 3, expected: -1},
		{input: []int{1, 2, 4, 5, 8}, target: 4, expected: 2},
		{input: []int{1, 2, 4, 5, 8}, target: 7, expected: -1},
		{input: []int{11, 12, 22, 25, 34, 64, 90}, target: 22, expected: 2},
		{input: []int{11, 12, 22, 25, 34, 64, 90}, target: 100, expected: -1},
	}

	algorithms := []struct {
		name  string
		search func([]int, int) int
	}{
		{name: "LinearSearch", search: LinearSearch},
		{name: "BinarySearch", search: BinarySearch},
	}

	for _, algo := range algorithms {
		for _, test := range tests {
			result := algo.search(test.input, test.target)
			if result != test.expected {
				t.Errorf("[%s] searching for %d in %v: expected %d, got %d", algo.name, test.target, test.input, test.expected, result)
			}
		}
	}
}

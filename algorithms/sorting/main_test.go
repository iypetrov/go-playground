package main

import (
	"reflect"
	"testing"
)

func TestSortingAlgorithms(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{input: []int{}, expected: []int{}},
		{input: []int{1}, expected: []int{1}},
		{input: []int{2, 1}, expected: []int{1, 2}},
		{input: []int{5, 1, 4, 2, 8}, expected: []int{1, 2, 4, 5, 8}},
		{input: []int{5, 1, 1, 2, 0, 0}, expected: []int{0, 0, 1, 1, 2, 5}},
		{input: []int{64, 34, 25, 12, 22, 11, 90}, expected: []int{11, 12, 22, 25, 34, 64, 90}},
	}

	algorithms := []struct {
		name string
		sort func([]int) []int
	}{
		{name: "BubbleSort", sort: BubbleSort},
		{name: "SelectionSort", sort: SelectionSort},
	}

	for _, algo := range algorithms {
		for _, test := range tests {
			original := make([]int, len(test.input))
			copy(original, test.input)

			sorted := algo.sort(test.input)
			if !reflect.DeepEqual(sorted, test.expected) {
				t.Errorf("[%s] expected sorted %v, got %v", algo.name, test.expected, sorted)
			}

			if !reflect.DeepEqual(test.input, test.expected) {
				t.Errorf("[%s] expected original %v to be modified to %v, got %v", algo.name, original, test.expected, test.input)
			}
		}
	}
}

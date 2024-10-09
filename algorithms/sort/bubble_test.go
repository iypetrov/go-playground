package main

import (
	"reflect"
	"testing"
)

func TestBubble(t *testing.T) {
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

	name := "Bubble"
	algo := Bubble[int]
	for _, test := range tests {
		result := algo(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("[%s] expected %v, got %v", name, test.expected, result)
		}
	}
}

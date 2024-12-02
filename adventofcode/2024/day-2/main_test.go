package main

import (
	"reflect"
	"testing"
)

func TestFinalState(t *testing.T) {
	tests := []struct {
		input    []int
		expected bool 
	}{
		{
			input:    []int{1, 1, 2},
			expected: false,
		},
		{
			input:    []int{1, 2, 2},
			expected: false,
		},
		{
			input:    []int{1, 2, 6},
			expected: false,
		},
		{
			input:    []int{1, 2, 3},
			expected: true,
		},
		{
			input:    []int{1, 2, 5},
			expected: true,
		},
		{
			input:    []int{2, 2, 1},
			expected: false,
		},
		{
			input:    []int{2, 1, 1},
			expected: false,
		},
		{
			input:    []int{6, 2, 1},
			expected: false,
		},
		{
			input:    []int{3, 2, 1},
			expected: true,
		},
		{
			input:    []int{5, 2, 1},
			expected: true,
		},
		{
			input:    []int{1, 2, 1},
			expected: false,
		},
	}

	name := "IsReportSafe"
	algo := IsReportSafe 
	for _, test := range tests {
		result := algo(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("[%s] expected %v, got %v", name, test.input, result)
		}
	}
}
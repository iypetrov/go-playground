package main

import (
	"reflect"
	"testing"
)

func TestFinalState(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{
			input:    0,
			expected: "_",
		},
		{
			input:    1,
			expected: "<_>",
		},
		{
			input:    2,
			expected: "<<_>>",
		},
		{
			input:    3,
			expected: "<<<_>>>",
		},
		{
			input:    4,
			expected: "<<<<_>>>>",
		},
	}

	name := "FinalState"
	algo := FinalState
	for _, test := range tests {
		result := algo(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("[%s] expected %v, got %v", name, test.input, result)
		}
	}
}

func TestGeneratePossiblePositions(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: ">_<",
			expected: []string{
				"><_",
				"_><",
			},
		},
		{
			input: ">>_<<",
			expected: []string{
				">_><<",
				">><_<",
				"_>><<",
				">><<_",
			},
		},
		{
			input: "_>><<",
			expected: []string{},
		},
	}

	name := "GeneratePossiblePositions"
	algo := GeneratePossiblePositions
	for _, test := range tests {
		result := algo(test.input)

		if len(result) != len(test.expected) {
			t.Errorf("[%s] expected length %d, got %d", name, len(test.expected), len(result))
			continue
		}

		expectedCount := make(map[string]int)
		for _, v := range test.expected {
			expectedCount[v]++
		}

		actualCount := make(map[string]int)
		for _, v := range result {
			actualCount[v]++
		}

		if !reflect.DeepEqual(expectedCount, actualCount) {
			t.Errorf("[%s] expected %v, got %v", name, test.expected, result)
		}
	}
}

func TestFrogLeap(t *testing.T) {
	tests := []struct {
		input    int
		expected []string
	}{
		{
			input: 0,
			expected: []string{
				"_",
			},
		},
		// {
		// 	input: 1,
		// 	expected: []string{
		// 		">_<",
		// 		"_><",
		// 		"<>_",
		// 		"<_>",
		// 	},
		// },
		{
			input: 2,
			expected: []string{
				">>_<<",
				">_><<",
				"><>_<",
				"><><_",
				"><_<>",
				"_<><>",
				"<_><>",
				"<<>_>",
				"<<_>>",
			},
		},
	}

	name := "FrogLeap"
	algo := FrogLeap
	for _, test := range tests {
		result := algo(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("[%s] expected %v, got %v", name, test.expected, result)
		}
	}
}

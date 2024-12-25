package main

import (
	"reflect"
	"testing"
)

func TestNPuzzle(t *testing.T) {
	tests := []struct {
		numberSlots int
		indexZero   int
		positions   []int
		bestPathLen int
		steps       []string
	}{
		{
			numberSlots: 8,
			indexZero:   0,
			positions:   []int{1, 2, 0, 3, 4, 5, 6, 7, 8},
			bestPathLen: 2,
			steps:       []string{"right", "right"},
		},
		{
			numberSlots: 8,
			indexZero:   4,
			positions:   []int{1, 0, 3, 4, 2, 5, 6, 7, 8},
			bestPathLen: 1,
			steps:       []string{"up"},
		},
		{
			numberSlots: 8,
			indexZero:   8,
			positions:   []int{1, 2, 3, 4, 5, 6, 0, 7, 8},
			bestPathLen: 2,
			steps:       []string{"left", "left"},
		},
		{
			numberSlots: 8,
			indexZero:   -1,
			positions:   []int{1, 2, 3, 4, 5, 6, 0, 7, 8},
			bestPathLen: 2,
			steps:       []string{"left", "left"},
		},
		{
			numberSlots: 15,
			indexZero:   0,
			positions:   []int{1, 2, 3, 0, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			bestPathLen: 3,
			steps:       []string{"right", "right", "right"},
		},
		{
			numberSlots: 15,
			indexZero:   15,
			positions:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0, 13, 14, 15},
			bestPathLen: 3,
			steps:       []string{"left", "left", "left"},
		},
		{
			numberSlots: 15,
			indexZero:   -1,
			positions:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0, 13, 14, 15},
			bestPathLen: 3,
			steps:       []string{"left", "left", "left"},
		},
	}

	for i, test := range tests {
		result, _, _ := NPuzzle(test.numberSlots, test.indexZero, test.positions)
		if !reflect.DeepEqual(result, test.steps) {
			t.Errorf("[NPuzzle %d] expected %v, got %v", i, test.steps, result)
		}
	}
}

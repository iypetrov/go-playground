package main

import (
	"reflect"
	"testing"
)

func TestNPuzzle(t *testing.T) {
	tests := []struct {
		number_slots  int
		index_zero    int
		positions     []int
		best_path_len int
		steps         []string
	}{
		{
			number_slots:  8,
			index_zero:    0,
			positions:     []int{1, 2, 0, 3, 4, 5, 6, 7, 8},
			best_path_len: 2,
			steps:         []string{"right", "right"},
		},
		{
			number_slots:  8,
			index_zero:    4,
			positions:     []int{1, 0, 3, 4, 2, 5, 6, 7, 8},
			best_path_len: 1,
			steps:         []string{"up"},
		},
		{
			number_slots:  8,
			index_zero:    8,
			positions:     []int{1, 2, 3, 4, 5, 6, 0, 7, 8},
			best_path_len: 2,
			steps:         []string{"left", "left"},
		},
		{
			number_slots:  8,
			index_zero:    -1,
			positions:     []int{1, 2, 3, 4, 5, 6, 0, 7, 8},
			best_path_len: 2,
			steps:         []string{"left", "left"},
		},
		{
			number_slots:  15,
			index_zero:    0,
			positions:     []int{1, 2, 3, 0, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			best_path_len: 3,
			steps:         []string{"right", "right", "right"},
		},
		{
			number_slots:  15,
			index_zero:    15,
			positions:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0, 13, 14, 15},
			best_path_len: 3,
			steps:         []string{"left", "left", "left"},
		},
		{
			number_slots:  15,
			index_zero:    -1,
			positions:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0, 13, 14, 15},
			best_path_len: 3,
			steps:         []string{"left", "left", "left"},
		},
	}

	name := "NPuzzle"
	algo := NPuzzle
	for _, test := range tests {
		len, result := algo(test.number_slots, test.index_zero, test.positions)
		if !reflect.DeepEqual(result, test.steps) {
			t.Errorf("[%s] expected %v, got %v", name, test.best_path_len, len)
			t.Errorf("[%s] expected %v, got %v", name, test.steps, result)
		}
	}
}

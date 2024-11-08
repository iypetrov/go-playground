package main

import (
	"reflect"
	"testing"
)

func TestNPuzzle(t *testing.T) {
	tests := []struct {
		number_slots int
		positions    []int
		steps        []string
	}{
		{
			number_slots: 8,
			positions:    []int{1, 2, 3, 4, 5, 6, 7, 8, -1},
			steps:        []string{"right", "right", "down", "down", "left", "left", "up", "up"},
		},
	}

	name := "NPuzzle"
	algo := NPuzzle
	for _, test := range tests {
		result := algo(test.number_slots, test.positions)
		if !reflect.DeepEqual(result, test.steps) {
			t.Errorf("[%s] expected %v, got %v", name, test.steps, result)
		}
	}
}

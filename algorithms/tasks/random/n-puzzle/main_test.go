package main

import (
	"reflect"
	"testing"
)

func TestEvaluatePosition (t *testing.T) {
	tests := []struct {
		currentPosition []int
		desiredPosition []int
		score int
	} {
		{
			currentPosition: []int{8, 7, 6, 5, 0, 1, 2, 3, 4},
			desiredPosition: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			score: 0,
		},
		{
			currentPosition: []int{0, 1, 2, 3, 4, 8, 7, 6, 5},
			desiredPosition: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			score: 5,
		},
		{
			currentPosition: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			desiredPosition: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			score: 9,
		},
	}

	name := "EvaluatePosition"
	algo := EvaluatePosition
	for _, test := range tests {
		result := algo(test.currentPosition, test.desiredPosition)
		if !reflect.DeepEqual(result, test.score) {
			t.Errorf("[%s] expected %v, got %v", name, test.score, result)
		}
	}
}

func TestFinalState(t *testing.T) {
	tests := []struct {
		numberSlots int
		indexZero   int
		positions   []int
	}{
		{
			numberSlots: 8,
			indexZero:   0,
			positions:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			numberSlots: 8,
			indexZero:   4,
			positions:   []int{1, 2, 3, 4, 0, 5, 6, 7, 8},
		},
		{
			numberSlots: 8,
			indexZero:   8,
			positions:   []int{1, 2, 3, 4, 5, 6, 7, 8, 0},
		},
	}

	name := "FinalPosition"
	algo := FinalPosition
	for _, test := range tests {
		result := algo(test.numberSlots, test.indexZero)
		if !reflect.DeepEqual(result, test.positions) {
			t.Errorf("[%s] expected %v, got %v", name, test.positions, result)
		}
	}
}

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

	name := "NPuzzle"
	algo := NPuzzle
	for _, test := range tests {
		len, result := algo(test.numberSlots, test.indexZero, test.positions)
		if !reflect.DeepEqual(result, test.steps) {
			t.Errorf("[%s] expected %v, got %v", name, test.bestPathLen, len)
			t.Errorf("[%s] expected %v, got %v", name, test.steps, result)
		}
	}
}

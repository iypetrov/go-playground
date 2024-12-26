package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBestMove(t *testing.T) {
	tests := []struct {
		board    [3][3]string
		expected [3][3]string
	}{
		{
			board: [3][3]string{
				{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
				{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
				{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
			},
			expected: [3][3]string{
				{X, EMPTY_CELL, EMPTY_CELL},
				{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
				{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
			},
		},
		{
			board: [3][3]string{
				{X, EMPTY_CELL, EMPTY_CELL},
				{EMPTY_CELL, O, EMPTY_CELL},
				{EMPTY_CELL, O, EMPTY_CELL},
			},
			expected: [3][3]string{
				{X, X, EMPTY_CELL},
				{EMPTY_CELL, O, EMPTY_CELL},
				{EMPTY_CELL, O, EMPTY_CELL},
			},
		},
		{
			board: [3][3]string{
				{X, O, X},
				{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
				{EMPTY_CELL, EMPTY_CELL, O},
			},
			expected: [3][3]string{
				{X, O, X},
				{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
				{X, EMPTY_CELL, O},
			},
		},
	}

	for _, test := range tests {
		game := NewGame(false)
		game.SetBotPrepState(test.board)
		game.Print()
		game.BestMove()
		fmt.Println("to")
		game.Print()
		fmt.Println("--------------------------")
		if !reflect.DeepEqual(game.Board, test.expected) {
			t.Errorf("Expected %v, but got %v", test.expected, game.Board)
		}
	}
}

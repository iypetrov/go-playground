package main

import (
	"fmt"
	"testing"
)

func isValidNQueensSolution(board [][]rune, N int) bool {
	rows := make([]bool, N)
	cols := make([]bool, N)
	diag1 := make([]bool, 2*N-1) 
	diag2 := make([]bool, 2*N-1) 

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if board[i][j] == '*' {
				if rows[i] || cols[j] || diag1[N-1+i-j] || diag2[i+j] {
					return false 
				}
				rows[i] = true
				cols[j] = true
				diag1[N-1+i-j] = true
				diag2[i+j] = true
			}
		}
	}
	return true
}

func TestNQueens(t *testing.T) {
	tests := []struct {
		N        int
		expected bool 
	}{
		{N: 4, expected: true},
		{N: 8, expected: true},
		{N: 1, expected: true},
		{N: 2, expected: false},
		{N: 3, expected: false}, 
		{N: 10, expected: true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("N=%d", test.N), func(t *testing.T) {
			board := NQueens(test.N)
			if test.expected {
				if board == nil {
					t.Errorf("Expected a solution for N=%d, but got none", test.N)
				} else if !isValidNQueensSolution(board, test.N) {
					t.Errorf("Invalid solution for N=%d", test.N)
				}
			} else {
				if board != nil {
					t.Errorf("Expected no solution for N=%d, but got a solution", test.N)
				}
			}
		})
	}
}

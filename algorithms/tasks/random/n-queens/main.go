package main

import (
	"fmt"
	"math/rand"
	"time"
)

const IterationLimit = 100000

type NQueensSolver struct {
	N                   int
	queenPositions      []int
	rowConflicts        []int
	diagonal1Conflicts  []int
	diagonal2Conflicts  []int
	queenConflicts      []int
}

func New(N int) *NQueensSolver {
	s := NQueensSolver{
		N:                   N,
		queenPositions:      make([]int, N),
		rowConflicts:        make([]int, N),
		diagonal1Conflicts:  make([]int, 2*N-1),
		diagonal2Conflicts:  make([]int, 2*N-1),
		queenConflicts:      make([]int, N),
	}

	row := 0
	for col := 0; col < s.N; col++ {
		s.queenPositions[col] = row
		row = (row + 2) % s.N
	}
	for col, row := range s.queenPositions {
		s.rowConflicts[row]++
		s.diagonal1Conflicts[s.N-1+col-row]++
		s.diagonal2Conflicts[col+row]++
	}

	return &s
}

func (s *NQueensSolver) updateConflicts() int {
	totalConflicts := 0
	for col := 0; col < s.N; col++ {
		s.queenConflicts[col] = s.getConflictsForQueen(col)
		totalConflicts += s.queenConflicts[col]
	}
	return totalConflicts
}

func (s *NQueensSolver) moveQueen(col, newRow int) {
	oldRow := s.queenPositions[col]
	s.rowConflicts[oldRow]--
	s.rowConflicts[newRow]++
	s.diagonal1Conflicts[s.N-1+col-oldRow]--
	s.diagonal2Conflicts[col+oldRow]--
	s.diagonal1Conflicts[s.N-1+col-newRow]++
	s.diagonal2Conflicts[col+newRow]++
	s.queenPositions[col] = newRow
}

func (s *NQueensSolver) getRowWithLeastConflicts(col int) int {
	conflicts := make([]int, s.N)
	minConflicts := s.N + 1

	for row := 0; row < s.N; row++ {
		conflicts[row] = s.rowConflicts[row] +
			s.diagonal1Conflicts[s.N-1+col-row] +
			s.diagonal2Conflicts[col+row]
		if conflicts[row] < minConflicts {
			minConflicts = conflicts[row]
		}
	}

	candidates := []int{}
	for row, c := range conflicts {
		if c == minConflicts {
			candidates = append(candidates, row)
		}
	}
	return candidates[randomInt(0, len(candidates)-1)]
}

func (s *NQueensSolver) getConflictsForQueen(col int) int {
	row := s.queenPositions[col]
	return s.rowConflicts[row] +
		s.diagonal1Conflicts[s.N-1+col-row] +
		s.diagonal2Conflicts[col+row] - 3
}

func (s *NQueensSolver) getColumnWithMostConflicts() int {
	maxConflicts := -1
	candidates := []int{}

	for col, c := range s.queenConflicts {
		if c > maxConflicts {
			maxConflicts = c
			candidates = []int{col}
		} else if c == maxConflicts {
			candidates = append(candidates, col)
		}
	}
	return candidates[randomInt(0, len(candidates)-1)]
}

func (s *NQueensSolver) buildBoard() [][]rune {
	board := make([][]rune, s.N)
	for i := range board {
		board[i] = make([]rune, s.N)
		for j := range board[i] {
			board[i][j] = '_'
		}
		board[i][s.queenPositions[i]] = '*'
	}
	return board
}

func NQueens(N int) [][]rune {
	solver := New(N)

	startTime := time.Now()
	totalConflicts := solver.updateConflicts()
	for i := 0; i < IterationLimit; i++ {
		if totalConflicts == 0 {
			return solver.buildBoard()
		}
		col := solver.getColumnWithMostConflicts()
		newRow := solver.getRowWithLeastConflicts(col)
		solver.moveQueen(col, newRow)
		totalConflicts = solver.updateConflicts()
	}
	duration := time.Since(startTime).Seconds()
	fmt.Printf("Solved in %.2f seconds\n", duration)

	return nil 
}

func randomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func main() {
	fmt.Println("hello n-queens")
}

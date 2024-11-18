package main

import (
	"fmt"
	"math"
)

// State represents the n-puzzle state
type State struct {
	board      []int
	zeroIndex  int
	cost       int // g: cost to reach this state
	heuristic  int // h: estimated cost to goal
	parentMove string
}

// ManhattanDistance calculates the heuristic for a given board
func ManhattanDistance(board []int, size int) int {
	distance := 0
	for i, tile := range board {
		if tile != 0 {
			targetRow := (tile - 1) / size
			targetCol := (tile - 1) % size
			currentRow := i / size
			currentCol := i % size
			distance += int(math.Abs(float64(targetRow-currentRow)) + math.Abs(float64(targetCol-currentCol)))
		}
	}
	return distance
}

// IsGoal checks if the current board matches the goal state
func IsGoal(board []int) bool {
	for i := 0; i < len(board)-1; i++ {
		if board[i] != i+1 {
			return false
		}
	}
	return board[len(board)-1] == 0
}

// GenerateSuccessors returns valid successor states
func GenerateSuccessors(state State, size int) []State {
	successors := []State{}
	directions := []struct {
		dx, dy   int
		moveName string
	}{
		{-1, 0, "up"}, {1, 0, "down"}, {0, -1, "left"}, {0, 1, "right"},
	}

	for _, dir := range directions {
		newRow := (state.zeroIndex / size) + dir.dx
		newCol := (state.zeroIndex % size) + dir.dy
		if newRow >= 0 && newRow < size && newCol >= 0 && newCol < size {
			newZeroIndex := newRow*size + newCol
			newBoard := make([]int, len(state.board))
			copy(newBoard, state.board)
			newBoard[state.zeroIndex], newBoard[newZeroIndex] = newBoard[newZeroIndex], newBoard[state.zeroIndex]
			successors = append(successors, State{
				board:      newBoard,
				zeroIndex:  newZeroIndex,
				cost:       state.cost + 1,
				heuristic:  ManhattanDistance(newBoard, size),
				parentMove: dir.moveName,
			})
		}
	}
	return successors
}

// IDAStar implements the IDA* search algorithm
func IDAStar(initial State, size int) (int, []string) {
	bound := initial.heuristic
	path := []State{initial}

	for {
		t, solution := search(path, 0, bound, size)
		if t == -1 {
			return len(solution), solution
		}
		if t == math.MaxInt {
			return -1, nil // No solution found
		}
		bound = t
	}
}

// search is a helper for the IDA* algorithm
func search(path []State, g, bound, size int) (int, []string) {
	node := path[len(path)-1]
	f := g + node.heuristic
	if f > bound {
		return f, nil
	}
	if IsGoal(node.board) {
		moves := []string{}
		for _, s := range path[1:] {
			moves = append(moves, s.parentMove)
		}
		return -1, moves
	}
	min := math.MaxInt
	for _, succ := range GenerateSuccessors(node, size) {
		exists := false
		for _, p := range path {
			if isEqual(p.board, succ.board) {
				exists = true
				break
			}
		}
		if !exists {
			path = append(path, succ)
			t, solution := search(path, g+1, bound, size)
			if t == -1 {
				return -1, solution
			}
			if t < min {
				min = t
			}
			path = path[:len(path)-1]
		}
	}
	return min, nil
}

// isEqual checks if two boards are identical
func isEqual(board1, board2 []int) bool {
	for i := range board1 {
		if board1[i] != board2[i] {
			return false
		}
	}
	return true
}

func main() {
	// Example puzzle: 8-puzzle
	initialBoard := []int{1, 2, 3, 4, 5, 6, 0, 7, 8}
	size := 3 // 3x3 board
	zeroIndex := 6

	initialState := State{
		board:     initialBoard,
		zeroIndex: zeroIndex,
		cost:      0,
		heuristic: ManhattanDistance(initialBoard, size),
	}

	fmt.Println("Solving n-puzzle with IDA*...")
	bestPathLen, solution := IDAStar(initialState, size)
	fmt.Printf("Solution found with %d moves: %v\n", bestPathLen, solution)
}

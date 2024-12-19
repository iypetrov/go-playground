package main

import (
	"fmt"
	"math"

	"github.com/deckarep/golang-set/v2"
)

type Tile int
type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Board struct {
	tiles          []Tile
	dimension      int
	emptyTileIndex int
}

func (b *Board) Equals(other Board) bool {
	if len(b.tiles) != len(other.tiles) {
		return false
	}
	for i := range b.tiles {
		if b.tiles[i] != other.tiles[i] {
			return false
		}
	}
	return true
}

type ManhattanState struct {
	board     Board
	direction Direction
	g         int // cost from the start state
	h         int // heuristic cost (Manhattan distance)
}

var goal Board
var visited = mapset.NewSet[Board]()

// Create a goal board based on the size and empty tile index
func goalBoard(size int, emptyTileIndex int) Board {
	board := Board{dimension: int(math.Sqrt(float64(size + 1)))}
	board.tiles = make([]Tile, size+1)
	board.emptyTileIndex = emptyTileIndex
	i := 0
	for ; i < emptyTileIndex; i++ {
		board.tiles[i] = Tile(i + 1)
	}
	for ; i <= size; i++ {
		board.tiles[i] = Tile(i)
	}
	board.tiles[emptyTileIndex] = 0
	return board
}

// Manhattan distance heuristic
func (b *Board) manhattan(goal Board) int {
	currentPositions := make(map[Tile][2]int)
	goalPositions := make(map[Tile][2]int)

	for i, tile := range b.tiles {
		currentPositions[tile] = [2]int{i / b.dimension, i % b.dimension}
		goalPositions[goal.tiles[i]] = [2]int{i / goal.dimension, i % goal.dimension}
	}

	dist := 0
	for tile, pos := range currentPositions {
		goalPos := goalPositions[tile]
		dist += int(math.Abs(float64(pos[0]-goalPos[0]))) + int(math.Abs(float64(pos[1]-goalPos[1])))
	}
	return dist
}

// Get the neighbors of a given board state
func (b *Board) getNeighbors() []ManhattanState {
	var neighbors []ManhattanState
	// Get potential moves for left, right, up, down
	emptyTileY := b.emptyTileIndex % b.dimension
	emptyTileX := b.emptyTileIndex / b.dimension

	// Left
	if emptyTileY > 0 {
		newBoard := *b
		newBoard.tiles[b.emptyTileIndex-1], newBoard.tiles[b.emptyTileIndex] = newBoard.tiles[b.emptyTileIndex], newBoard.tiles[b.emptyTileIndex-1]
		newBoard.emptyTileIndex = b.emptyTileIndex - 1
		neighbors = append(neighbors, ManhattanState{board: newBoard, direction: Left, g: 1, h: newBoard.manhattan(goal)})
	}

	// Right
	if emptyTileY < b.dimension-1 {
		newBoard := *b
		newBoard.tiles[b.emptyTileIndex+1], newBoard.tiles[b.emptyTileIndex] = newBoard.tiles[b.emptyTileIndex], newBoard.tiles[b.emptyTileIndex+1]
		newBoard.emptyTileIndex = b.emptyTileIndex + 1
		neighbors = append(neighbors, ManhattanState{board: newBoard, direction: Right, g: 1, h: newBoard.manhattan(goal)})
	}

	// Up
	if emptyTileX > 0 {
		newBoard := *b
		newBoard.tiles[b.emptyTileIndex-b.dimension], newBoard.tiles[b.emptyTileIndex] = newBoard.tiles[b.emptyTileIndex], newBoard.tiles[b.emptyTileIndex-b.dimension]
		newBoard.emptyTileIndex = b.emptyTileIndex - b.dimension
		neighbors = append(neighbors, ManhattanState{board: newBoard, direction: Up, g: 1, h: newBoard.manhattan(goal)})
	}

	// Down
	if emptyTileX < b.dimension-1 {
		newBoard := *b
		newBoard.tiles[b.emptyTileIndex+b.dimension], newBoard.tiles[b.emptyTileIndex] = newBoard.tiles[b.emptyTileIndex], newBoard.tiles[b.emptyTileIndex+b.dimension]
		newBoard.emptyTileIndex = b.emptyTileIndex + b.dimension
		neighbors = append(neighbors, ManhattanState{board: newBoard, direction: Down, g: 1, h: newBoard.manhattan(goal)})
	}

	return neighbors
}

// Check if the board is solvable
func (b *Board) isSolvable() bool {
	inversions := 0
	for i := 0; i < len(b.tiles)-1; i++ {
		for j := i + 1; j < len(b.tiles); j++ {
			if b.tiles[i] > 0 && b.tiles[j] > 0 && b.tiles[i] > b.tiles[j] {
				inversions++
			}
		}
	}
	emptyRow := b.emptyTileIndex / b.dimension
	if b.dimension%2 == 0 {
		return (inversions+emptyRow)%2 != 0
	}
	return inversions%2 == 0
}

// IDA* search algorithm
func search(root ManhattanState, bound int) (path []ManhattanState, result int) {
	f := root.g + root.h
	if f > bound {
		return nil, f
	}
	if root.board.Equals(goal) {
		return []ManhattanState{root}, 0
	}

	min := math.MaxInt
	visited.Add(root.board)

	neighbors := root.board.getNeighbors()
	for _, child := range neighbors {
		if !visited.Contains(child.board) {
			newPath, newResult := search(child, bound)
			if newResult == 0 {
				return append([]ManhattanState{root}, newPath...), 0
			}
			if newResult < min {
				min = newResult
			}
		}
	}

	visited.Remove(root.board)
	return nil, min
}

// Main function to run the IDA* search
func idaStar(root ManhattanState) (path []ManhattanState, result int) {
	bound := root.h
	for {
		_, result = search(root, bound)
		if result == 0 {
			return path, result
		}
		if result == math.MaxInt {
			return nil, -1
		}
		bound = result
	}
}

func main() {
	// Input
	size := 8
	emptyTileIndex := 8
	board := Board{dimension: int(math.Sqrt(float64(size + 1)))}
	board.tiles = make([]Tile, size+1)
	for i := 0; i < size+1; i++ {
		board.tiles[i] = Tile(i)
		if board.tiles[i] == 0 {
			board.emptyTileIndex = i
		}
	}

	goal = goalBoard(size, emptyTileIndex)
	if !board.isSolvable() {
		fmt.Println(-1)
		return
	}

	// Start the search
	start := ManhattanState{board: board, g: 0, h: board.manhattan(goal)}
	path, _ := idaStar(start)
	fmt.Println(path)

	// Output result
	// if path == nil {
	// 	fmt.Println(-1)
	// } else {
	// 	fmt.Println(len(path) - 1)
	// 	for _, state := range path[1:] {
	// 		fmt.Println(state.direction)
	// 	}
	// }
}

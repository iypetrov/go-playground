package main

import (
	"fmt"
	"math"

	"github.com/emirpasic/gods/sets/hashset"
)

type Direction int

const (
	LEFT Direction = iota
	RIGHT
	TOP
	DOWN
	NONE
)

func (d Direction) ToString() string {
	switch d {
	case LEFT:
		return "right"
	case RIGHT:
		return "left"
	case TOP:
		return "down"
	case DOWN:
		return "up"
	default:
		return ""
	}
}

type Board struct {
	Tiles          []int
	Dimension      float64
	EmptyTileIndex int
}

func NewBoard(size int, tiles []int) Board {
	size++
	return Board{
		Tiles:          tiles,
		Dimension:      math.Sqrt(float64(size)),
		EmptyTileIndex: findEmptyTileIndex(tiles),
	}
}

func findEmptyTileIndex(tiles []int) int {
	for i, tile := range tiles {
		if tile == 0 {
			return i
		}
	}
	return -1
}

func (b *Board) ManhattanDistance(goal Board) int {
	currentPositions := make(TilePosition)
	goalPositions := make(TilePosition)
	size := len(b.Tiles)

	for i := 0; i < size; i++ {
		currentPositions[b.Tiles[i]] = struct {
			Row    int
			Column int
		}{
			Row:    int(i / int(b.Dimension)),
			Column: i % int(b.Dimension),
		}

		goalPositions[goal.Tiles[i]] = struct {
			Row    int
			Column int
		}{
			Row:    int(i / int(goal.Dimension)),
			Column: i % int(goal.Dimension),
		}
	}

	sum := 0
	for tile, position := range currentPositions {
		goalPosition := goalPositions[tile]
		sum += int(math.Abs(float64(position.Row-goalPosition.Row)) + math.Abs(float64(position.Column-goalPosition.Column)))
	}

	return sum
}

func (b *Board) Goal(endEmptyTileIndex int) Board {
	size := len(b.Tiles)
	if endEmptyTileIndex == -1 {
		endEmptyTileIndex = size - 1
	}

	goalBoard := Board{
		Tiles:          make([]int, size+1),
		Dimension:      math.Sqrt(float64(size)),
		EmptyTileIndex: endEmptyTileIndex,
	}

	for i := 0; i < endEmptyTileIndex; i++ {
		goalBoard.Tiles[i] = i + 1
	}
	for i := endEmptyTileIndex; i < size; i++ {
		goalBoard.Tiles[i] = i
	}
	goalBoard.Tiles[endEmptyTileIndex] = 0

	return goalBoard
}

func (b *Board) Solvable() bool {
	tiles := b.Tiles
	dimension := int(b.Dimension)
	inversions := 0
	for i := 0; i < len(tiles)-1; i++ {
		for j := i + 1; j < len(tiles); j++ {
			if tiles[i] > 0 && tiles[j] > 0 && tiles[i] > tiles[j] {
				inversions++
			}
		}
	}

	if int(b.Dimension)%2 != 0 {
		return inversions%2 == 0
	} else {
		emptyRowIndex := int(math.Floor(float64(b.EmptyTileIndex) / float64(dimension)))
		return (emptyRowIndex+inversions)%2 != 0
	}
}

func (b *Board) Left() *Board {
	dimension := int(b.Dimension)
	emptyTileIndex := b.EmptyTileIndex
	emptyTileColumn := emptyTileIndex % dimension

	if emptyTileColumn-1 >= 0 {
		tempBoard := b
		tempBoard.Tiles = append([]int(nil), b.Tiles...)
		temp := tempBoard.Tiles[emptyTileIndex-1]
		tempBoard.Tiles[emptyTileIndex-1] = tempBoard.Tiles[emptyTileIndex]
		tempBoard.Tiles[emptyTileIndex] = temp
		tempBoard.EmptyTileIndex = emptyTileIndex - 1
		return tempBoard
	}

	return nil
}

func (b *Board) Right() *Board {
	dimension := int(b.Dimension)
	emptyTileIndex := b.EmptyTileIndex
	emptyTileColumn := emptyTileIndex % dimension

	if emptyTileColumn+1 < dimension {
		tempBoard := b
		tempBoard.Tiles = append([]int(nil), b.Tiles...)
		temp := tempBoard.Tiles[emptyTileIndex+1]
		tempBoard.Tiles[emptyTileIndex+1] = tempBoard.Tiles[emptyTileIndex]
		tempBoard.Tiles[emptyTileIndex] = temp
		tempBoard.EmptyTileIndex = emptyTileIndex + 1
		return tempBoard
	}

	return nil
}

func (b *Board) Top() *Board {
	dimension := int(b.Dimension)
	emptyTileIndex := b.EmptyTileIndex
	emptyTileRow := emptyTileIndex / dimension

	if emptyTileRow-1 >= 0 {
		tempBoard := b
		tempBoard.Tiles = append([]int(nil), b.Tiles...)
		temp := tempBoard.Tiles[emptyTileIndex-dimension]
		tempBoard.Tiles[emptyTileIndex-dimension] = tempBoard.Tiles[emptyTileIndex]
		tempBoard.Tiles[emptyTileIndex] = temp
		tempBoard.EmptyTileIndex = emptyTileIndex - dimension
		return tempBoard
	}

	return nil
}

func (b *Board) Bottom() *Board {
	dimension := int(b.Dimension)
	emptyTileIndex := b.EmptyTileIndex
	emptyTileRow := emptyTileIndex / dimension

	if emptyTileRow+1 < dimension {
		tempBoard := b
		tempBoard.Tiles = append([]int(nil), b.Tiles...)
		temp := tempBoard.Tiles[emptyTileIndex+dimension]
		tempBoard.Tiles[emptyTileIndex+dimension] = tempBoard.Tiles[emptyTileIndex]
		tempBoard.Tiles[emptyTileIndex] = temp
		tempBoard.EmptyTileIndex = emptyTileIndex + dimension
		return tempBoard
	}

	return nil
}

func (b *Board) Equal(board Board) bool {
	for i, tile := range b.Tiles {
		if tile != board.Tiles[i] {
			return false
		}
	}
	return true
}

func (b *Board) ToString() string {
	return fmt.Sprintf("%v", b.Tiles)
}

type TilePosition map[int]struct {
	Row    int
	Column int
}

type ManhattanState struct {
	Goal      Board
	Board     Board
	Direction Direction
	G         int
	H         int
}

func NewManhattanState(goal, board Board, direction Direction, g int) ManhattanState {
	return ManhattanState{
		Goal:      goal,
		Board:     board,
		Direction: direction,
		G:         g,
		H:         board.ManhattanDistance(goal),
	}
}

func (ms *ManhattanState) Neighbours() []ManhattanState {
	var neighbours []ManhattanState

	left := ms.Board.Left()
	if left != nil {
		neighbours = append(neighbours, NewManhattanState(ms.Goal, *left, LEFT, ms.G+1))
	}

	right := ms.Board.Right()
	if right != nil {
		neighbours = append(neighbours, NewManhattanState(ms.Goal, *right, RIGHT, ms.G+1))
	}

	top := ms.Board.Top()
	if top != nil {
		neighbours = append(neighbours, NewManhattanState(ms.Goal, *top, TOP, ms.G+1))
	}

	bottom := ms.Board.Bottom()
	if bottom != nil {
		neighbours = append(neighbours, NewManhattanState(ms.Goal, *bottom, DOWN, ms.G+1))
	}

	return neighbours
}

func (ms *ManhattanState) Search(goal Board) ([]string, int, error) {
	bound := ms.H
	var path []ManhattanState
	path = append(path, *ms)

	visited := hashset.New()
	visited.Add(ms.Board.ToString())

	for {
		t := search(goal, bound, visited, &path)
		if t == math.MaxInt {
			return nil, -1, fmt.Errorf("board is not solvable")
		}

		if t == 0 {
			var directions []string
			for _, state := range path {
				if state.Direction != NONE {
					directions = append(directions, state.Direction.ToString())
				}
			}
			return directions, len(path) - 1, nil
		}

		bound = t
	}
}

func search(goal Board, bound int, visited *hashset.Set, path *[]ManhattanState) int {
	node := (*path)[len(*path)-1]
	f := node.G + node.H

	if f > bound {
		return f
	}

	if node.Board.Equal(goal) {
		return 0
	}

	min := math.MaxInt
	for _, neighbour := range node.Neighbours() {
		neighbourBoardString := neighbour.Board.ToString()
		if !visited.Contains(neighbourBoardString) {
			*path = append(*path, neighbour)
			visited.Add(neighbourBoardString)
			t := search(goal, bound, visited, path)

			if t == 0 {
				return 0
			}

			if t < min {
				min = t
			}

			*path = (*path)[:len(*path)-1]
			visited.Remove(neighbourBoardString)
		}
	}

	return min
}

func NPuzzle(size, emptyTileIndex int, tiles []int) ([]string, int, error) {
	board := NewBoard(size, tiles)
	goal := board.Goal(emptyTileIndex)

	if !board.Solvable() {
		return []string{}, -1, fmt.Errorf("board is not solvable")
	}

	startState := NewManhattanState(goal, board, NONE, 0)

	path, steps, err := startState.Search(goal)
	if err != nil {
		return []string{}, -1, err
	}
	return path, steps, nil
}

func main() {
	fmt.Println("hello n-puzzle")
}

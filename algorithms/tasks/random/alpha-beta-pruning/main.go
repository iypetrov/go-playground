package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	X          = "X"
	O          = "0"
	BOARD_SIZE = 3
	EMPTY_CELL = "_"
)

type Game struct {
	Board          [BOARD_SIZE][BOARD_SIZE]string
	PlayerSymbol   string
	ComputerSymbol string
	Steps          int
}

func NewGame(isPlayerFirst bool) *Game {
	g := &Game{
		PlayerSymbol:   X,
		ComputerSymbol: O,
	}

	if !isPlayerFirst {
		g.PlayerSymbol, g.ComputerSymbol = O, X
	}

	for i := range g.Board {
		for j := range g.Board[i] {
			g.Board[i][j] = EMPTY_CELL
		}
	}

	return g
}

func (g *Game) SetBotPrepState(board [BOARD_SIZE][BOARD_SIZE]string) {
	g.PlayerSymbol = O
	g.ComputerSymbol = X
	g.Board = board

	numSteps := 0
	for i := range BOARD_SIZE {
		for j := range BOARD_SIZE {
			if board[i][j] != EMPTY_CELL {
				numSteps++
			}
		}
	}
	g.Steps = numSteps
}

func (g *Game) Print() {
	for i := 0; i < BOARD_SIZE; i++ {
		fmt.Println(strings.Join(g.Board[i][:], " "))
	}
}

func (g *Game) GameOver() bool {
	for i := 0; i < BOARD_SIZE; i++ {
		if g.Board[i][0] != EMPTY_CELL && g.Board[i][0] == g.Board[i][1] && g.Board[i][0] == g.Board[i][2] {
			return true
		}
		if g.Board[0][i] != EMPTY_CELL && g.Board[0][i] == g.Board[1][i] && g.Board[0][i] == g.Board[2][i] {
			return true
		}
	}

	if g.Board[0][0] != EMPTY_CELL && g.Board[0][0] == g.Board[1][1] && g.Board[0][0] == g.Board[2][2] {
		return true
	}
	if g.Board[0][2] != EMPTY_CELL && g.Board[0][2] == g.Board[1][1] && g.Board[0][2] == g.Board[2][0] {
		return true
	}

	for i := 0; i < BOARD_SIZE; i++ {
		for j := 0; j < BOARD_SIZE; j++ {
			if g.Board[i][j] == EMPTY_CELL {
				return false
			}
		}
	}

	return true
}

func (g *Game) BoardScore(depth int) int {
	for i := 0; i < BOARD_SIZE; i++ {
		if g.Board[i][0] != EMPTY_CELL && g.Board[i][0] == g.Board[i][1] && g.Board[i][0] == g.Board[i][2] {
			if g.Board[i][0] == g.ComputerSymbol {
				return 10 - depth
			}
			return depth - 10
		}
		if g.Board[0][i] != EMPTY_CELL && g.Board[0][i] == g.Board[1][i] && g.Board[0][i] == g.Board[2][i] {
			if g.Board[0][i] == g.ComputerSymbol {
				return 10 - depth
			}
			return depth - 10
		}
	}

	if g.Board[0][0] != EMPTY_CELL && g.Board[0][0] == g.Board[1][1] && g.Board[0][0] == g.Board[2][2] {
		if g.Board[0][0] == g.ComputerSymbol {
			return 10 - depth
		}
		return depth - 10
	}
	if g.Board[0][2] != EMPTY_CELL && g.Board[0][2] == g.Board[1][1] && g.Board[0][2] == g.Board[2][0] {
		if g.Board[0][2] == g.ComputerSymbol {
			return 10 - depth
		}
		return depth - 10
	}

	return 0
}

func (g *Game) MinMaxScore(depth int, alpha, beta int, isMaximizer bool) int {
	score := g.BoardScore(depth)
	if score != 0 || g.GameOver() {
		return score
	}

	if isMaximizer {
		best := -1000
		for i := 0; i < BOARD_SIZE; i++ {
			for j := 0; j < BOARD_SIZE; j++ {
				if g.Board[i][j] == EMPTY_CELL {
					g.Board[i][j] = g.ComputerSymbol
					best = max(best, g.MinMaxScore(depth+1, alpha, beta, false))
					g.Board[i][j] = EMPTY_CELL
					alpha = max(alpha, best)
					if beta <= alpha {
						break
					}
				}
			}
		}
		return best
	} else {
		best := 1000
		for i := 0; i < BOARD_SIZE; i++ {
			for j := 0; j < BOARD_SIZE; j++ {
				if g.Board[i][j] == EMPTY_CELL {
					g.Board[i][j] = g.PlayerSymbol
					best = min(best, g.MinMaxScore(depth+1, alpha, beta, true))
					g.Board[i][j] = EMPTY_CELL
					beta = min(beta, best)
					if beta <= alpha {
						break
					}
				}
			}
		}
		return best
	}
}

func (g *Game) BestMove() {
	bestVal := -1000
	bestRow, bestCol := -1, -1
	alpha, beta := -1000, 1000

	for i := 0; i < BOARD_SIZE; i++ {
		for j := 0; j < BOARD_SIZE; j++ {
			if g.Board[i][j] == EMPTY_CELL {
				g.Board[i][j] = g.ComputerSymbol
				moveVal := g.MinMaxScore(0, alpha, beta, false)
				g.Board[i][j] = EMPTY_CELL
				if moveVal > bestVal {
					bestVal = moveVal
					bestRow, bestCol = i, j
				}
			}
		}
	}

	if bestRow != -1 && bestCol != -1 {
		g.Board[bestRow][bestCol] = g.ComputerSymbol
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to go first? (y/n): ")
	input, _ := reader.ReadString('\n')
	isPlayerFirst := strings.TrimSpace(strings.ToLower(input)) == "y"

	game := NewGame(isPlayerFirst)

	if !isPlayerFirst {
		game.BestMove()
		game.Print()
	}

	for !game.GameOver() {
		fmt.Print("Enter your move (row and column): ")
		input, _ := reader.ReadString('\n')
		parts := strings.Fields(input)
		if len(parts) != 2 {
			fmt.Println("Invalid input. Enter row and column.")
			continue
		}
		row, err1 := strconv.Atoi(parts[0])
		col, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || row < 1 || col < 1 || row > BOARD_SIZE || col > BOARD_SIZE || game.Board[row-1][col-1] != EMPTY_CELL {
			fmt.Println("Invalid move. Try again.")
			continue
		}
		game.Board[row-1][col-1] = game.PlayerSymbol
		game.Steps++
		if game.GameOver() {
			break
		}
		game.BestMove()
		game.Print()
	}

	score := game.BoardScore(game.Steps)
	if score > 0 {
		fmt.Println("Computer wins!")
	} else if score < 0 {
		fmt.Println("Player wins!")
	} else {
		fmt.Println("It's a tie!")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

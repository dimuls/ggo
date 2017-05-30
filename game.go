package ggo

import (
	"errors"
)

type Parameters struct {
	BoardSize int `json:"boardSize"`
}

type Game struct {
	parameters     Parameters
	board          *board
	turnColor      Color
	turnID         int
	availableMoves map[int]bool
	boardHistory   map[boardHashSum]bool
}

func NewGame(parameters Parameters) *Game {
	return &Game{
		parameters: parameters,
		board:      newBoard(parameters.BoardSize),
		turnColor:  Black,
		turnID:     1,
	}
}

func (g *Game) Move(row int, column int, color Color) error {
	if g.turnColor != color {
		return errors.New("now turn of another color")
	}
	err := g.board.put(row, column, color)
	if err != nil {
		return err
	}
	if g.turnColor == Black {
		g.turnColor = White
	} else {
		g.turnColor = Black
	}
	g.turnID++
	g.boardHistory[g.board.getHashSum()] = true
	g.computeAvailableMoves()
}

func (g *Game) computeAvailableMoves() {

}

package ggo

import (
	"errors"

	"github.com/someanon/ggo/timer"
)

type Parameters struct {
	BoardSize  int               `json:"boardSize"`
	TimeSystem *timer.Parameters `json:"timeSystem"`
}

type Game struct {
	parameters       Parameters
	board            *board
	timer            *timer.Timer
	moveColor        Color
	moveID           int
	disallowedPlaces map[[2]int]nothing
}

func NewGame(parameters Parameters) *Game {
	g := &Game{
		parameters:       parameters,
		board:            newBoard(parameters.BoardSize),
		timer:            nil,
		moveColor:        Black,
		moveID:           1,
		disallowedPlaces: nil,
	}
	g.computeDisallowedMoves()
	return g
}

func (g *Game) Move(row int, column int, color Color) error {
	if g.moveColor != color {
		return errors.New("turn of another color")
	}
	if _, exists := g.disallowedPlaces[[2]int{row, column}]; exists {
		return errors.New("move is  disallowed")
	}
	err := g.board.put(row, column, color)
	if err != nil {
		return err
	}
	g.nextMove()
	g.computeDisallowedMoves()
}

func (g *Game) Pass(color Color) error {
	if g.moveColor != color {
		return errors.New("turn of another color")
	}
	g.nextMove()
}

func (g *Game) nextMove() {
	g.moveColor = g.nextColor(g.moveColor)
	g.moveID++
}

func (g *Game) nextColor(color Color) Color {
	if color == Black {
		return White
	}
	return Black
}

func (g *Game) computeDisallowedMoves() {
	g.disallowedPlaces = make(map[[2]int]nothing)
	for r := 0; r < g.board.size; r++ {
		for c := 0; c < g.board.size; c++ {
			p := g.board.places[r][c]
			if p.color == Empty {
				libertiesCount, _, dyingEnemyGroups := p.analyzeNeighbors(g.moveColor)
				if libertiesCount == 0 && len(dyingEnemyGroups) == 0 {
					g.disallowedPlaces[[2]int{r, c}] = nothing{}
				}
			}
		}
	}
	if g.board.koPlace != nil {
		g.disallowedPlaces[[2]int{g.board.koPlace.row, g.board.koPlace.column}] = nothing{}
	}
}

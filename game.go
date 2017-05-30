package ggo

import (
	"errors"
)

type Parameters struct {
	BoardSize int `json:"boardSize"`
}

type State struct {
	MoveID           int           `json:"moveID"`
	MoveColor        Color         `json:"moveColor"`
	DisallowedPlaces map[int]bool  `json:"disallowedPlaces"`
	Stones           map[int]Color `json:"stones"`
}

type Game struct {
	parameters       Parameters
	board            *board
	boardStates      map[hashSum]struct{}
	moveColor        Color
	moveID           int
	disallowedPlaces map[int]struct{}
}

func NewGame(parameters Parameters) *Game {
	g := &Game{
		parameters:       parameters,
		board:            newBoard(parameters.BoardSize),
		boardStates:      make(map[hashSum]struct{}),
		moveColor:        Black,
		moveID:           1,
		disallowedPlaces: nil,
	}
	g.computeAvailableMoves()
	return g
}

func (g *Game) GetState() State {
	s := State{
		MoveID:           g.moveID,
		MoveColor:        g.moveColor,
		DisallowedPlaces: make(map[int]bool, len(g.disallowedPlaces)),
		Stones:           make(map[int]Color),
	}
	for dp := range g.disallowedPlaces {
		s.DisallowedPlaces[dp] = false
	}
	for r := 0; r < g.board.size; r++ {
		for c := 0; c < g.board.size; c++ {
			if g.board.places[r][c].color != Empty {
				s.Stones[g.board.getPlaceID(r, c)] = g.board.places[r][c].color
			}
		}
	}
	return s
}

func (g *Game) Move(row int, column int, color Color) error {
	if g.moveColor != color {
		return errors.New("turn of another color")
	}
	if _, exists := g.disallowedPlaces[g.board.getPlaceID(row, column)]; exists {
		return errors.New("move is  disallowed")
	}
	err := g.board.put(row, column, color)
	if err != nil {
		return err
	}
	g.nextMove()
	g.boardStates[g.board.getHashSum()] = true
	g.computeAvailableMoves()
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

func (g *Game) computeAvailableMoves() {
	g.disallowedPlaces = make(map[int]struct{})
	for r := 0; r < g.board.size; r++ {
		for c := 0; c < g.board.size; c++ {
			pID := g.board.getPlaceID(r, c)
			p := g.board.places[r][c]
			noLiberties, dyingEnemyGroups := p.analyzeNeighbors(g.moveColor)
			if noLiberties && len(dyingEnemyGroups) == 0 {
				g.disallowedPlaces[pID] = struct{}{}
			} else if len(dyingEnemyGroups) == 1 && len(dyingEnemyGroups[0].places) == 1 {
				// we need to check for ko
				p.color = g.moveColor
				dyingEnemyGroups[0].places[0].color = Empty
				if _, exists := g.boardStates[g.board.getHashSum()]; exists {
					// ko place, disallow
					g.disallowedPlaces[pID] = struct{}{}
				}
				p.color = Empty
				dyingEnemyGroups[0].places[0].color = g.nextColor(g.moveColor)
			}
		}
	}
}

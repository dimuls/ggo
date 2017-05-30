package ggo

import (
	"errors"
)

type place struct {
	board  *board
	group  *group
	row    int
	column int
	color  Color
}

func newPlace(board *board, row int, column int) *place {
	return &place{
		board:  board,
		group:  nil,
		row:    row,
		column: column,
		color:  Empty,
	}
}

func (p *place) getNeighbors() []*place {
	neighbors := make([]*place, 0, 2)
	if p, err := p.board.getPlace(p.row-1, p.column); err == nil {
		neighbors = append(neighbors, p)
	}
	if p, err := p.board.getPlace(p.row, p.column-1); err == nil {
		neighbors = append(neighbors, p)
	}
	if p, err := p.board.getPlace(p.row+1, p.column); err == nil {
		neighbors = append(neighbors, p)
	}
	if p, err := p.board.getPlace(p.row, p.column+1); err == nil {
		neighbors = append(neighbors, p)
	}
	return neighbors
}

func (p *place) die() {
	enemyGroups := make(map[*group]bool)
	for _, n := range p.getNeighbors() {
		if n.color != Empty && n.color != p.color {
			enemyGroups[n.group] = true
		}
	}
	for eg := range enemyGroups {
		eg.liberties++
	}
	p.group = nil
	p.color = Empty
}

func (p *place) put(color Color) error {
	if color == Empty {
		return errors.New("color shouldn't be empty")
	}
	if p.color != Empty {
		return errors.New("place already occupied")
	}

	liberties := make(map[*place]bool)
	friendGroups := make(map[*group]bool)
	enemyGroups := make(map[*group]bool)

	for _, n := range p.getNeighbors() {
		if n.color == Empty {
			liberties[n] = true
		} else if n.color == color {
			friendGroups[n.group] = true
		} else {
			enemyGroups[n.group] = true
		}
	}

	for fg := range friendGroups {
		for _, fgp := range fg.places {
			for _, fgpn := range fgp.getNeighbors() {
				if fgpn.color == Empty {
					liberties[fgpn] = true
				}
			}
		}
	}

	dyingEnemyGroups := make([]*group, 0)
	for eg := range enemyGroups {
		if eg.liberties == 1 {
			dyingEnemyGroups = append(dyingEnemyGroups, eg)
		}
	}

	if len(liberties) == 0 && len(dyingEnemyGroups) == 0 {
		return errors.New("place haven't life and no neighbor enemy group is dying")
	}

	p.color = color

	if len(friendGroups) == 0 {
		p.group = &group{
			places: []*place{p},
		}
	} else {
		var baseGroup *group
		for fg := range friendGroups {
			if baseGroup == nil {
				baseGroup = fg
				continue
			}
			baseGroup.join(fg)
		}
		baseGroup.places = append(baseGroup.places, p)
		baseGroup.liberties = len(liberties)
	}

	for _, eg := range dyingEnemyGroups {
		eg.die()
	}

	return nil
}

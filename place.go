package ggo

import (
	"errors"
)

type nothing struct{}

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

func (p *place) neighbors() []*place {
	neighbors := make([]*place, 0, 2)
	if p, err := p.board.place(p.row-1, p.column); err == nil {
		neighbors = append(neighbors, p)
	}
	if p, err := p.board.place(p.row, p.column-1); err == nil {
		neighbors = append(neighbors, p)
	}
	if p, err := p.board.place(p.row+1, p.column); err == nil {
		neighbors = append(neighbors, p)
	}
	if p, err := p.board.place(p.row, p.column+1); err == nil {
		neighbors = append(neighbors, p)
	}
	return neighbors
}

func (p *place) analyzeNeighbors(color Color) (int, []*group, []*group) {
	libertiesMap := make(map[*place]nothing)
	friendGroupsMap := make(map[*group]nothing)
	enemyGroupsMap := make(map[*group]nothing)

	for _, n := range p.neighbors() {
		if n.color == Empty {
			libertiesMap[n] = nothing{}
		} else if n.color == color {
			friendGroupsMap[n.group] = nothing{}
		} else {
			enemyGroupsMap[n.group] = nothing{}
		}
	}

	for fg := range friendGroupsMap {
		for _, fgp := range fg.places {
			for _, fgpn := range fgp.neighbors() {
				if fgpn != p && fgpn.color == Empty {
					libertiesMap[fgpn] = nothing{}
				}
			}
		}
	}

	dyingEnemyGroups := make([]*group, 0)
	for eg := range enemyGroupsMap {
		if eg.liberties == 1 {
			dyingEnemyGroups = append(dyingEnemyGroups, eg)
		}
	}

	friendGroups := make([]*group, 0, len(friendGroupsMap))
	for fg := range friendGroupsMap {
		friendGroups = append(friendGroups, fg)
	}

	return len(libertiesMap), friendGroups, dyingEnemyGroups
}

func (p *place) put(color Color) error {
	if color == Empty {
		return errors.New("color shouldn't be empty")
	}
	if p.color != Empty {
		return errors.New("already occupied")
	}

	libertiesCount, friendGroups, dyingEnemyGroups := p.analyzeNeighbors(color)

	if libertiesCount == 0 && len(dyingEnemyGroups) == 0 {
		return errors.New("no liberties and no neighbor enemy group is dying")
	}

	p.color = color

	if len(friendGroups) == 0 {
		p.group = &group{
			places:    []*place{p},
			liberties: libertiesCount + len(dyingEnemyGroups),
		}
	} else {
		var baseGroup *group
		for _, fg := range friendGroups {
			if baseGroup == nil {
				baseGroup = fg
				continue
			}
			baseGroup.join(fg)
		}
		baseGroup.places = append(baseGroup.places, p)
		baseGroup.liberties = libertiesCount
	}

	for _, eg := range dyingEnemyGroups {
		eg.die()
	}

	if libertiesCount == 0 && len(friendGroups) == 0 &&
		len(dyingEnemyGroups) == 1 && len(dyingEnemyGroups[0].places) == 1 {
		p.board.koPlace = dyingEnemyGroups[0].places[0]
	} else {
		p.board.koPlace = nil
	}

	return nil
}

func (p *place) die() {
	enemyGroups := make(map[*group]struct{})
	for _, n := range p.neighbors() {
		if n.color != Empty && n.color != p.color {
			enemyGroups[n.group] = nothing{}
		}
	}
	for eg := range enemyGroups {
		eg.liberties++
	}
	p.group = nil
	p.color = Empty
}

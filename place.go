package ggo

import (
	"errors"
)

type place struct {
	board  *board
	group  *group
	id     int
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

func (p *place) analyzeNeighbors(color Color) (bool, []*group) {
	libertiesMap := make(map[*place]struct{})
	friendGroupsMap := make(map[*group]struct{})
	enemyGroupsMap := make(map[*group]struct{})

	for _, n := range p.getNeighbors() {
		if n.color == Empty {
			libertiesMap[n] = struct{}{}
		} else if n.color == color {
			friendGroupsMap[n.group] = struct{}{}
		} else {
			enemyGroupsMap[n.group] = struct{}{}
		}
	}

	noLiberties := true
	for fg := range friendGroupsMap {
		for _, fgp := range fg.places {
			for _, fgpn := range fgp.getNeighbors() {
				if fgpn != p && fgpn.color == Empty {
					noLiberties = false
					break
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

	return noLiberties, dyingEnemyGroups
}

func (p *place) analyzeNeighbors2(color Color) (int, []*group, []*group) {
	libertiesMap := make(map[*place]struct{})
	friendGroupsMap := make(map[*group]struct{})
	enemyGroupsMap := make(map[*group]struct{})

	for _, n := range p.getNeighbors() {
		if n.color == Empty {
			libertiesMap[n] = struct{}{}
		} else if n.color == color {
			friendGroupsMap[n.group] = struct{}{}
		} else {
			enemyGroupsMap[n.group] = struct{}{}
		}
	}

	for fg := range friendGroupsMap {
		for _, fgp := range fg.places {
			for _, fgpn := range fgp.getNeighbors() {
				if fgpn != p && fgpn.color == Empty {
					libertiesMap[fgpn] = struct{}{}
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
		return errors.New("place already occupied")
	}

	libertiesCount, friendGroups, dyingEnemyGroups := p.analyzeNeighbors2(color)

	if libertiesCount == 0 && len(dyingEnemyGroups) == 0 {
		return errors.New("place haven't life and no dying neighbor enemy group around")
	}

	p.color = color

	if len(friendGroups) == 0 {
		p.group = &group{
			places: []*place{p},
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

	return nil
}

func (p *place) die() {
	enemyGroups := make(map[*group]struct{})
	for _, n := range p.getNeighbors() {
		if n.color != Empty && n.color != p.color {
			enemyGroups[n.group] = struct{}{}
		}
	}
	for eg := range enemyGroups {
		eg.liberties++
	}
	p.group = nil
	p.color = Empty
}

package ggo

type group struct {
	places    []*place
	liberties int
}

func (g *group) join(joinGroup *group) {
	for _, p := range joinGroup.places {
		g.places = append(g.places, p)
		p.group = g
	}
}

func (g *group) die() {
	for _, p := range g.places {
		p.die()
	}
}

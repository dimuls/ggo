package ggo

import (
	"errors"
	"fmt"

	"github.com/onsi/gomega/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Board", func() {
	var b *board
	BeforeEach(func() {
		b = newBoard(3)
	})
	Describe("place", func() {
		Context("top left corner", func() {
			Specify("has right neighbors", func() {
				p := b.places[0][0]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(2))
				Expect(ns[0]).Should(EqualToPlace(b.places[1][0]))
				Expect(ns[1]).Should(EqualToPlace(b.places[0][1]))
			})
		})
		Context("top center", func() {
			Specify("has right neighbors", func() {
				p := b.places[0][1]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(3))
				Expect(ns[0]).Should(EqualToPlace(b.places[0][0]))
				Expect(ns[1]).Should(EqualToPlace(b.places[1][1]))
				Expect(ns[2]).Should(EqualToPlace(b.places[0][2]))
			})
		})
		Context("top right corner", func() {
			Specify("has right neighbors", func() {
				p := b.places[0][2]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(2))
				Expect(ns[0]).Should(EqualToPlace(b.places[0][1]))
				Expect(ns[1]).Should(EqualToPlace(b.places[1][2]))
			})
		})
		Context("center left", func() {
			Specify("has right neighbors", func() {
				p := b.places[1][0]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(3))
				Expect(ns[0]).Should(EqualToPlace(b.places[0][0]))
				Expect(ns[1]).Should(EqualToPlace(b.places[2][0]))
				Expect(ns[2]).Should(EqualToPlace(b.places[1][1]))
			})
		})
		Context("center center", func() {
			Specify("has right neighbors", func() {
				p := b.places[1][1]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(4))
				Expect(ns[0]).Should(EqualToPlace(b.places[0][1]))
				Expect(ns[1]).Should(EqualToPlace(b.places[1][0]))
				Expect(ns[2]).Should(EqualToPlace(b.places[2][1]))
				Expect(ns[3]).Should(EqualToPlace(b.places[1][2]))
			})
		})
		Context("center right", func() {
			Specify("has right neighbors", func() {
				p := b.places[1][2]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(3))
				Expect(ns[0]).Should(EqualToPlace(b.places[0][2]))
				Expect(ns[1]).Should(EqualToPlace(b.places[1][1]))
				Expect(ns[2]).Should(EqualToPlace(b.places[2][2]))
			})
		})
		Context("bottom left corner", func() {
			Specify("has right neighbors", func() {
				p := b.places[2][0]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(2))
				Expect(ns[0]).Should(EqualToPlace(b.places[1][0]))
				Expect(ns[1]).Should(EqualToPlace(b.places[2][1]))
			})
		})
		Context("bottom center", func() {
			Specify("has right neighbors", func() {
				p := b.places[2][1]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(3))
				Expect(ns[0]).Should(EqualToPlace(b.places[1][1]))
				Expect(ns[1]).Should(EqualToPlace(b.places[2][0]))
				Expect(ns[2]).Should(EqualToPlace(b.places[2][2]))
			})
		})
		Context("bottom right corner", func() {
			Specify("has right neighbors", func() {
				p := b.places[2][2]
				ns := p.neighbors()
				Expect(ns).Should(HaveLen(2))
				Expect(ns[0]).Should(EqualToPlace(b.places[1][2]))
				Expect(ns[1]).Should(EqualToPlace(b.places[2][1]))
			})
		})
	})
})

func EqualToPlace(expected *place) types.GomegaMatcher {
	return &equalToPlaceMatcher{
		expected: expected,
	}
}

type equalToPlaceMatcher struct {
	expected *place
}

func (m *equalToPlaceMatcher) Match(actual interface{}) (success bool,
	err error) {
	pa, ok := actual.(*place)
	if !ok {
		return false, errors.New("want actual to be *place")
	}

	pe := m.expected

	return pa.row == pe.row && pa.column == pe.column, nil
}

func (m *equalToPlaceMatcher) FailureMessage(actual interface{}) (
	message string) {
	pa := actual.(*place)
	pe := m.expected
	return fmt.Sprintf("expected place[%d][%d] but got place[%d][%d]",
		pe.row, pe.column, pa.row, pa.column)
}

func (m *equalToPlaceMatcher) NegatedFailureMessage(actual interface{}) (
	message string) {
	pe := m.expected
	return fmt.Sprintf("expected not place[%d][%d]",
		pe.row, pe.column)
}

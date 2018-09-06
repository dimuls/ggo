package ggo

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type hashSum [md5.Size]byte

type board struct {
	size    int
	places  [][]*place
	koPlace *place
}

func newBoard(size int) *board {
	b := &board{
		size: size,
	}
	b.places = make([][]*place, size)
	for r := 0; r < size; r++ {
		b.places[r] = make([]*place, size)
		for c := 0; c < size; c++ {
			b.places[r][c] = newPlace(b, r, c)
		}
	}
	return b
}

func (b *board) placeID(row int, column int) int {
	return row*b.size + column
}

func (b *board) place(row int, column int) (*place, error) {
	if row < 0 || column < 0 || row >= b.size || column >= b.size {
		return nil, errors.New(fmt.Sprintf("place not found at row=%d, column=%d", row, column))
	}
	return b.places[row][column], nil
}

func (b *board) put(row int, column int, color Color) error {
	p, err := b.place(row, column)
	if err != nil {
		return err
	}
	return p.put(color)
}

func (b *board) hashSum() hashSum {
	bytes := make([]byte, b.size*b.size)
	for r := 0; r < b.size; r++ {
		for c := 0; c < b.size; c++ {
			bytes[r*b.size+c] = byte(b.places[r][c].color)
		}
	}
	return hashSum(md5.Sum(bytes))
}

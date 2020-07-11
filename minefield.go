package main

import (
	"math/rand"
	"time"
)

type Miner interface {
	mineAt(int) bool
	adjacentMines(int) int
}

type Minefield struct {
	Miner
	context *GameContext
	mines   map[int]int
}

func createMinefield(context *GameContext) *Minefield {
	minefield := new(Minefield)
	minefield.context = context
	minefield.mines = make(map[int]int)
	minefield.placeMines()
	return minefield
}

func (minefield *Minefield) mineAt(index int) bool {
	_, hasMine := minefield.mines[index]
	return hasMine
}

func (minefield *Minefield) adjacentMines(index int) int {
	sum := 0
	fn := func(row int, column int) {
		index := minefield.context.layout.options.index(row, column)
		if minefield.mineAt(index) {
			sum++
		}
	}
	minefield.context.layout.options.forEachNeighbor(index, fn)
	return sum
}

func (minefield *Minefield) reset() {
	minefield.mines = make(map[int]int)
	minefield.placeMines()
}

func (minefield *Minefield) placeMines() {
	rows, columns := minefield.context.layout.rowsAndColumns()
	maxIndex := rows * columns
	mineCount := minefield.context.layout.options.mines

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for len(minefield.mines) < mineCount {
		index := r1.Intn(maxIndex)
		minefield.mines[index] = index
	}
}

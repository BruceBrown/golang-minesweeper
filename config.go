package main

import (
	"flag"
)

func skillLevel() string {
	level := flag.String("skillLevel", "beginner", "skillLevel {beginner|intermediate|expert}")
	flag.Parse()
	return *level
}

type Options struct {
	//func tiles() int { return rows * columns }
	//func blanks() int { return tiles() - mines }
	level   string
	rows    int
	columns int
	mines   int
}

func (options *Options) tiles() int {
	return options.rows * options.columns
}

func (options *Options) blanks() int {
	return options.tiles() - options.mines
}

func (options *Options) rowAndColumn(index int) (int, int) {
	row := index / options.columns
	column := index % options.columns

	return row, column
}

func (options *Options) index(row int, column int) int {
	return row*options.columns + column
}

// call function for each neighbor
func (options *Options) forEachNeighbor(index int, fn func(int, int)) {
	row, column := options.rowAndColumn(index)
	for r := row - 1; r <= row+1; r++ {
		for c := column - 1; c <= column+1; c++ {
			if r != row || c != column {
				if r >= 0 && r < options.rows && c >= 0 && c < options.columns {
					fn(r, c)
				}
			}
		}
	}
}

// option constants
const beginnerRows = 9
const beginnerColumns = 9
const beginnerMines = 10

const intermediateRows = 16
const intermediateColumns = 16
const intermediateMines = 40

const expertRows = 16
const expertColumns = 30
const expertMines = 99

func options() *Options {
	level := skillLevel()
	switch level {
	case "beginner":
		return &Options{level: level, rows: beginnerRows, columns: beginnerColumns, mines: beginnerMines}
	case "intermediate":
		return &Options{level: level, rows: intermediateRows, columns: intermediateColumns, mines: intermediateMines}
	case "expert":
		return &Options{level: level, rows: expertRows, columns: expertColumns, mines: expertMines}
	default:
		return &Options{level: "beginner", rows: beginnerRows, columns: beginnerColumns, mines: beginnerMines}
	}
}

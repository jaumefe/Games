package main

import (
	"image/color"
)

const (
	maxTileSize = 55
)

type tile struct {
	flag      bool
	nbhdMines int // Neighborhood mines
	isClicked bool
	isMine    bool
	colors    color.RGBA
}

func (t *tile) SetTileColor() {
	t.colors = nonClicked
	if t.isClicked && !t.isMine {
		switch t.nbhdMines {
		case 1:
			t.colors = blue
		case 2:
			t.colors = green
		case 3:
			t.colors = red
		case 4:
			t.colors = darkBlue
		case 5:
			t.colors = brown
		case 6:
			t.colors = cyan
		case 7:
			t.colors = black
		case 8:
			t.colors = gray
		default:
			t.colors = clicked
		}
	}
	if t.isMine {
		t.colors = red
	}
}

// Returns both index of the initial tile and all its neighbors
func LookForNeighborTiles(idx, totalRows, totalCols int) map[int]bool {
	row, col := SingleArrayToRowAndCol(idx, totalCols)
	excludeIndex := make(map[int]bool, 0)
	excludeIndex[RowAndColToSingleArray(row, col, totalCols)] = true
	if row > 0 {
		excludeIndex[RowAndColToSingleArray(row-1, col, totalCols)] = true
		if col > 0 {
			excludeIndex[RowAndColToSingleArray(row-1, col-1, totalCols)] = true
		}
		if col < totalCols-1 {
			excludeIndex[RowAndColToSingleArray(row-1, col+1, totalCols)] = true
		}
	}
	if row < totalRows-1 {
		excludeIndex[RowAndColToSingleArray(row+1, col, totalCols)] = true
		if col > 0 {
			excludeIndex[RowAndColToSingleArray(row+1, col-1, totalCols)] = true
		}
		if col < totalCols-1 {
			excludeIndex[RowAndColToSingleArray(row+1, col+1, totalCols)] = true
		}
	}
	if col > 0 {
		excludeIndex[RowAndColToSingleArray(row, col-1, totalCols)] = true
	}
	if col < totalCols-1 {
		excludeIndex[RowAndColToSingleArray(row, col+1, totalCols)] = true
	}
	return excludeIndex
}

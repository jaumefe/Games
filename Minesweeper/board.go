package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jaumefe/stats"
)

type Board struct {
	rows  int
	cols  int
	tiles [][]tile
	bombs int
}

func (b *Board) Grid(board *ebiten.Image) {
	for r := range b.rows {
		for c := range b.cols {
			t := b.tiles[r][c]
			x0 := t.x0
			y0 := t.y0
			x1 := x0 + maxTileSize
			y1 := y0 + maxTileSize
			colors := color.RGBA{R: 120, G: 120, B: 120, A: 255}
			if t.isMine {
				colors = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			}
			vector.DrawFilledRect(board, x0, y0, maxTileSize, maxTileSize, colors, false)
			vector.StrokeLine(board, x0, y0, x1, y0, float32(1), color.Black, false)
			vector.StrokeLine(board, x0, y0, x0, y1, float32(1), color.Black, false)
		}
	}

}

func (b *Board) Init(rows, cols, bombs int) {
	b.rows, b.cols, b.bombs = rows, cols, bombs
	b.tiles = make([][]tile, b.rows)
	for y := 0; y < b.rows; y++ {
		b.tiles[y] = make([]tile, b.cols)
		for x := 0; x < b.cols; x++ {
			x0 := float32(x * maxTileSize)
			y0 := float32(y * maxTileSize)
			b.tiles[y][x].x0 = x0
			b.tiles[y][x].y0 = y0
		}
	}
}

func (b *Board) CursorPositionToRowAndCol(x, y int) (row, col int) {
	row = int(y / maxTileSize)
	col = int(x / maxTileSize)
	return
}

func (b *Board) MineDistributionAfterFirstClick(row, col int) {
	length := b.rows * b.cols
	tilesArr := make([]any, 0, length)

	excludeIndex := make(map[int]bool, 0)
	// Initial tile
	excludeIndex[b.RowAndColToSingleArray(row, col)] = true
	fmt.Println(excludeIndex)
	if row > 0 {
		excludeIndex[b.RowAndColToSingleArray(row-1, col)] = true
		if col > 0 {
			excludeIndex[b.RowAndColToSingleArray(row-1, col-1)] = true
		}
		if col < b.cols-1 {
			excludeIndex[b.RowAndColToSingleArray(row-1, col+1)] = true
		}
	}
	if row < b.rows-1 {
		excludeIndex[b.RowAndColToSingleArray(row+1, col)] = true
		if col > 0 {
			excludeIndex[b.RowAndColToSingleArray(row+1, col-1)] = true
		}
		if col < b.cols-1 {
			excludeIndex[b.RowAndColToSingleArray(row+1, col+1)] = true
		}
	}
	if col > 0 {
		excludeIndex[b.RowAndColToSingleArray(row, col-1)] = true
	}
	if col < b.cols-1 {
		excludeIndex[b.RowAndColToSingleArray(row, col+1)] = true
	}

	nb := 0
	for r := 0; r < b.rows; r++ {
		for c := 0; c < b.cols; c++ {
			if nb < b.bombs {
				if !excludeIndex[b.RowAndColToSingleArray(r, c)] {
					b.tiles[r][c].isMine = true
					nb++
				}
			}
			tilesArr = append(tilesArr, b.tiles[r][c])
		}
	}

	fmt.Println(excludeIndex)

	fmt.Println("Before flushing")
	fmt.Println(tilesArr)
	stats.FisherYatesShuffleWithExclusion(tilesArr, excludeIndex)
	fmt.Println("After flushing")
	fmt.Println(tilesArr)
}

func (b *Board) RowAndColToSingleArray(row, col int) int {
	return col + row*b.cols
}

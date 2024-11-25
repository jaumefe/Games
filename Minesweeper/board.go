package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Board struct {
	rows  int
	cols  int
	tiles [][]tile
}

func (b *Board) Grid(board *ebiten.Image) {

	for c := range b.cols {
		for r := range b.rows {
			t := b.tiles[r][c]
			x0 := t.x0
			y0 := t.y0
			x1 := x0 + maxTileSize
			y1 := y0 + maxTileSize
			vector.DrawFilledRect(board, x0, y0, maxTileSize, maxTileSize, color.RGBA{R: 200, G: 200, B: 200, A: 255}, false)
			vector.StrokeLine(board, x0, y0, x1, y0, float32(1), color.Black, false)
			vector.StrokeLine(board, x0, y0, x0, y1, float32(1), color.Black, false)
		}
	}

}

func (b *Board) Init(rows, cols int) error {
	b.rows, b.cols = rows, cols

	b.tiles = make([][]tile, b.cols)

	for x := 0; x < b.cols; x++ {
		b.tiles[x] = make([]tile, b.rows)
		for y := 0; y < b.rows; y++ {
			x0 := float32(x * maxTileSize)
			y0 := float32(y * maxTileSize)
			b.tiles[x][y].x0 = x0
			b.tiles[x][y].y0 = y0
		}
	}

	return nil

}

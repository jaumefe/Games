package main

import "github.com/hajimehoshi/ebiten/v2"

const (
	maxTileSize = 55
)

type tile struct {
	flag      bool
	nbhdMines int // Neighborhood mines
	x0, y0    float32
	mine
}

func (t *tile) DrawTile(board *ebiten.Image) {

}

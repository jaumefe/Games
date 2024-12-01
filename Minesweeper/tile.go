package main

const (
	maxTileSize = 55
)

type tile struct {
	flag      bool
	nbhdMines int // Neighborhood mines
	x0, y0    float32
	isClicked bool
	isMine    bool
}

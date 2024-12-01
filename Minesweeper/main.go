package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidthMax  = 1840
	screenHeightMax = 1080
)

var initialClick = true

type Game struct {
	brd  Board
	diff string
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		row, col := g.brd.CursorPositionToRowAndCol(ebiten.CursorPosition())
		if row < g.brd.rows && col < g.brd.cols {
			g.brd.tiles[row][col].isClicked = true
			if initialClick {
				g.brd.MineDistributionAfterFirstClick(row, col)
				initialClick = false
			}
			fmt.Printf("Tile %d,%d: Clicked\n", row, col)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.brd.Grid(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidthMax, screenHeightMax
}

func (g *Game) Init(difficulty string) {
	var rows, cols, bombs int

	g.diff = difficulty
	if g.diff == "easy" {
		rows = 8
		cols = 8
		bombs = 10
	} else if g.diff == "medium" {
		rows = 16
		cols = 16
		bombs = 40
	} else if g.diff == "hard" {
		rows = 16
		cols = 30
		bombs = 99
	}

	g.brd.Init(rows, cols, bombs)

}

func main() {
	game := &Game{}
	game.Init("easy")
	ebiten.SetWindowSize(screenWidthMax, screenHeightMax)
	ebiten.SetWindowTitle("Minesweeper")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

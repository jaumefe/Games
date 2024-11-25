package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidthMax  = 1840
	screenHeightMax = 1080
)

type Game struct {
	brd  Board
	diff string
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.brd.Grid(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidthMax, screenHeightMax
}

func (g *Game) Init(difficulty string) error {
	var rows, cols int

	g.diff = difficulty
	if g.diff == "easy" {
		rows = 8
		cols = 8
	} else if g.diff == "medium" {
		rows = 16
		cols = 16
	} else if g.diff == "hard" {
		rows = 16
		cols = 30
	}

	if err := g.brd.Init(rows, cols); err != nil {
		return err
	}

	return nil
}

func main() {
	game := &Game{}
	if err := game.Init("hard"); err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(screenWidthMax, screenHeightMax)
	ebiten.SetWindowTitle("Minesweeper")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

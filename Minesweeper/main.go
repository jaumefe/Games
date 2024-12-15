package main

import (
	_ "embed"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/Roboto-Regular.ttf
var robotTFF []byte

const (
	screenWidthMax  = 1840
	screenHeightMax = 1080
)

var initialClick = true

type Game struct {
	brd  Board
	diff string
}

func RowAndColToSingleArray(row, col, totalCols int) int {
	return col + row*totalCols
}

func SingleArrayToRowAndCol(idx, totalCols int) (row, col int) {
	row = int(idx / totalCols)
	col = idx - totalCols*row
	return
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		row, col := g.brd.CursorPositionToRowAndCol(ebiten.CursorPosition())
		if row < g.brd.rows && col < g.brd.cols {
			idx := RowAndColToSingleArray(row, col, g.brd.cols)
			if initialClick {
				g.brd.MineDistributionAfterFirstClick(idx)
				g.brd.CountNeighborMinesAllBoard()
				initialClick = false
			}

			if !g.brd.tiles[idx].isClicked && g.brd.tiles[idx].nbhdMines == 0 && !g.brd.tiles[idx].isMine {
				excludeIndex := map[int]bool{idx: true}
				g.brd.NoMinesAutoShower(idx, excludeIndex)
			}
			g.brd.tiles[idx].isClicked = true

		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.brd.Grid(screen)
	if !initialClick {
		g.brd.ShowNeighborMines(screen)
	}
	g.brd.DrawFlag(screen)
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

	tt, err := opentype.Parse(robotTFF)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	mineFontOpts := &optsFontFace{
		size: 50,
		dpi:  72,
	}
	if err := minesFont.Init(tt, mineFontOpts); err != nil {
		log.Fatalf("failed to init font of mines text: %v", err)
	}

}

func main() {
	game := &Game{}
	game.Init("hard")
	ebiten.SetWindowSize(screenWidthMax, screenHeightMax)
	ebiten.SetWindowTitle("Minesweeper")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

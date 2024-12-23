package main

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/Roboto-Regular.ttf
var robotTFF []byte

const (
	screenWidthMax  = 1840
	screenHeightMax = 1080
)

var initialClick = true
var startTime time.Time
var minutes, seconds = 0, 0
var totalFlags = 0

type Game struct {
	brd      Board
	diff     string
	gameOver gameover
}

type gameover struct {
	Win  bool
	Lose bool
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
	if !g.gameOver.Win && !g.gameOver.Lose {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			row, col := g.brd.CursorPositionToRowAndCol(ebiten.CursorPosition())
			if row < g.brd.rows && col < g.brd.cols {
				idx := RowAndColToSingleArray(row, col, g.brd.cols)
				if !g.brd.tiles[idx].flag {
					if initialClick {
						startTime = time.Now()
						g.brd.MineDistributionAfterFirstClick(idx)
						g.brd.CountNeighborMinesAllBoard()
						initialClick = false
					}

					if !g.brd.tiles[idx].isClicked && g.brd.tiles[idx].nbhdMines == 0 && !g.brd.tiles[idx].isMine {
						excludeIndex := map[int]bool{idx: true}
						g.brd.NoMinesAutoShower(idx, excludeIndex)
					}
					g.brd.tiles[idx].isClicked = true
					if g.brd.tiles[idx].isMine {
						g.gameOver.Lose = true
					}
				}

			}
		} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
			row, col := g.brd.CursorPositionToRowAndCol(ebiten.CursorPosition())
			if row < g.brd.rows && col < g.brd.cols {
				idx := RowAndColToSingleArray(row, col, g.brd.cols)
				if !g.brd.tiles[idx].isClicked {
					g.brd.tiles[idx].flag = !g.brd.tiles[idx].flag
					if g.brd.tiles[idx].flag {
						totalFlags++
					} else {
						totalFlags--
					}
				}
			}
		}
		duration := time.Since(startTime)

		if !initialClick {
			minutes = int(duration.Minutes())
			seconds = int(duration.Seconds()) % 60
		}
	}
	g.gameOver.Win = g.checkIfEndGame()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.brd.Grid(screen)
	if !initialClick {
		g.brd.ShowNeighborMines(screen)
	}
	g.brd.DrawFlag(screen)
	printTime(screen)
	g.brd.ShowTotalFlags(screen)

	if g.gameOver.Lose {
		g.brd.PlotMines(screen)
	}
}

func printTime(screen *ebiten.Image) {
	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(screenWidthMax-100, 100)
	text.Draw(screen, timeStr, counterFont.txt, opts)
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

	counterFontOpts := &optsFontFace{
		size: 32,
		dpi:  48,
	}
	if err := counterFont.Init(tt, counterFontOpts); err != nil {
		log.Fatalf("failed to init font of mines text: %v", err)
	}
}

func (g *Game) checkIfEndGame() bool {
	totalClickedTiles := 0
	for i := 0; i < len(g.brd.tiles); i++ {
		if g.brd.tiles[i].isClicked {
			totalClickedTiles++
		}
	}

	return totalClickedTiles == len(g.brd.tiles)-g.brd.bombs
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

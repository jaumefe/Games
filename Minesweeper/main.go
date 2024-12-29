package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
	reset    bool
}

func (g *Game) Update() error {
	// Initializating the game
	if g.reset {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if x > bEasy.x0 && x < bEasy.x0+bEasy.width && y > bEasy.y0 && y < bEasy.y0+bEasy.height {
				g.diff = "easy"
			} else if x > bMedium.x0 && x < bMedium.x0+bMedium.width && y > bMedium.y0 && y < bMedium.y0+bMedium.height {
				g.diff = "medium"
			} else if x > bHard.x0 && x < bHard.x0+bHard.width && y > bHard.y0 && y < bHard.y0+bHard.height {
				g.diff = "hard"
			}
		}
		if g.diff != "" {
			g.Init()
		}
		return nil
	}

	// Reset button logic
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > bReset.x0 && x < bReset.x0+bReset.width && y > bReset.y0 && y < bReset.y0+bReset.height {
			g.reset = true
			g.diff = ""
			initialClick = true
			minutes, seconds = 0, 0
			totalFlags = 0
			g.gameOver.Win = false
			g.gameOver.Lose = false
			return nil
		}
	}

	// Game logic
	if !g.gameOver.Win && !g.gameOver.Lose {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
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
		} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
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
	if !g.reset {
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
		g.ResetButton(screen)
	} else {
		g.DifficultySelector(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidthMax, screenHeightMax
}

func main() {
	game := &Game{reset: true}
	if err := LoadFont(); err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(screenWidthMax, screenHeightMax)
	ebiten.SetWindowTitle("Minesweeper")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

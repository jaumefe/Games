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
var totalFlags = 0
var playerName []rune
var Name string
var popUp = false

var duration Duration = Duration{Minutes: 0, Seconds: 0}

type Game struct {
	brd      Board
	diff     string
	gameOver gameover
	reset    bool
	records  bool
}

func (g *Game) Update() error {
	// Initializating the game
	if g.reset {
		g.LogicDifficultySelector()

		if g.diff != "" {
			g.Init()
		}
		return nil
	}

	// Reset button logic
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > bReset.x0 && x < bReset.x0+bReset.width && y > bReset.y0 && y < bReset.y0+bReset.height {
			g.Reset()
			return nil
		}
	}

	// Game logic
	if !g.gameOver.Ended {
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
						g.gameOver.Ended = true
						g.records = true
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
		duration.Raw = time.Since(startTime)

		if !initialClick {
			duration.Minutes = int(duration.Raw.Minutes())
			duration.Seconds = int(duration.Raw.Seconds()) % 60
		}

	} else {
		if g.records {
			db, err := OpenDB()
			if err != nil {
				return err
			}

			// Save information on database
			if g.gameOver.Win {
				if popUp {
					// Getting the name of the player
					playerName = ebiten.AppendInputChars(playerName[:0])

					for _, r := range playerName {
						// Limiting the name to 10 characters
						if len(Name) < 10 {
							Name += string(r)
						}
					}
					// Delete character (backspace)
					if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(Name) > 0 {
						Name = Name[:len(Name)-1]
					}

					// Confirm player name (Enter)
					if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(Name) > 0 {
						popUp = false
					}

					if !popUp {
						err := SaveBestTime(g, Name, duration, db)
						if err != nil {
							return err
						}

						err = SaveStats(g, db, true)
						if err != nil {
							return err
						}
						g.records = false
					}
				}

			} else if g.gameOver.Lose {
				err := SaveStats(g, db, false)
				if err != nil {
					return err
				}
				g.records = false
			}

			if err := CloseDB(db); err != nil {
				return err
			}

		}
	}

	if !g.gameOver.Ended {
		g.gameOver.Win = g.checkIfEndGame()
		if g.gameOver.Win {
			popUp = true
			g.gameOver.Ended = true
			g.records = true
		}
	}

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

		if g.gameOver.Ended {
			if g.gameOver.Win && g.records {
				g.NamePopUp(screen)
			} else if !g.records {
				g.ShowStats(screen)
			}
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

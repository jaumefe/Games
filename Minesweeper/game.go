package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type gameover struct {
	Win   bool
	Lose  bool
	Ended bool
}

type button struct {
	x0, y0, width, height int
}

var bEasy = button{x0: 420, y0: screenHeightMax/2 - 200/2, width: 300, height: 200}
var bMedium = button{x0: 770, y0: screenHeightMax/2 - 200/2, width: 300, height: 200}
var bHard = button{x0: 1120, y0: screenHeightMax/2 - 200/2, width: 300, height: 200}
var bReset = button{x0: screenWidthMax - 170, y0: screenHeightMax - 300, width: 150, height: 100}

func RowAndColToSingleArray(row, col, totalCols int) int {
	return col + row*totalCols
}

func SingleArrayToRowAndCol(idx, totalCols int) (row, col int) {
	row = int(idx / totalCols)
	col = idx - totalCols*row
	return
}

func printTime(screen *ebiten.Image) {
	timeStr := fmt.Sprintf("%02d:%02d", duration.Minutes, duration.Seconds)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(screenWidthMax-125, 100)
	text.Draw(screen, timeStr, counterFont.txt, opts)
}

func (g *Game) Init() {
	var rows, cols, bombs int

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

	g.reset = false
	g.records = false
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

func (g *Game) DifficultySelector(screen *ebiten.Image) {
	var msg message
	var opts msgOpts
	msg.title = "8 x 8"
	msg.subtitle = "10 mines"
	opts.font = minesFont
	opts.marginX = 50
	opts.marginY = 50
	opts.subMarginX = opts.marginX + 50
	vector.DrawFilledRect(screen, 0, 0, screenWidthMax, screenHeightMax, gray, false)
	drawButton(screen, float32(bEasy.x0), float32(bEasy.y0), float32(bEasy.width), float32(bEasy.height), msg, &opts)
	msg.title = "16x16"
	msg.subtitle = "40 mines"
	opts.marginX = 70
	opts.subMarginX = opts.marginX + 35
	drawButton(screen, float32(bMedium.x0), float32(bMedium.y0), float32(bMedium.width), float32(bMedium.height), msg, &opts)
	msg.title = "16x30"
	msg.subtitle = "99 mines"
	drawButton(screen, float32(bHard.x0), float32(bHard.y0), float32(bHard.width), float32(bHard.height), msg, &opts)
}

func (g *Game) ResetButton(screen *ebiten.Image) {
	var msg message
	var opts msgOpts
	msg.title = "New game"
	opts.font = counterFont
	opts.marginX = 50
	opts.marginY = float32(bReset.height)/2 - 35
	drawButton(screen, float32(bReset.x0), float32(bReset.y0), float32(bReset.width), float32(bReset.height), msg, &opts)
}

func drawButton(screen *ebiten.Image, x0, y0, width, height float32, msg message, mOpts *msgOpts) {
	var lineWidth float32 = 5
	marginX, marginY := mOpts.marginX, mOpts.marginY
	vector.DrawFilledRect(screen, x0, y0, width, height, nonClicked, false)
	x0i := x0 - lineWidth/2
	x0f := x0 + width + lineWidth/2
	vector.StrokeLine(screen, x0i, y0, x0f, y0, lineWidth, clicked, false)
	y0i := y0 - lineWidth/2
	y0f := y0 + height + lineWidth/2
	vector.StrokeLine(screen, x0, y0i, x0, y0f, lineWidth, clicked, false)
	y0i = y0 - lineWidth/2
	y0f = y0 + height + lineWidth/2
	vector.StrokeLine(screen, x0+width, y0i, x0+width, y0f, lineWidth, clicked, false)
	x0i = x0 - lineWidth/2
	x0f = x0 + width + lineWidth/2
	vector.StrokeLine(screen, x0i, y0+height, x0f, y0+height, lineWidth, clicked, false)

	font := mOpts.font
	txt := msg.title
	opts := &text.DrawOptions{}
	offsetX0 := float64(x0 + width/2 - marginX)
	offsetY0 := float64(y0 + height/2 - marginY)
	opts.GeoM.Translate(offsetX0, offsetY0)
	opts.ColorScale.SetR(0)
	opts.ColorScale.SetG(0)
	opts.ColorScale.SetB(0)
	text.Draw(screen, txt, font.txt, opts)

	if msg.subtitle != "" {
		marginX, marginY = mOpts.subMarginX, mOpts.subMarginY
		txt = msg.subtitle
		opts = &text.DrawOptions{}
		offsetX0 = float64(x0 + width/2 - marginX)
		offsetY0 = float64(y0 + height/2 - marginY)
		opts.GeoM.Translate(offsetX0, offsetY0)
		opts.ColorScale.SetR(0)
		opts.ColorScale.SetG(0)
		opts.ColorScale.SetB(0)
		text.Draw(screen, txt, font.txt, opts)
	}
}

func (g *Game) Reset() {
	g.reset = true
	g.diff = ""
	initialClick = true
	totalFlags = 0
	duration = Duration{Minutes: 0, Seconds: 0}
	g.gameOver.Win = false
	g.gameOver.Lose = false
	g.gameOver.Ended = false
	playerName = nil
	g.records = false
	Name = ""
}

func (g *Game) LogicDifficultySelector() {
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
}

func (g *Game) NamePopUp(screen *ebiten.Image) {
	popupWidth, popupHeight := 300, 150

	popup := ebiten.NewImage(popupWidth, popupHeight)
	popup.Fill(color.RGBA{50, 50, 50, 255})

	txt := "Enter your name: " + Name
	opts := &text.DrawOptions{}
	opts.ColorScale.SetR(255)
	opts.ColorScale.SetG(255)
	opts.ColorScale.SetB(255)
	opts.ColorScale.SetA(255)
	opts.GeoM.Translate(10, float64(popupHeight)/2-counterFont.txt.Metrics().HAscent/2)
	text.Draw(popup, txt, counterFont.txt, opts)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidthMax/2-float64(popupWidth)/2, screenHeightMax/2-float64(popupHeight)/2)
	screen.DrawImage(popup, op)

}

func (g *Game) ShowStats(screen *ebiten.Image) {
	popupWidth, popupHeight := 500, 300
	popup := ebiten.NewImage(popupWidth, popupHeight)
	popup.Fill(color.RGBA{100, 100, 100, 255})

	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer CloseDB(db)

	diff, err := DiffStringToId(g, db)
	if err != nil {
		log.Fatal(err)
	}

	bts, err := QueryBestTimes(db, diff)
	if err != nil {
		log.Fatal(err)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidthMax/2-float64(popupWidth)/2, screenHeightMax/2-float64(popupHeight)/2)

	for i := 0; i < len(bts); i++ {
		txt := fmt.Sprintf("%d.", i+1)
		opts := &text.DrawOptions{}
		opts.ColorScale.SetR(0)
		opts.ColorScale.SetG(0)
		opts.ColorScale.SetB(0)
		opts.ColorScale.SetA(255)
		opts.GeoM.Translate(50, 25+float64(i)*20)
		text.Draw(popup, txt, counterFont.txt, opts)

		txt = bts[i].PlayerName
		opts.GeoM.Translate(100, 0)
		text.Draw(popup, txt, counterFont.txt, opts)

		txt = fmt.Sprintf("%02d:%02d", bts[i].Time.Minutes, bts[i].Time.Seconds)
		opts.GeoM.Translate(200, 0)
		text.Draw(popup, txt, counterFont.txt, opts)
	}

	won, played, err := QueryStats(db, diff)
	if err != nil {
		log.Fatal(err)
	}

	percentage := float32(won) / float32(played) * 100
	txt := fmt.Sprintf("Games won: %d/%d (%.2f%%)", won, played, percentage)
	opts := &text.DrawOptions{}
	opts.ColorScale.SetR(0)
	opts.ColorScale.SetG(0)
	opts.ColorScale.SetB(0)
	opts.ColorScale.SetA(255)
	opts.GeoM.Translate(50, 250)
	text.Draw(popup, txt, counterFont.txt, opts)

	screen.DrawImage(popup, op)

}

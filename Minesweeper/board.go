package main

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jaumefe/stats"
)

type Board struct {
	rows  int
	cols  int
	tiles []tile
	bombs int
	coord [][]coord
}

type coord struct {
	x0, y0 float32
}

func (b *Board) Grid(board *ebiten.Image) {
	vector.DrawFilledRect(board, 0, 0, screenWidthMax, screenHeightMax, color.RGBA{R: 80, G: 80, B: 80, A: 0}, false)
	for r := range b.rows {
		for c := range b.cols {
			idx := RowAndColToSingleArray(r, c, b.cols)
			cds := b.coord[r][c]
			x0 := cds.x0
			y0 := cds.y0
			x1 := x0 + maxTileSize
			y1 := y0 + maxTileSize
			b.tiles[idx].SetTileColor()
			vector.DrawFilledRect(board, x0, y0, maxTileSize, maxTileSize, b.tiles[idx].colors, false)
			vector.StrokeLine(board, x0, y0, x1, y0, float32(1), color.Black, false)
			vector.StrokeLine(board, x0, y0, x0, y1, float32(1), color.Black, false)
			vector.StrokeLine(board, x1, y0, x1, y1, float32(1), color.Black, false)
			vector.StrokeLine(board, x0, y1, x1, y1, float32(1), color.Black, false)
		}
	}
}

func (b *Board) Init(rows, cols, bombs int) {
	b.rows, b.cols, b.bombs = rows, cols, bombs
	b.coord = make([][]coord, b.rows)
	length := b.rows * b.cols
	b.tiles = make([]tile, length)
	for y := 0; y < b.rows; y++ {
		b.coord[y] = make([]coord, b.cols)
		for x := 0; x < b.cols; x++ {
			x0 := float32(x * maxTileSize)
			y0 := float32(y * maxTileSize)
			b.coord[y][x].x0 = x0
			b.coord[y][x].y0 = y0
		}
	}
}

func (b *Board) CursorPositionToRowAndCol(x, y int) (row, col int) {
	row = int(y / maxTileSize)
	col = int(x / maxTileSize)
	return
}

func (b *Board) MineDistributionAfterFirstClick(idx int) {

	excludeIndex := LookForNeighborTiles(idx, b.rows, b.cols)

	nb := 0
	for r := 0; r < b.rows; r++ {
		for c := 0; c < b.cols; c++ {
			if nb < b.bombs {
				idx := RowAndColToSingleArray(r, c, b.cols)
				if !excludeIndex[idx] {
					b.tiles[idx].isMine = true
					nb++
				}
			}
		}
	}

	opts := &stats.ShuffleOptions{Seed: time.Now().UnixNano(), ExcludeIndices: excludeIndex}
	stats.FisherYatesShuffle(b.tiles, opts)
}

func (b *Board) countNeighborMinesSingleTile(row, col int) {
	idx := RowAndColToSingleArray(row, col, b.cols)
	neighborIdx := LookForNeighborTiles(idx, b.rows, b.cols)

	for i := range neighborIdx {
		if i != idx {
			if b.tiles[i].isMine {
				b.tiles[idx].nbhdMines++
			}
		}
	}

}

func (b *Board) CountNeighborMinesAllBoard() {
	for r := 0; r < b.rows; r++ {
		for c := 0; c < b.cols; c++ {
			b.countNeighborMinesSingleTile(r, c)
		}
	}
}

func (b *Board) ShowNeighborMines(screen *ebiten.Image) {
	for x := 0; x < len(b.coord); x++ {
		for y := 0; y < len(b.coord[x]); y++ {
			idx := RowAndColToSingleArray(x, y, b.cols)
			nbhdMinesStr := strconv.Itoa(b.tiles[idx].nbhdMines)
			if nbhdMinesStr != "0" && !b.tiles[idx].isMine && b.tiles[idx].isClicked {
				opts := &text.DrawOptions{}
				offsetX0 := float64(b.coord[x][y].x0 + 15)
				offsetY0 := float64(b.coord[x][y].y0)
				opts.GeoM.Translate(offsetX0, offsetY0)
				opts.ColorScale.SetR(0)
				opts.ColorScale.SetG(0)
				opts.ColorScale.SetB(0)
				text.Draw(screen, nbhdMinesStr, minesFont.txt, opts)
			}
		}
	}

}

func (b *Board) NoMinesAutoShower(idx int, exclude map[int]bool) {
	neighbors := LookForNeighborTiles(idx, b.rows, b.cols)
	for n := range neighbors {
		if !exclude[n] && !b.tiles[n].isMine {
			b.tiles[n].isClicked = true
			b.tiles[n].flag = false
			if b.tiles[n].nbhdMines == 0 {
				exclude[n] = true
				b.NoMinesAutoShower(n, exclude)
			}
		}
	}
}

func (b *Board) DrawFlag(board *ebiten.Image) {
	for x := 0; x < len(b.coord); x++ {
		for y := 0; y < len(b.coord[x]); y++ {
			idx := RowAndColToSingleArray(x, y, b.cols)
			if b.tiles[idx].flag {
				x0, y0 := b.coord[x][y].x0, b.coord[x][y].y0
				vector.StrokeLine(board, float32(x0+10), float32(y0+45), float32(x0+45), float32(y0+45), 5, color.Black, false)
				vector.StrokeLine(board, float32(x0)+27.5, float32(y0+45), float32(x0)+27.5, float32(y0+10), 5, color.Black, false)
				vector.DrawFilledRect(board, float32(x0+10), float32(y0+10), 19, 17.5, color.RGBA{R: 255, A: 255}, false)
			}
		}
	}

}

func (b *Board) ShowTotalFlags(screen *ebiten.Image) {
	msg := fmt.Sprintf("%02d/%d", totalFlags, b.bombs)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(screenWidthMax-125, 50)
	text.Draw(screen, msg, counterFont.txt, opts)
}

func (b *Board) PlotMines(screen *ebiten.Image) {
	for x := 0; x < len(b.coord); x++ {
		for y := 0; y < len(b.coord[x]); y++ {
			idx := RowAndColToSingleArray(x, y, b.cols)
			if b.tiles[idx].isMine && !b.tiles[idx].flag {
				x0, y0 := b.coord[x][y].x0, b.coord[x][y].y0
				vector.DrawFilledCircle(screen, x0+maxTileSize/2, y0+maxTileSize/2, 15, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)
				vector.StrokeLine(screen, x0+5, y0+maxTileSize/2, x0+50, y0+maxTileSize/2, 5, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)
				vector.StrokeLine(screen, x0+maxTileSize/2, y0+5, x0+maxTileSize/2, y0+50, 5, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)
			}
		}
	}
}

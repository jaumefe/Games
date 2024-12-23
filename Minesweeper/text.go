package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type fontFace struct {
	face font.Face
	txt  *text.GoXFace
}

type optsFontFace struct {
	size float64
	dpi  float64
}

var minesFont, counterFont fontFace

var (
	blue       = color.RGBA{R: 0, G: 170, B: 228, A: 255}
	green      = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	darkRed    = color.RGBA{R: 125, G: 0, B: 0, A: 255}
	darkBlue   = color.RGBA{R: 0, G: 80, B: 139, A: 255}
	brown      = color.RGBA{R: 225, G: 127, B: 50, A: 255}
	cyan       = color.RGBA{R: 0, G: 125, B: 125, A: 255}
	black      = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	gray       = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	nonClicked = color.RGBA{R: 120, G: 120, B: 120, A: 255}
	clicked    = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	red        = color.RGBA{R: 255, G: 0, B: 0, A: 255}
)

func (ff *fontFace) Init(tt *sfnt.Font, optsff *optsFontFace) error {
	var err error
	opts := &opentype.FaceOptions{
		Size:    optsff.size,
		DPI:     optsff.dpi,
		Hinting: font.HintingFull,
	}
	ff.face, err = opentype.NewFace(tt, opts)
	if err != nil {
		return err
	}

	ff.txt = text.NewGoXFace(ff.face)

	return nil
}

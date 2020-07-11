package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// layout constants
const (
	windowLeft       = 100
	windowTop        = 100
	gridLeft         = 15
	gridTop          = 81
	tileSide         = 20
	digitPanelWidth  = 65
	digitPanelHeight = 37
	timerTop         = 21
	flagTop          = 21
	faceTop          = 19
	faceWidth        = 42
	faceHeight       = 42
)

const (
	digitWidth           = 19
	digitHeight          = 33
	digitPanelHorzMargin = (digitPanelWidth - (3 * digitWidth)) / 4
	digitPanelVertMargin = (digitPanelHeight - digitHeight) / 2
)

const (
	beginnerLayoutHeight     = 276
	beginnerLayoutWidth      = 210
	intermediateLayoutHeight = 416
	intermediateLayoutWidth  = 350
	expertLayoutHeight       = 416
	expertLayoutWidth        = 630
)

const (
	beginnerDigitPanelOffset = 16
	defaultDigitPanelOffset  = 20
)

// Layout provides rect parameters based upon the game options
type Layout struct {
	options *Options
}

func (layout *Layout) timerDigitPanel() *sdl.Rect {
	left := layout.width() - layout.digitPanelOffset() - digitPanelWidth
	var top int32 = timerTop

	return layout.digitPanel(left, top)
}

func (layout *Layout) timerDigit(position int32) *sdl.Rect {
	left := layout.width() - layout.digitPanelOffset() - digitPanelWidth
	left += digitPanelHorzMargin*(position+1) + digitWidth*position
	top := int32(timerTop + digitPanelVertMargin)
	return &sdl.Rect{X: left, Y: top, W: digitWidth, H: digitHeight}
}

func (layout *Layout) flagDigit(position int32) *sdl.Rect {
	left := layout.digitPanelOffset()
	left += digitPanelHorzMargin*(position+1) + digitWidth*position
	top := int32(flagTop + digitPanelVertMargin)
	return &sdl.Rect{X: left, Y: top, W: digitWidth, H: digitHeight}
}

func (layout *Layout) grid() *sdl.Rect {
	width := int32(layout.options.columns * tileSide)
	height := int32(layout.options.rows * tileSide)

	return &sdl.Rect{X: gridLeft, Y: gridTop, W: width, H: height}
}

func (layout *Layout) digitPanel(left int32, top int32) *sdl.Rect {
	return &sdl.Rect{X: left, Y: top, W: digitPanelWidth, H: digitPanelHeight}
}

func (layout *Layout) tile(boundingBox *sdl.Rect, index int) *sdl.Rect {
	row, column := layout.options.rowAndColumn(index)
	left := boundingBox.X + int32(column*tileSide)
	top := boundingBox.Y + int32(row*tileSide)

	return &sdl.Rect{X: left, Y: top, W: tileSide, H: tileSide}
}

func (layout *Layout) height() int32 {
	switch layout.options.level {
	case "beginner":
		return beginnerLayoutHeight
	case "intermediate":
		return intermediateLayoutHeight
	case "expert":
		return expertLayoutHeight
	default:
		return beginnerLayoutHeight
	}
}

func (layout *Layout) rowsAndColumns() (int, int) {
	return layout.options.rows, layout.options.columns
}

func (layout *Layout) width() int32 {
	switch layout.options.level {
	case "beginner":
		return beginnerLayoutWidth
	case "intermediate":
		return intermediateLayoutWidth
	case "expert":
		return expertLayoutWidth
	default:
		return beginnerLayoutWidth
	}
}

func (layout *Layout) digitPanelOffset() int32 {
	switch layout.options.level {
	case "beginner":
		return beginnerDigitPanelOffset
	case "intermediate":
		return defaultDigitPanelOffset
	case "expert":
		return defaultDigitPanelOffset
	default:
		return defaultDigitPanelOffset
	}
}

func (layout *Layout) flagDigitPanel() *sdl.Rect {
	left := layout.digitPanelOffset()
	var top int32 = flagTop
	return layout.digitPanel(left, top)
}

func (layout *Layout) face() *sdl.Rect {
	left := layout.width()/2 - faceWidth/2
	return &sdl.Rect{X: left, Y: faceTop, W: faceWidth, H: faceHeight}
}

func windowRect(layout *Layout) sdl.Rect {
	width := layout.width()
	height := layout.height()

	return sdl.Rect{X: windowLeft, Y: windowTop, W: width, H: height}
}

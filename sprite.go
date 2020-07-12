package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Renderer interface {
	render(surface *sdl.Surface)
}

type MouseHandler interface {
	hitTest(event *sdl.MouseButtonEvent) bool
	handleEvent(event *sdl.MouseButtonEvent)
	//boundingBox() *sdl.Rect
}

type TileListener interface {
	reveal(mine bool, adjacentMines bool)
	clear()
	flag(flagged bool)
}

type BaseTileListener struct {
	TileListener
}

func (listener *BaseTileListener) reveal(mine bool, adjacentMines bool) {
}

func (listener *BaseTileListener) clear() {
}

func (listener *BaseTileListener) flag(flagged bool) {
}

type FlagStateListener interface {
	flagStateChanged(exhausted bool)
}

const (
	gameStateInit    = "init"
	gameStatePlaying = "playing"
	gameStateWin     = "win"
	gameStateLose    = "lose"
)

type GameStateListener interface {
	// valid strings are init, playing, won, lost
	gameStateChanged(state string)
}

type SpriteHandler interface {
	Renderer
	MouseHandler
}

type Sprite struct {
	SpriteHandler

	context *GameContext
}

// convience functions
func ptInRect(event *sdl.MouseButtonEvent, rect *sdl.Rect) bool {
	x, y := int32(event.X), int32(event.Y)
	inside := x >= rect.X && x < rect.X+rect.W && y >= rect.Y && y < rect.Y+rect.H
	return inside
}

func (sprite *Sprite) hitTest(event *sdl.MouseButtonEvent) bool {
	return false
}

func (sprite *Sprite) handleEvent(event *sdl.MouseButtonEvent) {
}

type Background struct {
	Sprite
}

func (panel *Background) render(surface *sdl.Surface) {
	name := "bg_" + panel.context.layout.options.level
	image := panel.context.imageRepo.imageForName(name)
	image.Blit(nil, surface, panel.boundingBox())
}

func (background *Background) boundingBox() *sdl.Rect {
	width := background.context.layout.width()
	height := background.context.layout.height()
	return &sdl.Rect{X: 0, Y: 0, W: width, H: height}
}

func createBackground(context *GameContext) *Background {
	background := new(Background)
	background.context = context

	return background
}

type Button struct {
	BaseTileListener
	Sprite

	state      string
	revealed   int
	listeners  []GameStateListener
	imageNames map[string]string
}

func createButton(context *GameContext) *Button {
	button := new(Button)
	button.context = context
	button.state = gameStateInit
	button.imageNames = map[string]string{
		gameStateInit: "face_playing", gameStatePlaying: "face_playing",
		gameStateWin: "face_win", gameStateLose: "face_lose"}

	return button
}

func (button *Button) boundingBox() *sdl.Rect {
	return button.context.layout.face()
}

func (button *Button) render(surface *sdl.Surface) {
	name := button.imageNames[button.state]
	image := button.context.imageRepo.imageForName(name)
	image.Blit(nil, surface, button.boundingBox())
}

func (button *Button) reveal(mine bool, adjacentMines bool) {
	if mine {
		button.state = gameStateLose
		button.notifyListeners()
	} else {
		if button.state == gameStateInit {
			button.state = gameStatePlaying
			button.notifyListeners()
		}
		button.revealed++
		if button.revealed == button.context.layout.options.blanks() {
			button.state = gameStateWin
			button.notifyListeners()
		}
	}
}

func (button *Button) notifyListeners() {
	for _, listener := range button.listeners {
		listener.gameStateChanged(button.state)
	}
}

func (button *Button) hitTest(event *sdl.MouseButtonEvent) bool {
	return ptInRect(event, button.boundingBox())
}

func (button *Button) handleEvent(event *sdl.MouseButtonEvent) {
	if event.Button == sdl.BUTTON_LEFT {
		button.state = gameStateInit
		button.revealed = 0
		button.notifyListeners()
	}
}

type Grid struct {
	GameStateListener
	FlagStateListener

	Sprite

	tiles     []*Tile
	minefield *Minefield
}

func createGrid(context *GameContext) *Grid {
	grid := new(Grid)
	grid.context = context
	boundingBox := context.layout.grid()
	grid.minefield = createMinefield(context)
	tiles := context.layout.options.tiles()
	for i := 0; i < tiles; i++ {
		tile := createTile(context)
		grid.tiles = append(grid.tiles, tile)
	}

	for index, tile := range grid.tiles {
		adjacentMines := grid.minefield.adjacentMines(index)
		isMine := grid.minefield.mineAt(index)
		rect := context.layout.tile(boundingBox, index)
		tile.boundingBox = rect
		tile.isMine = isMine
		tile.adjacentMines = adjacentMines
	}
	return grid
}

func (panel *Grid) render(surface *sdl.Surface) {
	for _, tile := range panel.tiles {
		tile.render(surface)
	}
}

func (grid *Grid) boundingBox() *sdl.Rect {
	return grid.context.layout.grid()
}

func (grid *Grid) hitTest(event *sdl.MouseButtonEvent) bool {
	return ptInRect(event, grid.boundingBox())
}

func (grid *Grid) gameStateChanged(state string) {
	if state == gameStateInit {
		grid.minefield.reset()
		for index, tile := range grid.tiles {
			isMine := grid.minefield.mineAt(index)
			adjacentMines := grid.minefield.adjacentMines(index)
			tile.reset(isMine, adjacentMines)
		}
	}
	for _, tile := range grid.tiles {
		tile.gameStateChanged(state)
	}
}

func (grid *Grid) flagStateChanged(exhausted bool) {
	for _, tile := range grid.tiles {
		tile.flagStateChanged(exhausted)
	}
}

func (grid *Grid) handleEvent(event *sdl.MouseButtonEvent) {
	box := grid.boundingBox()
	column := (event.X - box.X) / tileSide
	row := (event.Y - box.Y) / tileSide
	index := grid.context.layout.options.index(int(row), int(column))
	grid.tiles[index].handleEvent(event)
}

func (grid *Grid) setListeners(topLevelListeners []TileListener) {
	options := grid.context.layout.options
	// build a mesh of adjacent tile listeners
	for index, tile := range grid.tiles {
		listeners := topLevelListeners
		fn := func(r int, c int) {
			index := options.index(r, c)
			tile := grid.tiles[index]
			listeners = append(listeners, tile)
		}
		options.forEachNeighbor(index, fn)
		tile.listeners = listeners
	}
}

type Tile struct {
	TileListener
	GameStateListener
	FlagStateListener
	Sprite

	isMine        bool
	adjacentMines int
	adjacentFlags int
	flagged       bool
	revealed      bool
	gameOver      bool
	flagRemaining bool
	boundingBox   *sdl.Rect
	listeners     []TileListener
}

func createTile(context *GameContext) *Tile {
	tile := new(Tile)
	tile.context = context
	tile.flagRemaining = true
	return tile
}

func (tile *Tile) reveal(hasMine bool, hasAdjacentMines bool) {
	if !hasMine && !hasAdjacentMines {
		tile.tryReveal()
	}
}

func (tile *Tile) clear() {
	tile.tryReveal()
}

func (tile *Tile) flag(isFlagged bool) {
	if isFlagged {
		tile.adjacentFlags++
	} else {
		tile.adjacentFlags--
	}
}

func (tile *Tile) gameStateChanged(state string) {
	switch state {
	case gameStateInit:
		tile.flagged = false
		tile.revealed = false
		tile.adjacentFlags = 0
		tile.gameOver = false
		tile.flagRemaining = true
	case gameStateLose, gameStateWin:
		tile.gameOver = true
	}
}

func (tile *Tile) flagStateChanged(exhausted bool) {
	tile.flagRemaining = !exhausted
}

func (tile *Tile) render(surface *sdl.Surface) {
	if tile.revealed {
		if tile.isMine {
			image := tile.context.imageRepo.imageForName("tile_mine")
			image.BlitScaled(nil, surface, tile.boundingBox)
		} else {
			image := tile.context.imageRepo.imageForRevealedTile(tile.adjacentMines)
			image.BlitScaled(nil, surface, tile.boundingBox)
		}
	} else if tile.flagged {
		image := tile.context.imageRepo.imageForName("tile_flag")
		image.BlitScaled(nil, surface, tile.boundingBox)
	} else {
		image := tile.context.imageRepo.imageForName("tile")
		image.BlitScaled(nil, surface, tile.boundingBox)
	}
}

func (tile *Tile) handleEvent(event *sdl.MouseButtonEvent) {
	if event.Button == sdl.BUTTON_LEFT {
		if tile.revealed {
			tile.tryClear()
		} else {
			tile.tryReveal()
		}
	} else if event.Button == sdl.BUTTON_RIGHT {
		tile.tryToggleFlag()
	}
}

func (tile *Tile) tryClear() {
	if tile.adjacentFlags == tile.adjacentMines {
		for _, listener := range tile.listeners {
			listener.clear()
		}
	}
}

func (tile *Tile) tryReveal() {
	if tile.gameOver || tile.flagged || tile.revealed {
		return
	}
	tile.revealed = true
	for _, listener := range tile.listeners {
		listener.reveal(tile.isMine, tile.adjacentMines > 0)
	}
}

func (tile *Tile) tryToggleFlag() {
	if tile.gameOver || tile.revealed {
		return
	}
	if !tile.flagged && !tile.flagRemaining {
		return
	}
	tile.flagged = !tile.flagged
	for _, listener := range tile.listeners {
		listener.flag(tile.flagged)
	}
}

func (tile *Tile) reset(isMine bool, adjacentMines int) {
	tile.isMine = isMine
	tile.adjacentMines = adjacentMines
}

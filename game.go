package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Runner interface {
	run()
}

type GameContext struct {
	layout    *Layout
	imageRepo *ImageRepo
}

func (context *GameContext) withLayout(layout *Layout) *GameContext {
	context.layout = layout
	return context
}

func (context *GameContext) withImageRepo(imageRepo *ImageRepo) *GameContext {
	context.imageRepo = imageRepo
	return context
}

type Game struct {
	Runner

	context *GameContext
	window  *sdl.Window
	sprites []SpriteHandler
}

func (game *Game) addSprite(sprite SpriteHandler) {
	game.sprites = append(game.sprites, sprite)
}

func CreateGame(context *GameContext, window *sdl.Window) *Game {
	game := new(Game)
	game.context = context
	game.window = window

	// add sprites, building up the z-order from back to front
	background := createBackground(context)
	game.addSprite(background)

	timer := createTimer(context)
	game.addSprite(timer)

	flagCounter := createFlagCounter(context)
	game.addSprite(flagCounter)

	button := createButton(context)
	game.addSprite(button)

	grid := createGrid(context)
	game.addSprite(grid)

	// wire in listeners
	gameStateListeners := []GameStateListener{grid, timer, flagCounter}
	button.listeners = gameStateListeners

	tileListeners := []TileListener{button, flagCounter}
	grid.setListeners(tileListeners)

	flagStateListeners := []FlagStateListener{grid}
	flagCounter.flagStateListeners = flagStateListeners
	return game
}

func (game *Game) run() {
	game.render()
	running := true
	for running {
		event := sdl.WaitEventTimeout(100)
		if event != nil {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					game.handleEvent(t)
					game.render()
				}
			}
		} else {
			// render on timeout
			game.render()
		}
	}
}

func (game *Game) render() {
	surface, err := game.window.GetSurface()
	if err != nil {
		panic(err)
	}

	for _, renderer := range game.sprites {
		renderer.render(surface)
	}
	game.window.UpdateSurface()
}

func (game *Game) handleEvent(mouseEvent *sdl.MouseButtonEvent) {
	for _, mouseHandler := range game.sprites {
		if mouseHandler.hitTest(mouseEvent) {
			mouseHandler.handleEvent(mouseEvent)
		}
	}
}

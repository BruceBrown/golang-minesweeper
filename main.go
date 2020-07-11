package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	options := options()
	layout := Layout{options}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	rect := windowRect(&layout)
	window, err := sdl.CreateWindow("Minesweeper", rect.X, rect.Y, rect.W, rect.H, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	imageRepo := CreateImageRepo("images/")

	context := &GameContext{}
	context = context.withLayout(&layout).withImageRepo(imageRepo).withSurface(surface)
	game := CreateGame(context, window)
	game.run()
}

package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	// parse command line options and construct the layout based upon them
	options := options()
	layout := Layout{options}

	// initialize the UI
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// create the window and make it visible
	rect := windowRect(&layout)
	window, err := sdl.CreateWindow("Minesweeper", rect.X, rect.Y, rect.W, rect.H, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// initialize the image repo. Should find a way to imclude these in the binary as a resource
	imageRepo := CreateImageRepo("minesweeper/images")

	// build up the game context, which gets injected into a lot of places
	context := new(GameContext)
	context = context.withLayout(&layout).withImageRepo(imageRepo)
	game := CreateGame(context, window)

	// go run the game...
	game.run()
}

// This is how we'd like to build the images into the binary...
//go:embedglob Assets images/*.bmp

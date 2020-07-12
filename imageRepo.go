package main

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type ImageRepo struct {
	box               *packr.Box
	images            map[string]*sdl.Surface
	adjacentTileNames []string
	digitNames        []string
}

func CreateImageRepo(folder string) *ImageRepo {
	repo := new(ImageRepo)
	repo.box = packr.New("myBox", folder)
	repo.images = make(map[string]*sdl.Surface)
	repo.adjacentTileNames = []string{
		"tile_none", "tile_one", "tile_two", "tile_three", "tile_four",
		"tile_five", "tile_six", "tile_seven", "tile_eight"}
	repo.digitNames = []string{
		"digit_zero", "digit_one", "digit_two", "digit_three", "digit_four",
		"digit_five", "digit_six", "digit_seven", "digit_eight", "digit_nine"}

	return repo
}

func (repo *ImageRepo) imageForName(name string) *sdl.Surface {
	image := repo.images[name]
	if image == nil {
		image = load(repo, name)
		repo.images[name] = image
	}
	return image
}

func (repo *ImageRepo) imageForRevealedTile(adjacentMines int) *sdl.Surface {
	return repo.imageForName(repo.adjacentTileNames[adjacentMines])
}

func (repo *ImageRepo) imageForDigit(digit int) *sdl.Surface {
	return repo.imageForName(repo.digitNames[digit])
}

func load(repo *ImageRepo, name string) *sdl.Surface {
	/*
		file := repo.folder
		file += "minesweeper_"
		file += name
		file += ".bmp"
		image, err := img.Load(file)
	*/
	file := "images/minesweeper_" + name + ".png"

	bytes, err := repo.box.Find(file)
	if err != nil {
		panic(err)
	}

	ops, err := sdl.RWFromMem(bytes)
	if err != nil {
		panic(err)
	}

	image, err := img.LoadPNGRW(ops)
	if err != nil {
		panic(err)
	}

	return image
}

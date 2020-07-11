package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type FlagCounter struct {
	GameStateListener

	BaseTileListener
	Sprite

	flagStateListeners []FlagStateListener
	flags              int
}

func createFlagCounter(context *GameContext) *FlagCounter {
	flagCounter := new(FlagCounter)
	flagCounter.context = context
	flagCounter.flags = context.layout.options.mines
	return flagCounter
}

func (flagCounter *FlagCounter) hitTest(event *sdl.MouseButtonEvent) bool {
	return false
}

func (flagCounter *FlagCounter) render(surface *sdl.Surface) {
	name := "digit_panel"
	image := flagCounter.context.imageRepo.imageForName(name)
	boundingBox := flagCounter.context.layout.flagDigitPanel()
	image.Blit(nil, surface, boundingBox)

	onesDigit := flagCounter.flags % 10
	tensDigit := flagCounter.flags / 10 % 10
	hundredsDigit := flagCounter.flags / 100 % 10

	rect := flagCounter.context.layout.flagDigit(2)
	image = flagCounter.context.imageRepo.imageForDigit(onesDigit)
	image.Blit(nil, surface, rect)

	rect = flagCounter.context.layout.flagDigit(1)
	image = flagCounter.context.imageRepo.imageForDigit(tensDigit)
	image.Blit(nil, surface, rect)

	rect = flagCounter.context.layout.flagDigit(0)
	image = flagCounter.context.imageRepo.imageForDigit(hundredsDigit)
	image.Blit(nil, surface, rect)
}

func (flagCounter *FlagCounter) gameStateChanged(state string) {
	if state == gameStateInit {
		flagCounter.flags = flagCounter.context.layout.options.mines
	}
}

func (flagCounter *FlagCounter) flag(flagged bool) {
	if flagged {
		flagCounter.flags--
		if flagCounter.flags == 0 {
			flagCounter.notifyListeners(true)
		}
	} else {
		flagCounter.flags++
		if flagCounter.flags == 1 {
			flagCounter.notifyListeners(false)
		}
	}
}

func (flagCounter *FlagCounter) notifyListeners(exhausted bool) {
	for _, listener := range flagCounter.flagStateListeners {
		listener.flagStateChanged(exhausted)
	}
}

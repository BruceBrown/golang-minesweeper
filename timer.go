package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Timer struct {
	GameStateListener
	Sprite

	elapsed int
	running bool
	start   time.Time
}

func createTimer(context *GameContext) *Timer {
	timer := new(Timer)
	timer.context = context

	return timer
}

func (timer *Timer) gameStateChanged(state string) {
	switch state {
	case gameStateInit:
		timer.running = false
		timer.elapsed = 0
	case gameStatePlaying:
		timer.running = true
		timer.start = time.Now()
	case gameStateWin:
		timer.running = false
		timer.elapsed = int(time.Since(timer.start).Seconds())
	case gameStateLose:
		timer.running = false
		timer.elapsed = int(time.Since(timer.start).Seconds())
	}
}

func (timer *Timer) render(surface *sdl.Surface) {
	elapsed := timer.elapsed
	if timer.running {
		elapsed = int(time.Since(timer.start).Seconds())
	}
	name := "digit_panel"
	image := timer.context.imageRepo.imageForName(name)
	boundingBox := timer.context.layout.timerDigitPanel()
	image.Blit(nil, surface, boundingBox)

	onesDigit := elapsed % 10
	tensDigit := elapsed / 10 % 10
	hundredsDigit := elapsed / 100 % 10

	rect := timer.context.layout.timerDigit(2)
	image = timer.context.imageRepo.imageForDigit(onesDigit)
	image.Blit(nil, surface, rect)

	rect = timer.context.layout.timerDigit(1)
	image = timer.context.imageRepo.imageForDigit(tensDigit)
	image.Blit(nil, surface, rect)

	rect = timer.context.layout.timerDigit(0)
	image = timer.context.imageRepo.imageForDigit(hundredsDigit)
	image.Blit(nil, surface, rect)
}

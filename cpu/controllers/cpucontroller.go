package controllers

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/go-sdl-lib/game"
	"github.com/hculpan/kcpu/cpu/executor"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/hculpan/kcpu/cpu/pages"
	"github.com/veandco/go-sdl2/sdl"
)

type CpuController struct {
	game.GameController
}

func NewCpuController(gameWidth, gameHeight int32, config executor.CpuConfig) CpuController {
	result := CpuController{}

	windowBackground := sdl.Color{R: 0, G: 0, B: 0, A: 0}

	result.Game = model.NewCpuGame(gameWidth, gameHeight, 0.1, config)
	result.Window = component.NewWindow(gameWidth, gameHeight, "KCPU", windowBackground)

	result.RegisterPages()

	return result
}

func (c *CpuController) RegisterPages() {
	component.RegisterPage(pages.NewStartPage("StartPage", 0, 0, c.Window.Width, c.Window.Height))
	component.RegisterPage(pages.NewMainPage("MainPage", 0, 0, c.Window.Width, c.Window.Height))
	component.RegisterPage(pages.NewErrorPage("ErrorPage", 0, 0, c.Window.Width, c.Window.Height))
}

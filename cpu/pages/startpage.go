package pages

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/veandco/go-sdl2/sdl"
)

type StartPage struct {
	component.BasePage
}

func NewStartPage(name string, x, y, width, height int32) *StartPage {
	p := StartPage{}
	p.Initialize()
	p.Name = "StartPage"
	p.SetPosition(0, 0)
	p.SetSize(width, height)

	yloc := height / 6
	config := component.NewLabelConfig()
	config.FontSize = 64
	config.Justify = component.JUSTIFY_CENTER
	config.Color = sdl.Color{R: 89, G: 183, B: 56, A: 0}
	p.AddChild(component.NewLabelComponentWithConfig(0, yloc, width, 60, func() string {
		return "Start KCPU"
	}, config))
	p.AddChild(component.NewButtonComponent(
		width/2-width/8-width/16,
		yloc+200,
		width/8,
		75,
		"Run",
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		sdl.Color{R: 89, G: 183, B: 56, A: 255},
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		func() {
			model.Game.SingleStepStatus = START_NORMAL
			component.SwitchPage("MainPage")
		},
	))
	p.AddChild(component.NewButtonComponent(
		width/2+width/16,
		yloc+200,
		width/8,
		75,
		"Debug",
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		sdl.Color{R: 89, G: 183, B: 56, A: 255},
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		func() {
			model.Game.SingleStepStatus = START_SINGLE_STEP
			component.SwitchPage("MainPage")
		},
	))

	return &p
}

func (m *StartPage) KeyEvent(event *sdl.KeyboardEvent) bool {
	keycode := sdl.GetKeyFromScancode(event.Keysym.Scancode)
	if keycode == sdl.K_r {
		model.Game.Reset()
		return true
	}

	return component.PassKeyEventToChildren(event, m.Children)
}

func (m *StartPage) Draw(r *sdl.Renderer) error {
	return component.DrawParentAndChildren(r, m)
}

func (m *StartPage) MouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	if m.IsPointInComponent(event.X, event.Y) {
		return component.PassMouseButtonEventToChildren(event, m.Children)
	}

	return false
}

package pages

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/kcpu/cpu/components"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/veandco/go-sdl2/sdl"
)

type MainPage struct {
	component.BasePage
}

func NewMainPage(name string, x, y, width, height int32) *MainPage {
	p := MainPage{}
	p.Name = "MainPage"
	p.SetPosition(0, 0)
	p.SetSize(width, height)

	p.AddChild(components.NewTerminalComponent(0, 0, width, height))
	//	p.AddChild(components.NewHeaderComponent(0, 0, width, 40))

	return &p
}

func (m *MainPage) KeyEvent(event *sdl.KeyboardEvent) bool {
	keycode := sdl.GetKeyFromScancode(event.Keysym.Scancode)
	if keycode == sdl.K_r {
		model.Game.Reset()
		return true
	}

	return component.PassKeyEventToChildren(event, m.Children)
}

func (m *MainPage) Draw(r *sdl.Renderer) error {
	return component.DrawParentAndChildren(r, m)
}

func (m *MainPage) MouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	if m.IsPointInComponent(event.X, event.Y) {
		return component.PassMouseButtonEventToChildren(event, m.Children)
	}

	return false
}

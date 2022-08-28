package pages

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/veandco/go-sdl2/sdl"
)

type ErrorPage struct {
	component.BasePage
}

func NewErrorPage(name string, x, y, width, height int32) *ErrorPage {
	p := ErrorPage{}
	p.Initialize()
	p.Name = "ErrorPage"
	p.SetPosition(0, 0)
	p.SetSize(width, height)

	yloc := height / 6
	config := component.NewLabelConfig()
	config.FontSize = 64
	config.Justify = component.JUSTIFY_CENTER
	config.Color = sdl.Color{R: 255, G: 0, B: 0, A: 0}
	p.AddChild(component.NewLabelComponentWithConfig(0, yloc, width, 60, func() string {
		return "CPU ERROR"
	}, config))
	p.AddChild(component.NewLabelComponentWithConfig(0, yloc+200, width, 60, func() string {
		if model.Game.Error != nil {
			return model.Game.Error.Error()
		} else {
			return "unknown"
		}
	}, config))

	return &p
}

func (m *ErrorPage) KeyEvent(event *sdl.KeyboardEvent) bool {
	keycode := sdl.GetKeyFromScancode(event.Keysym.Scancode)
	if keycode == sdl.K_r {
		model.Game.Reset()
		return true
	}

	return component.PassKeyEventToChildren(event, m.Children)
}

func (m *ErrorPage) Draw(r *sdl.Renderer) error {
	return component.DrawParentAndChildren(r, m)
}

func (m *ErrorPage) MouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	if m.IsPointInComponent(event.X, event.Y) {
		return component.PassMouseButtonEventToChildren(event, m.Children)
	}

	return false
}

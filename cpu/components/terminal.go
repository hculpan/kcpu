package components

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/go-sdl-lib/resources"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/veandco/go-sdl2/sdl"
)

type TerminalComponent struct {
	component.BaseComponent
}

func NewTerminalComponent(x, y, width, height int32) *TerminalComponent {
	result := &TerminalComponent{}

	result.SetPosition(x, y)
	result.SetSize(width, height)

	return result
}

func (c *TerminalComponent) DrawComponent(r *sdl.Renderer) error {
	for i := 0; i < 24; i++ {
		bytes := model.Game.Cpu.GetVideoCharacterLine(i)
		text, err := resources.Fonts.CreateTexture(string(bytes), sdl.Color{R: 50, G: 255, B: 50, A: 255}, "CourierNew-24", r)
		if err != nil {
			return err
		}
		_, _, w, h, err := text.Query()
		if err != nil {
			return err
		}
		r.Copy(text, &sdl.Rect{X: 0, Y: 0, W: w, H: h}, &sdl.Rect{X: c.X + 5, Y: c.Y + int32(i*60), W: int32(w), H: int32(h)})
	}

	return nil
}

func (c *TerminalComponent) Draw(r *sdl.Renderer) error {
	if err := component.DrawParentAndChildren(r, c); err != nil {
		return err
	}

	return nil
}

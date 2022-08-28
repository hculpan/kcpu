package components

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/veandco/go-sdl2/sdl"
)

type LeftBarComponent struct {
	component.BaseComponent
}

func NewLeftBarComponent(x, y, width, height int32) *LeftBarComponent {
	result := &LeftBarComponent{}
	result.Initialize()

	result.SetPosition(x, y)
	result.SetSize(width, height)

	return result
}

func (c *LeftBarComponent) DrawComponent(r *sdl.Renderer) error {
	r.SetDrawColor(0, 0, 0, 255)
	rect := sdl.Rect{X: c.X, Y: c.Y, W: c.Width, H: c.Height}
	r.FillRect(&rect)

	return nil
}

func (c *LeftBarComponent) Draw(r *sdl.Renderer) error {
	if err := component.DrawParentAndChildren(r, c); err != nil {
		return err
	}

	return nil
}

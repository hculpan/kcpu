package components

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/veandco/go-sdl2/sdl"
)

type BottomBarComponent struct {
	component.BaseComponent

	CpuStatus component.Component
}

func NewBottomBar(x, y, width, height int32) *BottomBarComponent {
	result := &BottomBarComponent{}
	result.Initialize()

	result.SetPosition(x, y)
	result.SetSize(width, height)

	result.CpuStatus = NewCpuStatusComponent(x+10, y+10, width-20, height-20)
	result.AddChild(result.CpuStatus)

	return result
}

func (c *BottomBarComponent) DrawComponent(r *sdl.Renderer) error {
	r.SetDrawColor(50, 50, 50, 255)
	rect := sdl.Rect{X: c.X, Y: c.Y, W: c.Width, H: c.Height}
	r.FillRect(&rect)

	//	r.SetDrawColor(0, 0, 0, 255)
	//	r.DrawLine(c.X, c.Y, c.Width, c.Y)

	return nil
}

func (c *BottomBarComponent) Draw(r *sdl.Renderer) error {
	c.CpuStatus.SetVisible(model.Game.SingleStepStatus == model.SSTEP_WAITING)

	if err := component.DrawParentAndChildren(r, c); err != nil {
		return err
	}

	return nil
}

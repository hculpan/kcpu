package components

import (
	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/veandco/go-sdl2/sdl"
)

type RightBarComponent struct {
	component.BaseComponent

	StepButton component.Component
	RunButton  component.Component
}

func NewRightBarComponent(x, y, width, height int32) *RightBarComponent {
	result := &RightBarComponent{}
	result.Initialize()

	result.SetPosition(x, y)
	result.SetSize(width, height)

	result.StepButton = component.NewButtonComponent(
		x+width/2-90,
		y+height-300,
		180,
		65,
		"Step",
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		sdl.Color{R: 89, G: 183, B: 56, A: 255},
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		func() {
			model.Game.SingleStepStatus = model.SSTEP_RUN
		},
	)

	result.RunButton = component.NewButtonComponent(
		x+width/2-90,
		y+height-225,
		180,
		65,
		"Run",
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		sdl.Color{R: 89, G: 183, B: 56, A: 255},
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		func() {
			model.Game.SingleStepStatus = model.SSTEP_OFF
		},
	)

	result.AddChild(component.NewButtonComponent(
		x+width/2-90,
		y+height-150,
		180,
		65,
		"Reset",
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		sdl.Color{R: 89, G: 183, B: 56, A: 255},
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		func() {
			model.Game.Reset()
		},
	))

	result.AddChild(component.NewButtonComponent(
		x+width/2-90,
		y+height-75,
		180,
		65,
		"Reset Step",
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		sdl.Color{R: 89, G: 183, B: 56, A: 255},
		sdl.Color{R: 0, G: 0, B: 0, A: 255},
		func() {
			model.Game.SingleStepStatus = model.SSTEP_WAITING
			model.Game.Reset()
		},
	))

	result.AddChild(result.StepButton)
	result.AddChild(result.RunButton)

	return result
}

func (c *RightBarComponent) DrawComponent(r *sdl.Renderer) error {
	r.SetDrawColor(50, 50, 50, 255)
	rect := sdl.Rect{X: c.X, Y: c.Y, W: c.Width, H: c.Height}
	r.FillRect(&rect)

	return nil
}

func (c *RightBarComponent) Draw(r *sdl.Renderer) error {
	if model.Game.SingleStepStatus != model.SSTEP_WAITING {
		c.StepButton.SetVisible(false)
		c.RunButton.SetVisible(false)
	} else {
		c.StepButton.SetVisible(true)
		c.RunButton.SetVisible(true)
	}

	if err := component.DrawParentAndChildren(r, c); err != nil {
		return err
	}

	return nil
}

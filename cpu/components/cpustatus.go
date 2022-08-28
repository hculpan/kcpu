package components

import (
	"strings"

	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/kcpu/cpu/model"
	"github.com/veandco/go-sdl2/sdl"
)

type CpuStatusComponent struct {
	component.BaseComponent
}

func NewCpuStatusComponent(x, y, width, height int32) *CpuStatusComponent {
	result := &CpuStatusComponent{}
	result.Initialize()

	result.SetPosition(x, y)
	result.SetSize(width, height)

	config := component.NewLabelConfig()
	config.FontFile = "resources/fonts/Courier-New-Regular.ttf"
	config.FontName = "CourierNew"
	config.FontSize = 24
	config.Justify = component.JUSTIFY_LEFT
	config.Color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	result.AddChild(component.NewLabelComponentWithConfig(x+10, y+10, width-20, 35, func() string {
		lines := strings.Split(model.Game.Cpu.ToString(), "\n")
		//		fmt.Println(lines[0])
		return lines[0]
	}, config))
	result.AddChild(component.NewLabelComponentWithConfig(x+10, y+45, width-20, 35, func() string {
		lines := strings.Split(model.Game.Cpu.ToString(), "\n")
		return strings.Trim(lines[1], " ")
	}, config))

	return result
}

func (c *CpuStatusComponent) DrawComponent(r *sdl.Renderer) error {
	r.SetDrawColor(0, 0, 0, 255)
	rect := sdl.Rect{X: c.X, Y: c.Y, W: c.Width, H: c.Height}
	r.FillRect(&rect)

	return nil
}

func (c *CpuStatusComponent) Draw(r *sdl.Renderer) error {
	if err := component.DrawParentAndChildren(r, c); err != nil {
		return err
	}

	return nil
}

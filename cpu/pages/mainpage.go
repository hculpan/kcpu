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

/*
type MainPage struct {
	Cycle int
}

func NewMainPage() *MainPage {
	return &MainPage{Cycle: 0}
}

func (m *MainPage) Draw(screen *ebiten.Image, g *common.Game) {
	const border = 20

	// Draw the sample text
	for i := 0; i < 24; i++ {
		bytes := g.Cpu.GetVideoCharacterLine(i)
		text.Draw(screen, string(bytes), g.NormalFont, border, (26*i)+(border*2), color.RGBA{R: 50, G: 200, B: 50, A: 255})
	}
}

func (m *MainPage) Update(g *common.Game) error {
	if m.Cycle == 0 {
		g.Cpu.SoftReset()
	}

	for i := 0; i < 5000; i++ {
		err := g.Cpu.ExecuteSingle()
		if err != nil {
			log.Fatal(err)
		}
		m.Cycle++
	}

	return nil
}
*/

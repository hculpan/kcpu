package pages

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hculpan/kcpu/cpu/common"
)

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
	/*	for i := 0; i < 24; i++ {
		g.Cpu.CursorToPos(30, i)
		g.Cpu.SetCharacterAtCursor('H')
		g.Cpu.SetCharacterAtCursor('e')
		g.Cpu.SetCharacterAtCursor('l')
		g.Cpu.SetCharacterAtCursor('l')
		g.Cpu.SetCharacterAtCursor('o')
		g.Cpu.SetCharacterAtCursor(' ')
		g.Cpu.SetCharacterAtCursor('w')
		g.Cpu.SetCharacterAtCursor('o')
		g.Cpu.SetCharacterAtCursor('r')
		g.Cpu.SetCharacterAtCursor('l')
		g.Cpu.SetCharacterAtCursor('d')
	}*/

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

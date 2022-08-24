package common

import (
	"bufio"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hculpan/kcpu/cpu/executor"
	"github.com/hculpan/kcpu/cpu/resources"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func NewGame(screenWidth int, screenHeight int, execFile string) *Game {
	tt, err := opentype.Parse(resources.Courier_ttf)
	if err != nil {
		log.Fatal(err)
	}

	c := executor.NewCpu()
	result := &Game{
		Cpu:          &c,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
	if err := result.loadProgram(execFile); err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	result.NormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return result
}

type Game struct {
	Cpu          *executor.Cpu
	ScreenWidth  int
	ScreenHeight int

	NormalFont font.Face
}

func (g *Game) loadProgram(f string) error {
	file, err := os.Open(f)

	if err != nil {
		return err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	// calculate the bytes size
	var size int64 = info.Size()
	bytes := make([]byte, size)

	// read into buffer
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		return err
	}

	copy(g.Cpu.Memory, bytes)

	return nil
}

func (g *Game) Update() error {
	return CurrentPage.Update(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	CurrentPage.Draw(screen, g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}

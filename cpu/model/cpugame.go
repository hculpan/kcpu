package model

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/kcpu/cpu/executor"
)

const (
	SSTEP_OFF = iota
	SSTEP_WAITING
	SSTEP_RUN
	CPU_HALTED
)

type CpuGame struct {
	component.BaseGame

	Cpu              *executor.Cpu
	Error            error
	SingleStepStatus int
}

var Game *CpuGame

func NewCpuGame(gameWidth, gameHeight int32, livingRatio float32, config executor.CpuConfig) *CpuGame {
	rand.Seed(time.Now().UnixNano())

	cpu := executor.NewCpu(config)
	Game = &CpuGame{
		Cpu: &cpu,
	}

	Game.Reset()

	Game.SingleStepStatus = CPU_HALTED
	Game.Initialize(gameWidth, gameHeight)
	if err := Game.loadProgram(config.ProgramFilename, config.StartingAddress); err != nil {
		log.Fatal(err)
	}

	return Game
}

func (g *CpuGame) Update() error {
	if g.SingleStepStatus == CPU_HALTED {
		return nil
	}

	if g.Cpu.Cycle == 0 {
		g.Cpu.SoftReset()
	}

	switch g.SingleStepStatus {
	case SSTEP_OFF:
		for i := 0; i < 5000; i++ {
			err := g.Cpu.ExecuteSingle()
			if err != nil {
				g.SingleStepStatus = CPU_HALTED
				g.Error = err
				component.SwitchPage("ErrorPage")
				break
			}
		}
	case SSTEP_WAITING:
		// do nothing
	case SSTEP_RUN:
		err := g.Cpu.ExecuteSingle()
		if err != nil {
			g.Error = err
			component.SwitchPage("ErrorPage")
		}
		g.SingleStepStatus = SSTEP_WAITING
	default:
		// do nothing
	}

	return nil
}

func (g *CpuGame) Reset() error {
	g.Cpu.SoftReset()
	//	g.SingleStepStatus = SSTEP_OFF
	return nil
}

func (g *CpuGame) loadProgram(f string, addr uint16) error {
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

	for i := 0; i < int(size); i++ {
		g.Cpu.Memory[i+int(addr)] = bytes[i]
	}

	return nil
}

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

type CpuGame struct {
	component.BaseGame

	Cpu *executor.Cpu
}

var Game *CpuGame

func NewCpuGame(gameWidth, gameHeight int32, livingRatio float32, programFilename string) *CpuGame {
	rand.Seed(time.Now().UnixNano())

	cpu := executor.NewCpu()
	Game = &CpuGame{
		Cpu: &cpu,
	}

	Game.Reset()

	Game.Initialize(gameWidth, gameHeight)
	if err := Game.loadProgram(programFilename); err != nil {
		log.Fatal(err)
	}

	return Game
}

func (g *CpuGame) Update() error {
	if g.Cpu.Cycle == 0 {
		g.Cpu.SoftReset()
	}

	for i := 0; i < 5000; i++ {
		err := g.Cpu.ExecuteSingle()
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (g *CpuGame) Reset() error {
	return nil
}

func (g *CpuGame) loadProgram(f string) error {
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

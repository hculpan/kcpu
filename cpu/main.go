package main

/*func main() {
	cmd.Execute()
}*/

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/go-sdl-lib/resources"
	"github.com/hculpan/kcpu/cpu/controllers"
	"github.com/hculpan/kcpu/cpu/executor"
)

//go:embed resources/fonts/*
var appFonts embed.FS

func main() {
	// variables declaration
	var programFilename string
	var startingAddr string

	// flags declaration using flag package
	flag.StringVar(&programFilename, "p", "", "KCPU file to execute")
	flag.StringVar(&startingAddr, "a", "", "Address to load and execute program")

	flag.Parse() // after declaring flags we need to call it

	// check if cli params match
	if len(programFilename) == 0 {
		fmt.Println("KCPU file to execute is required")
		return
	}

	config := executor.CpuConfig{
		ProgramFilename: programFilename,
	}
	config.StartingAddress = 0
	if len(startingAddr) > 0 {
		if strings.HasPrefix(startingAddr, "0x") {
			n, err := strconv.ParseInt(startingAddr[2:], 16, 32)
			if err != nil || n < 0 || n > 65535 {
				fmt.Printf("Invalid starting address: %s\n", startingAddr)
				return
			}
			config.StartingAddress = uint16(n)
		} else {
			n, err := strconv.Atoi(startingAddr)
			if err != nil || n < 0 || n > 65535 {
				fmt.Printf("Invalid starting address: %s\n", startingAddr)
				return
			}
			config.StartingAddress = uint16(n)
		}
	}

	component.SetupSDL()

	if err := resources.FontsInit(appFonts); err != nil {
		fmt.Println(err)
		return
	}

	if err := resources.Fonts.RegisterFont("TruenoLight-64", "built-in-fonts/TruenoLight.otf", 64); err != nil {
		log.Fatal(err)
	}
	if err := resources.Fonts.RegisterFont("CourierNew-64", "resources/fonts/Courier-New-Regular.ttf", 64); err != nil {
		log.Fatal(err)
	}

	// Since our cells are all 3 pixels with a 1 pixel barrier
	// around them, we want to make sure our widht/height is
	// a divisor of 4
	var gameWidth int32 = 867
	var gameHeight int32 = 806 // 706 terminal + 100 bottom bar

	gamecontroller := controllers.NewCpuController(gameWidth, gameHeight, config)
	if err := gamecontroller.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

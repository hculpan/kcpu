package main

/*func main() {
	cmd.Execute()
}*/

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/hculpan/go-sdl-lib/component"
	"github.com/hculpan/go-sdl-lib/resources"
	"github.com/hculpan/kcpu/cpu/controllers"
)

//go:embed resources/fonts/*

var appFonts embed.FS

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Must provide program to execute as second parameter")
		return
	}

	component.SetupSDL()

	if err := resources.FontsInit(appFonts); err != nil {
		fmt.Println(err)
		return
	}

	if err := resources.Fonts.RegisterFont("HackBold-24", "built-in-fonts/TruenoLight.otf", 24); err != nil {
		log.Fatal(err)
	}
	if err := resources.Fonts.RegisterFont("CourierNew-24", "resources/fonts/Courier-New-Regular.ttf", 24); err != nil {
		log.Fatal(err)
	}

	// Since our cells are all 3 pixels with a 1 pixel barrier
	// around them, we want to make sure our widht/height is
	// a divisor of 4
	var gameWidth int32 = 570
	var gameHeight int32 = 625

	gamecontroller := controllers.NewCpuController(gameWidth, gameHeight, os.Args[1])
	if err := gamecontroller.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

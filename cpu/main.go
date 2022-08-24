package main

/*func main() {
	cmd.Execute()
}*/

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hculpan/kcpu/cpu/common"
	"github.com/hculpan/kcpu/cpu/pages"
)

const (
	screenWidth  = 805
	screenHeight = (25 * 24) + 60
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Must specify a file to execute")
		return
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("KCPU")
	common.RegisterPage("main", pages.NewMainPage())
	common.SwitchPage("main")
	if err := ebiten.RunGame(common.NewGame(screenWidth, screenHeight, os.Args[1])); err != nil {
		log.Fatal(err)
	}
}

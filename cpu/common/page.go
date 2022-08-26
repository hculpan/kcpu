package common

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Page interface {
	Draw(screen *ebiten.Image, g *Game)
	Update(g *Game) error
}

var CurrentPage Page
var pagesMap map[string]Page = make(map[string]Page)

func RegisterPage(pageName string, page Page) {
	pagesMap[pageName] = page
}

func SwitchPage(pageName string) error {
	if v, ok := pagesMap[pageName]; ok {
		CurrentPage = v
		return nil
	}

	return fmt.Errorf("invalid page identifier '%s'", pageName)
}

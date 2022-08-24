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
var pages map[string]Page = make(map[string]Page)

func RegisterPage(pageName string, page Page) {
	pages[pageName] = page
}

func SwitchPage(pageName string) error {
	if v, ok := pages[pageName]; ok {
		CurrentPage = v
		return nil
	}

	return fmt.Errorf("invalid page identifier '%s'", pageName)
}

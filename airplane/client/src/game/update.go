package game

import (
	apcui "airplaneClient/src/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Update() error {
	apcui.GUIMgr.Update()

	g.debugMsg = ""
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.debugMsg = "Left"
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.debugMsg = "Right"
	}
	return nil
}

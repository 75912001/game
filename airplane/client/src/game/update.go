package game

import "github.com/hajimehoshi/ebiten"

func (g *Game) Update(screen *ebiten.Image) error {
	g.debugMsg = ""
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.debugMsg = "Left"
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.debugMsg = "Right"
	}
	return nil
}

package game

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

func (p *Game) Draw(screen *ebitenv2.Image) {
	p.user.Draw(screen)
}

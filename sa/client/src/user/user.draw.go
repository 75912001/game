package user

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/ui"
)

func (p *User) Draw(screen *ebitenv2.Image) {
	if p.role != nil {
		p.role.DrawAll(screen)
	}
	ui.GUIMgr.Draw(screen)
}

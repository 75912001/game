package user

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/ui"
)

func (p *User) Draw(screen *ebitenv2.Image) {
	// 填充灰色背景
	//screen.Fill(color.RGBA{
	//	R: 128,
	//	G: 128,
	//	B: 128,
	//	A: 255,
	//})

	p.role.Draw(screen)
	ui.GUIMgr.Draw(screen)
}

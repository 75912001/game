package game

import (
	"saClient/src/cfg"
	"saClient/src/ui"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

func (p *Game) Draw(screen *ebitenv2.Image) {
	p.user.Draw(screen)

	// 在屏幕右上角显示帧率
	fps := ebitenv2.ActualFPS()
	x := float32(cfg.GCommon.ScreenMaxWidth - 100)
	y := float32(10)
	ui.Printf(screen, x, y, "FPS: %.1f", fps)
}

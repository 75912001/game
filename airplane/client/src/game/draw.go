package game

import (
	apcui "airplaneClient/src/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

func (g *Game) Draw(screen *ebiten.Image) {
	// 填充灰色背景
	screen.Fill(color.RGBA{
		R: 128,
		G: 128,
		B: 128,
		A: 255,
	})

	apcui.GUIMgr.Draw(screen)

	// 显示调试信息
	if g.debugMsg != "" {
		ebitenutil.DebugPrint(screen, g.debugMsg)
	}
}

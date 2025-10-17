package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

func (g *Game) Draw(screen *ebiten.Image) {
	// 填充灰色背景
	_ = screen.Fill(color.RGBA{
		R: 128,
		G: 128,
		B: 128,
		A: 255,
	})

	// 显示调试信息
	if g.debugMsg != "" {
		_ = ebitenutil.DebugPrint(screen, g.debugMsg)
	}
}

package game

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"saClient/src/ui"
)

func (p *Game) Draw(screen *ebitenv2.Image) {
	// 填充绿色背景
	screen.Fill(color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})

	// 根据游戏状态绘制不同的内容
	switch p.state {
	case StateMenu: // 菜单状态下绘制 UI
		ui.GUIMgr.Draw(screen)
	case StateScene: // 场景状态,绘制场景背景等
	case StateBattling: // 战斗状态,绘制战斗相关内容
	case StateGameOver: // 游戏结束,显示游戏结束信息
		ui.Printf(screen, 280, 280, "*** 游戏结束 ***")
	}
	// 显示调试信息
	if p.debugMsg != "" {
		ui.Printf(screen, 0, 0, p.debugMsg)
	}
}

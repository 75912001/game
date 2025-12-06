package game

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"saClient/src/ui"
)

func (p *Game) Draw(screen *ebitenv2.Image) {
	// 填充草地一样的绿色背景
	screen.Fill(color.RGBA{
		R: 34,
		G: 139,
		B: 34,
		A: 255,
	})

	// 根据游戏状态绘制不同的内容
	switch p.state {
	case State_StartMenu: // 开始菜单状态下绘制 UI
		ui.GUIMgr.Draw(screen)
	case State_Scene: // 场景状态,绘制场景背景等
	case State_Battling: // 战斗状态,绘制战斗相关内容
	case State_GameOver: // 游戏结束,显示游戏结束信息
		ui.Printf(screen, 280, 280, "*** 游戏结束 ***")
	}
}

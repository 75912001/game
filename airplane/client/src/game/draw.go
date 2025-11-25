package game

import (
	"airplaneClient/src/ui"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

func (p *Game) Draw(screen *ebitenv2.Image) {
	// 填充灰色背景
	screen.Fill(color.RGBA{
		R: 128,
		G: 128,
		B: 128,
		A: 255,
	})

	// 根据游戏状态绘制不同的内容
	switch p.state {
	case StateMenu: // 菜单状态下绘制 UI
		ui.GUIMgr.Draw(screen)
	case StateBattling: // 游戏进行中,绘制飞机等游戏对象
		// 绘制-用户-子弹
		for _, bullet := range p.userBullets {
			bullet.Draw(screen)
		}
		// 绘制-用户-飞机
		if p.userPlane != nil {
			p.userPlane.Draw(screen)
		}
	case StatePaused: // 暂停状态,绘制游戏对象和暂停提示
		if p.userPlane != nil {
			p.userPlane.Draw(screen)
		}
		// 使用中文字体显示暂停信息
		ui.Printf(screen, 280, 280, "*** 游戏暂停 ***")
		ui.Printf(screen, 300, 310, "按 ESC 继续")
	case StateGameOver: // 游戏结束,显示游戏结束信息
		ui.Printf(screen, 280, 280, "*** 游戏结束 ***")
	}
	// 显示调试信息
	if p.debugMsg != "" {
		ui.Printf(screen, 0, 0, p.debugMsg)
	}
}

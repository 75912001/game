package game

import (
	"airplaneClient/src/battle"
	"airplaneClient/src/common"
	"airplaneClient/src/ui"
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (p *Game) Update() error {
	ui.GUIMgr.Update()

	// 根据游戏状态执行不同的更新逻辑
	switch p.state {
	case StateMenu: // 菜单状态下的更新逻辑
		p.updateMenu()
	case StateBattling: // 游戏进行中的更新逻辑
		p.updateBattling()
	case StatePaused: // 暂停状态下的更新逻辑
		p.updatePaused()
	case StateGameOver: // 游戏结束状态下的更新逻辑
		p.updateGameOver()
	}

	return nil
}

// updateMenu 菜单状态的更新
func (p *Game) updateMenu() {
	p.debugMsg = "菜单状态 - 点击开始按钮开始游戏"
}

// updateBattling 游戏进行中的更新
func (p *Game) updateBattling() {
	// 首次进入游戏状态时,初始化飞机
	if p.playerPlane == nil {
		p.startGame()
	}

	if p.playerPlane != nil {
		// 更新飞机
		p.playerPlane.Update()
		// 用户按空格键发射子弹
		if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeySpace) {
			bullet := p.playerPlane.Fire()
			p.bullets = append(p.bullets, bullet)
		}
	}

	// 更新所有子弹
	p.updateBullets()

	p.debugMsg = fmt.Sprintf("游戏进行中 - 按空格发射子弹, 当前子弹数: %v", len(p.bullets))

	// 按 ESC 暂停游戏
	if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeyEscape) {
		p.state = StatePaused
	}
}

// startGame 开始游戏,初始化游戏对象
func (p *Game) startGame() {
	// 创建玩家飞机,初始位置在屏幕下方中央
	p.playerPlane = battle.NewPlane(1, 1, 375, 500, 2)
}

// updatePaused 暂停状态的更新
func (p *Game) updatePaused() {
	p.debugMsg = "游戏暂停 - 按 ESC 继续"

	// 按 ESC 继续游戏
	if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeyEscape) {
		p.state = StateBattling
	}
}

// updateGameOver 游戏结束状态的更新
func (p *Game) updateGameOver() {
	p.debugMsg = "游戏结束"
}

// updateBullets 更新所有子弹
func (p *Game) updateBullets() {
	// 更新每颗子弹的位置
	for _, bullet := range p.bullets {
		bullet.Update()
	}

	// 移除飞出屏幕的子弹
	validBullets := make([]*battle.Bullet, 0, len(p.bullets))
	for _, bullet := range p.bullets {
		if !bullet.IsOutOfScreen(common.ScreenWidth, common.ScreenHeight) {
			validBullets = append(validBullets, bullet)
		}
	}
	p.bullets = validBullets
}

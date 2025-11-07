package game

import (
	apbattle "airplaneClient/src/battle"
	apcui "airplaneClient/src/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (p *Game) Update() error {
	apcui.GUIMgr.Update()

	// 根据游戏状态执行不同的更新逻辑
	switch p.state {
	case StateMenu:
		// 菜单状态下的更新逻辑
		p.updateMenu()
	case StateBattling:
		// 游戏进行中的更新逻辑
		p.updateBattling()
	case StatePaused:
		// 暂停状态下的更新逻辑
		p.updatePaused()
	case StateGameOver:
		// 游戏结束状态下的更新逻辑
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

	// 更新飞机
	if p.playerPlane != nil {
		p.playerPlane.Update()
	}

	p.debugMsg = "游戏进行中"

	// 按 ESC 暂停游戏
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		p.state = StatePaused
	}
}

// startGame 开始游戏,初始化游戏对象
func (p *Game) startGame() {
	// 创建玩家飞机,初始位置在屏幕下方中央
	p.playerPlane = apbattle.NewPlane(375, 500, 2)
}

// updatePaused 暂停状态的更新
func (p *Game) updatePaused() {
	p.debugMsg = "游戏暂停 - 按 ESC 继续"

	// 按 ESC 继续游戏
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		p.state = StateBattling
	}
}

// updateGameOver 游戏结束状态的更新
func (p *Game) updateGameOver() {
	p.debugMsg = "游戏结束"
}

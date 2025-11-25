package battle

import ebitenv2 "github.com/hajimehoshi/ebiten/v2"

// UpdateEnemy 更新敌机状态（自动向下移动）
func (p *Plane) UpdateEnemy() {
	// 敌机自动向下移动
	p.y += p.speed
}

// IsOutOfScreen 判断敌机是否飞出屏幕
func (p *Plane) IsOutOfScreen(screenHeight float64) bool {
	return screenHeight < p.y
}

// DrawEnemy 绘制敌机
func (p *Plane) DrawEnemy(screen *ebitenv2.Image) {
	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.frames[0], op)

	// 绘制调试边界
	p.DrawDebugBounds(screen)
}

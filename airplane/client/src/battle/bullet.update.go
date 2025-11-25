package battle

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"math"
)

// Update 更新子弹位置
func (b *Bullet) Update() {
	// 根据方向更新位置
	b.x += math.Cos(b.direction) * b.speed
	b.y += math.Sin(b.direction) * b.speed
}

// Draw 绘制子弹
func (b *Bullet) Draw(screen *ebitenv2.Image) {
	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.frames[0], op)
}

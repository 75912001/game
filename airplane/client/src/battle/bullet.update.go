package battle

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"math"
)

// Update 更新子弹位置
func (b *Bullet) Update() {
	// 根据方向更新位置
	b.SetX(b.GetX() + math.Cos(b.Orientation)*b.GetSpeed())
	b.SetY(b.GetY() + math.Sin(b.Orientation)*b.GetSpeed())
}

// Draw 绘制子弹
func (b *Bullet) Draw(screen *ebitenv2.Image) {
	op := &ebitenv2.DrawImageOptions{}
	// 应用缩放
	op.GeoM.Scale(b.GetScale(), b.GetScale())
	op.GeoM.Translate(b.GetX(), b.GetY())
	screen.DrawImage(b.GetFrames()[0], op)
	// 绘制调试边界
	b.DrawDebugBounds(screen)
}

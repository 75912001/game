package common

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2vector "github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// DrawDebugBounds 绘制调试边界
func (p *Object) DrawDebugBounds(screen *ebitenv2.Image) {
	// 绘制图像边界（蓝色）
	if p.debugDrawImageBounds {
		bounds := p.GetBounds()
		ebitenv2vector.StrokeRect(
			screen,
			float32(bounds.Min.X),
			float32(bounds.Min.Y),
			float32(bounds.Dx()),
			float32(bounds.Dy()),
			1,                                      // 线宽
			color.RGBA{R: 0, G: 0, B: 255, A: 255}, // 蓝色
			false,
		)
	}

	// 绘制碰撞体边界（红色）
	if p.debugDrawCollider {
		colliderBounds := p.GetColliderBounds()
		ebitenv2vector.StrokeRect(
			screen,
			float32(colliderBounds.Min.X),
			float32(colliderBounds.Min.Y),
			float32(colliderBounds.Dx()),
			float32(colliderBounds.Dy()),
			1,                                      // 线宽
			color.RGBA{R: 255, G: 0, B: 0, A: 255}, // 红色
			false,
		)
	}
}

package common

import (
	"image"
)

// GetColliderBounds 获取对象的碰撞体边界矩形
func (p *Object) GetColliderBounds() image.Rectangle {
	// 计算碰撞体的偏移量，使其居中在缩放后的图像中
	scaledImageWidth := p.GetScaleImageWidth()
	scaledImageHeight := p.GetScaleImageHeight()
	scaledColliderWidth := p.GetScaleColliderWidth()
	scaledColliderHeight := p.GetScaleColliderHeight()
	offsetX := (scaledImageWidth - scaledColliderWidth) / 2
	offsetY := (scaledImageHeight - scaledColliderHeight) / 2

	return image.Rect(
		int(p.x+offsetX),
		int(p.y+offsetY),
		int(p.x+offsetX+scaledColliderWidth),
		int(p.y+offsetY+scaledColliderHeight),
	)
}

// CollidesWith 检测是否与另一个对象发生碰撞（使用碰撞体）
func (p *Object) CollidesWith(other *Object) bool {
	return p.GetColliderBounds().Overlaps(other.GetColliderBounds())
}

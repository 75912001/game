package common

import (
	"image"
)

func (p *Object) GetColliderWidth() float64 {
	return p.colliderWidth * p.scale
}

func (p *Object) GetColliderHeight() float64 {
	return p.colliderHeight * p.scale
}

// GetColliderBounds 获取对象的碰撞体边界矩形-缩放后
func (p *Object) GetColliderBounds() image.Rectangle {
	// 计算碰撞体的偏移量，使其居中在缩放后的图像中
	scaledImageWidth := p.GetImageWidth()
	scaledImageHeight := p.GetImageHeight()
	scaledColliderWidth := p.GetColliderWidth()
	scaledColliderHeight := p.GetColliderHeight()
	offsetX := (scaledImageWidth - scaledColliderWidth) / 2
	offsetY := (scaledImageHeight - scaledColliderHeight) / 2

	return image.Rect(
		int(p.x+offsetX),
		int(p.y+offsetY),
		int(p.x+offsetX+scaledColliderWidth),
		int(p.y+offsetY+scaledColliderHeight),
	)
}

// CollidesWith 检测是否与另一个对象发生碰撞-使用碰撞体
func (p *Object) CollidesWith(other *Object) bool {
	return p.GetColliderBounds().Overlaps(other.GetColliderBounds())
}

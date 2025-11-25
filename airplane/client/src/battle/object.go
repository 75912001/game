package battle

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2vector "github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
)

type Object struct {
	id    uint32 // ID
	level uint32 // 等级

	imageWidth  float64 // 图像-宽
	imageHeight float64 // 图像-高

	colliderWidth  float64 // 碰撞体-宽
	colliderHeight float64 // 碰撞体-高

	x float64 // x 坐标
	y float64 // y 坐标

	speed float64 // 移动速度

	frames []*ebitenv2.Image // 动画帧

	//------------------------------------------------------------
	// 是否画出碰撞体边界(调试用) 红色
	debugDrawCollider bool
	// 是否画出图像边界(调试用) 蓝色
	debugDrawImageBounds bool
}

func newObject(id, level uint32,
	imageWidth, imageHeight, colliderWidth, colliderHeight float64,
	x, y, speed float64, frames []*ebitenv2.Image) *Object {
	return &Object{
		id:    id,
		level: level,

		imageWidth:  imageWidth,
		imageHeight: imageHeight,

		colliderWidth:  colliderWidth,
		colliderHeight: colliderHeight,

		x: x,
		y: y,

		speed: speed,

		frames:               frames,
		debugDrawCollider:    true,
		debugDrawImageBounds: true,
	}
}

// GetX 获取X坐标
func (p *Object) GetX() float64 {
	return p.x
}

// GetY 获取Y坐标
func (p *Object) GetY() float64 {
	return p.y
}

// GetImageWidth 获取图像宽度
func (p *Object) GetImageWidth() float64 {
	return p.imageWidth
}

// GetImageHeight 获取图像高度
func (p *Object) GetImageHeight() float64 {
	return p.imageHeight
}

// GetBounds 获取对象的边界矩形（基于图像大小）
func (p *Object) GetBounds() image.Rectangle {
	return image.Rect(
		int(p.x),
		int(p.y),
		int(p.x+p.imageWidth),
		int(p.y+p.imageHeight),
	)
}

// GetColliderBounds 获取对象的碰撞体边界矩形
func (p *Object) GetColliderBounds() image.Rectangle {
	// 计算碰撞体的偏移量，使其居中
	offsetX := (p.imageWidth - p.colliderWidth) / 2
	offsetY := (p.imageHeight - p.colliderHeight) / 2

	return image.Rect(
		int(p.x+offsetX),
		int(p.y+offsetY),
		int(p.x+offsetX+p.colliderWidth),
		int(p.y+offsetY+p.colliderHeight),
	)
}

// CollidesWith 检测是否与另一个对象发生碰撞（使用碰撞体）
func (p *Object) CollidesWith(other *Object) bool {
	return p.GetColliderBounds().Overlaps(other.GetColliderBounds())
}

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

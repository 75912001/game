package common

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Object struct {
	id    uint32 // ID
	level uint32 // 等级

	isEnemy bool // 是否敌人

	imageWidth  float64 // 图像-宽
	imageHeight float64 // 图像-高

	colliderWidth  float64 // 碰撞体-宽
	colliderHeight float64 // 碰撞体-高

	x float64 // x 坐标
	y float64 // y 坐标

	speed   float64 // 移动速度
	scale   float64 // 缩放比例(1.0为原始大小)
	hp      uint32  // 生命值
	defense uint32  // 防御力

	frames []*ebitenv2.Image // 动画帧

	//------------------------------------------------------------
	// 是否画出碰撞体边界(调试用) 红色
	debugDrawCollider bool
	// 是否画出图像边界(调试用) 蓝色
	debugDrawImageBounds bool
}

func NewObject(id, level uint32, isEnemy bool,
	imageWidth, imageHeight, colliderWidth, colliderHeight float64,
	x, y, speed, scale float64, frames []*ebitenv2.Image) *Object {
	return &Object{
		id:    id,
		level: level,

		isEnemy: isEnemy,

		imageWidth:  imageWidth,
		imageHeight: imageHeight,

		colliderWidth:  colliderWidth,
		colliderHeight: colliderHeight,

		x: x,
		y: y,

		speed:   speed,
		scale:   scale,
		hp:      100,
		defense: 0,

		frames:               frames,
		debugDrawCollider:    true,
		debugDrawImageBounds: true,
	}
}

func (p *Object) IsEnemy() bool {
	return p.isEnemy
}

// GetX 获取X坐标
func (p *Object) GetX() float64 {
	return p.x
}

func (p *Object) SetX(x float64) {
	p.x = x
}

// GetY 获取Y坐标
func (p *Object) GetY() float64 {
	return p.y
}

func (p *Object) SetY(y float64) {
	p.y = y
}

func (p *Object) GetSpeed() float64 {
	return p.speed
}

// GetScale 获取缩放比例
func (p *Object) GetScale() float64 {
	return p.scale
}

// GetImageWidth 获取图像宽度-缩放后
func (p *Object) GetImageWidth() float64 {
	return p.imageWidth * p.scale
}

// GetRawImageWidth 获取图像宽度-原始大小
func (p *Object) GetRawImageWidth() float64 {
	return p.imageWidth
}

// GetImageHeight 获取图像高度-缩放后
func (p *Object) GetImageHeight() float64 {
	return p.imageHeight * p.scale
}

// GetRawImageHeight 获取图像高度-原始大小
func (p *Object) GetRawImageHeight() float64 {
	return p.imageHeight
}

func (p *Object) GetHp() uint32 {
	return p.hp
}

func (p *Object) SetHp(hp uint32) {
	p.hp = hp
}

// GetFrames 获取帧图像
func (p *Object) GetFrames() []*ebitenv2.Image {
	return p.frames
}

// GetBounds 获取对象的边界矩形-缩放后
func (p *Object) GetBounds() image.Rectangle {
	return image.Rect(
		int(p.x),
		int(p.y),
		int(p.x+p.GetImageWidth()),
		int(p.y+p.GetImageHeight()),
	)
}

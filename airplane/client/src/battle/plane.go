package battle

import (
	"airplaneClient/src/common"
	"airplaneClient/src/config"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

// Plane 飞机结构体
type Plane struct {
	id      uint32  // ID
	level   uint32  // 等级
	isEnemy bool    // 是否敌人
	speed   float64 // 移动速度
	scale   float64 // 缩放比例(1.0为原始大小)
	defense float64 // 防御力
	part    [common.PlanePartTypeMax]*PlanePart
}

// NewPlane 创建-飞机
func NewPlane(id, level uint32, isEnemy bool, x, y, speed, scale float64, defense float64) *Plane {
	plane := &Plane{
		id:      id,
		level:   level,
		isEnemy: isEnemy,
		speed:   speed,
		scale:   scale,
		defense: defense,
	}
	planeCfg := config.GPlaneMgr.GetPlane(common.PlaneKey{id, level})
	for part := range plane.part {
		planePartFrameData := planeCfg.PartsFramesData[part]
		if len(planePartFrameData) == 0 { // 没有这个部件
			continue
		}
		// 部件帧信息
		planePartFramesInfo := planeCfg.PartsFramesInfo[part]
		if len(planePartFramesInfo) == 0 { // 没有这个部件
			continue
		}
		// 机身0帧信息
		body0FrameInfo := planeCfg.PartsFramesInfo[common.PlanePartTypeBody][0].Frame
		// 当前部件0帧信息
		part0FrameInfo := planePartFramesInfo[0].Frame
		var partX float64
		var partY float64
		switch common.PlanePartType(part) { // 不同部件的偏移位置不同
		case common.PlanePartTypeNose: // 机头位于机身顶部中央
			if body0FrameInfo.Width < part0FrameInfo.Width { // 机身宽度 < 机头的宽度
				partX = x - (float64(part0FrameInfo.Width)/2 - float64(body0FrameInfo.Width)/2)
			} else if part0FrameInfo.Width < body0FrameInfo.Width { // 机头的宽度 < 机身宽度
				partX = x + (float64(body0FrameInfo.Width)/2 - float64(part0FrameInfo.Width)/2)
			} else { // 机身宽度 == 机头的宽度
				partX = x
			}
			partY = y - float64(part0FrameInfo.Height)
		case common.PlanePartTypeBody: // 使用机身部件的坐标作为飞机的坐标
			partX = x
			partY = y
		case common.PlanePartTypeLeftWing: // 左翼位于机身左侧顶部
			partX = x - float64(part0FrameInfo.Width)
			partY = y
		case common.PlanePartTypeRightWing: // 右翼位于机身右侧顶部
			partX = x + float64(body0FrameInfo.Width)
			partY = y
		default: // 不存在的部件
			return nil
		}
		plane.part[part] = NewPlanePart(
			plane,
			partX,
			partY,
			100,
			common.PlanePartType(part),
			float64(part0FrameInfo.Width),
			float64(part0FrameInfo.Height),
			float64(part0FrameInfo.Width)*0.8,
			float64(part0FrameInfo.Height)*0.8,
			planePartFrameData,
		)
	}
	return plane
}

// IsOutOfScreen 判断敌机是否飞出屏幕
func (p *Plane) IsOutOfScreen() bool {
	if !p.IsEnemy() { // 用户
		return false
	}
	// todo menglc 使用 全屏幕判断

	return screenHeight < p.GetY()-p.part[common.PlanePartTypeBody].imageHeight
}

func (p *Plane) IsEnemy() bool {
	return p.isEnemy
}

// GetX 获取X坐标-使用机身
func (p *Plane) GetX() float64 {
	return p.part[common.PlanePartTypeBody].x
}

func (p *Plane) SetX(x float64) {
	p.x = x
}

// GetY 获取Y坐标-使用机身
func (p *Plane) GetY() float64 {
	return p.part[common.PlanePartTypeBody].y
}

func (p *Plane) SetY(y float64) {
	p.y = y
}

func (p *Plane) GetSpeed() float64 {
	return p.speed
}

// GetScale 获取缩放比例
func (p *Plane) GetScale() float64 {
	return p.scale
}

// GetImageWidth 获取图像宽度-缩放后
func (p *Plane) GetImageWidth() float64 {
	return p.imageWidth * p.scale
}

// GetRawImageWidth 获取图像宽度-原始大小
func (p *Plane) GetRawImageWidth() float64 {
	return p.imageWidth
}

// GetImageHeight 获取图像高度-缩放后
func (p *Plane) GetImageHeight() float64 {
	return p.imageHeight * p.scale
}

// GetRawImageHeight 获取图像高度-原始大小
func (p *Plane) GetRawImageHeight() float64 {
	return p.imageHeight
}

func (p *Plane) GetHp() uint32 {
	return p.hp
}

func (p *Plane) SetHp(hp uint32) {
	p.hp = hp
}

// GetFrames 获取帧图像
func (p *Plane) GetFrames() []*ebitenv2.Image {
	return p.frames
}

// GetBounds 获取对象的边界矩形-缩放后
func (p *Plane) GetBounds() image.Rectangle {
	return image.Rect(
		int(p.x),
		int(p.y),
		int(p.x+p.GetImageWidth()),
		int(p.y+p.GetImageHeight()),
	)
}

package battle

import (
	"airplaneClient/src/common"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// PlanePart 飞机部件
type PlanePart struct {
	plane          *Plane               // 所属 Plane
	x              float64              // x 坐标
	y              float64              // y 坐标
	hp             uint32               // 生命值
	part           common.PlanePartType // 部件类型
	imageWidth     float64              // 图像-宽
	imageHeight    float64              // 图像-高
	colliderWidth  float64              // 碰撞体-宽
	colliderHeight float64              // 碰撞体-高
	frames         []*ebitenv2.Image    // 动画帧
	framesIdx      int                  // 当前帧索引
	//------------------------------------------------------------
	// 是否画出碰撞体边界(调试用) 红色
	debugDrawCollider bool
	// 是否画出图像边界(调试用) 蓝色
	debugDrawImageBounds bool
}

func NewPlanePart(plane *Plane,
	x, y float64,
	hp uint32,
	part common.PlanePartType,
	imageWidth, imageHeight float64,
	colliderWidth, colliderHeight float64,
	frames []*ebitenv2.Image) *PlanePart {
	return &PlanePart{
		plane:                plane,
		x:                    x,
		y:                    y,
		hp:                   hp,
		part:                 part,
		imageWidth:           imageWidth,
		imageHeight:          imageHeight,
		colliderWidth:        colliderWidth,
		colliderHeight:       colliderHeight,
		frames:               frames,
		framesIdx:            0,
		debugDrawCollider:    true,
		debugDrawImageBounds: true,
	}
}

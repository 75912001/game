package battle

import (
	"airplaneClient/src/common"
	"airplaneClient/src/resources"
	resourcescommon "airplaneClient/src/resources/common"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// Plane 飞机结构体
type Plane struct {
	*common.Object

	hp uint32

	currentFrameType resourcescommon.PlaneFrameType // 当前帧类型
	flipHorizontal   bool                           // 是否水平镜像(用于左倾)

	frameCounter int // 帧计数器，用于控制动画速度
}

// NewPlane 创建一个新飞机
func NewPlane(id, level uint32, isBelongsToUser bool, x, y, speed float64) *Plane {
	// 加载飞机图片的4帧动画
	// 001.001.png 包含4帧:
	// 帧0-直飞(正), 帧1-微右倾, 帧2-右倾, 帧3-大右倾
	// 左倾通过镜像右倾帧实现
	var planeFrameType resourcescommon.PlaneFrameType
	frames, imageWidth, imageHeight, err := resources.LoadPlaneFrames(id, level, planeFrameType.GetFrameTypeCount())
	if err != nil {
		panic(err)
	}
	return &Plane{
		Object: common.NewObject(id, level,
			isBelongsToUser,
			imageWidth,
			imageHeight,
			imageWidth*0.6,  // 碰撞体宽度为图像宽度的60%
			imageHeight*0.6, // 碰撞体高度为图像高度的60%
			x,
			y,
			speed,
			frames,
		),
		hp:               100,
		currentFrameType: resourcescommon.PlaneFrameTypeStraight,
		flipHorizontal:   false, // 默认不镜像
		frameCounter:     0,
	}
}

// NewEnemyPlane 创建一个敌机（单帧，向下移动）
func NewEnemyPlane(id, level uint32, isBelongsToUser bool, x, y, speed float64) *Plane {
	// 敌机只加载1帧图像
	frames, imageWidth, imageHeight, err := resources.LoadEnemyPlaneFrames(id, level, 1)
	if err != nil {
		panic(err)
	}
	return &Plane{
		Object: common.NewObject(id, level,
			isBelongsToUser,
			imageWidth,
			imageHeight,
			imageWidth*0.6,  // 碰撞体宽度为图像宽度的60%
			imageHeight*0.6, // 碰撞体高度为图像高度的60%
			x,
			y,
			speed,
			frames,
		),
		hp:               10, // 敌机HP为10
		currentFrameType: 0,  // 只有一帧
		flipHorizontal:   false,
		frameCounter:     0,
	}
}

// Fire 发射一颗子弹
func (p *Plane) Fire() *Bullet {
	// 计算子弹初始位置 - 在飞机顶端中央
	bulletX := p.GetX() + p.GetImageWidth()/2
	bulletY := p.GetY()

	// 创建一颗向上移动的子弹 (id=1, level=1, speed=5)
	return NewBullet(1, 1, !p.IsEnemy(), bulletX, bulletY, 1, BulletDirectionUp, p)
}

// GetCurrentImage 获取当前飞机图像
func (p *Plane) GetCurrentImage() *ebitenv2.Image {
	return p.GetFrames()[p.currentFrameType]
}

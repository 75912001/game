package battle

import (
	"airplaneClient/src/common"
	resourcescommon "airplaneClient/src/resources/common"
	"airplaneClient/src/ui"
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// Update 更新飞机状态
func (p *Plane) Update() {
	// 检测移动方向
	currentMoveDir := 0 // x轴方向上: -1:左移, 0:不动, 1:右移
	if ebitenv2.IsKeyPressed(ebitenv2.KeyLeft) || ebitenv2.IsKeyPressed(ebitenv2.KeyA) {
		p.x -= p.speed
		currentMoveDir = -1
	}
	if ebitenv2.IsKeyPressed(ebitenv2.KeyRight) || ebitenv2.IsKeyPressed(ebitenv2.KeyD) {
		p.x += p.speed
		currentMoveDir = 1
	}
	if ebitenv2.IsKeyPressed(ebitenv2.KeyUp) || ebitenv2.IsKeyPressed(ebitenv2.KeyW) {
		p.y -= p.speed
	}
	if ebitenv2.IsKeyPressed(ebitenv2.KeyDown) || ebitenv2.IsKeyPressed(ebitenv2.KeyS) {
		p.y += p.speed
	}

	// 边界检测 - x
	if p.x < 0 {
		p.x = 0
	}
	if (common.ScreenWidth - p.imageWidth) < p.x {
		p.x = common.ScreenWidth - p.imageWidth
	}

	// 边界检测 - y
	if p.y < 0 {
		p.y = 0
	}
	if (common.ScreenHeight - p.imageHeight) < p.y {
		p.y = common.ScreenHeight - p.imageHeight
	}

	// 更新动画帧
	p.updateFrame(currentMoveDir)
}

// updateFrame 更新动画帧
func (p *Plane) updateFrame(currentMoveDir int) {
	targetFrame := resourcescommon.PlaneFrameTypeStraight
	targetFlip := false // 默认不镜像

	if currentMoveDir < 0 { // 左移：使用帧3并镜像（最大倾斜）
		targetFrame = resourcescommon.PlaneFrameTypeSharpRight
		targetFlip = true
	} else if 0 < currentMoveDir { // 右移：使用帧3（最大倾斜）
		targetFrame = resourcescommon.PlaneFrameTypeSharpRight
		targetFlip = false
	} else {
		// 不移动：回到直飞状态
		targetFrame = resourcescommon.PlaneFrameTypeStraight
		targetFlip = false
	}

	// 平滑过渡到目标帧
	p.frameCounter++
	if resourcescommon.PlaneFrameDelay <= p.frameCounter {
		p.frameCounter = 0

		// 如果改变方向，先回到帧0再转向
		if targetFlip != p.flipHorizontal && p.currentFrameType != resourcescommon.PlaneFrameTypeStraight { // 先回到直飞
			p.currentFrameType--
		} else { // 正常过渡
			if p.currentFrameType < targetFrame {
				p.currentFrameType++
			} else if targetFrame < p.currentFrameType {
				p.currentFrameType--
			}
		}
		// 确保帧索引在有效范围内
		if p.currentFrameType < 0 {
			p.currentFrameType = resourcescommon.PlaneFrameTypeStraight
		} else if resourcescommon.PlaneFrameTypeMax <= p.currentFrameType {
			p.currentFrameType = resourcescommon.PlaneFrameTypeMax - 1
		}
		// 只有在帧0时才允许改变镜像方向
		if p.currentFrameType == resourcescommon.PlaneFrameTypeStraight {
			p.flipHorizontal = targetFlip
		}
	}
}

// Draw 绘制飞机
func (p *Plane) Draw(screen *ebitenv2.Image) {
	op := &ebitenv2.DrawImageOptions{}

	// 如果需要水平镜像（向左倾斜）
	if p.flipHorizontal {
		// 先水平翻转
		op.GeoM.Scale(-1, 1)
		// 翻转后需要调整X坐标（因为翻转会改变原点位置）
		op.GeoM.Translate(p.x+p.imageWidth, p.y)
	} else {
		// 正常绘制
		op.GeoM.Translate(p.x, p.y)
	}

	screen.DrawImage(p.GetCurrentImage(), op)

	// 绘制调试边界
	p.DrawDebugBounds(screen)

	ui.Printf(screen, 10, 10, fmt.Sprintf("使用方向键或WASD移动飞机 x:%v y:%v", p.x, p.y))
}

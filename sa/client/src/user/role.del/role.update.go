package role_del

import (
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/ui"
	"saClient/src/user/role"
)

// Update 更新角色状态
func (p *role.Role) Update() {
	// 检测移动方向
	var direction proto.RoleDirection = proto.RoleDirection_RoleDirection_Unknow

	// 检测按键输入
	up := ebitenv2.IsKeyPressed(ebitenv2.KeyUp) || ebitenv2.IsKeyPressed(ebitenv2.KeyW)
	down := ebitenv2.IsKeyPressed(ebitenv2.KeyDown) || ebitenv2.IsKeyPressed(ebitenv2.KeyS)
	left := ebitenv2.IsKeyPressed(ebitenv2.KeyLeft) || ebitenv2.IsKeyPressed(ebitenv2.KeyA)
	right := ebitenv2.IsKeyPressed(ebitenv2.KeyRight) || ebitenv2.IsKeyPressed(ebitenv2.KeyD)

	// 根据按键组合确定方向
	if up && left {
		direction = proto.RoleDirection_RoleDirection_UpLeft
	} else if up && right {
		direction = proto.RoleDirection_RoleDirection_UpRight
	} else if down && left {
		direction = proto.RoleDirection_RoleDirection_DownLeft
	} else if down && right {
		direction = proto.RoleDirection_RoleDirection_DownRight
	} else if up {
		direction = proto.RoleDirection_RoleDirection_Up
	} else if down {
		direction = proto.RoleDirection_RoleDirection_Down
	} else if left {
		direction = proto.RoleDirection_RoleDirection_Left
	} else if right {
		direction = proto.RoleDirection_RoleDirection_Right
	}

	// 移动角色
	if direction != proto.RoleDirection_RoleDirection_Unknow {
		p.move(direction)
	}
}

// move 根据方向移动角色
func (p *role.Role) move(direction proto.RoleDirection) {
	speed := SpeedDefault
	x := p.GetX()
	y := p.GetY()

	switch direction {
	case proto.RoleDirection_RoleDirection_Up:
		y -= speed
	case proto.RoleDirection_RoleDirection_Down:
		y += speed
	case proto.RoleDirection_RoleDirection_Left:
		x -= speed
	case proto.RoleDirection_RoleDirection_Right:
		x += speed
	case proto.RoleDirection_RoleDirection_UpLeft:
		x -= speed * 0.707 // 斜向移动时速度分量
		y -= speed * 0.707
	case proto.RoleDirection_RoleDirection_UpRight:
		x += speed * 0.707
		y -= speed * 0.707
	case proto.RoleDirection_RoleDirection_DownLeft:
		x -= speed * 0.707
		y += speed * 0.707
	case proto.RoleDirection_RoleDirection_DownRight:
		x += speed * 0.707
		y += speed * 0.707
	}

	// 边界检测
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	maxX := float64(common.ScreenWidth) - p.GetImageWidth()
	maxY := float64(common.ScreenHeight) - p.GetImageHeight()
	if x > maxX {
		x = maxX
	}
	if y > maxY {
		y = maxY
	}

	p.SetX(x)
	p.SetY(y)
}

// Draw 绘制角色
func (p *role.Role) Draw(screen *ebitenv2.Image) {
	op := &ebitenv2.DrawImageOptions{}
	// 应用缩放
	op.GeoM.Scale(p.GetScale(), p.GetScale())
	// 应用位置
	op.GeoM.Translate(p.GetX(), p.GetY())

	// 绘制角色图像
	if p.imageMgr != nil {
		// TODO: 绘制当前帧图像
	}

	// 显示调试信息
	ui.Printf(screen, 10, 10, fmt.Sprintf("使用方向键或WASD移动角色 x:%.0f y:%.0f", p.GetX(), p.GetY()))

	// 绘制调试边界
	if p.debugDrawImageBounds {
		p.DrawDebugBounds(screen)
	}
}

// GetX 获取X坐标
func (p *role.Role) GetX() int {
	if p.object != nil && p.object.GetImage() != nil {
		return p.object.GetImage().GetX()
	}
	return 0
}

// SetX 设置X坐标
func (p *role.Role) SetX(x int) {
	if p.object != nil && p.object.GetImage() != nil {
		p.object.GetImage().SetX(x)
	}
}

// GetY 获取Y坐标
func (p *role.Role) GetY() int {
	if p.object != nil && p.object.GetImage() != nil {
		return p.object.GetImage().GetY()
	}
	return 0
}

// SetY 设置Y坐标
func (p *role.Role) SetY(y int) {
	if p.object != nil && p.object.GetImage() != nil {
		p.object.GetImage().SetY(y)
	}
}

// GetScale 获取缩放比例
func (p *role.Role) GetScale() float64 {
	return ScaleDefault
}

// GetImageWidth 获取图像宽度
func (p *role.Role) GetImageWidth() float64 {
	if p.object != nil && p.object.GetImage() != nil {
		return float64(p.object.GetImage().GetWidth())
	}
	return 0
}

// GetImageHeight 获取图像高度
func (p *role.Role) GetImageHeight() float64 {
	if p.object != nil && p.object.GetImage() != nil {
		return float64(p.object.GetImage().GetHeight())
	}
	return 0
}

// DrawDebugBounds 绘制调试边界
func (p *role.Role) DrawDebugBounds(screen *ebitenv2.Image) {
	// TODO: 实现调试边界绘制
}

package role

import (
	"saClient/src/cfg"
	"saClient/src/proto"

	"github.com/hajimehoshi/ebiten/v2"
)

// HandleInput 处理键盘输入
func (p *Role) HandleInput() { // 获取当前位置
	// 检测 WASD 按键
	up := ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyD)

	// 无按键则不处理
	if !up && !down && !left && !right {
		return
	}

	moveSpeed := int(cfg.GCommon.RoleDefaultMoveSpeed)
	// 获取当前位置
	x := p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_BottomCenterX)
	y := p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_BottomCenterY)

	// 计算方向和位移
	var direction proto.AssetDirection
	dx, dy := 0, 0

	if up && right {
		direction = proto.AssetDirection_AssetDirection_UpRight
		dx, dy = moveSpeed, -moveSpeed
	} else if up && left {
		direction = proto.AssetDirection_AssetDirection_UpLeft
		dx, dy = -moveSpeed, -moveSpeed
	} else if down && right {
		direction = proto.AssetDirection_AssetDirection_DownRight
		dx, dy = moveSpeed, moveSpeed
	} else if down && left {
		direction = proto.AssetDirection_AssetDirection_DownLeft
		dx, dy = -moveSpeed, moveSpeed
	} else if up {
		direction = proto.AssetDirection_AssetDirection_Up
		dy = -moveSpeed
	} else if down {
		direction = proto.AssetDirection_AssetDirection_Down
		dy = moveSpeed
	} else if left {
		direction = proto.AssetDirection_AssetDirection_Left
		dx = -moveSpeed
	} else { // right
		direction = proto.AssetDirection_AssetDirection_Right
		dx = moveSpeed
	}

	// 更新位置和方向
	p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_BottomCenterX, uint64(x+dx))
	p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_BottomCenterY, uint64(y+dy))
	p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_Direction, uint64(direction))

	// 更新动画帧
	p.frameIdx++
}

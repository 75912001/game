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

	// 计算新位置
	newX := x + dx
	newY := y + dy

	// 获取角色尺寸(使用当前帧的尺寸)
	roleHalfWidth := p.roleSprite.roleImageSprite.Frame.Width / 2
	roleHeight := p.roleSprite.roleImageSprite.Frame.Height

	// 限制在地图范围内(考虑角色尺寸)
	mapWidth, mapHeight := p.scene.GetMapSize()
	// 左边界(角色中心点至少要离左边缘半个角色宽度
	if newX < roleHalfWidth {
		newX = roleHalfWidth
	}
	// 右边界:角色中心点至少要离右边缘半个角色宽度
	if newX > mapWidth-roleHalfWidth {
		newX = mapWidth - roleHalfWidth
	}
	if true { // 只露脚即可
		if newY < 20 {
			newY = 20
		}
	} else { // 上边界:脚底Y坐标至少要离上边缘一个角色高度
		if newY < roleHeight {
			newY = roleHeight
		}
	}
	// 下边界:脚底Y坐标不能超过地图底部
	if newY > mapHeight {
		newY = mapHeight
	}

	// 更新位置和方向
	p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_BottomCenterX, uint64(newX))
	p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_BottomCenterY, uint64(newY))
	p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_Direction, uint64(direction))

	// 更新动画帧（每 6 tick 切换一帧，60 TPS 下约 10 FPS 动画）
	p.frameTick++
	if p.frameTick >= 6 {
		p.frameTick = 0
		p.frameIdx++
	}
}

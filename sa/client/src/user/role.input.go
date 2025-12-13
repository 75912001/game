package user

import (
	"log"
	"saClient/src/cfg"
	"saClient/src/proto"

	xutil "github.com/75912001/xlib/util"
	"github.com/hajimehoshi/ebiten/v2"
)

// HandleInput 处理键盘输入
// 移动按屏幕方向(World坐标系)，而非Tile坐标系
func (p *Role) HandleInput() {
	// 检测 WASD 按键
	up := ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyD)

	// 无按键则不处理
	if !up && !down && !left && !right {
		return
	}

	// 计算屏幕方向的移动向量 (World 坐标系)
	var dx, dy float32 = 0, 0
	var orientation proto.AssetOrientation

	// 屏幕方向: W=上, S=下, A=左, D=右
	if up && right {
		orientation = proto.AssetOrientation_AssetOrientation_UpRight
		dx, dy = 1, -1
	} else if up && left {
		orientation = proto.AssetOrientation_AssetOrientation_UpLeft
		dx, dy = -1, -1
	} else if down && right {
		orientation = proto.AssetOrientation_AssetOrientation_DownRight
		dx, dy = 1, 1
	} else if down && left {
		orientation = proto.AssetOrientation_AssetOrientation_DownLeft
		dx, dy = -1, 1
	} else if up {
		orientation = proto.AssetOrientation_AssetOrientation_Up
		dy = -1
	} else if down {
		orientation = proto.AssetOrientation_AssetOrientation_Down
		dy = 1
	} else if left {
		orientation = proto.AssetOrientation_AssetOrientation_Left
		dx = -1
	} else { // right
		orientation = proto.AssetOrientation_AssetOrientation_Right
		dx = 1
	}

	// 归一化对角线移动 (避免对角线移动更快)
	if dx != 0 && dy != 0 {
		dx *= 0.707 // 1/sqrt(2) ≈ 0.707
		dy *= 0.707
	}

	// 更新方向
	p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_Orientation, uint64(orientation))

	// 移动速度 (像素/帧)
	moveSpeed := cfg.GCommon.RoleDefaultMoveSpeed

	// 计算新的 World 坐标
	newWorldX := p.sprite.bottomCenterWorldX + dx*moveSpeed
	newWorldY := p.sprite.bottomCenterWorldY + dy*moveSpeed

	// 转换为 Tile 坐标进行边界检测
	tileX, tileY := p.scene.WorldToTile(newWorldX, newWorldY)

	// 限制在地图边界内
	clampedTX, clampedTY := p.scene.ClampTileBounds(tileX, tileY)

	// 如果被限制了，需要重新计算 World 坐标
	if !xutil.Float32Equal(clampedTX, tileX) || !xutil.Float32Equal(clampedTY, tileY) {
		newWorldX, newWorldY = p.scene.TileToWorld(clampedTX, clampedTY)
	}

	// 更新 World 坐标
	p.sprite.bottomCenterWorldX = newWorldX
	p.sprite.bottomCenterWorldY = newWorldY

	// 同步 Tile 坐标到 proto (用于记录/网络同步)
	p.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TX, clampedTX)
	p.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TY, clampedTY)

	// 更新动画帧
	p.frameTick++
	if p.frameTick >= 6 {
		p.frameTick = 0
		p.frameIdx++
	}

	// 判断角色 wx, wy 是否触发了碰撞
	isCollision := p.scene._map.cfg.FindCollisionObject(p.sprite.bottomCenterWorldX, p.sprite.bottomCenterWorldY)
	if isCollision {
		log.Printf("Role HandleInput collision at world (%.2f, %.2f)", p.sprite.bottomCenterWorldX, p.sprite.bottomCenterWorldY)
	}
}

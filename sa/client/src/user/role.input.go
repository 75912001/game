package user

import (
	"log"
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/user/camera"

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
	if up || down || left || right { // 角色移动
		p.handleKeyMove(up, down, left, right)
		return
	}
}

// 处理角色移动
func (p *Role) handleKeyMove(up, down, left, right bool) {
	newRoleSprite := p.sprite // 复制当前角色状态
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
	newRoleSprite.orientation = uint32(orientation)
	// 移动速度 (像素/帧)
	moveSpeed := cfg.GCommon.RoleDefaultMoveSpeed
	// 计算新的 World 坐标
	newWorldX := newRoleSprite.bottomCenterWX + dx*moveSpeed
	newWorldY := newRoleSprite.bottomCenterWY + dy*moveSpeed
	mapCfg := p.scene._map.tiledMapCfg
	// 限制在地图边界内
	newRoleSprite.bottomCenterTX, newRoleSprite.bottomCenterTY = mapCfg.IsometricCT.W2T(newWorldX, newWorldY)
	clampedTX, clampedTY := mapCfg.IsometricCT.ClampTileBounds(newRoleSprite.bottomCenterTX, newRoleSprite.bottomCenterTY)
	if !xutil.Float32Equal(newRoleSprite.bottomCenterTX, clampedTX) || !xutil.Float32Equal(newRoleSprite.bottomCenterTY, clampedTY) { // 需要限制
		newWorldX, newWorldY = mapCfg.IsometricCT.T2W(clampedTX, clampedTY)
		newRoleSprite.bottomCenterTX, newRoleSprite.bottomCenterTY = clampedTX, clampedTY
	}
	// 更新 World 坐标
	newRoleSprite.bottomCenterWX = newWorldX
	newRoleSprite.bottomCenterWY = newWorldY
	// 更新动画帧
	p.frameTick++
	if p.frameTick >= 6 {
		p.frameTick = 0
		p.frameIdx++
	}
	// 是否使用预设的角色状态(用于回退)
	var rollBack = false
	if unreachableObject, ok := mapCfg.FindUnreachableObject(newRoleSprite.bottomCenterWX, newRoleSprite.bottomCenterWY); ok { // 判断角色 wx, wy 是否触发了-不可达区域
		log.Printf("Role HandleInput unreachable at world (%.2f, %.2f) unreachableObject:%+v", newRoleSprite.bottomCenterWX, newRoleSprite.bottomCenterWY, unreachableObject)
		// 回退到之前的位置
		rollBack = true
	} else if collisionObject, ok := mapCfg.FindCollisionObject(newRoleSprite.bottomCenterWX, newRoleSprite.bottomCenterWY); ok { // 判断角色 wx, wy 是否触发了碰撞
		log.Printf("Role HandleInput collision at world (%.2f, %.2f) collisionObject:%+v", newRoleSprite.bottomCenterWX, newRoleSprite.bottomCenterWY, collisionObject)
		switch collisionObject.Type {
		case cfg.TiledObjectType_Portal: // 传送点
			p.SwitchScene(collisionObject.PortalCfg.MapID, collisionObject.PortalCfg.TX, collisionObject.PortalCfg.TY)
		case cfg.TiledObjectType_ArrivalPortal: // 传送点-到达
		}
	}
	if !rollBack { // 正常移动，更新动画帧索引
		p.sprite = newRoleSprite
		p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_Orientation, uint64(newRoleSprite.orientation))
		p.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TX, newRoleSprite.bottomCenterTX)
		p.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TY, newRoleSprite.bottomCenterTY)
	}

	p.UpdateWithAction()
}

// SwitchScene 切换场景 (带过渡动画)
func (p *Role) SwitchScene(mapID common.AssetID, tx, ty float32) {
	// 如果已经在过渡中，忽略新的切换请求
	if p.scene.transition.IsActive() || p.pendingScene != nil {
		return
	}
	// 获取目标地图的过渡配置
	targetMapCfg := cfg.GMapMgr.Maps.Get(mapID)
	// 预创建新场景 (使用目标地图的过渡配置)
	p.pendingScene = NewScene(mapID)
	p.pendingTX = tx
	p.pendingTY = ty
	// 当前场景转场退出 (使用目标地图的过渡配置)
	p.scene.transition.Out(targetMapCfg.SceneTransition)
}

// doSwitchScene 执行实际的场景切换 (无过渡动画)
func (p *Role) doSwitchScene(mapID common.AssetID, tx, ty float32) {
	p.scene = NewScene(mapID)
	p.camera = camera.NewCamera()
	// 初始化角色的 World 坐标
	p.sprite.bottomCenterWX, p.sprite.bottomCenterWY = p.scene._map.tiledMapCfg.IsometricCT.T2W(tx, ty)
}

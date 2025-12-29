package user

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"saClient/src/cfg"
	"saClient/src/common"
	commoncamera "saClient/src/common/camera"
	"saClient/src/proto"
)

// HandleInput 处理键盘输入
// 移动按屏幕方向(World坐标系)，而非Tile坐标系
func (p *Role) HandleInput() {
	// 检测 WASD 按键
	up := ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyD)

	// 移动输入优先级最高，立即打断其他动作
	if up || down || left || right {
		p.SetAction(proto.RoleAction_RoleAction_Move)
		p.handleKeyMove(up, down, left, right)
		return
	}

	if p.scene._map.IsArpgMap() { // 仅在 ARPG 地图中处理自动攻击
		p.handleArpgAutoAttack()
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
	newRoleSprite.orientation = orientation
	// 移动速度 (像素/帧)
	moveSpeed := cfg.GCommon.RoleDefMoveSpeed
	// 计算新的 World 坐标
	newWorldX := newRoleSprite.actionAnchorPointWX + dx*moveSpeed
	newWorldY := newRoleSprite.actionAnchorPointWY + dy*moveSpeed
	mapCfg := p.scene._map.tiledMapCfg

	{ // 限制在菱形地图边界内 (World 坐标)
		clampedX, clampedY, blocked := mapCfg.Boundary.ClampWithW(newWorldX, newWorldY)
		if blocked { // 坐标被限制了
			log.Printf("Role HandleInput boundary hit at world (%.3f, %.3f)", newWorldX, newWorldY)
		}
		newWorldX, newWorldY = clampedX, clampedY
	}

	newRoleSprite.actionAnchorPointTX, newRoleSprite.actionAnchorPointTY = mapCfg.IsometricCT.W2T(newWorldX, newWorldY)
	// 更新 World 坐标
	newRoleSprite.actionAnchorPointWX = newWorldX
	newRoleSprite.actionAnchorPointWY = newWorldY
	// 更新动画帧
	p.frameTick++
	if common.FrameTickPerChange <= p.frameTick {
		p.frameTick = 0
		p.frameIdx++
	}

	// 是否使用预设的角色状态(用于回退)
	var rollBack = false
	if portalObject, ok := mapCfg.FindPortalByObject(newRoleSprite.actionAnchorPointWX, newRoleSprite.actionAnchorPointWY); ok { // 判断角色 wx, wy 是否触发了-传送
		log.Printf("Role HandleInput portal at world (%.3f, %.3f) portalObject:%+v", newRoleSprite.actionAnchorPointWX, newRoleSprite.actionAnchorPointWY, portalObject)
		p.SwitchScene(portalObject.PortalCfg.MapID, portalObject.PortalCfg.WX, portalObject.PortalCfg.WY)
	} else if mapCfg.TileBlocked.IsBlockedWithTF(newRoleSprite.actionAnchorPointTX, newRoleSprite.actionAnchorPointTY) { // 使用 tile 图层 图块 阻挡检测
		log.Printf("Role HandleInput tile blocked at world (%.2f, %.2f)", newRoleSprite.actionAnchorPointWX, newRoleSprite.actionAnchorPointWY)
		rollBack = true
	} else if blockedObject, ok := mapCfg.FindBlockedByObject(newRoleSprite.actionAnchorPointWX, newRoleSprite.actionAnchorPointWY); ok { // 判断角色 wx, wy 是否触发了-阻挡
		log.Printf("Role HandleInput blocked at world (%.2f, %.2f) blockedObject:%+v", newRoleSprite.actionAnchorPointWX, newRoleSprite.actionAnchorPointWY, blockedObject)
		// 回退到之前的位置
		rollBack = true
	}
	if !rollBack { // 正常移动，更新动画帧索引
		p.sprite = newRoleSprite
		p.SetValueU64(proto.AssetIDRecord_AssetIDRecord_Orientation, uint64(newRoleSprite.orientation))
		p.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WX, newRoleSprite.actionAnchorPointWX)
		p.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WY, newRoleSprite.actionAnchorPointWY)
	}
	p.UpdateWithAction()
}

// SwitchScene 切换场景 (带过渡动画)
func (p *Role) SwitchScene(mapID common.AssetID, wx, wy float32) {
	// 如果已经在过渡中，忽略新的切换请求
	if p.scene.transition.IsActive() || p.pendingScene != nil {
		return
	}
	// 获取目标地图的过渡配置
	targetMapCfg := cfg.GMapMgr.Maps.Get(mapID)
	// 预创建新场景 (使用目标地图的过渡配置)
	p.pendingScene = NewScene(mapID)
	p.pendingWX = wx
	p.pendingWY = wy
	// 当前场景转场退出 (使用目标地图的过渡配置)
	p.scene.transition.Out(targetMapCfg.SceneTransition)
}

// doSwitchScene 执行实际的场景切换 (无过渡动画)
func (p *Role) doSwitchScene(mapID common.AssetID, wx, wy float32) {
	p.scene = NewScene(mapID)
	p.camera = commoncamera.NewCamera()
	// 初始化角色的 World 坐标
	p.sprite.actionAnchorPointWX, p.sprite.actionAnchorPointWY = wx, wy
}

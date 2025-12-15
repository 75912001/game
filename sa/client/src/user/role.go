package user

import (
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/user/camera"
)

type Role struct {
	roleRecord *proto.RoleRecord // 角色-记录
	cfgRole    *cfg.Role         // 角色-配置

	frameIdx  uint32     // 当前帧索引
	frameTick uint32     // 帧计数器（用于控制动画速度）
	sprite    RoleSprite // 角色-精灵

	scene  *Scene         // 角色所在场景
	camera *camera.Camera // 角色摄像机

	// 场景切换相关
	pendingScene *Scene  // 待切换的新场景
	pendingTX    float32 // 待切换的目标 Tile X
	pendingTY    float32 // 待切换的目标 Tile Y
}

func NewRole(roleRecord *proto.RoleRecord) *Role {
	role := &Role{
		roleRecord: roleRecord,
	}
	roleAssetID := role.GetAssetID()
	role.cfgRole = cfg.GRoleMgr.Roles.Get(roleAssetID)
	if role.cfgRole == nil {
		// todo menglc  日志记录错误: 角色配置不存在 roleAssetID
		return nil
	}

	mapID := role.GetValueU32(proto.AssetIDRecord_AssetIDRecord_MapID)
	tx := role.GetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TX)
	ty := role.GetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TY)
	role.doSwitchScene(common.AssetID(mapID), tx, ty)

	role.UpdateWithAction()

	return role
}

// UpdateWithAction 更新角色状态和摄像机位置
func (p *Role) UpdateWithAction() {
	// 设置方向
	p.sprite.orientation = p.GetValueU32(proto.AssetIDRecord_AssetIDRecord_Orientation)

	// 设置图像
	images := p.cfgRole.ResRole.Move.Frames[p.sprite.orientation]
	p.sprite.image = images[p.frameIdx%uint32(len(images))]

	frames := p.cfgRole.ResRole.Move.FrameInfo[p.sprite.orientation]
	p.sprite.roleImageSprite = frames[p.frameIdx%uint32(len(frames))]

	// 更新角色中心 World 坐标 (脚底中心向上偏移半个角色高度)
	p.sprite.centerWX = p.sprite.bottomCenterWX
	p.sprite.centerWY = p.sprite.bottomCenterWY - float32(p.sprite.roleImageSprite.Frame.Height/2)

	// 更新摄像机跟随点 (World 坐标)
	p.camera.FollowX = int(p.sprite.centerWX)
	p.camera.FollowY = int(p.sprite.centerWY)

	// 计算摄像机视口左上角 (World 坐标)
	viewportX := p.camera.FollowX - cfg.GCommon.ScreenMaxWidth/2
	viewportY := p.camera.FollowY - cfg.GCommon.ScreenMaxHeight/2

	// 获取地图尺寸
	mapWidth, mapHeight := p.scene._map.tiledMapCfg.PixelW, p.scene._map.tiledMapCfg.PixelH

	// 限制摄像机不显示地图边界之外
	if mapWidth <= cfg.GCommon.ScreenMaxWidth {
		viewportX = -(cfg.GCommon.ScreenMaxWidth - mapWidth) / 2
	} else {
		if viewportX < 0 {
			viewportX = 0
		}
		if viewportX > mapWidth-cfg.GCommon.ScreenMaxWidth {
			viewportX = mapWidth - cfg.GCommon.ScreenMaxWidth
		}
	}

	if mapHeight <= cfg.GCommon.ScreenMaxHeight {
		viewportY = -(cfg.GCommon.ScreenMaxHeight - mapHeight) / 2
	} else {
		if viewportY < 0 {
			viewportY = 0
		}
		if viewportY > mapHeight-cfg.GCommon.ScreenMaxHeight {
			viewportY = mapHeight - cfg.GCommon.ScreenMaxHeight
		}
	}

	p.camera.ViewportX = viewportX
	p.camera.ViewportY = viewportY
}

// GetUUID 获取-角色UUID
func (p *Role) GetUUID() common.RoleUUID {
	return common.RoleUUID(p.roleRecord.UUID)
}

// GetAssetID 获取-角色资产ID
func (p *Role) GetAssetID() common.AssetID {
	v, _ := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_AssetID)]
	return common.AssetID(v)
}

// GetRoleNick 获取-角色昵称
func (p *Role) GetRoleNick() string {
	return p.roleRecord.Nick
}

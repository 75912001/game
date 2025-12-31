package user

import (
	"fmt"
	"saClient/src/cfg"
	"saClient/src/common"
	commoncamera "saClient/src/common/camera"
	"saClient/src/proto"
)

type Role struct {
	roleRecord *proto.RoleRecord // 角色-记录
	cfgRole    *cfg.Role         // 角色-配置

	animationFrame common.AnimationFrame
	sprite         RoleSprite // 角色-精灵

	scene  *Scene               // 角色所在场景
	camera *commoncamera.Camera // 角色摄像机

	// 场景切换相关
	pendingScene *Scene  // 待切换的新场景
	pendingWX    float32 // 待切换的目标 World X
	pendingWY    float32 // 待切换的目标 World Y

	BattleStats *RoleBattleStats

	ArpgAI *ArpgRoleAI // 角色战斗AI
}

func NewRole(roleRecord *proto.RoleRecord) *Role {
	role := &Role{
		roleRecord: roleRecord,
	}
	role.ArpgAI = NewArpgRoleAI(role)                     // 初始化战斗AI
	role.sprite.action = proto.RoleAction_RoleAction_Move // 默认动作-移动
	role.BattleStats = NewRoleBattleStats(role)
	roleAssetID := role.GetAssetID()
	role.cfgRole = cfg.GRoleMgr.Roles.Get(roleAssetID)
	if role.cfgRole == nil {
		// todo menglc  日志记录错误: 角色配置不存在 roleAssetID
		return nil
	}

	mapID := role.GetValueU32(proto.AssetIDRecord_AssetIDRecord_MapID)
	tx := role.GetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WX)
	ty := role.GetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WY)
	role.doSwitchScene(common.AssetID(mapID), tx, ty)
	role.sprite.orientation = proto.AssetOrientation(role.GetValueU64(proto.AssetIDRecord_AssetIDRecord_Orientation))
	role.UpdateWithAction()

	return role
}

// UpdateWithAction 更新角色状态和摄像机位置
// 根据当前动作和方向更新角色显示的帧图像
func (p *Role) UpdateWithAction() {
	// 获取当前动作的数据
	actionData := p.cfgRole.ResRole.Actions[p.sprite.action]
	if actionData == nil {
		panic(fmt.Sprintf("角色动作数据不存在 %v %v", p.GetAssetID(), p.sprite.action))
	}

	// 设置图像 (性能: 数组访问 O(1))
	images := actionData.Frames[p.sprite.orientation]
	p.sprite.image = images[p.animationFrame.FrameIdx%uint32(len(images))]

	frames := actionData.FrameInfo[p.sprite.orientation]
	p.sprite.roleImageSprite = frames[p.animationFrame.FrameIdx%uint32(len(frames))]

	// 更新角色中心 World 坐标 (脚底中心向上偏移半个角色高度)
	p.sprite.centerWX = p.sprite.wx
	p.sprite.centerWY = p.sprite.wy - float32(p.sprite.roleImageSprite.Frame.H/2)

	// 更新摄像机跟随点 (World 坐标)
	p.camera.FollowWX = int(p.sprite.centerWX)
	p.camera.FollowWY = int(p.sprite.centerWY)

	// 计算摄像机视口左上角 (World 坐标)
	viewportX := p.camera.FollowWX - cfg.GCommon.ScreenMaxWidth/2
	viewportY := p.camera.FollowWY - cfg.GCommon.ScreenMaxHeight/2

	// 获取地图尺寸
	mapWidth, mapHeight := p.scene._map.tiledMapCfg.WPixel, p.scene._map.tiledMapCfg.HPixel

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

	p.camera.ViewportWX = viewportX
	p.camera.ViewportWY = viewportY
}

// GetUUID 获取-角色UUID
func (p *Role) GetUUID() common.RoleUUID {
	return common.RoleUUID(p.roleRecord.UUID)
}

// GetAssetID 获取-角色资产ID
func (p *Role) GetAssetID() common.AssetID {
	v, _ := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AssetID)]
	return common.AssetID(v)
}

// GetRoleNick 获取-角色昵称
func (p *Role) GetRoleNick() string {
	return p.roleRecord.Nick
}

// GetWX 获取-角色世界坐标X
func (p *Role) GetWX() float32 {
	return p.sprite.wx
}

// SetAction 设置角色当前动作
// 切换动作时自动重置帧索引, 从头播放新动作
// 参数:
//   - action: 要切换到的动作类型
//
// 性能: O(1) 赋值操作
func (p *Role) SetAction(action proto.RoleAction) {
	if p.sprite.action != action {
		p.sprite.action = action
		p.animationFrame.Reset()
	}
}

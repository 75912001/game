package role

import (
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/user/scene"
)

type Role struct {
	roleRecord *proto.RoleRecord // 角色-记录
	cfgRole    *cfg.Role         // 角色-配置

	frameIdx uint32 // 当前帧索引

	debugDrawImageBounds bool // 是否画出图像边界(调试用)

	scene *scene.Scene // 角色所在场景
}

func NewRole(roleRecord *proto.RoleRecord) *Role {
	role := &Role{
		roleRecord:           roleRecord,
		debugDrawImageBounds: true,
	}
	roleAssetID := role.GetAssetID()
	role.cfgRole = cfg.GRoleMgr.Roles.Get(roleAssetID)
	if role.cfgRole == nil {
		// todo menglc  日志记录错误: 角色配置不存在 roleAssetID
		return nil
	}
	role.scene = scene.NewScene(common.AssetID(role.GetValueU32(proto.AssetIDRecord_AssetIDRecord_MapID)))
	if role.scene == nil {
		// todo menglc 日志记录错误: 场景创建失败 mapID
		return nil
	}
	return role
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

package role

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/user/role.del"
)

type Role struct {
	roleRecord *proto.RoleRecord // 角色-记录
	cfgRole    *cfg.Role         // 角色-配置

	frames []*ebitenv2.Image // 动画帧

	frameIdx uint32 // 当前帧索引

	imageMgr *role_del.ImageMgr // 图像-管理器

	debugDrawImageBounds bool // 是否画出图像边界(调试用)
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
	return role
}

// GetID 获取-角色ID
func (p *Role) GetID() common.UserID {
	return common.UserID(p.roleRecord.RoleID)
}

// GetAssetID 获取-角色资产ID
func (p *Role) GetAssetID() common.AssetID {
	v, _ := p.roleRecord.AssetRoleBaseMap[uint32(proto.AssetIDRoleBase_AssetIDRoleBase_RoleID)]
	return common.AssetID(v)
}

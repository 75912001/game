package user

import (
	"fmt"
	"time"

	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/user/role"
)

type User struct {
	userRecord *proto.UserRecord
	roleMgr    *role.RoleMgr
	role       *role.Role // 当前选择的角色
}

func NewUser() *User {
	user := &User{
		userRecord: &proto.UserRecord{
			RoleRecordMap: make(map[uint64]*proto.RoleRecord),
		},
		roleMgr: role.NewRoleMgr(),
	}
	return user
}

func (p *User) Login(account string, password string) error {
	_ = account
	_ = password
	// 初始化用户数据 - 模拟数据
	if true {
		var roleUUID common.RoleUUID = 1
		now := uint64(time.Now().Unix())
		p.userRecord.RoleRecordMap[uint64(roleUUID)] = &proto.RoleRecord{
			UUID:             uint64(roleUUID),
			Nick:             "role.uuid.1",
			AssetIDRecordMap: make(map[uint32]uint64),
			RecordMap:        map[uint64]*proto.RecordPrimary{},
			PetRecordMap:     map[uint64]*proto.PetRecord{},
		}
		roleRecord := p.userRecord.RoleRecordMap[uint64(roleUUID)]

		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Exp)] = 0                                                            // 经验值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = 9999                                                          // 生命值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_MP)] = uint64(cfg.GCommon.RoleMPMax)                                 // 魔法值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_CreateTimestamp)] = now                                              // 创建时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_LastLoginTimestamp)] = now                                           // 上次登录时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_LastLogoutTimestamp)] = 0                                            // 上次登出时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_MapID)] = 2000001                                                    // 所在地图ID (地图ID范围起始)
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_BottomCenterX)] = 2700                                               // 当前位置x
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_BottomCenterY)] = 1620                                               // 当前位置y
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Orientation)] = uint64(proto.AssetOrientation_AssetOrientation_Down) // 当前朝向-下
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Pose)] = uint64(proto.RoleAction_RoleAction_Stand)                   // 姿势 0:站立
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_AvailablePoint)] = uint64(cfg.GCommon.RoleInitialAvailablePoint)     // 可用点数
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_RebirthCount)] = 0                                                   // 转生次数
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_AssetID)] = 1000101                                                  // 资产ID (角色ID范围起始)

		// 元素属性 (101-104)
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalEarth)] = 0 // 元素属性-土
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalWater)] = 5 // 元素属性-水
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalFire)] = 5  // 元素属性-火
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalWind)] = 0  // 元素属性-风

		// 角色属性 (201-204)
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Strength)] = 10  // 腕力
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Endurance)] = 10 // 耐力
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Agility)] = 10   // 速度
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Stamina)] = 10   // 体力

		roleObject := role.NewRole(roleRecord)
		if roleObject == nil {
			return fmt.Errorf("new role failed, roleUUID %d", roleUUID)
		}
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = uint64(roleObject.GetHpMax()) // 将生命值设为最大值
		ok := p.roleMgr.Roles.AddIfNotExist(roleUUID, roleObject)
		if !ok {
			return fmt.Errorf("add role failed, roleUUID %d already exists", roleUUID)
		}
		p.role = roleObject
	}
	if false {
		var roleUUID common.RoleUUID = 2
		now := uint64(time.Now().Unix())
		p.userRecord.RoleRecordMap[uint64(roleUUID)] = &proto.RoleRecord{
			UUID:             uint64(roleUUID),
			Nick:             "role.uuid.2",
			AssetIDRecordMap: make(map[uint32]uint64),
			RecordMap:        map[uint64]*proto.RecordPrimary{},
			PetRecordMap:     map[uint64]*proto.PetRecord{},
		}
		roleRecord := p.userRecord.RoleRecordMap[uint64(roleUUID)]

		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Exp)] = 0                                                            // 经验值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = 9999                                                          // 生命值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_MP)] = uint64(cfg.GCommon.RoleMPMax)                                 // 魔法值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_CreateTimestamp)] = now                                              // 创建时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_LastLoginTimestamp)] = now                                           // 上次登录时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_LastLogoutTimestamp)] = 0                                            // 上次登出时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_MapID)] = 2000001                                                    // 所在地图ID (地图ID范围起始)
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_BottomCenterX)] = 200                                                // 当前位置x
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_BottomCenterY)] = 200                                                // 当前位置y
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Orientation)] = uint64(proto.AssetOrientation_AssetOrientation_Down) // 当前朝向-下
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Pose)] = uint64(proto.RoleAction_RoleAction_Stand)                   // 姿势 0:站立
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_AvailablePoint)] = uint64(cfg.GCommon.RoleInitialAvailablePoint)     // 可用点数
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_RebirthCount)] = 0                                                   // 转生次数
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_AssetID)] = 1000101                                                  // 资产ID (角色ID范围起始)

		// 元素属性 (101-104)
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalEarth)] = 0 // 元素属性-土
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalWater)] = 5 // 元素属性-水
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalFire)] = 5  // 元素属性-火
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalWind)] = 0  // 元素属性-风

		// 角色属性 (201-204)
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Strength)] = 10  // 腕力
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Endurance)] = 10 // 耐力
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Agility)] = 10   // 速度
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Stamina)] = 10   // 体力

		roleObject := role.NewRole(roleRecord)
		if roleObject == nil {
			return fmt.Errorf("new role failed, roleUUID %d", roleUUID)
		}
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = uint64(roleObject.GetHpMax()) // 将生命值设为最大值
		ok := p.roleMgr.Roles.AddIfNotExist(roleUUID, roleObject)
		if !ok {
			return fmt.Errorf("add role failed, roleUUID %d already exists", roleUUID)
		}
		p.role = roleObject
	}
	return nil
}

func (p *User) SelectRole(roleUUID common.RoleUUID) error {
	_ = roleUUID
	return nil
}

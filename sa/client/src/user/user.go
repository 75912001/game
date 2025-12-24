package user

import (
	"fmt"
	"time"

	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/proto"
)

type User struct {
	userRecord *proto.UserRecord
	roleMgr    *RoleMgr
	role       *Role // 当前选择的角色
}

var GUser *User

func NewUser() *User {
	user := &User{
		userRecord: &proto.UserRecord{
			RoleRecordMap: make(map[uint64]*proto.RoleRecord),
		},
		roleMgr: NewRoleMgr(),
	}
	GUser = user
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

		// 通用
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Exp)] = 0                                                            // 经验值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = 9999                                                          // 生命值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_MP)] = uint64(cfg.GCommon.RoleMPMax)                                 // 魔法值
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_CreateTimestamp)] = now                                              // 创建时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_MapID)] = 2000001                                                    // 所在地图ID (地图ID范围起始)
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WX)] = uint64(1600 * float32(common.Float32Ratio))      // 当前位置tx
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WY)] = uint64(800 * float32(common.Float32Ratio))       // 当前位置ty
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Orientation)] = uint64(proto.AssetOrientation_AssetOrientation_Down) // 当前朝向-下
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Pose)] = uint64(proto.RoleAction_RoleAction_Stand)                   // 姿势 0:站立
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_RebirthCount)] = 0                                                   // 转生次数

		// 元素
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalEarth)] = 0 // 元素属性-土
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalWater)] = 5 // 元素属性-水
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalFire)] = 5  // 元素属性-火
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_ElementalWind)] = 0  // 元素属性-风

		// 角色
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_LastLoginTimestamp)] = now                                       // 上次登录时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_LastLogoutTimestamp)] = 0                                        // 上次登出时间戳
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AvailablePoint)] = uint64(cfg.GCommon.RoleInitialAvailablePoint) // 可用点数
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AssetID)] = 1000101                                              // 资产ID (角色ID范围起始)

		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStrength)] = 10  // 腕力
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesEndurance)] = 10 // 耐力
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesAgility)] = 10   // 速度
		roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStamina)] = 10   // 体力

		roleObject := NewRole(roleRecord)
		if roleObject == nil {
			return fmt.Errorf("new role failed, roleUUID %d", roleUUID)
		}

		hpMax := roleObject.BattleStats.GetHpMax()
		roleObject.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = uint64(hpMax) // 将生命值设为最大值

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

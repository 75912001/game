package cfg

import (
	"os"
	"path/filepath"
	"saClient/src/common"
	"saClient/src/proto"

	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Common 通用配置
type Common struct {
	RoleCountMax                uint32  `yaml:"roleCountMax"`                // 允许创建的最大角色数
	RoleRebirthMax              uint32  `yaml:"roleRebirthCountMax"`         // 允许转生的最大次数
	RoleDefMoveSpeed            float32 `yaml:"roleDefMoveSpeed"`            // 角色默认移动速度
	RoleArpgDefViewRange        float32 `yaml:"roleArpgDefViewRange"`        // 角色-arpg-默认-视野范围
	RoleArpgDefAttackRangeAxe   float32 `yaml:"roleArpgDefAttackRangeAxe"`   // 角色-arpg-默认-攻击范围-斧头
	RoleArpgDefAttackRangeBow   float32 `yaml:"roleArpgDefAttackRangeBow"`   // 角色-arpg-默认-攻击范围-弓箭
	RoleArpgDefAttackRangeStaff float32 `yaml:"roleArpgDefAttackRangeStaff"` // 角色-arpg-默认-攻击范围-法杖
	RoleArpgDefCdTimeMsAxe      int64   `yaml:"roleArpgDefCdTimeMsAxe"`      // 角色-arpg-默认-cd时间-毫秒-斧头
	RoleArpgDefCdTimeMsBow      int64   `yaml:"roleArpgDefCdTimeMsBow"`      // 角色-arpg-默认-cd时间-毫秒-弓箭
	RoleArpgDefCdTimeMsStaff    int64   `yaml:"roleArpgDefCdTimeMsStaff"`    // 角色-arpg-默认-cd时间-毫秒-法杖
	RoleDefScale                float32 `yaml:"roleDefScale"`                // 角色默认缩放比例
	RoleLevelUpAvailablePoint   uint32  `yaml:"roleLevelUpAvailablePoint"`   // 角色每升一级,获得可用点数
	RoleMPMax                   uint32  `yaml:"roleMPMax"`                   // 角色属性-魔法值最大值
	RoleInitialAvailablePoint   uint32  `yaml:"roleInitialAvailablePoint"`   // 角色初始可分配点数
	RoleDefCritRate             float32 `yaml:"roleDefCritRate"`             // 角色-默认-暴击率
	RoleDefCounterRate          float32 `yaml:"roleDefCounterRate"`          // 角色-默认-反击率
	RoleDefDodgeRate            float32 `yaml:"roleDefDodgeRate"`            // 角色-默认-闪避率
	RoleDefHitRate              float32 `yaml:"roleDefHitRate"`              // 角色-默认-命中率
	RoleDefCritDamageBonusRate  float32 `yaml:"roleDefCritDamageBonusRate"`  // 角色-默认-暴击伤害加成倍数
	RoleDefStatusResistRate     float32 `yaml:"roleDefStatusResistRate"`     // 角色-默认-异常状态抗性比率

	PetRebirthCountMax              uint32  `yaml:"petRebirthCountMax"`              // 宠物最大转生次数
	PetDefMoveSpeed                 float32 `yaml:"petDefMoveSpeed"`                 // 宠物默认移动速度
	PetMPMax                        uint32  `yaml:"petMPMax"`                        // 宠物属性-MP最大值
	PetDefArpgViewRange             float32 `yaml:"petDefArpgViewRange"`             // 宠物-默认-Arpg视野范围
	PetDefArpgMoveSpeed             float32 `yaml:"petDefArpgMoveSpeed"`             // 宠物-默认-Arpg战斗中移动速度
	PetDefArpgChaseRadiusMultiplier float32 `yaml:"petDefArpgChaseRadiusMultiplier"` // 宠物-默认-arpg 追击半径乘数 [1.0, ∞) 追击半径 = 巡逻半径 * 追击半径乘数
	PetArpgDefAttackRange           float32 `yaml:"petArpgDefAttackRange"`           // 宠物-arpg-默认-攻击范围
	PetArpgDefCdTimeMs              int64   `yaml:"petArpgDefCdTimeMs"`              // 宠物-arpg-默认-攻击冷却时间(毫秒)
	PetDefScale                     float32 `yaml:"petDefScale"`                     // 宠物默认缩放比例
	PetDefCritRate                  float32 `yaml:"petDefCritRate"`                  // 宠物-默认-暴击率
	PetDefCounterRate               float32 `yaml:"petDefCounterRate"`               // 宠物-默认-反击率
	PetDefDodgeRate                 float32 `yaml:"petDefDodgeRate"`                 // 宠物-默认-闪避率
	PetDefHitRate                   float32 `yaml:"petDefHitRate"`                   // 宠物-默认-命中率
	PetDefCritDamageBonusRate       float32 `yaml:"petDefCritDamageBonusRate"`       // 宠物-默认-暴击伤害加成倍数
	PetDefStatusResistRate          float32 `yaml:"petDefStatusResistRate"`          // 宠物-默认-异常状态抗性比率

	WindowDefWidth  int `yaml:"windowDefWidth"`  // 窗口默认宽度
	WindowDefHeight int `yaml:"windowDefHeight"` // 窗口默认高度
	ScreenMaxWidth  int `yaml:"screenMaxWidth"`  // 屏幕最大宽度
	ScreenMaxHeight int `yaml:"screenMaxHeight"` // 屏幕最大高度
}

var GCommon = &Common{}

// Load 加载通用配置文件
func (p *Common) Load() error {
	// 构建配置文件路径
	cfgPath := filepath.Join(common.AppCfgDir, "common.yaml")
	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取通用配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	// 解析YAML
	err = yaml.Unmarshal(data, p)
	if err != nil {
		return errors.WithMessagef(err, "解析通用配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	return nil
}

// Check 检查配置
func (p *Common) Check() error {
	return nil
}

// Assemble 组装配置
func (p *Common) Assemble() error {
	return nil
}

// GetRoleArpgDefAttackRangeByWeaponType 获取-角色-arpg-默认-攻击范围-根据武器类型
func (p *Common) GetRoleArpgDefAttackRangeByWeaponType(weaponType proto.RoleWeaponType) float32 {
	switch weaponType {
	case proto.RoleWeaponType_RoleWeaponType_Axe:
		return p.RoleArpgDefAttackRangeAxe
	case proto.RoleWeaponType_RoleWeaponType_Bow:
		return p.RoleArpgDefAttackRangeBow
	case proto.RoleWeaponType_RoleWeaponType_Staff:
		return p.RoleArpgDefAttackRangeStaff
	default:
		return p.RoleArpgDefAttackRangeAxe
	}
}

// GetRoleArpgDefCdTimeMsByWeaponType 获取-角色-arpg-默认-cd时间-毫秒-根据武器类型
func (p *Common) GetRoleArpgDefCdTimeMsByWeaponType(weaponType proto.RoleWeaponType) int64 {
	switch weaponType {
	case proto.RoleWeaponType_RoleWeaponType_Axe:
		return p.RoleArpgDefCdTimeMsAxe
	case proto.RoleWeaponType_RoleWeaponType_Bow:
		return p.RoleArpgDefCdTimeMsBow
	case proto.RoleWeaponType_RoleWeaponType_Staff:
		return p.RoleArpgDefCdTimeMsStaff
	default:
		return p.RoleArpgDefCdTimeMsAxe
	}
}

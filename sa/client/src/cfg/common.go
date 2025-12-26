package cfg

import (
	"os"
	"path/filepath"
	"saClient/src/common"

	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Common 通用配置
type Common struct {
	RoleCountMax               uint32  `yaml:"roleCountMax"`               // 允许创建的最大角色数
	RoleRebirthMax             uint32  `yaml:"roleRebirthCountMax"`        // 允许转生的最大次数
	RoleDefMoveSpeed           float32 `yaml:"roleDefMoveSpeed"`           // 角色默认移动速度
	RoleDefScale               float32 `yaml:"roleDefScale"`               // 角色默认缩放比例
	RoleLevelUpAvailablePoint  uint32  `yaml:"roleLevelUpAvailablePoint"`  // 角色每升一级,获得可用点数
	RoleMPMax                  uint32  `yaml:"roleMPMax"`                  // 角色属性-魔法值最大值
	RoleInitialAvailablePoint  uint32  `yaml:"roleInitialAvailablePoint"`  // 角色初始可分配点数
	RoleDefCritRate            float32 `yaml:"roleDefCritRate"`            // 角色-默认-暴击率
	RoleDefCounterRate         float32 `yaml:"roleDefCounterRate"`         // 角色-默认-反击率
	RoleDefDodgeRate           float32 `yaml:"roleDefDodgeRate"`           // 角色-默认-闪避率
	RoleDefHitRate             float32 `yaml:"roleDefHitRate"`             // 角色-默认-命中率
	RoleDefCritDamageBonusRate float32 `yaml:"roleDefCritDamageBonusRate"` // 角色-默认-暴击伤害加成倍数
	RoleDefStatusResistRate    float32 `yaml:"roleDefStatusResistRate"`    // 角色-默认-异常状态抗性比率

	PetRebirthCountMax        uint32  `yaml:"petRebirthCountMax"`        // 宠物最大转生次数
	PetDefMoveSpeed           float32 `yaml:"petDefMoveSpeed"`           // 宠物默认移动速度
	PetMPMax                  uint32  `yaml:"petMPMax"`                  // 宠物属性-MP最大值
	PetDefViewRange           float32 `yaml:"petDefViewRange"`           // 宠物-默认-视野范围
	PetDefArpgMoveSpeed       float32 `yaml:"petDefArpgMoveSpeed"`       // 宠物-默认-战斗中移动速度
	PetDefCritRate            float32 `yaml:"petDefCritRate"`            // 宠物-默认-暴击率
	PetDefCounterRate         float32 `yaml:"petDefCounterRate"`         // 宠物-默认-反击率
	PetDefDodgeRate           float32 `yaml:"petDefDodgeRate"`           // 宠物-默认-闪避率
	PetDefHitRate             float32 `yaml:"petDefHitRate"`             // 宠物-默认-命中率
	PetDefCritDamageBonusRate float32 `yaml:"petDefCritDamageBonusRate"` // 宠物-默认-暴击伤害加成倍数
	PetDefStatusResistRate    float32 `yaml:"petDefStatusResistRate"`    // 宠物-默认-异常状态抗性比率

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

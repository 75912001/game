package cfg

import (
	"fmt"
	xutil "github.com/75912001/xlib/util"
	"os"
	"path/filepath"
	"saClient/src/common"
	commonelemental "saClient/src/common/elemental"
	"saClient/src/proto"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// PetElemental 宠物元素属性
type PetElemental struct {
	Earth uint32 `yaml:"earth"` // 土
	Water uint32 `yaml:"water"` // 水
	Fire  uint32 `yaml:"fire"`  // 火
	Wind  uint32 `yaml:"wind"`  // 风
}

// PetAttributes 宠物基础属性
type PetAttributes struct {
	Attack          uint32  `yaml:"attack"`            // 攻击
	Defense         uint32  `yaml:"defense"`           // 防御
	Agility         uint32  `yaml:"agility"`           // 敏捷
	HP              uint32  `yaml:"hp"`                // 生命
	CritRate        float32 `yaml:"crit_rate"`         // 暴击率
	CounterRate     float32 `yaml:"counter_rate"`      // 反击率
	DodgeRate       float32 `yaml:"dodge_rate"`        // 闪避率
	HitRate         float32 `yaml:"hit_rate"`          // 命中率
	CritDamageBonus float32 `yaml:"crit_damage_bonus"` // 暴击伤害加成倍数
	StatusResist    float32 `yaml:"status_resist"`     // 异常状态抗性比率
}

// GrowthRange 成长范围
type GrowthRange struct {
	Min float32 `yaml:"min"` // 最小值
	Max float32 `yaml:"max"` // 最大值
}

// PetGrowth 宠物成长配置
type PetGrowth struct {
	Attack  GrowthRange `yaml:"attack"`  // 攻击成长
	Defense GrowthRange `yaml:"defense"` // 防御成长
	Agility GrowthRange `yaml:"agility"` // 敏捷成长
	HP      GrowthRange `yaml:"hp"`      // 生命成长
}

// Pet 宠物信息
type Pet struct {
	ID          common.AssetID  `yaml:"id"`          // 宠物ID
	Name        string          `yaml:"name"`        // 名称
	Rarity      proto.PetRarity `yaml:"rarity"`      // 稀有度: 1-普通, 2-稀有, 3-史诗, 4-传说, 5-神话
	Description string          `yaml:"description"` // 描述
	Elemental   PetElemental    `yaml:"elemental"`   // 元素属性
	Attributes  PetAttributes   `yaml:"attributes"`  // 基础属性
	Growth      PetGrowth       `yaml:"growth"`      // 成长配置

	// todo menglc 装配出宠物的其他配置
}

var GPetMgr = newPetMgr()

// PetMgr 配置管理器
type PetMgr struct {
	Pets *xmap.MapMgr[common.AssetID, *Pet] // key: 宠物ID
}

func newPetMgr() *PetMgr {
	return &PetMgr{
		Pets: xmap.NewMapMgr[common.AssetID, *Pet](),
	}
}

// Load 加载配置文件
func (p *PetMgr) Load() error {
	// 构建配置文件路径
	cfgPath := filepath.Join(common.AppCfgDir, "pet.yaml")
	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取宠物配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	// 定义中间结构用于YAML解析
	var config struct {
		Pets []*Pet `yaml:"pets"`
	}
	// 解析YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析宠物配置文件失败: %v %v", cfgPath, err)
	}
	// 加载每个宠物
	for _, pet := range config.Pets {
		if pet.ID < common.AssetID(proto.AssetIDRange_AssetIDRange_Pet_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Pet_End) < pet.ID {
			return fmt.Errorf("宠物ID超出范围: %d %v", pet.ID, xruntime.Location())
		}
		if pet.Rarity <= proto.PetRarity_PetRarity_Unknow || proto.PetRarity_PetRarity_Max <= pet.Rarity { // 稀有度错误
			return fmt.Errorf("宠物稀有度配置错误: %v %v %v", pet.ID, pet.Rarity, xruntime.Location())
		}
		// 设置默认值

		if xutil.Float32Equal(pet.Attributes.CritRate, 0) {
			pet.Attributes.CritRate = 0.1
		}
		if xutil.Float32Equal(pet.Attributes.CounterRate, 0) {
			pet.Attributes.CounterRate = 0.1
		}
		if xutil.Float32Equal(pet.Attributes.DodgeRate, 0) {
			pet.Attributes.DodgeRate = 0.1
		}
		if xutil.Float32Equal(pet.Attributes.HitRate, 0) {
			pet.Attributes.HitRate = 0.9
		}
		if xutil.Float32Equal(pet.Attributes.CritDamageBonus, 0) {
			pet.Attributes.CritDamageBonus = 0.5
		}
		if xutil.Float32Equal(pet.Attributes.StatusResist, 0) {
			pet.Attributes.StatusResist = 0.2
		}
		// 验证战斗属性范围 [0,1]
		if xutil.Float32Less(pet.Attributes.CritRate, 0) || xutil.Float32Less(1, pet.Attributes.CritRate) {
			return fmt.Errorf("宠物ID %d 暴击率超出范围[0,1]: %v %v", pet.ID, pet.Attributes.CritRate, xruntime.Location())
		}
		if xutil.Float32Less(pet.Attributes.CounterRate, 0) || xutil.Float32Less(1, pet.Attributes.CounterRate) {
			return fmt.Errorf("宠物ID %d 反击率超出范围[0,1]: %v %v", pet.ID, pet.Attributes.CounterRate, xruntime.Location())
		}
		if xutil.Float32Less(pet.Attributes.DodgeRate, 0) || xutil.Float32Less(1, pet.Attributes.DodgeRate) {
			return fmt.Errorf("宠物ID %d 闪避率超出范围[0,1]: %v %v", pet.ID, pet.Attributes.DodgeRate, xruntime.Location())
		}
		if xutil.Float32Less(pet.Attributes.HitRate, 0) || xutil.Float32Less(1, pet.Attributes.HitRate) {
			return fmt.Errorf("宠物ID %d 命中率超出范围[0,1]: %v %v", pet.ID, pet.Attributes.HitRate, xruntime.Location())
		}
		if xutil.Float32Less(pet.Attributes.CritDamageBonus, 0) || xutil.Float32Less(1, pet.Attributes.CritDamageBonus) {
			return fmt.Errorf("宠物ID %d 暴击伤害加成超出范围[0,1]: %v %v", pet.ID, pet.Attributes.CritDamageBonus, xruntime.Location())
		}
		if xutil.Float32Less(pet.Attributes.StatusResist, 0) || xutil.Float32Less(1, pet.Attributes.StatusResist) {
			return fmt.Errorf("宠物ID %d 异常状态抗性超出范围[0,1]: %v %v", pet.ID, pet.Attributes.StatusResist, xruntime.Location())
		}
		ok := p.Pets.AddIfNotExist(pet.ID, pet)
		if !ok { // 添加失败
			return fmt.Errorf("添加宠物失败,宠物已存在: %v %v", pet.ID, xruntime.Location())
		}
	}
	p.Pets.Foreach(
		func(id common.AssetID, pet *Pet) bool { // 验证元素属性是否合法
			err = commonelemental.Validate(common.AssetQuantity(pet.Elemental.Earth), common.AssetQuantity(pet.Elemental.Water), common.AssetQuantity(pet.Elemental.Fire), common.AssetQuantity(pet.Elemental.Wind))
			if err != nil {
				err = fmt.Errorf("宠物ID %d 元素属性验证失败: %v", pet.ID, err)
				return false
			}
			return true
		},
	)
	return nil
}

// Check 检查配置
func (p *PetMgr) Check() error {
	return nil
}

// Assemble 组装配置
func (p *PetMgr) Assemble() error {
	return nil
}

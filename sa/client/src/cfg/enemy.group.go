package cfg

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"saClient/src/common"
	"saClient/src/proto"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// EnemyGroupEnemy 敌人组中的敌人条目
type EnemyGroupEnemy struct {
	ID     common.AssetID `yaml:"id"`     // 引用 pet.yaml 中的宠物ID
	Weight uint32         `yaml:"weight"` // 权重(加权池) [默认:0] 0:必定出现
	Level  uint32         `yaml:"level"`  // 指定敌人等级 [1,140] 不配置则按组等级规则计算

	Pet *Pet // 运行时字段,引用的宠物配置
}

// EnemyGroup 敌人组
type EnemyGroup struct {
	ID              uint32             `yaml:"id"`              // 敌人组ID
	Name            string             `yaml:"name"`            // 名称
	IsBoss          bool               `yaml:"isBoss"`          // 是否是Boss组 [默认:false]
	CountRange      []uint32           `yaml:"countRange"`      // 敌人数量范围 [min, max] [1,10]之内
	LevelRange      []uint32           `yaml:"levelRange"`      // 敌人等级范围 [min, max] [1,140]之内
	RoleLevelOffset []int              `yaml:"roleLevelOffset"` // 角色等级偏移 [min, max]
	Captured        *bool              `yaml:"captured"`        // 是否允许捕获 [默认:true]
	BabyRate        uint32             `yaml:"babyRate"`        // 宝宝概率(十万分率) [默认:0]
	Enemies         []*EnemyGroupEnemy `yaml:"enemies"`         // 敌人列表

	// 运行时字段
	MustAppearEnemies []*EnemyGroupEnemy // weight=0 的必定出现敌人
	WeightedEnemies   []*EnemyGroupEnemy // weight>0 的加权池敌人
	TotalWeight       uint32             // 加权池总权重
}

var GEnemyGroupMgr = newEnemyGroupMgr()

// EnemyGroupMgr 敌人组配置管理器
type EnemyGroupMgr struct {
	EnemyGroups *xmap.MapMgr[uint32, *EnemyGroup] // key: 敌人组ID
}

func newEnemyGroupMgr() *EnemyGroupMgr {
	return &EnemyGroupMgr{
		EnemyGroups: xmap.NewMapMgr[uint32, *EnemyGroup](),
	}
}

// Load 加载配置文件
func (p *EnemyGroupMgr) Load() error {
	cfgPath := filepath.Join(common.AppCfgDir, "enemy.group.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取敌人组配置文件失败: %v %v", cfgPath, xruntime.Location())
	}

	var config struct {
		EnemyGroups []*EnemyGroup `yaml:"enemyGroups"`
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析敌人组配置文件失败: %v %v", cfgPath, xruntime.Location())
	}

	for _, group := range config.EnemyGroups {
		if err := p.validateGroup(group); err != nil {
			return err
		}
		p.processDef(group)
		p.classifyEnemies(group)

		ok := p.EnemyGroups.AddIfNotExist(group.ID, group)
		if !ok {
			return fmt.Errorf("添加敌人组失败,敌人组已存在: %v %v", group.ID, xruntime.Location())
		}
	}
	return nil
}

// validateGroup 验证敌人组配置
func (p *EnemyGroupMgr) validateGroup(group *EnemyGroup) error {
	// 验证敌人列表非空
	if len(group.Enemies) == 0 {
		return fmt.Errorf("敌人组ID %d 的敌人列表为空 %v", group.ID, xruntime.Location())
	}

	if group.IsBoss {
		return p.validateBossGroup(group)
	}
	return p.validateNormalGroup(group)
}

// validateBossGroup 验证Boss组配置
func (p *EnemyGroupMgr) validateBossGroup(group *EnemyGroup) error {
	// Boss组: 每个敌人必须指定等级
	for i, enemy := range group.Enemies {
		if enemy.Level == 0 {
			return fmt.Errorf("敌人组ID %d (Boss) 的敌人[%d]必须指定等级 %v", group.ID, i, xruntime.Location())
		}
		if enemy.Level < uint32(proto.LevelRange_LevelRange_Min) || uint32(proto.LevelRange_LevelRange_Max) < enemy.Level {
			return fmt.Errorf("敌人组ID %d (Boss) 的敌人[%d]等级超出范围[%d,%d]: %d %v",
				group.ID, i, proto.LevelRange_LevelRange_Min, proto.LevelRange_LevelRange_Max, enemy.Level, xruntime.Location())
		}
	}
	return nil
}

// validateNormalGroup 验证普通组配置
func (p *EnemyGroupMgr) validateNormalGroup(group *EnemyGroup) error {
	// 验证 countRange
	if len(group.CountRange) != 2 {
		return fmt.Errorf("敌人组ID %d 的countRange必须是[min,max]格式 %v", group.ID, xruntime.Location())
	}
	if group.CountRange[0] < uint32(proto.CombatEnemyGroupEnemyCountRange_CombatEnemyGroupEnemyCountRange_Min) ||
		uint32(proto.CombatEnemyGroupEnemyCountRange_CombatEnemyGroupEnemyCountRange_Max) < group.CountRange[1] {
		return fmt.Errorf("敌人组ID %d 的countRange超出范围[%d,%d]: %v %v",
			group.ID, proto.CombatEnemyGroupEnemyCountRange_CombatEnemyGroupEnemyCountRange_Min, proto.CombatEnemyGroupEnemyCountRange_CombatEnemyGroupEnemyCountRange_Max, group.CountRange, xruntime.Location())
	}
	if group.CountRange[1] < group.CountRange[0] {
		return fmt.Errorf("敌人组ID %d 的countRange[0]不能大于countRange[1]: %v %v",
			group.ID, group.CountRange, xruntime.Location())
	}

	// 验证 levelRange/roleLevelOffset 二选一
	hasLevelRange := len(group.LevelRange) == 2
	hasRoleLevelOffset := len(group.RoleLevelOffset) == 2
	if !hasLevelRange && !hasRoleLevelOffset {
		return fmt.Errorf("敌人组ID %d 必须配置levelRange或roleLevelOffset %v", group.ID, xruntime.Location())
	}

	// 验证 levelRange 范围
	if hasLevelRange {
		if group.LevelRange[0] < uint32(proto.LevelRange_LevelRange_Min) || uint32(proto.LevelRange_LevelRange_Max) < group.LevelRange[1] {
			return fmt.Errorf("敌人组ID %d 的levelRange超出范围[%d,%d]: %v %v",
				group.ID, proto.LevelRange_LevelRange_Min, proto.LevelRange_LevelRange_Max, group.LevelRange, xruntime.Location())
		}
		if group.LevelRange[1] < group.LevelRange[0] {
			return fmt.Errorf("敌人组ID %d 的levelRange[0]不能大于levelRange[1]: %v %v",
				group.ID, group.LevelRange, xruntime.Location())
		}
	}

	// 验证 babyRate 范围
	if group.BabyRate < uint32(proto.CombatEnemyGroupBabyRate_CombatEnemyGroupBabyRate_Min) || uint32(proto.CombatEnemyGroupBabyRate_CombatEnemyGroupBabyRate_Max) < group.BabyRate {
		return fmt.Errorf("敌人组ID %d 的babyRate超出范围[%d,%d]: %d %v",
			group.ID, proto.CombatEnemyGroupBabyRate_CombatEnemyGroupBabyRate_Min, proto.CombatEnemyGroupBabyRate_CombatEnemyGroupBabyRate_Max, group.BabyRate, xruntime.Location())
	}

	// 验证敌人条目
	for i, enemy := range group.Enemies {
		if enemy.Level != 0 && (enemy.Level < uint32(proto.LevelRange_LevelRange_Min) || uint32(proto.LevelRange_LevelRange_Max) < enemy.Level) {
			return fmt.Errorf("敌人组ID %d 的敌人[%d]等级超出范围[%d,%d]: %d %v",
				group.ID, i, proto.LevelRange_LevelRange_Min, proto.LevelRange_LevelRange_Max, enemy.Level, xruntime.Location())
		}
	}
	return nil
}

// processDef 处理默认值
func (p *EnemyGroupMgr) processDef(group *EnemyGroup) {
	// captured 默认为 true
	if group.Captured == nil {
		var defCaptured = true
		group.Captured = &defCaptured
	}

	// Boss组强制不可捕获
	if group.IsBoss {
		var defCaptured = false
		group.Captured = &defCaptured
		group.BabyRate = 0
	}
}

// classifyEnemies 分类敌人(必定出现/加权池)
func (p *EnemyGroupMgr) classifyEnemies(group *EnemyGroup) {
	group.MustAppearEnemies = make([]*EnemyGroupEnemy, 0)
	group.WeightedEnemies = make([]*EnemyGroupEnemy, 0)
	group.TotalWeight = 0

	for _, enemy := range group.Enemies {
		if enemy.Weight == 0 {
			group.MustAppearEnemies = append(group.MustAppearEnemies, enemy)
		} else {
			group.WeightedEnemies = append(group.WeightedEnemies, enemy)
			group.TotalWeight += enemy.Weight
		}
	}
}

// Check 检查配置
func (p *EnemyGroupMgr) Check() error {
	var err error
	p.EnemyGroups.Foreach(func(id uint32, group *EnemyGroup) bool {
		// 检查敌人ID是否在pet配置中存在
		for i, enemy := range group.Enemies {
			enemy.Pet = GPetMgr.Pets.Get(enemy.ID)
			if enemy.Pet == nil {
				err = fmt.Errorf("敌人组ID %d 的敌人[%d]引用的宠物ID %d 不存在 %v",
					group.ID, i, enemy.ID, xruntime.Location())
				return false
			}
		}
		if !group.IsBoss { // 普通组
			// 如果随机的最大数量,大于必选出现的数量,且非必选数量为0,则代表无法选够所需数量,报错
			if len(group.MustAppearEnemies) < int(group.CountRange[1]) && len(group.WeightedEnemies) == 0 {
				err = fmt.Errorf("敌人组ID %d 的必定出现敌人数量 %d 少于最大敌人数量 %d,且无加权池敌人无法补足 %v",
					group.ID, len(group.MustAppearEnemies), group.CountRange[1], xruntime.Location())
				return false
			}
		}

		return true
	})
	return err
}

// Assemble 组装配置
func (p *EnemyGroupMgr) Assemble() error {
	return nil
}

// GeneratedEnemy 生成的敌人实例
type GeneratedEnemy struct {
	Config *EnemyGroupEnemy // 配置引用
	Level  uint32           // 实际等级(1级即为宝宝)
}

// Generate 根据敌人组配置生成敌人数组
// group: 敌人组配置
// roleLevel: 角色等级
// 返回: 生成的敌人数组
func (group *EnemyGroup) Generate(roleLevel uint32) []*GeneratedEnemy {
	if group.IsBoss {
		return group.generateBoss()
	}
	return group.generateNormal(roleLevel)
}

// generateBoss 生成Boss组敌人
func (group *EnemyGroup) generateBoss() []*GeneratedEnemy {
	result := make([]*GeneratedEnemy, 0, len(group.Enemies))
	for _, enemy := range group.Enemies {
		result = append(result, &GeneratedEnemy{
			Config: enemy,
			Level:  enemy.Level,
		})
	}
	return result
}

// generateNormal 生成普通组敌人
func (group *EnemyGroup) generateNormal(roleLevel uint32) []*GeneratedEnemy {
	// 计算敌人数量
	count := xutil.RandomInt(int(group.CountRange[0]), int(group.CountRange[1]))
	result := make([]*GeneratedEnemy, 0, count)

	// 1. 先放置必定出现的敌人(weight=0)
	mustAppearIdx := 0
	for len(result) < count && mustAppearIdx < len(group.MustAppearEnemies) {
		enemy := group.MustAppearEnemies[mustAppearIdx]
		level := group.calculateLevel(enemy, roleLevel)
		if group.rollBaby() {
			level = 1
		}
		result = append(result, &GeneratedEnemy{
			Config: enemy,
			Level:  level,
		})
		mustAppearIdx++
	}

	// 2. 剩余位置从加权池随机选择
	for len(result) < count {
		enemy := group.randomWeightedEnemy()
		level := group.calculateLevel(enemy, roleLevel)
		if group.rollBaby() {
			level = 1
		}
		result = append(result, &GeneratedEnemy{
			Config: enemy,
			Level:  level,
		})
	}

	return result
}

// calculateLevel 计算敌人等级
func (group *EnemyGroup) calculateLevel(enemy *EnemyGroupEnemy, roleLevel uint32) uint32 {
	// 如果敌人配置了固定等级，使用固定等级
	if enemy.Level != 0 {
		return enemy.Level
	}
	if len(group.LevelRange) == 2 { // levelRange
		return uint32(xutil.RandomInt(int(group.LevelRange[0]), int(group.LevelRange[1])))
	} else { // roleLevelOffset
		offset := xutil.RandomInt(group.RoleLevelOffset[0], group.RoleLevelOffset[1])
		calcLevel := int(roleLevel) + offset
		if calcLevel < int(proto.LevelRange_LevelRange_Min) {
			calcLevel = int(proto.LevelRange_LevelRange_Min)
		}
		if int(proto.LevelRange_LevelRange_Max) < calcLevel {
			calcLevel = int(proto.LevelRange_LevelRange_Max)
		}
		return uint32(calcLevel)
	}
}

// randomWeightedEnemy 从加权池随机选择敌人
func (group *EnemyGroup) randomWeightedEnemy() *EnemyGroupEnemy {
	if group.TotalWeight == 0 || len(group.WeightedEnemies) == 0 {
		return nil
	}

	roll := uint32(rand.Intn(int(group.TotalWeight)))
	var cumulative uint32 = 0
	for _, enemy := range group.WeightedEnemies {
		cumulative += enemy.Weight
		if roll < cumulative {
			return enemy
		}
	}
	return group.WeightedEnemies[len(group.WeightedEnemies)-1]
}

// rollBaby 判定是否成为宝宝
func (group *EnemyGroup) rollBaby() bool {
	if group.BabyRate <= 0 {
		return false
	}
	if uint32(proto.CombatEnemyGroupBabyRate_CombatEnemyGroupBabyRate_Max) <= group.BabyRate {
		return true
	}
	return rand.Intn(int(proto.CombatEnemyGroupBabyRate_CombatEnemyGroupBabyRate_Max)) < int(group.BabyRate)
}

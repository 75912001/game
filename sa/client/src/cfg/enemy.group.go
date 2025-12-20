package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"saClient/src/common"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// todo menglc 放入 proto 中
const (
	EnemyGroupCountMin    = 1      // 敌人数量最小值
	EnemyGroupCountMax    = 10     // 敌人数量最大值
	EnemyGroupLevelMin    = 1      // 敌人等级最小值
	EnemyGroupLevelMax    = 140    // 敌人等级最大值
	EnemyGroupBabyRateMin = 0      // 宝宝概率最小值
	EnemyGroupBabyRateMax = 100000 // 宝宝概率最大值(十万分率,100000=100%)
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
	BabyRate        int                `yaml:"babyRate"`        // 宝宝概率(十万分率) [默认:0]
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
func (m *EnemyGroupMgr) Load() error {
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
		if err := m.validateGroup(group); err != nil {
			return err
		}
		m.processDefaults(group)
		m.classifyEnemies(group)

		ok := m.EnemyGroups.AddIfNotExist(group.ID, group)
		if !ok {
			return fmt.Errorf("添加敌人组失败,敌人组已存在: %v %v", group.ID, xruntime.Location())
		}
	}
	return nil
}

// validateGroup 验证敌人组配置
func (m *EnemyGroupMgr) validateGroup(group *EnemyGroup) error {
	// 验证敌人列表非空
	if len(group.Enemies) == 0 {
		return fmt.Errorf("敌人组ID %d 的敌人列表为空 %v", group.ID, xruntime.Location())
	}

	if group.IsBoss {
		return m.validateBossGroup(group)
	}
	return m.validateNormalGroup(group)
}

// validateBossGroup 验证Boss组配置
func (m *EnemyGroupMgr) validateBossGroup(group *EnemyGroup) error {
	// Boss组: 每个敌人必须指定等级
	for i, enemy := range group.Enemies {
		if enemy.Level == 0 {
			return fmt.Errorf("敌人组ID %d (Boss) 的敌人[%d]必须指定等级 %v", group.ID, i, xruntime.Location())
		}
		if enemy.Level < EnemyGroupLevelMin || EnemyGroupLevelMax < enemy.Level {
			return fmt.Errorf("敌人组ID %d (Boss) 的敌人[%d]等级超出范围[%d,%d]: %d %v",
				group.ID, i, EnemyGroupLevelMin, EnemyGroupLevelMax, enemy.Level, xruntime.Location())
		}
	}
	return nil
}

// validateNormalGroup 验证普通组配置
func (m *EnemyGroupMgr) validateNormalGroup(group *EnemyGroup) error {
	// 验证 countRange
	if len(group.CountRange) != 2 {
		return fmt.Errorf("敌人组ID %d 的countRange必须是[min,max]格式 %v", group.ID, xruntime.Location())
	}
	if group.CountRange[0] < EnemyGroupCountMin || EnemyGroupCountMax < group.CountRange[1] {
		return fmt.Errorf("敌人组ID %d 的countRange超出范围[%d,%d]: %v %v",
			group.ID, EnemyGroupCountMin, EnemyGroupCountMax, group.CountRange, xruntime.Location())
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
		if group.LevelRange[0] < EnemyGroupLevelMin || EnemyGroupLevelMax < group.LevelRange[1] {
			return fmt.Errorf("敌人组ID %d 的levelRange超出范围[%d,%d]: %v %v",
				group.ID, EnemyGroupLevelMin, EnemyGroupLevelMax, group.LevelRange, xruntime.Location())
		}
		if group.LevelRange[1] < group.LevelRange[0] {
			return fmt.Errorf("敌人组ID %d 的levelRange[0]不能大于levelRange[1]: %v %v",
				group.ID, group.LevelRange, xruntime.Location())
		}
	}

	// 验证 babyRate 范围
	if group.BabyRate < EnemyGroupBabyRateMin || EnemyGroupBabyRateMax < group.BabyRate {
		return fmt.Errorf("敌人组ID %d 的babyRate超出范围[%d,%d]: %d %v",
			group.ID, EnemyGroupBabyRateMin, EnemyGroupBabyRateMax, group.BabyRate, xruntime.Location())
	}

	// 验证敌人条目
	for i, enemy := range group.Enemies {
		if enemy.Level != 0 && (enemy.Level < EnemyGroupLevelMin || EnemyGroupLevelMax < enemy.Level) {
			return fmt.Errorf("敌人组ID %d 的敌人[%d]等级超出范围[%d,%d]: %d %v",
				group.ID, i, EnemyGroupLevelMin, EnemyGroupLevelMax, enemy.Level, xruntime.Location())
		}
	}
	return nil
}

// processDefaults 处理默认值
func (m *EnemyGroupMgr) processDefaults(group *EnemyGroup) {
	// captured 默认为 true
	if group.Captured == nil {
		var defaultCaptured = true
		group.Captured = &defaultCaptured
	}

	// Boss组强制不可捕获
	if group.IsBoss {
		var defaultCaptured = false
		group.Captured = &defaultCaptured
		group.BabyRate = 0
	}
}

// classifyEnemies 分类敌人(必定出现/加权池)
func (m *EnemyGroupMgr) classifyEnemies(group *EnemyGroup) {
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
func (m *EnemyGroupMgr) Check() error {
	var err error
	m.EnemyGroups.Foreach(func(id uint32, group *EnemyGroup) bool {
		// 检查敌人ID是否在pet配置中存在
		for i, enemy := range group.Enemies {
			enemy.Pet = GPetMgr.Pets.Get(enemy.ID)
			if enemy.Pet == nil {
				err = fmt.Errorf("敌人组ID %d 的敌人[%d]引用的宠物ID %d 不存在 %v",
					group.ID, i, enemy.ID, xruntime.Location())
				return false
			}
		}
		return true
	})
	return err
}

// Assemble 组装配置
func (m *EnemyGroupMgr) Assemble() error {
	return nil
}

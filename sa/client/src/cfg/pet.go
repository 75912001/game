package cfg

import (
	"encoding/json"
	"fmt"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"

	"os"
	"path/filepath"
	"saClient/src/common"
	"saClient/src/elemental"
	"saClient/src/proto"
)

// PetElemental 宠物元素属性
type PetElemental struct {
	Earth common.AssetQuantity `json:"earth"` // 土
	Water common.AssetQuantity `json:"water"` // 水
	Fire  common.AssetQuantity `json:"fire"`  // 火
	Wind  common.AssetQuantity `json:"wind"`  // 风
}

// PetBasic 宠物基础属性
type PetBasic struct {
	Attack  common.AssetQuantity `json:"attack"`  // 攻击
	Defense common.AssetQuantity `json:"defense"` // 防御
	Agility common.AssetQuantity `json:"agility"` // 敏捷
	HP      common.AssetQuantity `json:"hp"`      // 生命
}

// Pet 宠物信息
type Pet struct {
	ID          common.AssetID  `json:"id"`          // 宠物ID
	Name        string          `json:"name"`        // 名称
	Rarity      proto.PetRarity `json:"rarity"`      // 稀有度: 1-普通, 2-稀有, 3-史诗, 4-传说, 5-神话
	Description string          `json:"description"` // 描述
	Elemental   PetElemental    `json:"elemental"`   // 元素属性
	Basic       PetBasic        `json:"basic"`       // 基础属性

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
	cfgPath := filepath.Join(common.AppCfgDir, "pet.json")
	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取宠物配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	// 定义中间结构用于JSON解析
	var config struct {
		Pets []*Pet `json:"pets"`
	}
	// 解析JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析宠物配置文件失败: %v %v", cfgPath, err)
	}
	// 加载每个宠物
	for _, pet := range config.Pets {
		if pet.ID < common.AssetID(proto.AssetIDRange_AssetIDRange_Pet_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Pet_End) < pet.ID {
			return fmt.Errorf("宠物ID超出范围: %d %v", pet.ID, xruntime.Location())
		}
		// todo menglc 验证宠物其他配置是否合法
		//  rarity, basic

		ok := p.Pets.AddIfNotExist(pet.ID, pet)
		if !ok { // 添加失败
			return fmt.Errorf("添加宠物失败,宠物已存在: %v %v", pet.ID, xruntime.Location())
		}
	}
	return nil
}

// Check 检查配置
func (p *PetMgr) Check() error {
	var err error
	p.Pets.Foreach(
		func(id common.AssetID, pet *Pet) bool { // 验证元素属性是否合法
			err = elemental.Validate(pet.Elemental.Earth, pet.Elemental.Water, pet.Elemental.Fire, pet.Elemental.Wind)
			if err != nil {
				err = fmt.Errorf("宠物ID %d 元素属性验证失败: %v", pet.ID, err)
				return false
			}
			return true
		},
	)
	return err
}

// Assemble 组装配置
func (p *PetMgr) Assemble() error {
	return nil
}

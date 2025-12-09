package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/res"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// MapExitRect 出口矩形
type MapExitRect struct {
	X int `yaml:"x"` // 出口矩形x
	Y int `yaml:"y"` // 出口矩形y
	W int `yaml:"w"` // 出口矩形宽度
	H int `yaml:"h"` // 出口矩形高度
}

// MapExit 地图出口
type MapExit struct {
	TargetTeleportID TeleportPointID `yaml:"targetTeleportID"` // 目标传送点ID
	Rect             MapExitRect     `yaml:"rect"`             // 出口矩形

	TargetTeleportPoint *TeleportPoint // 目标传送点 (Assemble时填充)
}

// Map 地图信息
type Map struct {
	AssetID common.AssetID `yaml:"assetID"` // 地图资产ID
	Name    string         `yaml:"name"`    // 地图名称
	Width   int            `yaml:"width"`   // 地图宽度
	Height  int            `yaml:"height"`  // 地图高度
	Exits   []*MapExit     `yaml:"exits"`   // 出口数组

	ResMap *res.Map // 资源-地图
}

var GMapMgr = newMapMgr()

// MapMgr 地图配置管理器
type MapMgr struct {
	Maps *xmap.MapMgr[common.AssetID, *Map] // key: 地图资产ID
}

func newMapMgr() *MapMgr {
	return &MapMgr{
		Maps: xmap.NewMapMgr[common.AssetID, *Map](),
	}
}

// Load 加载地图配置文件
func (p *MapMgr) Load() error {
	cfgPath := filepath.Join(common.AppCfgDir, "map.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取地图配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	var config struct {
		Maps []*Map `yaml:"maps"`
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析地图配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	for _, m := range config.Maps {
		if m.AssetID < common.AssetID(proto.AssetIDRange_AssetIDRange_Map_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Map_End) < m.AssetID {
			return fmt.Errorf("地图ID超出范围: %d %v", m.AssetID, xruntime.Location())
		}
		ok := p.Maps.AddIfNotExist(m.AssetID, m)
		if !ok {
			return fmt.Errorf("添加地图失败,地图已存在: %v %v", m.AssetID, xruntime.Location())
		}
	}
	return nil
}

// Check 检查地图配置
func (p *MapMgr) Check() error {
	var err error
	p.Maps.Foreach(func(id common.AssetID, m *Map) bool {
		// 检查宽高
		if m.Width <= 0 || m.Height <= 0 {
			err = fmt.Errorf("地图 %d 宽高无效: width=%d, height=%d", m.AssetID, m.Width, m.Height)
			return false
		}
		// 检查出口
		for _, exit := range m.Exits {
			// 检查目标传送点是否存在
			targetPoint, exists := GTeleportMgr.Points.Find(exit.TargetTeleportID)
			if !exists {
				err = fmt.Errorf("地图 %d 的出口目标传送点 %d 不存在", m.AssetID, exit.TargetTeleportID)
				return false
			}
			// 检查目标传送点所在地图是否存在
			if _, exists := p.Maps.Find(targetPoint.MapID); !exists {
				err = fmt.Errorf("地图 %d 的出口目标传送点 %d 所在地图 %d 不存在", m.AssetID, exit.TargetTeleportID, targetPoint.MapID)
				return false
			}
			// 检查出口矩形是否在地图范围内
			if exit.Rect.X < 0 || exit.Rect.Y < 0 ||
				exit.Rect.X+exit.Rect.W > m.Width ||
				exit.Rect.Y+exit.Rect.H > m.Height {
				err = fmt.Errorf("地图 %d 的出口矩形超出地图范围", m.AssetID)
				return false
			}
		}
		// 检查资源
		resMap := res.GMapMgr.Maps.Get(m.AssetID)
		if resMap == nil {
			err = fmt.Errorf("地图 %d 的资源不存在", m.AssetID)
			return false
		}
		return true
	})
	return err
}

// Assemble 组装地图配置
func (p *MapMgr) Assemble() error {
	var err error
	p.Maps.Foreach(func(id common.AssetID, m *Map) bool {
		for _, exit := range m.Exits {
			// 关联目标传送点
			exit.TargetTeleportPoint = GTeleportMgr.Points.Get(exit.TargetTeleportID)
		}
		return true
	})
	// 地图资源
	p.Maps.Foreach(
		func(id common.AssetID, m *Map) bool {
			// 资源
			m.ResMap = res.GMapMgr.Maps.Get(m.AssetID)
			return true
		},
	)

	return err
}

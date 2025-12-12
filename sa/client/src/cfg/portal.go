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

// PortalID 传送ID类型
type PortalID uint32

// PortalPoint 传送点
type PortalPoint struct {
	ID    PortalID       `yaml:"id"`    // 传送点ID (地图ID * 100 + 序号)
	Name  string         `yaml:"name"`  // 传送点名称
	MapID common.AssetID `yaml:"mapID"` // 所在地图ID
	X     int            `yaml:"x"`     // x坐标
	Y     int            `yaml:"y"`     // y坐标
}

var GPortalMgr = newPortalMgr()

// PortalMgr 传送点管理器
type PortalMgr struct {
	Points *xmap.MapMgr[PortalID, *PortalPoint] // key: 传送点ID
}

func newPortalMgr() *PortalMgr {
	return &PortalMgr{
		Points: xmap.NewMapMgr[PortalID, *PortalPoint](),
	}
}

// Load 加载传送点配置文件
func (p *PortalMgr) Load() error {
	cfgPath := filepath.Join(common.AppCfgDir, "portal.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取传送点配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	var config struct {
		TargetPoints []*PortalPoint `yaml:"targetPoints"`
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析传送点配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	for _, point := range config.TargetPoints {
		ok := p.Points.AddIfNotExist(point.ID, point)
		if !ok {
			return fmt.Errorf("添加传送点失败,传送点已存在: %v %v", point.ID, xruntime.Location())
		}
	}
	return nil
}

// Check 检查传送点配置
func (p *PortalMgr) Check() error {
	var err error
	p.Points.Foreach(func(id PortalID, point *PortalPoint) bool {
		// 检查所在地图是否存在
		targetMap, exists := GMapMgr.Maps.Find(point.MapID)
		if !exists {
			err = fmt.Errorf("传送点 %d 所在地图 %d 不存在", point.ID, point.MapID)
			return false
		}
		// 检查坐标是否在地图范围内 todo menglc 后续可扩展为检查是否在可行走区域内
		_ = targetMap
		// 检查传送点ID是否符合规则 (地图ID * 100 + 序号)
		expectedMapID := common.AssetID(point.ID / 100)
		if expectedMapID != point.MapID {
			err = fmt.Errorf("传送点 %d 的ID与所在地图 %d 不匹配 (期望地图ID: %d)", point.ID, point.MapID, expectedMapID)
			return false
		}
		return true
	})
	return err
}

// Assemble 组装传送点配置
func (p *PortalMgr) Assemble() error {
	return nil
}

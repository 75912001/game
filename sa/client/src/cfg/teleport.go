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

// TeleportPointID 传送点ID类型
type TeleportPointID uint32

// TeleportPoint 传送点
type TeleportPoint struct {
	ID    TeleportPointID `yaml:"id"`    // 传送点ID (地图ID * 100 + 序号)
	MapID common.AssetID  `yaml:"mapID"` // 所在地图ID
	X     int             `yaml:"x"`     // x坐标
	Y     int             `yaml:"y"`     // y坐标
	Name  string          `yaml:"name"`  // 传送点名称
}

var GTeleportMgr = newTeleportMgr()

// TeleportMgr 传送点配置管理器
type TeleportMgr struct {
	Points *xmap.MapMgr[TeleportPointID, *TeleportPoint] // key: 传送点ID
}

func newTeleportMgr() *TeleportMgr {
	return &TeleportMgr{
		Points: xmap.NewMapMgr[TeleportPointID, *TeleportPoint](),
	}
}

// Load 加载传送点配置文件
func (p *TeleportMgr) Load() error {
	cfgPath := filepath.Join(common.AppCfgDir, "teleport.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取传送点配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	var config struct {
		TeleportPoints []*TeleportPoint `yaml:"teleportPoints"`
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析传送点配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	for _, point := range config.TeleportPoints {
		ok := p.Points.AddIfNotExist(point.ID, point)
		if !ok {
			return fmt.Errorf("添加传送点失败,传送点已存在: %v %v", point.ID, xruntime.Location())
		}
	}
	return nil
}

// Check 检查传送点配置
func (p *TeleportMgr) Check() error {
	var err error
	p.Points.Foreach(func(id TeleportPointID, point *TeleportPoint) bool {
		// 检查所在地图是否存在
		targetMap, exists := GMapMgr.Maps.Find(point.MapID)
		if !exists {
			err = fmt.Errorf("传送点 %d 所在地图 %d 不存在", point.ID, point.MapID)
			return false
		}
		// 检查坐标是否在地图范围内
		if point.X < 0 || point.Y < 0 ||
			point.X >= targetMap.Width ||
			point.Y >= targetMap.Height {
			err = fmt.Errorf("传送点 %d 坐标(%d,%d)超出地图 %d 范围", point.ID, point.X, point.Y, point.MapID)
			return false
		}
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
func (p *TeleportMgr) Assemble() error {
	return nil
}

package cfg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"saClient/src/common"
	"saClient/src/proto"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Map 地图信息
type Map struct {
	AssetID      common.AssetID `yaml:"assetID"`      // 地图资产ID
	Name         string         `yaml:"name"`         // 地图名称
	TransitionID uint32         `yaml:"transitionID"` // 过渡效果ID

	ResMap          *TiledMap        // 资源-地图
	SceneTransition *SceneTransition // 过渡效果配置
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
		// 检查资源
		exist := GTiledMapMgr.Maps.IsExist(id)
		if !exist {
			err = fmt.Errorf("地图资源不存在: %d %v", id, xruntime.Location())
			return false
		}
		// 检查过渡效果配置
		switch m.TransitionID {
		case SceneTransitionIDVertical, SceneTransitionIDHorizontal, SceneTransitionIDFade:
			if !GSceneTransitionMgr.Transitions.IsExist(m.TransitionID) {
				err = fmt.Errorf("地图 %d 引用的过渡效果不存在: %d %v", id, m.TransitionID, xruntime.Location())
				return false
			}
		default:
			log.Printf("地图 %d 未指定过渡效果ID%v", id, xruntime.Location())
		}
		return true
	})
	return err
}

// Assemble 组装地图配置
func (p *MapMgr) Assemble() error {
	p.Maps.Foreach(func(id common.AssetID, m *Map) bool {
		// 关联-地图资源
		m.ResMap = GTiledMapMgr.Maps.Get(m.AssetID)
		// 关联-过渡效果配置
		switch m.TransitionID {
		case SceneTransitionIDVertical, SceneTransitionIDHorizontal, SceneTransitionIDFade:
			m.SceneTransition = GSceneTransitionMgr.Transitions.Get(m.TransitionID)
		default:
			m.SceneTransition = DefSceneTransition()
			log.Printf("地图 %d 引用的场景过渡效果不存在: %d, 使用默认过渡效果 %d %v", id, m.TransitionID, m.SceneTransition.ID, xruntime.Location())
		}
		return true
	})
	return nil
}

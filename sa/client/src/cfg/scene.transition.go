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

// 场景-过渡效果ID
const (
	SceneTransitionIDVertical   uint32 = 1 // 垂直 (上下合拢/展开)
	SceneTransitionIDHorizontal uint32 = 2 // 水平 (左右合拢/展开)
	SceneTransitionIDFade       uint32 = 3 // 淡入淡出
)

// SceneTransition 场景-过渡效果配置
type SceneTransition struct {
	ID    uint32  `yaml:"id"`    // 过渡效果ID
	Speed float32 `yaml:"speed"` // 过渡速度 (每帧进度增量)
}

func DefaultSceneTransition() *SceneTransition {
	return &SceneTransition{
		ID:    SceneTransitionIDVertical,
		Speed: 0.1,
	}
}

var GSceneTransitionMgr = newSceneTransitionMgr()

// SceneTransitionMgr 过渡效果配置管理器
type SceneTransitionMgr struct {
	Transitions *xmap.MapMgr[uint32, *SceneTransition] // key: 过渡效果ID
}

func newSceneTransitionMgr() *SceneTransitionMgr {
	return &SceneTransitionMgr{
		Transitions: xmap.NewMapMgr[uint32, *SceneTransition](),
	}
}

// Load 加载过渡效果配置文件
func (p *SceneTransitionMgr) Load() error {
	cfgPath := filepath.Join(common.AppCfgDir, "scene.transition.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取过渡效果配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	var config struct {
		Transitions []*SceneTransition `yaml:"transitions"`
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析过渡效果配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	for _, t := range config.Transitions {
		if t.ID == 0 {
			return fmt.Errorf("过渡效果ID不能为0: %v", xruntime.Location())
		}
		switch t.ID {
		case SceneTransitionIDVertical:
		case SceneTransitionIDHorizontal:
		case SceneTransitionIDFade:
		default:
			return fmt.Errorf("无效的过渡效果ID: %d %v", t.ID, xruntime.Location())
		}
		// 补充默认值
		if t.Speed <= 0 {
			t.Speed = 0.1
		}
		ok := p.Transitions.AddIfNotExist(t.ID, t)
		if !ok {
			return fmt.Errorf("添加过渡效果失败,ID已存在: %v %v", t.ID, xruntime.Location())
		}
	}
	return nil
}

// Get 获取过渡效果配置
func (p *SceneTransitionMgr) Get(id uint32) *SceneTransition {
	return p.Transitions.Get(id)
}

// Check 检查地图配置
func (p *SceneTransitionMgr) Check() error {
	return nil
}

// Assemble 组装地图配置
func (p *SceneTransitionMgr) Assemble() error {
	return nil
}

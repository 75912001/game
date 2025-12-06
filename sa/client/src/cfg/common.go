package cfg

import (
	"encoding/json"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"saClient/src/common"
)

// Common 通用配置
type Common struct {
	RoleCountMax   int `json:"roleCountMax"`        // 允许创建的最大角色数
	RoleRebirthMax int `json:"roleRebirthCountMax"` // 允许转生的最大次数
}

var GCommon = &Common{}

// Load 加载通用配置文件
func (p *Common) Load() error {
	// 构建配置文件路径
	cfgPath := filepath.Join(common.AppCfgDir, "common.json")
	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取通用配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	// 解析JSON
	err = json.Unmarshal(data, p)
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

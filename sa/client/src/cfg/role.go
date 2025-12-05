package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"saClient/src/common"
	"saClient/src/proto"
)

// Role 角色信息
type Role struct {
	ID          common.AssetID `json:"id"`          // 角色ID
	Name        string         `json:"name"`        // 英文名称
	NameCN      string         `json:"nameCN"`      // 中文名称
	Gender      string         `json:"gender"`      // 性别
	Color       string         `json:"color"`       // 颜色
	Description string         `json:"description"` // 描述
}

var GRoleMgr = newRoleMgr()

// RoleMgr 角色配置管理器
type RoleMgr struct {
	roles map[common.AssetID]*Role // key: 角色ID
}

func newRoleMgr() *RoleMgr {
	return &RoleMgr{
		roles: make(map[common.AssetID]*Role),
	}
}

// Load 加载角色配置文件
func (m *RoleMgr) Load() error {
	// 构建配置文件路径
	cfgPath := filepath.Join(common.AppCfgDir, "role.json")

	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("读取角色配置文件失败: %v", err)
	}
	// 定义中间结构用于JSON解析
	var config struct {
		Roles []*Role `json:"roles"`
	}
	// 解析JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("解析角色配置文件失败: %v", err)
	}
	// 加载每个角色
	for _, role := range config.Roles {
		if role.ID < common.AssetID(proto.AssetIDRange_AssetIDRange_Role_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Role_End) < role.ID {
			return fmt.Errorf("角色ID超出范围: %d", role.ID)
		}
		// 检查是否重复
		if _, exists := m.roles[role.ID]; exists {
			return fmt.Errorf("角色ID重复: %d", role.ID)
		}
		m.roles[role.ID] = role
	}
	return nil
}

func (m *RoleMgr) Check() error {
	return nil
}

// GetRole 获取角色信息
func (m *RoleMgr) GetRole(id common.AssetID) *Role {
	return m.roles[id]
}

// GetRoleCount 获取角色数量
func (m *RoleMgr) GetRoleCount() uint32 {
	return uint32(len(m.roles))
}

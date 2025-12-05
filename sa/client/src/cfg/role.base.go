package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"saClient/src/common"
	"saClient/src/proto"
)

// RoleBase 角色基础属性
type RoleBase struct {
	ID          common.AssetID `json:"id"`          // 属性ID
	Name        string         `json:"name"`        // 英文名称
	NameCN      string         `json:"nameCN"`      // 中文名称
	Category    string         `json:"category"`    // 类别 primary/basic
	Description string         `json:"description"` // 描述
}

var GRoleBaseMgr = newRoleBaseMgr()

// RoleBaseMgr 角色基础属性配置管理器
type RoleBaseMgr struct {
	roleBases map[common.AssetID]*RoleBase // key: 角色基础属性ID
}

func newRoleBaseMgr() *RoleBaseMgr {
	return &RoleBaseMgr{
		roleBases: make(map[common.AssetID]*RoleBase),
	}
}

// Load 加载角色基础属性配置文件
func (m *RoleBaseMgr) Load() error {
	// 构建配置文件路径
	cfgPath := filepath.Join(common.AppCfgDir, "role.base.json")
	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("读取角色基础属性配置文件失败: %v", err)
	}
	// 定义中间结构用于JSON解析
	var config struct {
		RoleBase []*RoleBase `json:"roleBase"`
	}
	// 解析JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("解析角色基础属性配置文件失败: %v", err)
	}

	// 加载每个基础属性
	for _, roleBase := range config.RoleBase {
		if roleBase.ID < common.AssetID(proto.AssetIDRange_AssetIDRange_RoleBase_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_RoleBase_End) < roleBase.ID {
			return fmt.Errorf("角色基础属性ID超出范围: %d", roleBase.ID)
		}
		// 检查是否重复
		if _, exists := m.roleBases[roleBase.ID]; exists {
			return fmt.Errorf("角色基础属性ID重复: %d", roleBase.ID)
		}
		m.roleBases[roleBase.ID] = roleBase
	}
	return nil
}

func (m *RoleBaseMgr) Check() error {
	return nil
}

// GetRoleBase 获取角色基础属性信息
func (m *RoleBaseMgr) GetRoleBase(id common.AssetID) *RoleBase {
	return m.roleBases[id]
}

// GetRoleBaseCount 获取角色基础属性数量
func (m *RoleBaseMgr) GetRoleBaseCount() uint32 {
	return uint32(len(m.roleBases))
}

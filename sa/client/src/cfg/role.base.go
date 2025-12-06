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
	RoleBases *xmap.MapMgr[common.AssetID, *RoleBase] // key: 角色基础属性ID
}

func newRoleBaseMgr() *RoleBaseMgr {
	return &RoleBaseMgr{
		RoleBases: xmap.NewMapMgr[common.AssetID, *RoleBase](),
	}
}

// Load 加载基础属性配置文件
func (p *RoleBaseMgr) Load() error {
	// 构建配置文件路径
	cfgPath := filepath.Join(common.AppCfgDir, "role.base.json")
	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取角色基础属性配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	// 定义中间结构用于JSON解析
	var config struct {
		RoleBase []*RoleBase `json:"roleBase"`
	}
	// 解析JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析角色基础属性配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	// 加载每个基础属性
	for _, roleBase := range config.RoleBase {
		if roleBase.ID < common.AssetID(proto.AssetIDRange_AssetIDRange_RoleBase_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_RoleBase_End) < roleBase.ID {
			return fmt.Errorf("角色基础属性ID超出范围: %d %v", roleBase.ID, xruntime.Location())
		}
		ok := p.RoleBases.AddIfNotExist(roleBase.ID, roleBase)
		if !ok { // 添加失败
			return fmt.Errorf("添加角色基础属性失败,属性已存在: %d %v", roleBase.ID, xruntime.Location())
		}
	}
	return nil
}

// Check 检查数据
func (p *RoleBaseMgr) Check() error {
	return nil
}

// Assemble 组装数据
func (p *RoleBaseMgr) Assemble() error {
	return nil
}

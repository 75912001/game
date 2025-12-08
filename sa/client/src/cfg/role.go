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

// Role 角色信息
type Role struct {
	ID          common.AssetID `yaml:"id"`          // 角色ID
	Name        string         `yaml:"name"`        // 名称
	Gender      string         `yaml:"gender"`      // 性别
	Color       string         `yaml:"color"`       // 颜色
	Description string         `yaml:"description"` // 描述

	roleBase *RoleBase // 基础属性
	resRole  *res.Role // 角色资源
}

var GRoleMgr = newRoleMgr()

// RoleMgr 配置管理器
type RoleMgr struct {
	Roles *xmap.MapMgr[common.AssetID, *Role] // key: 角色ID
}

func newRoleMgr() *RoleMgr {
	return &RoleMgr{
		Roles: xmap.NewMapMgr[common.AssetID, *Role](),
	}
}

// Load 加载配置文件
func (p *RoleMgr) Load() error {
	// 构建配置文件路径
	cfgPath := filepath.Join(common.AppCfgDir, "role.yaml")
	// 读取配置文件
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取角色配置文件失败: %v %v", cfgPath, xruntime.Location())
	}
	// 定义中间结构用于YAML解析
	var config struct {
		Roles []*Role `yaml:"roles"`
	}
	// 解析YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析角色配置文件失败: %v %v", cfgPath, err)
	}
	// 加载每个角色
	for _, role := range config.Roles {
		if 1000101 != role.ID { // todo menglc 仅加载测试角色
			continue
		}
		if role.ID < common.AssetID(proto.AssetIDRange_AssetIDRange_Role_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Role_End) < role.ID {
			return fmt.Errorf("角色ID超出范围: %d %v", role.ID, xruntime.Location())
		}
		ok := p.Roles.AddIfNotExist(role.ID, role)
		if !ok { // 添加失败
			return fmt.Errorf("添加角色失败,角色已存在: %v %v", role.ID, xruntime.Location())
		}
	}
	return nil
}

// Check 检查配置
func (p *RoleMgr) Check() error {
	var err error
	p.Roles.Foreach( // 检查每个角色的基础属性和资源是否存在
		func(id common.AssetID, role *Role) bool {
			// 基础属性
			roleBase := GRoleBaseMgr.RoleBases.Get(role.ID)
			if roleBase == nil {
				err = fmt.Errorf("角色ID %d 的基础属性不存在", role.ID)
				return false
			}
			// 资源
			resRole := res.GRoleMgr.Roles.Get(role.ID)
			if resRole == nil {
				err = fmt.Errorf("角色ID %d 的资源不存在", role.ID)
				return false
			}
			return true
		},
	)
	return err
}

// Assemble 组装配置
func (p *RoleMgr) Assemble() error {
	p.Roles.Foreach(
		func(id common.AssetID, role *Role) bool {
			// 基础属性
			role.roleBase = GRoleBaseMgr.RoleBases.Get(role.ID)
			// 资源
			role.resRole = res.GRoleMgr.Roles.Get(role.ID)
			return true
		},
	)
	return nil
}

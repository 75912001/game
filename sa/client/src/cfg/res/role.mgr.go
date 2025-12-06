package res

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"saClient/src/common"
	"saClient/src/proto"
	"strconv"
)

var GRoleMgr = newRoleMgr()

// RoleMgr 角色配置管理器
type RoleMgr struct {
	Roles map[common.AssetID]*Role // key: 角色ID
}

func newRoleMgr() *RoleMgr {
	return &RoleMgr{
		Roles: make(map[common.AssetID]*Role),
	}
}

// Load 加载角色res配置文件
func (p *RoleMgr) Load() error {
	// 读取目录下所有文件
	roleDirectories, err := os.ReadDir(common.AppResRoleDir)
	if err != nil {
		return fmt.Errorf("读取角色配置目录失败: %v", err)
	}
	for _, roleDir := range roleDirectories {
		if !roleDir.IsDir() { // 非目录
			continue
		}
		roleDirName := roleDir.Name() // 目录名称-角色ID
		var roleID common.AssetID     // 角色ID
		{
			id64, err := strconv.ParseUint(roleDirName, 10, 32)
			if err != nil {
				return fmt.Errorf("解析角色目录名称 %s 为 ID 失败: %v", roleDirName, err)
			}
			roleID = common.AssetID(id64)
		}
		// 读取 dirName 目录下的所有文件
		roleDirPath := filepath.Join(common.AppResRoleDir, roleDirName)
		roleFiles, err := os.ReadDir(roleDirPath)
		if err != nil {
			return fmt.Errorf("读取角色目录 %s 失败: %v", roleDirName, err)
		}
		// 匹配 role.${roleID}.${action}.${arg}.json 格式的文件 例如 `role.1000101.move.up.json`
		pattern := regexp.MustCompile(fmt.Sprintf(`^%v\.%v\.([^.]+)\.([^.]+)\.json$`, common.GetNameByAssetType(proto.AssetType_AssetType_Role), roleID))
		for _, roleFile := range roleFiles {
			if roleFile.IsDir() { // 目录
				continue
			}
			jsonFileName := roleFile.Name()
			if !pattern.MatchString(jsonFileName) { // 不匹配格式
				continue
			}
			imageFilePath := filepath.Join(roleDirPath, jsonFileName[:len(jsonFileName)-len(filepath.Ext(jsonFileName))]+".png")
			// 加载配置文件
			jsonFilePath := filepath.Join(roleDirPath, jsonFileName)
			roleJson, err := loadRoleJson(jsonFilePath)
			if err != nil {
				return fmt.Errorf("加载配置文件 %s 失败: %v", jsonFileName, err)
			}
			err = loadRoleImage(roleID, imageFilePath, roleJson)
			if err != nil {
				return fmt.Errorf("加载角色 %s 失败: %v", jsonFileName, err)
			}
		}
	}
	return nil
}

func (p *RoleMgr) Check() error {
	return nil
}

func (p *RoleMgr) Assemble() error {
	return nil
}

// Get 获取角色信息
func (p *RoleMgr) Get(roleID common.AssetID) *Role {
	return p.Roles[roleID]
}

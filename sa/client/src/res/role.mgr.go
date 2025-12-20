package res

import (
	"fmt"
	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
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
	Roles *xmap.MapMgr[common.AssetID, *Role]
}

func newRoleMgr() *RoleMgr {
	return &RoleMgr{
		Roles: xmap.NewMapMgr[common.AssetID, *Role](),
	}
}

// Load 加载角色res配置文件
func (p *RoleMgr) Load() error {
	// 读取目录下所有文件-role
	roleDirectories, err := os.ReadDir(common.AppResRoleDir)
	if err != nil {
		return errors.WithMessagef(err, "读取角色配置目录失败: %v %v", common.AppResRoleDir, xruntime.Location())
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
				return errors.WithMessagef(err, "解析角色目录名称 %s 为 ID 失败: %v", roleDirName, xruntime.Location())
			}
			roleID = common.AssetID(id64)
		}
		// 读取 dirName 目录下的所有文件
		roleDirPath := filepath.Join(common.AppResRoleDir, roleDirName)
		roleFiles, err := os.ReadDir(roleDirPath)
		if err != nil {
			return errors.WithMessagef(err, "读取角色目录 %s 失败: %v", roleDirPath, xruntime.Location())
		}
		// 匹配 role.${roleID}.${arg}.${data}.json 格式的文件 例如 `role.1000101.move.up.json`
		rolePattern := regexp.MustCompile(fmt.Sprintf(`^%v\.%v\.([^.]+)\.([^.]+)\.json$`, common.GetNameByAssetType(proto.AssetType_AssetType_Role), roleID))
		for _, roleFile := range roleFiles {
			if roleFile.IsDir() { // 目录
				continue
			}
			jsonFileName := roleFile.Name()
			if !rolePattern.MatchString(jsonFileName) { // 不匹配格式
				continue
			}
			imageFilePath := filepath.Join(roleDirPath, jsonFileName[:len(jsonFileName)-len(filepath.Ext(jsonFileName))]+".png")
			// 加载配置文件
			jsonFilePath := filepath.Join(roleDirPath, jsonFileName)
			roleJson, err := loadRoleJson(jsonFilePath)
			if err != nil {
				return errors.WithMessagef(err, "加载配置文件失败 %s %v", jsonFileName, xruntime.Location())
			}
			err = loadRoleImage(roleID, imageFilePath, roleJson)
			if err != nil {
				return errors.WithMessagef(err, "加载角色失败 %v %v", roleID, xruntime.Location())
			}
		}
	}
	return nil
}

func (p *RoleMgr) Check() error {
	// 验证每个角色的8个方向是否完整
	var checkErr error
	p.Roles.Foreach(func(key common.AssetID, role *Role) bool {
		if err := checkRole8Directions(role); err != nil {
			checkErr = err
			return false // 停止遍历
		}
		return true // 继续遍历
	})
	return checkErr
}

// checkRole8Directions 验证角色8个方向是否完整
func checkRole8Directions(role *Role) error {
	// 需要验证的8个方向
	directions := []proto.AssetOrientation{
		proto.AssetOrientation_AssetOrientation_Up,
		proto.AssetOrientation_AssetOrientation_UpRight,
		proto.AssetOrientation_AssetOrientation_Right,
		proto.AssetOrientation_AssetOrientation_DownRight,
		proto.AssetOrientation_AssetOrientation_Down,
		proto.AssetOrientation_AssetOrientation_DownLeft,
		proto.AssetOrientation_AssetOrientation_Left,
		proto.AssetOrientation_AssetOrientation_UpLeft,
	}

	// 验证 Move 动画的8个方向
	for _, dir := range directions {
		if len(role.Move.Frames[dir]) == 0 {
			return fmt.Errorf("角色 %v 移动动画缺少方向 %v %v", role.ID, GetNameByAssetOrientation(dir), xruntime.Location())
		}
	}

	return nil
}

func (p *RoleMgr) Assemble() error {
	return nil
}

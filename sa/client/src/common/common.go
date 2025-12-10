package common

import (
	"log"
	"os"
	"path/filepath"
)

const (
	AppVersion     = "0.0.1"          // 版本号
	AppVersionDesc = "SA Game v0.0.1" // 版本描述
	AppName        = "SA"             // 应用名称
	AppWindowTitle = "SA Game"        // 窗口标题
)

var (
	AppExeDir       string // 可执行文件目录
	AppDataDir      string // 数据目录
	AppLogDir       string // 日志目录
	AppCfgDir       string // 配置目录
	AppTempDir      string // 临时目录
	AppBinDir       string // 可执行文件目录
	AppResDir       string // 资源目录
	AppResSystemDir string // 资源目录-系统
	AppResFontDir   string // 资源目录-字体
	AppResMapDir    string // 资源目录-地图
	AppResRoleDir   string // 资源目录-角色
)

type AssetID uint32       // 资产ID
type AssetQuantity uint64 // 资产数量
type UserUUID uint64      // 用户ID
type RoleUUID uint64      // 角色ID 每个用户的角色唯一ID
type PetUUID uint64       // 宠物唯一ID 每个角色的宠物唯一ID

type AssetMap map[AssetID]AssetQuantity // 资产映射

func init() {
	var err error
	// 获取可执行文件所在目录
	AppExeDir, err = os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	AppBinDir = filepath.Dir(AppExeDir)
	clientDir := filepath.Dir(AppBinDir)
	// 配置
	AppCfgDir = filepath.Join(clientDir, "cfg")
	// 资源
	AppResDir = filepath.Join(clientDir, "res")
	AppResSystemDir = filepath.Join(AppResDir, "system")
	AppResFontDir = filepath.Join(AppResDir, "font")
	AppResMapDir = filepath.Join(AppResDir, "map")
	AppResRoleDir = filepath.Join(AppResDir, "role")
}

package common

import (
	"log"
	"os"
	"path/filepath"
)

const (
	AppVersion     = "0.0.1"
	AppName        = "saClient"
	AppWindowTitle = "SA Game"
)

var (
	AppExeDir     string // 可执行文件目录
	AppDataDir    string // 数据目录
	AppLogDir     string // 日志目录
	AppCfgDir     string // 配置目录
	AppTempDir    string // 临时目录
	AppBinDir     string // 可执行文件目录
	AppResDir     string // 资源目录
	AppResFontDir string // 资源目录-字体
	AppResMapDir  string // 资源目录-地图
	AppResRoleDir string // 资源目录-角色

)

const (
	WindowWidth  = 1024
	WindowHeight = 800

	ScreenWidth  = 1024
	ScreenHeight = 800
)

func init() {
	var err error
	// 获取可执行文件所在目录
	AppExeDir, err = os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	AppBinDir = filepath.Dir(AppExeDir)
	clientDir := filepath.Dir(AppBinDir)
	AppCfgDir = filepath.Join(clientDir, "cfg")
	AppResDir = filepath.Join(clientDir, "res")
	AppResFontDir = filepath.Join(AppResDir, "font")
	AppResMapDir = filepath.Join(AppResDir, "map")
	AppResRoleDir = filepath.Join(AppResDir, "role")
}

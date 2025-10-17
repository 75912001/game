package common

import (
	"log"
	"os"
	"path/filepath"
)

const (
	AppVersion     = "0.0.1"
	AppName        = "airplane"
	AppWindowTitle = "Airplane Game"
)

var (
	AppExeDir      string // 可执行文件目录
	AppDataDir     string // 数据目录
	AppLogDir      string // 日志目录
	AppConfDir     string // 配置目录
	AppTempDir     string // 临时目录
	AppResourceDir string // 资源目录
	AppBinDir      string // 可执行文件目录
)

const (
	WindowWidth  = 800
	WindowHeight = 600

	ScreenWidth  = 800
	ScreenHeight = 600
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
	AppResourceDir = filepath.Join(clientDir, "resource")
}

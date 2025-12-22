package main

import (
	xpacket "github.com/75912001/xlib/packet"
	xtime "github.com/75912001/xlib/time"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"os"
	"path/filepath"
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/game"
)

func main() {
	xtime.GTimeMgr.Update() // 初始化时间管理器
	xpacket.SetEndianMode(xpacket.LittleEndian)
	// todo menglc 初始化程序, 需要 加载 app 配置文件
	err := LoadCfg()
	if err != nil {
		panic(err)
	}
	err = CheckCfg()
	if err != nil {
		panic(err)
	}
	err = AssembleCfg()
	if err != nil {
		panic(err)
	}
	// 设置窗口大小和标题
	ebitenv2.SetWindowSize(cfg.GCommon.WindowDefaultWidth, cfg.GCommon.WindowDefaultHeight)
	ebitenv2.SetWindowTitle(common.AppWindowTitle)

	// 使用标准库加载图标
	windowIconPath := filepath.Join(common.AppResSystemDir, "window.icon.png")
	iconFile, err := os.Open(windowIconPath)
	if err != nil {
		panic(err)
	}
	icon, _, err := image.Decode(iconFile)
	_ = iconFile.Close()
	if err != nil {
		panic(err)
	}
	ebitenv2.SetWindowIcon([]image.Image{icon})

	// 创建并初始化游戏
	gameObject := game.NewGame()
	gameObject.Init()

	// 运行游戏
	err = ebitenv2.RunGame(gameObject)
	if err != nil {
		log.Fatal(err)
	}
}

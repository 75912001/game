package main

import (
	apccommon "airplaneClient/src/common"
	apcgame "airplaneClient/src/game"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// 设置窗口大小和标题
	ebiten.SetWindowSize(apccommon.WindowWidth, apccommon.WindowHeight)
	ebiten.SetWindowTitle(apccommon.AppWindowTitle)

	// 使用标准库加载图标
	windowIconPath := filepath.Join(apccommon.AppResourcesDir, "window.icon.png")
	iconFile, err := os.Open(windowIconPath)
	if err != nil {
		log.Fatal(err)
	}
	icon, _, err := image.Decode(iconFile)
	_ = iconFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowIcon([]image.Image{icon})

	// 创建并初始化游戏
	game := &apcgame.Game{}
	game.Init()

	// 运行游戏
	err = ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}

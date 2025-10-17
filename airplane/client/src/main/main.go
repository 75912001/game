package main

import (
	apccommon "airplaneClient/src/common"
	apcgame "airplaneClient/src/game"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"image"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fmt.Printf("Hello World!\n")

	// 设置窗口大小和标题
	ebiten.SetWindowSize(apccommon.WindowWidth, apccommon.WindowHeight)
	ebiten.SetWindowTitle(apccommon.AppWindowTitle)

	// 使用标准库加载图标
	windowIconPath := filepath.Join(apccommon.AppResourceDir, "window.icon.png")
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

	// 运行游戏
	err = ebiten.RunGame(&apcgame.Game{})
	if err != nil {
		log.Fatal(err)
	}
}

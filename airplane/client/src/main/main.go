package main

import (
	apccommon "airplaneClient/src/common"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
)

type Game struct {
}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

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
	err = ebiten.RunGame(&Game{})
	if err != nil {
		log.Fatal(err)
	}
}

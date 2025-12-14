package user

import (
	"image/color"
	"saClient/src/ui"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (p *Role) Draw(screen *ebitenv2.Image) {
	// 绘制场景 (包含过渡效果)
	p.scene.Draw(screen, p.camera)

	// 绘制角色
	// 角色屏幕位置 = 角色 World 坐标 - 摄像机视口 World 坐标 - 角色图片偏移
	screenX := p.sprite.centerWorldX - float32(p.camera.ViewportX) - float32(p.sprite.roleImageSprite.Frame.Width/2)
	screenY := p.sprite.centerWorldY - float32(p.camera.ViewportY) - float32(p.sprite.roleImageSprite.Frame.Height/2)

	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(p.sprite.image, op)

	// 绘制调试信息
	p.drawDebugInfo(screen)
}

// drawDebugInfo 绘制调试信息
func (p *Role) drawDebugInfo(screen *ebitenv2.Image) {
	mapCfg := p.scene._map.cfg
	// 获取角色 Tile 坐标
	tileX, tileY := mapCfg.IsometricCT.W2T(p.sprite.bottomCenterWorldX, p.sprite.bottomCenterWorldY)

	// 获取角色 World 坐标
	worldX := p.sprite.bottomCenterWorldX
	worldY := p.sprite.bottomCenterWorldY

	// 获取角色 Screen 坐标
	roleScreenX := p.sprite.centerWorldX - float32(p.camera.ViewportX)
	roleScreenY := p.sprite.centerWorldY - float32(p.camera.ViewportY)

	if true { // 绘制脚底中心点 (红色圆形)
		bottomCenterScreenX := p.sprite.bottomCenterWorldX - float32(p.camera.ViewportX)
		bottomCenterScreenY := p.sprite.bottomCenterWorldY - float32(p.camera.ViewportY)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		vector.FillCircle(screen, bottomCenterScreenX, bottomCenterScreenY, 5, red, false)
	}

	// 显示地图信息
	y := float32(10.0)
	ui.Printf(screen, 10, y, "=== Map Info ===")
	y += 20
	ui.Printf(screen, 10, y, "Map ID: %d", mapCfg.ID)
	y += 20
	ui.Printf(screen, 10, y, "Tile Size: %d x %d", mapCfg.Width, mapCfg.Height)
	y += 20
	ui.Printf(screen, 10, y, "Pixel Size: %d x %d", mapCfg.PixelW, mapCfg.PixelH)

	// 显示角色坐标信息
	y += 30
	ui.Printf(screen, 10, y, "=== Role Position ===")
	y += 20
	ui.Printf(screen, 10, y, "Tile:   (%.2f, %.2f)", tileX, tileY)
	y += 20
	ui.Printf(screen, 10, y, "World:  (%.1f, %.1f)", worldX, worldY)
	y += 20
	ui.Printf(screen, 10, y, "Screen: (%.1f, %.1f)", roleScreenX, roleScreenY)

	// 显示摄像机信息
	y += 30
	ui.Printf(screen, 10, y, "=== Camera ===")
	y += 20
	ui.Printf(screen, 10, y, "Viewport: (%d, %d)", p.camera.ViewportX, p.camera.ViewportY)
}

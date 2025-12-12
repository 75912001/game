package user

import (
	"saClient/src/ui"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

func (p *Role) Draw(screen *ebitenv2.Image) {
	// 绘制场景
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
	// 获取地图信息
	mapID := p.scene.GetMapID()
	tileW, tileH := p.scene.GetMapTileSize()
	pixelW, pixelH := p.scene.GetMapPixeSize()

	// 获取角色 Tile 坐标
	tileX, tileY := p.scene.WorldToTile(p.sprite.bottomCenterWorldX, p.sprite.bottomCenterWorldY)

	// 获取角色 World 坐标
	worldX := p.sprite.bottomCenterWorldX
	worldY := p.sprite.bottomCenterWorldY

	// 获取角色 Screen 坐标
	roleScreenX := p.sprite.centerWorldX - float32(p.camera.ViewportX)
	roleScreenY := p.sprite.centerWorldY - float32(p.camera.ViewportY)

	// 显示地图信息
	y := float32(10.0)
	ui.Printf(screen, 10, y, "=== Map Info ===")
	y += 20
	ui.Printf(screen, 10, y, "Map ID: %d", mapID)
	y += 20
	ui.Printf(screen, 10, y, "Tile Size: %d x %d", tileW, tileH)
	y += 20
	ui.Printf(screen, 10, y, "Pixel Size: %d x %d", pixelW, pixelH)

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

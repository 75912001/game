package user

import (
	"image/color"
	"saClient/src/cfg"
	"saClient/src/common"
	commoncamera "saClient/src/common/camera"
	"saClient/src/ui"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawDebugInfo 绘制调试信息
func (p *Role) drawDebugInfo(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	mapCfg := p.scene._map.tiledMapCfg
	// 获取角色 Tile 坐标
	tileX, tileY := mapCfg.IsometricCT.W2T(p.sprite.actionAnchorPointWX, p.sprite.actionAnchorPointWY)

	// 获取角色 World 坐标
	worldX := p.sprite.actionAnchorPointWX
	worldY := p.sprite.actionAnchorPointWY

	// 获取角色 Screen 坐标
	roleScreenX := p.sprite.cameraAnchorPointWX - float32(p.camera.ViewportWX)
	roleScreenY := p.sprite.cameraAnchorPointWY - float32(p.camera.ViewportWY)

	if true { // 绘制脚底中心点 (红色圆形)
		bottomCenterScreenX := p.sprite.actionAnchorPointWX - float32(p.camera.ViewportWX)
		bottomCenterScreenY := p.sprite.actionAnchorPointWY - float32(p.camera.ViewportWY)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		vector.FillCircle(screen, bottomCenterScreenX, bottomCenterScreenY, 5, red, false)

		// 绘制视野范围 (红色虚线圆)
		common.DrawDashedCircle(screen, bottomCenterScreenX, bottomCenterScreenY, cfg.GCommon.RoleArpgDefViewRange, common.Colors_Red, 48, 0.5, 2.0)
	}

	if true { // 绘制调试边界
		p.scene._map.drawBorder(screen, cam)
		p.scene._map.drawPortal(screen, cam)
		p.scene._map.drawBlocked(screen, cam)
		p.scene._map.drawSpawnPointDebug(screen, cam) // 绘制刷怪点调试信息
	}

	// 绘制敌人调试信息 (视野范围)
	for _, enemy := range p.scene._map.spawnManager.GetAllEnemies() {
		enemy.DrawDebug(screen, p.camera)
	}

	// 显示地图信息
	y := float32(10.0)
	ui.Printf(screen, 10, y, "=== Map Info ===")
	y += 20
	ui.Printf(screen, 10, y, "Map ID: %d", mapCfg.ID)
	y += 20
	ui.Printf(screen, 10, y, "Tile Size: %d x %d", mapCfg.TileCountW, mapCfg.TileCountH)
	y += 20
	ui.Printf(screen, 10, y, "Pixel Size: %d x %d", mapCfg.WPixel, mapCfg.HPixel)

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
	ui.Printf(screen, 10, y, "Viewport: (%d, %d)", p.camera.ViewportWX, p.camera.ViewportWY)
}

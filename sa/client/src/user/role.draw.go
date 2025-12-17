package user

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	commoncamera "saClient/src/common/camera"
	commonrenderable "saClient/src/common/renderable"
	"saClient/src/ui"
)

func (p *Role) DrawAll(screen *ebitenv2.Image) {
	// 填充草地一样的绿色背景
	screen.Fill(color.RGBA{
		R: 34,
		G: 139,
		B: 34,
		A: 255,
	})

	p.scene._map.DrawCollision(screen, p.camera)
	p.scene._map.DrawGround(screen, p.camera)

	p.scene._map.DrawBuilding(screen, p.camera)
	renderables := p.collectRenderables()
	commonrenderable.SortYAndDraw(screen, p.camera, renderables)

	{
		// todo menglc 绘制剩余地图元素
		p.scene._map.DrawObjects(screen, p.camera)
	}

	p.scene._map.DrawOverhead(screen, p.camera)

	// 绘制过渡效果 (最上层)
	p.scene.transition.Draw(screen)

	// 绘制调试信息
	p.drawDebugInfo(screen)
}

// GetWY 返回角色动作锚点的 World Y 坐标（用于Y-Sorting）
func (p *Role) GetWY() float32 {
	return p.sprite.actionAnchorPointWY
}

// Draw 仅绘制角色本身(实现 Renderable 接口)
func (p *Role) Draw(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	screenX := p.sprite.cameraAnchorPointWX - float32(cam.ViewportWX) - float32(p.sprite.roleImageSprite.Frame.Width/2)
	screenY := p.sprite.cameraAnchorPointWY - float32(cam.ViewportWY) - float32(p.sprite.roleImageSprite.Frame.Height/2)

	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(p.sprite.image, op)
}

// collectRenderables 收集所有需要Y-Sorting的对象
func (p *Role) collectRenderables() []commonrenderable.IRenderable {
	var list []commonrenderable.IRenderable
	list = append(list, p) // 角色自己

	//list = append(list, p.scene.buildingMgr.GetRenderables()...)
	//list = append(list, p.scene.plantMgr.GetRenderables()...)
	//list = append(list, p.scene.decorationMgr.GetRenderables()...)
	//list = append(list, p.scene.itemMgr.GetRenderables()...)

	return list
}

// drawDebugInfo 绘制调试信息
func (p *Role) drawDebugInfo(screen *ebitenv2.Image) {
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
	ui.Printf(screen, 10, y, "Viewport: (%d, %d)", p.camera.ViewportWX, p.camera.ViewportWY)
}

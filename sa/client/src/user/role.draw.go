package user

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/cfg"
	commoncamera "saClient/src/common/camera"
	commonrenderable "saClient/src/common/renderable"
)

func (p *Role) DrawAll(screen *ebitenv2.Image) {
	p.scene._map.DrawGround(screen, p.camera)

	renderables := p.collectRenderables()
	commonrenderable.SortYAndDraw(screen, p.camera, renderables)
	p.scene._map.RecycleTileSortInfoSlice(renderables) // 回收对象池

	p.scene._map.DrawOverhead(screen, p.camera)
	p.scene._map.DrawCollision(screen, p.camera)

	// 绘制过渡效果 (最上层)
	p.scene.transition.Draw(screen)

	// 绘制调试信息
	p.drawDebugInfo(screen, p.camera)
}

// GetWY 返回角色动作锚点的 World Y 坐标（用于Y-Sorting）
func (p *Role) GetWY() float32 {
	return p.sprite.actionAnchorPointWY
}

// Draw 仅绘制角色本身(实现 Renderable 接口)
func (p *Role) Draw(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	screenX := p.sprite.cameraAnchorPointWX - float32(cam.ViewportWX) - float32(p.sprite.roleImageSprite.Frame.W/2)
	screenY := p.sprite.cameraAnchorPointWY - float32(cam.ViewportWY) - float32(p.sprite.roleImageSprite.Frame.H/2)

	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(p.sprite.image, op)
}

// collectRenderables 收集所有需要Y-Sorting的对象
func (p *Role) collectRenderables() []commonrenderable.IRenderable {
	// Y-Sorting 层
	renderables := make([]commonrenderable.IRenderable, 0, 1024)
	renderables = append(renderables, p) // 角色自己

	// 添加怪物
	for _, enemy := range p.scene._map.spawnManager.GetAllEnemies() {
		renderables = append(renderables, enemy)
	}

	renderables = p.scene._map.GetTileSortInfoSlice(p.camera, cfg.TiledLayerType_Building, renderables)
	renderables = p.scene._map.GetTileSortInfoSlice(p.camera, cfg.TiledLayerType_Objects, renderables)

	return renderables
}

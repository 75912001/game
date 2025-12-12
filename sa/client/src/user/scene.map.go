package user

import (
	"image"
	"image/color"
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/user/camera"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Map Tiled 地图场景
type Map struct {
	id  common.AssetID // 地图ID
	cfg *cfg.TiledMap  // Tiled 地图资源
}

// NewMap 创建 Tiled 地图场景
func NewMap(mapID common.AssetID) *Map {
	m := &Map{
		id: mapID,
	}
	m.cfg = cfg.GTiledMapMgr.Maps.Get(mapID)
	return m
}

// Update 更新
func (p *Map) Update() {
}

// Draw 绘制 Tiled 地图
func (p *Map) Draw(screen *ebitenv2.Image, cam *camera.Camera) {
	// 遍历所有图层
	for _, layer := range p.cfg.Layers {
		if !layer.Visible {
			continue
		}
		if layer.Type != cfg.TiledLayerType_TileLayer {
			continue
		}
		p.drawLayer(screen, cam, layer)
	}

	if true { // 绘制调试边界
		p.drawBorder(screen, cam)
		p.drawCollision(screen, cam)
	}
}

// drawBorder 绘制地图边界(调试用)-红色加粗线条
func (p *Map) drawBorder(screen *ebitenv2.Image, cam *camera.Camera) {
	// 获取菱形四角 (World 坐标)
	topX, topY, rightX, rightY, bottomX, bottomY, leftX, leftY := p.cfg.CT.GetDiamondCorners()

	// World -> Screen
	camX := float32(cam.ViewportX)
	camY := float32(cam.ViewportY)

	sTopX, sTopY := p.cfg.CT.WorldToScreen(topX, topY, camX, camY)
	sRightX, sRightY := p.cfg.CT.WorldToScreen(rightX, rightY, camX, camY)
	sBottomX, sBottomY := p.cfg.CT.WorldToScreen(bottomX, bottomY, camX, camY)
	sLeftX, sLeftY := p.cfg.CT.WorldToScreen(leftX, leftY, camX, camY)

	// 绘制四条边界线(红色加粗)
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	strokeWidth := float32(3.0)

	vector.StrokeLine(screen, float32(sTopX), float32(sTopY), float32(sRightX), float32(sRightY), strokeWidth, red, false)
	vector.StrokeLine(screen, float32(sRightX), float32(sRightY), float32(sBottomX), float32(sBottomY), strokeWidth, red, false)
	vector.StrokeLine(screen, float32(sBottomX), float32(sBottomY), float32(sLeftX), float32(sLeftY), strokeWidth, red, false)
	vector.StrokeLine(screen, float32(sLeftX), float32(sLeftY), float32(sTopX), float32(sTopY), strokeWidth, red, false)
}

// drawCollision 绘制碰撞区域(调试用)-黄色加粗线条
func (p *Map) drawCollision(screen *ebitenv2.Image, cam *camera.Camera) {
	yellow := color.RGBA{R: 255, G: 255, B: 0, A: 255}
	strokeWidth := float32(3.0)

	camX := float32(cam.ViewportX)
	camY := float32(cam.ViewportY)

	// 遍历所有图层，查找对象图层中的碰撞对象
	for _, layer := range p.cfg.Layers {
		if layer.Type != cfg.TiledLayerType_ObjectLayer {
			continue
		}
		for _, obj := range layer.Objects {
			if !obj.Collision {
				continue
			}
			// 绘制多边形碰撞区域
			if len(obj.Polygon) > 0 {
				for i := 0; i < len(obj.Polygon); i++ {
					// 获取当前点和下一个点（形成闭合多边形）
					curr := obj.Polygon[i]
					next := obj.Polygon[(i+1)%len(obj.Polygon)]

					// 多边形顶点是相对于对象位置的坐标，转换为 World 坐标
					worldX1 := obj.X + curr.X
					worldY1 := obj.Y + curr.Y
					worldX2 := obj.X + next.X
					worldY2 := obj.Y + next.Y

					// World -> Screen
					screenX1, screenY1 := p.cfg.CT.WorldToScreen(worldX1, worldY1, camX, camY)
					screenX2, screenY2 := p.cfg.CT.WorldToScreen(worldX2, worldY2, camX, camY)

					vector.StrokeLine(screen, screenX1, screenY1, screenX2, screenY2, strokeWidth, yellow, false)
				}
			} else {
				// 绘制矩形碰撞区域
				// obj.X, obj.Y 是 Tiled 中的像素坐标，需要转换为 Tile 坐标
				// Tiled 等距地图中，像素坐标除以 TileHeight 得到 Tile 坐标
				tileX := obj.X / float32(p.cfg.TileHeight)
				tileY := obj.Y / float32(p.cfg.TileHeight)
				tileW := obj.Width / float32(p.cfg.TileHeight)
				tileH := obj.Height / float32(p.cfg.TileHeight)

				// 四个角点的 Tile 坐标，直接转换为 Screen 坐标
				x1, y1 := p.cfg.CT.TileToScreen(tileX, tileY, camX, camY)
				x2, y2 := p.cfg.CT.TileToScreen(tileX+tileW, tileY, camX, camY)
				x3, y3 := p.cfg.CT.TileToScreen(tileX+tileW, tileY+tileH, camX, camY)
				x4, y4 := p.cfg.CT.TileToScreen(tileX, tileY+tileH, camX, camY)

				vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, yellow, false)
				vector.StrokeLine(screen, x2, y2, x3, y3, strokeWidth, yellow, false)
				vector.StrokeLine(screen, x3, y3, x4, y4, strokeWidth, yellow, false)
				vector.StrokeLine(screen, x4, y4, x1, y1, strokeWidth, yellow, false)
			}
		}
	}
}

// drawLayer 绘制单个图层
func (p *Map) drawLayer(screen *ebitenv2.Image, cam *camera.Camera, layer *cfg.TiledLayer) {
	if 0 < len(layer.Data) {
		p.drawData(screen, cam, layer.Data, layer.Width, layer.Height)
	}
}

// drawData 绘制 tile 数据
func (p *Map) drawData(screen *ebitenv2.Image, cam *camera.Camera, data []int, width, height int) {
	for i, gid := range data {
		if gid == 0 {
			continue
		}

		tileX := i % width
		tileY := i / width

		// 使用 CT 获取 tile 图像的屏幕位置
		screenX, screenY := p.cfg.CT.TileToImageScreenPos(tileX, tileY, cam.ViewportX, cam.ViewportY)

		// 裁剪：跳过屏幕外的 tile
		if screenX < -p.cfg.TileWidth || cfg.GCommon.ScreenMaxWidth < screenX ||
			screenY < -p.cfg.TileHeight || cfg.GCommon.ScreenMaxHeight < screenY {
			continue
		}

		// 获取 tile 图像
		// todo menglc 优化：可以缓存 tile 图像，避免重复获取
		tileImg := p.getTileImage(gid)
		if tileImg == nil {
			continue
		}

		// 绘制
		op := &ebitenv2.DrawImageOptions{}
		op.GeoM.Translate(float64(screenX), float64(screenY))
		screen.DrawImage(tileImg, op)
	}
}

// getTileImage 根据 GID 获取 tile 图像
func (p *Map) getTileImage(gid int) *ebitenv2.Image {
	var tileset *cfg.TiledTileset
	for i := len(p.cfg.Tilesets) - 1; i >= 0; i-- {
		ts := p.cfg.Tilesets[i]
		if gid >= ts.FirstGID {
			tileset = ts
			break
		}
	}
	if tileset == nil || tileset.Image == nil {
		return nil
	}

	localID := gid - tileset.FirstGID
	if localID < 0 || localID >= tileset.TileCount {
		return nil
	}

	col := localID % tileset.Columns
	row := localID / tileset.Columns
	x := col * tileset.TileWidth
	y := row * tileset.TileHeight

	return tileset.Image.SubImage(image.Rect(x, y, x+tileset.TileWidth, y+tileset.TileHeight)).(*ebitenv2.Image)
}

// ClampTileBounds 将 tile 坐标限制在地图边界内
func (p *Map) ClampTileBounds(tileX, tileY float32) (clampedTX, clampedTY float32) {
	return p.cfg.CT.ClampTile(tileX, tileY)
}

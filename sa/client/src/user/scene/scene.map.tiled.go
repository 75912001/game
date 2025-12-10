package scene

import (
	"image"
	"saClient/src/cfg"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/common"
	"saClient/src/res"
	"saClient/src/user/camera"
)

// TiledMap Tiled 地图场景
type TiledMap struct {
	id       common.AssetID // 地图ID
	tiledMap *res.TiledMap  // Tiled 地图资源
}

// NewTiledMap 创建 Tiled 地图场景
func NewTiledMap(mapID common.AssetID) *TiledMap {
	m := &TiledMap{
		id: mapID,
	}
	m.tiledMap = res.GTiledMapMgr.Maps.Get(mapID)
	return m
}

// Update 更新
func (p *TiledMap) Update() {
}

// Draw 绘制 Tiled 地图
func (p *TiledMap) Draw(screen *ebitenv2.Image, cam *camera.Camera) {
	// 遍历所有图层
	for _, layer := range p.tiledMap.Layers {
		if !layer.Visible || layer.Type != "tilelayer" {
			continue
		}
		p.drawLayer(screen, cam, layer)
	}
}

// drawLayer 绘制单个图层
func (p *TiledMap) drawLayer(screen *ebitenv2.Image, cam *camera.Camera, layer *res.TiledLayer) {
	// 处理 infinite 地图的 chunks
	if len(layer.Chunks) > 0 {
		for _, chunk := range layer.Chunks {
			p.drawChunk(screen, cam, chunk)
		}
		return
	}

	// 处理有限地图的 data
	if len(layer.Data) > 0 {
		p.drawData(screen, cam, layer.Data, 0, 0, layer.Width, layer.Height)
	}
}

// drawChunk 绘制数据块
func (p *TiledMap) drawChunk(screen *ebitenv2.Image, cam *camera.Camera, chunk *res.TiledChunk) {
	p.drawData(screen, cam, chunk.Data, chunk.X, chunk.Y, chunk.Width, chunk.Height)
}

// drawData 绘制 tile 数据
func (p *TiledMap) drawData(screen *ebitenv2.Image, cam *camera.Camera, data []int, startX, startY, width, height int) {
	for i, gid := range data {
		if gid == 0 {
			continue // 空 tile
		}

		// 计算 tile 在 chunk 中的位置
		localX := i % width
		localY := i / width
		tileX := startX + localX
		tileY := startY + localY

		// 计算屏幕位置
		screenX, screenY := p.getTileScreenPos(tileX, tileY)
		screenX -= cam.ScreenX
		screenY -= cam.ScreenY

		// 裁剪：跳过屏幕外的 tile
		if screenX < -p.tiledMap.TileWidth || screenX > cfg.GCommon.ScreenMaxWidth ||
			screenY < -p.tiledMap.TileHeight || screenY > cfg.GCommon.ScreenMaxHeight {
			continue
		}

		// 获取 tile 图像
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

// getTileScreenPos 获取 tile 的屏幕位置（世界坐标）
func (p *TiledMap) getTileScreenPos(tileX, tileY int) (screenX, screenY int) {
	switch p.tiledMap.Orientation {
	case "staggered":
		// 等距交错地图
		screenX = tileX * p.tiledMap.TileWidth
		if p.tiledMap.StaggerAxis == "y" {
			// Y 轴交错
			if p.tiledMap.StaggerIndex == "odd" {
				if tileY%2 == 1 {
					screenX += p.tiledMap.TileWidth / 2
				}
			} else { // even
				if tileY%2 == 0 {
					screenX += p.tiledMap.TileWidth / 2
				}
			}
			screenY = tileY * (p.tiledMap.TileHeight / 2)
		} else {
			// X 轴交错（较少使用）
			screenY = tileY * p.tiledMap.TileHeight
			if p.tiledMap.StaggerIndex == "odd" {
				if tileX%2 == 1 {
					screenY += p.tiledMap.TileHeight / 2
				}
			} else {
				if tileX%2 == 0 {
					screenY += p.tiledMap.TileHeight / 2
				}
			}
		}
	case "isometric":
		// 等距地图（菱形）
		// 原始公式会让左半边地图出现负坐标，需要加上偏移
		// 偏移量 = Height * TileWidth / 2，让最左边的 tile(0, Height-1) 的 X 坐标从 0 开始
		offsetX := p.tiledMap.Height * (p.tiledMap.TileWidth / 2)
		screenX = (tileX-tileY)*(p.tiledMap.TileWidth/2) + offsetX
		screenY = (tileX + tileY) * (p.tiledMap.TileHeight / 2)
	default:
		// 正交地图
		screenX = tileX * p.tiledMap.TileWidth
		screenY = tileY * p.tiledMap.TileHeight
	}
	return
}

// getTileImage 根据 GID 获取 tile 图像
func (p *TiledMap) getTileImage(gid int) *ebitenv2.Image {
	// 查找对应的 tileset
	var tileset *res.TiledTileset
	for i := len(p.tiledMap.Tilesets) - 1; i >= 0; i-- {
		ts := p.tiledMap.Tilesets[i]
		if gid >= ts.FirstGID {
			tileset = ts
			break
		}
	}
	if tileset == nil || tileset.Image == nil {
		return nil
	}

	// 计算 tile 在 tileset 中的位置
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

// GetMapSize 获取地图尺寸（像素）
func (p *TiledMap) GetMapSize() (width, height int) {
	return p.tiledMap.GetPixelWidth(), p.tiledMap.GetPixelHeight()
}

// GetTiledMap 获取 Tiled 地图资源
func (p *TiledMap) GetTiledMap() *res.TiledMap {
	return p.tiledMap
}

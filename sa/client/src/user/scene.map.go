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
	id          common.AssetID // 地图ID
	tiledMapCfg *cfg.TiledMap  // Tiled 地图资源
	mapCfg      *cfg.Map       // 地图配置
}

// NewMap 创建 Tiled 地图场景
func NewMap(mapID common.AssetID) *Map {
	m := &Map{
		id: mapID,
	}
	m.tiledMapCfg = cfg.GTiledMapMgr.Maps.Get(mapID)
	m.mapCfg = cfg.GMapMgr.Maps.Get(mapID)
	return m
}

// Update 更新
func (p *Map) Update() {
}

// Draw 绘制 Tiled 地图
func (p *Map) Draw(screen *ebitenv2.Image, cam *camera.Camera) {
	// 遍历所有图层
	for _, layer := range p.tiledMapCfg.Layers {
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
		p.drawPortal(screen, cam)
		p.drawBlocked(screen, cam)
	}
}

// drawBorder 绘制地图边界(调试用)
func (p *Map) drawBorder(screen *ebitenv2.Image, cam *camera.Camera) {
	// World -> Screen
	camX := float32(cam.ViewportX)
	camY := float32(cam.ViewportY)

	sTopX, sTopY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.TopWX, p.tiledMapCfg.TopWY, camX, camY)
	sRightX, sRightY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.RightWX, p.tiledMapCfg.RightWY, camX, camY)
	sBottomX, sBottomY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.BottomWX, p.tiledMapCfg.BottomWY, camX, camY)
	sLeftX, sLeftY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.LeftWX, p.tiledMapCfg.LeftWY, camX, camY)

	// 绘制四条边界线
	drawDiamond(screen, sTopX, sTopY, sRightX, sRightY, sBottomX, sBottomY, sLeftX, sLeftY, common.Colors_Red, 3.0)
}

// drawPortal 绘制传送区域(调试用)
func (p *Map) drawPortal(screen *ebitenv2.Image, cam *camera.Camera) {
	camX := float32(cam.ViewportX)
	camY := float32(cam.ViewportY)

	for _, layer := range p.tiledMapCfg.Layers {
		if layer.Type != cfg.TiledLayerType_ObjectLayer {
			continue
		}
		for _, obj := range layer.Objects {
			if obj.PortalCfg == nil { // 非-传送点
				continue
			}
			p.drawTiledObjectRect(screen, obj, camX, camY, common.Colors_Yellow, 3.0)
		}
	}
}

// drawBlocked 绘制阻挡(调试用)
func (p *Map) drawBlocked(screen *ebitenv2.Image, cam *camera.Camera) {
	camX := float32(cam.ViewportX)
	camY := float32(cam.ViewportY)

	for _, layer := range p.tiledMapCfg.Layers {
		if layer.Type != cfg.TiledLayerType_ObjectLayer {
			continue
		}
		for _, obj := range layer.Objects {
			if !obj.Blocked { // 非-阻挡
				continue
			}
			p.drawTiledObjectRect(screen, obj, camX, camY, common.Colors_Red, 3.0)
		}
	}
}

// drawTiledObjectRect 绘制 TiledObject 的等距矩形边界
func (p *Map) drawTiledObjectRect(screen *ebitenv2.Image, obj *cfg.TiledObject, camX, camY float32, clr color.RGBA, strokeWidth float32) {
	// obj.X, obj.Y 是 Tiled 中的像素坐标，转换为 Tile 坐标
	th := float32(p.tiledMapCfg.TileHeight)
	tileX := obj.X/th - 0.5
	tileY := obj.Y/th - 0.5
	tileW := obj.Width / th
	tileH := obj.Height / th

	// 四个角点的 Tile 坐标 -> Screen 坐标
	x1, y1 := p.tiledMapCfg.IsometricCT.T2S(tileX, tileY, camX, camY)
	x2, y2 := p.tiledMapCfg.IsometricCT.T2S(tileX+tileW, tileY, camX, camY)
	x3, y3 := p.tiledMapCfg.IsometricCT.T2S(tileX+tileW, tileY+tileH, camX, camY)
	x4, y4 := p.tiledMapCfg.IsometricCT.T2S(tileX, tileY+tileH, camX, camY)

	drawDiamond(screen, x1, y1, x2, y2, x3, y3, x4, y4, clr, strokeWidth)
}

// drawDiamond 绘制菱形边框
func drawDiamond(screen *ebitenv2.Image, x1, y1, x2, y2, x3, y3, x4, y4 float32, clr color.RGBA, strokeWidth float32) {
	vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, clr, false)
	vector.StrokeLine(screen, x2, y2, x3, y3, strokeWidth, clr, false)
	vector.StrokeLine(screen, x3, y3, x4, y4, strokeWidth, clr, false)
	vector.StrokeLine(screen, x4, y4, x1, y1, strokeWidth, clr, false)
}

// drawLayer 绘制单个图层
func (p *Map) drawLayer(screen *ebitenv2.Image, cam *camera.Camera, layer *cfg.TiledLayer) {
	if 0 < len(layer.Data) {
		p.drawData(screen, cam, layer.Data, layer.Width, layer.Height)
	}
}

// drawData 绘制 tile 数据
// todo menglc 研究逻辑, 优化计算, 让相关的瓦片挂载在地图上, 不用每次绘制的时候都放在内存中. 每次加载地图.都将该地图资源全量加载, 离开地图时,卸载地图所占资源
func (p *Map) drawData(screen *ebitenv2.Image, cam *camera.Camera, data []int, width, height int) {
	for i, gid := range data {
		if gid == 0 {
			continue
		}

		tileX := i % width
		tileY := i / width

		// 使用 IsometricCT 获取 tile 图像的屏幕位置
		screenX, screenY := p.tiledMapCfg.IsometricCT.TileImageScreenPos(tileX, tileY, cam.ViewportX, cam.ViewportY)

		// 获取 tile 图像
		// todo menglc 优化：可以缓存 tile 图像，避免重复获取
		tileImg, tileset := p.getTileImage(gid)
		if tileImg == nil {
			continue
		}

		// 修正 Y 坐标：tileset 高度大于 map tile 高度时，图像底部需对齐到 tile 基准位置
		screenY -= tileset.TileHeight - p.tiledMapCfg.TileHeight

		// 裁剪：跳过屏幕外的 tile
		if screenX < -tileset.TileWidth || cfg.GCommon.ScreenMaxWidth < screenX ||
			screenY < -tileset.TileHeight || cfg.GCommon.ScreenMaxHeight < screenY {
			continue
		}

		// 绘制
		op := &ebitenv2.DrawImageOptions{}
		op.GeoM.Translate(float64(screenX), float64(screenY))
		screen.DrawImage(tileImg, op)
	}
}

// getTileImage 根据 GID 获取 tile 图像和所属的 tileset
func (p *Map) getTileImage(gid int) (*ebitenv2.Image, *cfg.TiledTileset) {
	var tileset *cfg.TiledTileset
	for i := len(p.tiledMapCfg.Tilesets) - 1; i >= 0; i-- {
		ts := p.tiledMapCfg.Tilesets[i]
		if gid >= ts.FirstGID {
			tileset = ts
			break
		}
	}
	if tileset == nil || tileset.Image == nil {
		return nil, nil
	}

	localID := gid - tileset.FirstGID
	if localID < 0 || localID >= tileset.TileCount {
		return nil, nil
	}

	col := localID % tileset.Columns
	row := localID / tileset.Columns
	x := col * tileset.TileWidth
	y := row * tileset.TileHeight

	img := tileset.Image.SubImage(image.Rect(x, y, x+tileset.TileWidth, y+tileset.TileHeight)).(*ebitenv2.Image)
	return img, tileset
}

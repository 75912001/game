package scene

import (
	"image"
	"image/color"
	"saClient/src/cfg"
	restiled "saClient/src/res/tiled"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"saClient/src/common"
	"saClient/src/user/camera"
)

// TiledMap Tiled 地图场景
type TiledMap struct {
	id       common.AssetID     // 地图ID
	tiledMap *restiled.TiledMap // Tiled 地图资源

	DebugDrawBorder bool // 是否绘制地图边界(调试用)
}

// NewTiledMap 创建 Tiled 地图场景
func NewTiledMap(mapID common.AssetID) *TiledMap {
	m := &TiledMap{
		id:              mapID,
		DebugDrawBorder: true,
	}
	m.tiledMap = restiled.GMapMgr.Maps.Get(mapID)
	return m
}

// Update 更新
func (p *TiledMap) Update() {
}

// Draw 绘制 Tiled 地图
func (p *TiledMap) Draw(screen *ebitenv2.Image, cam *camera.Camera) {
	// 遍历所有图层
	for _, layer := range p.tiledMap.Layers {
		if !layer.Visible { // 不可见图层跳过
			continue
		}
		if layer.Type != restiled.LayerType_TileLayer { // 仅处理瓦片图层
			continue
		}
		p.drawLayer(screen, cam, layer)
	}

	// 绘制调试边界
	if p.DebugDrawBorder {
		p.drawBorder(screen, cam)
	}
}

// drawBorder 绘制地图边界(调试用)-红色加粗线条
func (p *TiledMap) drawBorder(screen *ebitenv2.Image, cam *camera.Camera) {
	w := p.tiledMap.Width
	h := p.tiledMap.Height
	tw := p.tiledMap.TileWidth
	th := p.tiledMap.TileHeight

	// 等距地图的四个角点(世界坐标)
	// 顶点: tile(0,0) 的顶部
	// 右点: tile(W-1,0) 的右侧
	// 底点: tile(W-1,H-1) 的底部
	// 左点: tile(0,H-1) 的左侧
	offsetX := h * (tw / 2)

	// 顶点 - tile(0,0) 的顶部中心
	topX := float32(offsetX + tw/2)
	topY := float32(0)

	// 右点 - tile(W-1,0) 的右侧
	rightX := float32((w-1-0)*(tw/2) + offsetX + tw)
	rightY := float32((w-1+0)*(th/2) + th/2)

	// 底点 - tile(W-1,H-1) 的底部
	bottomX := float32((w-1-(h-1))*(tw/2) + offsetX + tw/2)
	bottomY := float32((w-1+(h-1))*(th/2) + th)

	// 左点 - tile(0,H-1) 的左侧
	leftX := float32((0-(h-1))*(tw/2) + offsetX)
	leftY := float32((0+(h-1))*(th/2) + th/2)

	// 转换为屏幕坐标
	camX := float32(cam.ScreenX)
	camY := float32(cam.ScreenY)
	topX -= camX
	topY -= camY
	rightX -= camX
	rightY -= camY
	bottomX -= camX
	bottomY -= camY
	leftX -= camX
	leftY -= camY

	// 绘制四条边界线(红色加粗)
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	strokeWidth := float32(3.0)

	// 顶点 -> 右点
	vector.StrokeLine(screen, topX, topY, rightX, rightY, strokeWidth, red, false)
	// 右点 -> 底点
	vector.StrokeLine(screen, rightX, rightY, bottomX, bottomY, strokeWidth, red, false)
	// 底点 -> 左点
	vector.StrokeLine(screen, bottomX, bottomY, leftX, leftY, strokeWidth, red, false)
	// 左点 -> 顶点
	vector.StrokeLine(screen, leftX, leftY, topX, topY, strokeWidth, red, false)
}

// drawLayer 绘制单个图层
func (p *TiledMap) drawLayer(screen *ebitenv2.Image, cam *camera.Camera, layer *restiled.TiledLayer) {
	// 处理有限地图的 data
	if len(layer.Data) > 0 {
		p.drawData(screen, cam, layer.Data, layer.Width, layer.Height)
	}
}

// drawData 绘制 tile 数据
func (p *TiledMap) drawData(screen *ebitenv2.Image, cam *camera.Camera, data []int, width, height int) {
	for i, gid := range data {
		if gid == 0 {
			continue // 空 tile
		}

		// 计算 tile 在 chunk 中的位置
		localX := i % width
		localY := i / width
		tileX := localX
		tileY := localY

		// 计算屏幕位置
		screenX, screenY := p.getTileScreenPos(tileX, tileY)
		screenX -= cam.ScreenX
		screenY -= cam.ScreenY

		// 裁剪：跳过屏幕外的 tile
		if screenX < -p.tiledMap.TileWidth || cfg.GCommon.ScreenMaxWidth < screenX ||
			screenY < -p.tiledMap.TileHeight || cfg.GCommon.ScreenMaxHeight < screenY {
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
	var tileset *restiled.TiledTileset
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

// GetMapSize 获取地图尺寸(像素)
func (p *TiledMap) GetMapSize() (width, height int) {
	return p.tiledMap.GetPixelWidth(), p.tiledMap.GetPixelHeight()
}

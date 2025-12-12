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
	halfTW := tw / 2
	halfTH := th / 2

	// 边界与 GetMapSize 一致，从 X=0 开始到 X=GetPixelWidth 结束
	// GetPixelWidth = (w+h) * halfTW
	// GetPixelHeight = (w+h) * halfTH
	//
	// 菱形边界的四个角点:
	// 顶点: X = h*halfTW, Y = 0
	// 右点: X = (w+h)*halfTW, Y = w*halfTH
	// 底点: X = w*halfTW, Y = (w+h)*halfTH
	// 左点: X = 0, Y = h*halfTH

	// 顶点
	topX := float32(h * halfTW)
	topY := float32(0)

	// 右点
	rightX := float32((w + h) * halfTW)
	rightY := float32(w * halfTH)

	// 底点
	bottomX := float32(w * halfTW)
	bottomY := float32((w + h) * halfTH)

	// 左点
	leftX := float32(0)
	leftY := float32(h * halfTH)

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
		// 偏移量 = (Height-1) * TileWidth / 2，让最左边的 tile(0, Height-1) 的 X 坐标从 0 开始
		offsetX := (p.tiledMap.Height - 1) * (p.tiledMap.TileWidth / 2)
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

// WorldToTile 世界坐标转换为 tile 坐标
// worldX, worldY: 世界坐标(像素),基于 tile 菱形中心
// 返回: tileX, tileY (浮点数,可用于精确位置判断)
func (p *TiledMap) WorldToTile(worldX, worldY float64) (tileX, tileY float64) {
	tw := float64(p.tiledMap.TileWidth)
	th := float64(p.tiledMap.TileHeight)
	halfTW := tw / 2
	halfTH := th / 2
	offsetX := float64(p.tiledMap.Height-1) * halfTW

	// getTileScreenPos 返回 tile 图像左上角位置
	// tile 菱形中心相对于左上角的偏移是 (halfTW, halfTH)
	// 所以需要先将世界坐标转换为 tile 图像左上角坐标
	screenX := worldX - halfTW
	screenY := worldY - halfTH

	// 逆变换公式:
	// screenX = (tileX - tileY) * halfTW + offsetX
	// screenY = (tileX + tileY) * halfTH
	// 解方程组得:
	sx := screenX - offsetX
	tileX = (sx/halfTW + screenY/halfTH) / 2
	tileY = (screenY/halfTH - sx/halfTW) / 2
	return
}

// TileToWorld tile 坐标转换为世界坐标
// tileX, tileY: tile 坐标
// 返回: worldX, worldY (世界坐标,像素,基于 tile 菱形中心)
func (p *TiledMap) TileToWorld(tileX, tileY float64) (worldX, worldY float64) {
	tw := float64(p.tiledMap.TileWidth)
	th := float64(p.tiledMap.TileHeight)
	halfTW := tw / 2
	halfTH := th / 2
	offsetX := float64(p.tiledMap.Height-1) * halfTW

	// 先计算 tile 图像左上角位置
	screenX := (tileX-tileY)*halfTW + offsetX
	screenY := (tileX + tileY) * halfTH

	// 加上菱形中心偏移得到世界坐标
	worldX = screenX + halfTW
	worldY = screenY + halfTH
	return
}

// IsInMapBounds 检查世界坐标是否在地图边界内
func (p *TiledMap) IsInMapBounds(worldX, worldY float64) bool {
	tileX, tileY := p.WorldToTile(worldX, worldY)
	w := float64(p.tiledMap.Width)
	h := float64(p.tiledMap.Height)
	return tileX >= 0 && tileX < w && tileY >= 0 && tileY < h
}

// ClampToMapBounds 将世界坐标限制在地图边界内
// 如果坐标超出边界,将其限制到最近的边界位置
func (p *TiledMap) ClampToMapBounds(worldX, worldY float64) (clampedX, clampedY float64) {
	tileX, tileY := p.WorldToTile(worldX, worldY)
	w := float64(p.tiledMap.Width)
	h := float64(p.tiledMap.Height)

	// 限制 tile 坐标到有效范围
	// 使用小于 Width/Height 的最大值,避免刚好在边界外
	maxTileX := w - 0.01
	maxTileY := h - 0.01

	if tileX < 0 {
		tileX = 0
	} else if tileX > maxTileX {
		tileX = maxTileX
	}

	if tileY < 0 {
		tileY = 0
	} else if tileY > maxTileY {
		tileY = maxTileY
	}

	// 转换回世界坐标
	return p.TileToWorld(tileX, tileY)
}

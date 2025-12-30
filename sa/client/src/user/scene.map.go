package user

import (
	xmap "github.com/75912001/xlib/map"
	xpool "github.com/75912001/xlib/pool"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"saClient/src/cfg"
	"saClient/src/common"
	commoncamera "saClient/src/common/camera"
	commonrenderable "saClient/src/common/renderable"
)

// TileCacheInfo 缓存的瓦片信息
type TileCacheInfo struct {
	Image  *ebitenv2.Image // 裁剪后的图像
	Width  int             // tileset 的 tile 宽度
	Height int             // tileset 的 tile 高度

	SlicedHorizontal *cfg.TileSlicedHorizontal // 水平切片-瓦片 (如树木, 雕像建筑)
	SlicedVertical   *cfg.TileSlicedVertical   // 垂直切片-瓦片 (如大型建筑)
}

// UmbrellaShapedCanopyDrawInfo 冠部-绘制信息(用于收集待绘制的冠部)
type UmbrellaShapedCanopyDrawInfo struct {
	Image   *ebitenv2.Image // 冠部-图像
	ScreenX int             // 屏幕 X 坐标
	ScreenY int             // 屏幕 Y 坐标
}

func (p *UmbrellaShapedCanopyDrawInfo) reset() {
	p.Image = nil
	p.ScreenX = 0
	p.ScreenY = 0
}

// TileSortInfo 待排序的瓦片绘制信息（用于 Y-Sorting）
type TileSortInfo struct {
	Image   *ebitenv2.Image // 瓦片图像
	ScreenX int             // 屏幕 X 坐标
	ScreenY int             // 屏幕 Y 坐标
	WY      float32         // World Y 坐标（用于排序）
}

func (p *TileSortInfo) reset() {
	p.Image = nil
	p.ScreenX = 0
	p.ScreenY = 0
	p.WY = 0
}

// GetWY 实现 IRenderable 接口
func (p *TileSortInfo) GetWY() float32 {
	return p.WY
}

// Draw 实现 IRenderable 接口
func (p *TileSortInfo) Draw(screen *ebitenv2.Image, camera *commoncamera.Camera) {
	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(p.ScreenX), float64(p.ScreenY))
	screen.DrawImage(p.Image, op)
}

// Map Tiled 地图场景
type Map struct {
	id           common.AssetID     // 地图ID
	tiledMapCfg  *cfg.TiledMap      // Tiled 地图资源
	mapCfg       *cfg.Map           // 地图配置
	spawnManager *ArpgEnemySpawnMgr // 刷怪点管理器

	tileCache    *xmap.MapMgr[int, *TileCacheInfo] // key:GID val:TileCacheInfo 缓存 (位于 多个图层中 相同 GID 只存一份)
	tileSortPool *xpool.Pool[*TileSortInfo]        // TileSortInfo 对象池

	umbrellaShapedCanopyDrawPool *xpool.Pool[*UmbrellaShapedCanopyDrawInfo] // UmbrellaShapedCanopyDrawInfo 对象池
	umbrellaShapedCanopyDrawList []*UmbrellaShapedCanopyDrawInfo            // 当前帧待绘制-上部列表
}

// NewMap 创建 Tiled 地图场景
func NewMap(mapID common.AssetID) *Map {
	m := &Map{
		id:        mapID,
		tileCache: xmap.NewMapMgr[int, *TileCacheInfo](),
		tileSortPool: xpool.NewPool(
			func() *TileSortInfo { return &TileSortInfo{} },
			func(t *TileSortInfo) {
				t.reset()
			},
		),
		umbrellaShapedCanopyDrawPool: xpool.NewPool(
			func() *UmbrellaShapedCanopyDrawInfo { return &UmbrellaShapedCanopyDrawInfo{} },
			func(t *UmbrellaShapedCanopyDrawInfo) {
				t.reset()
			},
		),
		umbrellaShapedCanopyDrawList: make([]*UmbrellaShapedCanopyDrawInfo, 0, 256),
	}
	m.tiledMapCfg = cfg.GTiledMapMgr.Maps.Get(mapID)
	m.mapCfg = cfg.GMapMgr.Maps.Get(mapID)
	m.spawnManager = NewArpgEnemySpawnMgr(m.tiledMapCfg)
	m.buildTileCache()
	return m
}

// buildTileCache 预构建瓦片缓存
func (p *Map) buildTileCache() {
	// 收集所有使用的 GID
	usedGIDs := make(map[int]struct{})
	for _, layer := range p.tiledMapCfg.Layers {
		for _, gid := range layer.Data {
			if gid != 0 {
				usedGIDs[gid] = struct{}{}
			}
		}
	}

	// 为每个 GID 创建缓存
	for gid := range usedGIDs {
		img, tileset := p.getTileImage(gid)
		if img != nil {
			cacheInfo := &TileCacheInfo{
				Image:  img,
				Width:  tileset.TileWidth,
				Height: tileset.TileHeight,
			}
			// 处理水平切分的 tile
			if 0 < tileset.HorizontalSlicePixel {
				tileSliceHorizontal := cfg.NewTileSliceHorizontal()
				tileSliceHorizontal.Split(img, tileset)
				cacheInfo.SlicedHorizontal = tileSliceHorizontal // 存入缓存
			}
			// 处理垂直切分的 tile
			if 0 < tileset.VerticalSlicePixel {
				verticalSliced := cfg.NewTileSlicedVertical()
				verticalSliced.Split(img, tileset)
				cacheInfo.SlicedVertical = verticalSliced // 存入缓存
			}
			p.tileCache.Add(gid, cacheInfo)
		}
	}
}

// getTileImageCached 从缓存获取瓦片信息
func (p *Map) getTileImageCached(gid int) *TileCacheInfo {
	return p.tileCache.Get(gid)
}

// Update 更新
func (p *Map) Update() {
	p.spawnManager.Update()
}

// DrawCollision 碰撞层
func (p *Map) DrawCollision(screen *ebitenv2.Image, cam *commoncamera.Camera) {
}

// DrawGround 地面
func (p *Map) DrawGround(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	// 遍历所有图层
	for _, layer := range p.tiledMapCfg.Layers {
		if !layer.Visible {
			continue
		}
		if layer.LayerType != cfg.TiledLayerType_Ground {
			continue
		}
		p.drawLayer(screen, cam, layer)
	}
}

// DrawOverhead 头顶
func (p *Map) DrawOverhead(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	// 绘制拆分 tile 的冠部
	p.drawSplitUmbrellaShapedCanopy(screen)
	// 遍历所有图层
	for _, layer := range p.tiledMapCfg.Layers {
		if !layer.Visible {
			continue
		}
		if layer.LayerType != cfg.TiledLayerType_Overhead {
			continue
		}
		p.drawLayer(screen, cam, layer)
	}
}

// drawSplitUmbrellaShapedCanopy 绘制拆分 tile 的冠部
func (p *Map) drawSplitUmbrellaShapedCanopy(screen *ebitenv2.Image) {
	for _, canopy := range p.umbrellaShapedCanopyDrawList {
		op := &ebitenv2.DrawImageOptions{}
		op.GeoM.Translate(float64(canopy.ScreenX), float64(canopy.ScreenY))
		screen.DrawImage(canopy.Image, op)
	}

	// 回收所有冠部信息
	for _, canopy := range p.umbrellaShapedCanopyDrawList {
		p.umbrellaShapedCanopyDrawPool.Put(canopy)
	}
	p.umbrellaShapedCanopyDrawList = p.umbrellaShapedCanopyDrawList[:0]
}

// drawBorder 绘制地图边界(调试用)
func (p *Map) drawBorder(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	// World -> Screen
	camX := float32(cam.ViewportWX)
	camY := float32(cam.ViewportWY)

	sTopX, sTopY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.Boundary.TopWX, p.tiledMapCfg.Boundary.TopWY, camX, camY)
	sRightX, sRightY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.Boundary.RightWX, p.tiledMapCfg.Boundary.RightWY, camX, camY)
	sBottomX, sBottomY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.Boundary.BottomWX, p.tiledMapCfg.Boundary.BottomWY, camX, camY)
	sLeftX, sLeftY := p.tiledMapCfg.IsometricCT.W2S(p.tiledMapCfg.Boundary.LeftWX, p.tiledMapCfg.Boundary.LeftWY, camX, camY)

	// 绘制四条边界线
	drawDiamond(screen, sTopX, sTopY, sRightX, sRightY, sBottomX, sBottomY, sLeftX, sLeftY, common.Colors_Red, 3.0)
}

// drawPortal 绘制传送区域(调试用)
func (p *Map) drawPortal(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	camX := float32(cam.ViewportWX)
	camY := float32(cam.ViewportWY)

	for _, layer := range p.tiledMapCfg.Layers {
		if layer.LayerType != cfg.TiledLayerType_Collision {
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
func (p *Map) drawBlocked(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	camX := float32(cam.ViewportWX)
	camY := float32(cam.ViewportWY)

	for _, layer := range p.tiledMapCfg.Layers {
		if layer.LayerType != cfg.TiledLayerType_Collision {
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

// drawSpawnPointDebug 绘制刷怪点调试信息(调试用)
func (p *Map) drawSpawnPointDebug(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	camX := float32(cam.ViewportWX)
	camY := float32(cam.ViewportWY)

	// 颜色定义
	colorSpawn := color.RGBA{255, 0, 0, 200}    // 红色 - 生成半径
	colorPatrol := color.RGBA{255, 165, 0, 200} // 橙色 - 巡逻半径
	colorChase := color.RGBA{255, 255, 0, 200}  // 黄色 - 追击半径

	for _, sp := range p.spawnManager.GetSpawnPoints() {
		obj := sp.Object
		// 转换圆心到屏幕坐标
		cx, cy := p.tiledMapCfg.IsometricCT.W2S(obj.WX, obj.WY, camX, camY)

		// 绘制生成半径 (红色)
		if obj.SpawnRadius > 0 {
			drawDashedCircle(screen, cx, cy, obj.SpawnRadius, colorSpawn, 48, 0.5, 2.0)
		}

		// 绘制巡逻半径 (橙色)
		if obj.PatrolRadius > 0 {
			drawDashedCircle(screen, cx, cy, obj.PatrolRadius, colorPatrol, 48, 0.5, 2.0)
		}

		// 绘制追击半径 (黄色)
		chaseRadius := obj.PatrolRadius * cfg.GCommon.PetDefArpgChaseRadiusMultiplier
		if chaseRadius > 0 {
			drawDashedCircle(screen, cx, cy, chaseRadius, colorChase, 48, 0.5, 2.0)
		}
	}
}

// drawDashedCircle 绘制虚线圆
// cx, cy: 圆心屏幕坐标
// radius: 半径（世界坐标单位）
// clr: 颜色
// dashCount: 虚线段数量
// dashRatio: 实线占比（0.5 表示实线和空白各占一半）
// strokeWidth: 线宽
func drawDashedCircle(screen *ebitenv2.Image, cx, cy, radius float32, clr color.RGBA, dashCount int, dashRatio float32, strokeWidth float32) {
	if dashCount <= 0 || radius <= 0 {
		return
	}

	angleStep := 2 * math.Pi / float64(dashCount)
	dashAngle := angleStep * float64(dashRatio)

	for i := 0; i < dashCount; i++ {
		startAngle := float64(i) * angleStep
		endAngle := startAngle + dashAngle

		// 计算起点和终点
		x1 := cx + radius*float32(math.Cos(startAngle))
		y1 := cy + radius*float32(math.Sin(startAngle))
		x2 := cx + radius*float32(math.Cos(endAngle))
		y2 := cy + radius*float32(math.Sin(endAngle))

		vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, clr, false)
	}
}

// drawTiledObjectRect 绘制 TiledObject 的矩形边界
func (p *Map) drawTiledObjectRect(screen *ebitenv2.Image, obj *cfg.TiledObject, camX, camY float32, clr color.RGBA, strokeWidth float32) {
	if obj.IsWorldCoordinateSystem { // World 坐标系：直接绘制屏幕
		screenX, screenY := p.tiledMapCfg.IsometricCT.W2S(obj.X, obj.Y, camX, camY)
		vector.StrokeRect(screen, screenX, screenY, obj.Width, obj.Height, strokeWidth, clr, false)
		return
	}

	// Tiled 等距像素坐标：转换为 Tile 坐标绘制菱形
	th := float32(p.tiledMapCfg.TileHPixel)
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
func (p *Map) drawLayer(screen *ebitenv2.Image, cam *commoncamera.Camera, layer *cfg.TiledLayer) {
	if 0 < len(layer.Data) {
		p.drawData(screen, cam, layer.Data, layer.Width, layer.Height)
	}
}

// drawData 绘制 tile 数据
func (p *Map) drawData(screen *ebitenv2.Image, cam *commoncamera.Camera, data []int, width, height int) {
	for i, gid := range data {
		if gid == 0 {
			continue
		}

		tileX := i % width
		tileY := i / width

		// 使用 IsometricCT 获取 tile 图像的屏幕位置
		screenX, screenY := p.tiledMapCfg.IsometricCT.TileImageScreenPos(tileX, tileY, cam.ViewportWX, cam.ViewportWY)

		// 从缓存获取 tile 信息
		tileInfo := p.getTileImageCached(gid)
		if tileInfo == nil {
			continue
		}

		// 修正 Y 坐标：tileset 高度大于 map tile 高度时，图像底部需对齐到 tile 基准位置
		screenY -= tileInfo.Height - p.tiledMapCfg.TileHPixel

		// 裁剪：跳过屏幕外的 tile
		if screenX < -tileInfo.Width || cfg.GCommon.ScreenMaxWidth < screenX ||
			screenY < -tileInfo.Height || cfg.GCommon.ScreenMaxHeight < screenY {
			continue
		}

		// 绘制
		op := &ebitenv2.DrawImageOptions{}
		op.GeoM.Translate(float64(screenX), float64(screenY))
		screen.DrawImage(tileInfo.Image, op)
	}
}

// getTileImage 根据 GID 获取 tile 图像和所属的 tileset
func (p *Map) getTileImage(gid int) (*ebitenv2.Image, *cfg.TiledTileset) {
	var tileset *cfg.TiledTileset
	for i := len(p.tiledMapCfg.Tiles) - 1; i >= 0; i-- {
		ts := p.tiledMapCfg.Tiles[i]
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

// GetTileSortInfoSlice 收集指定图层类型的可见瓦片信息（用于 Y-Sorting）
// 参数 result: 调用方传入的 slice，避免每帧分配
// 返回: 追加后的 slice
func (p *Map) GetTileSortInfoSlice(cam *commoncamera.Camera, layerType cfg.TiledLayerType, result []commonrenderable.IRenderable) []commonrenderable.IRenderable {
	for _, layer := range p.tiledMapCfg.Layers {
		if !layer.Visible || layer.LayerType != layerType {
			continue
		}
		result = p.collectLayerTiles(cam, layer, result)
	}
	return result
}

// collectLayerTiles 收集单个图层的可见瓦片
func (p *Map) collectLayerTiles(cam *commoncamera.Camera, layer *cfg.TiledLayer, result []commonrenderable.IRenderable) []commonrenderable.IRenderable {
	data := layer.Data
	width := layer.Width

	for i, gid := range data {
		if gid == 0 {
			continue
		}

		tileX := i % width
		tileY := i / width

		// 获取屏幕位置
		screenX, screenY := p.tiledMapCfg.IsometricCT.TileImageScreenPos(tileX, tileY, cam.ViewportWX, cam.ViewportWY)

		// 获取缓存的 tile 信息
		tileInfo := p.getTileImageCached(gid)
		if tileInfo == nil {
			continue
		}

		// 修正 Y 坐标
		screenY -= tileInfo.Height - p.tiledMapCfg.TileHPixel

		// 裁剪：跳过屏幕外的 tile
		if screenX < -tileInfo.Width || cfg.GCommon.ScreenMaxWidth < screenX ||
			screenY < -tileInfo.Height || cfg.GCommon.ScreenMaxHeight < screenY {
			continue
		}

		// 计算 World Y(瓦片底部中心点，用于排序)
		_, wy := p.tiledMapCfg.IsometricCT.T2W(float32(tileX)+0.5, float32(tileY)+0.5)

		// 处理垂直切片的 tile (大型建筑)
		if tileInfo.SlicedVertical != nil {
			for _, slice := range tileInfo.SlicedVertical.Slices {
				// 每个切片的屏幕位置
				sliceScreenX := screenX + slice.OffsetX
				sliceScreenY := screenY

				// 每个切片的 WY: 锚点在切片最底部像素
				// 切片底部相对于 tile 图像底部的偏移 = BottomPixelY - (Height - 1)
				// 偏移为负 = 切片底部在图像底部之上 = 更远离观察者 = WY 更低
				sliceWY := wy + float32(slice.BottomPixelY-tileInfo.Height+1)

				sortInfo := p.tileSortPool.Get()
				sortInfo.Image = slice.Image
				sortInfo.ScreenX = sliceScreenX
				sortInfo.ScreenY = sliceScreenY

				var offsetWY float32 = -float32(tileInfo.SlicedVertical.SlicePixel / 2) // todo menglc [微调] 为了使遮挡不穿模. 根据垂直切割宽度来设置偏移量. 避免角色在下方时,遮挡角色
				sortInfo.WY = sliceWY + offsetWY
				result = append(result, sortInfo)
			}
		} else if tileInfo.SlicedHorizontal != nil {
			// 下部分参与 Y-Sorting
			trunkScreenY := screenY + tileInfo.SlicedHorizontal.CanopyHeight // 树干起始Y = 整体Y + 树冠高度
			sortInfo := p.tileSortPool.Get()
			sortInfo.Image = tileInfo.SlicedHorizontal.TrunkImage
			sortInfo.ScreenX = screenX
			sortInfo.ScreenY = trunkScreenY
			sortInfo.WY = wy
			result = append(result, sortInfo)

			// 树冠部分收集到列表，稍后在 DrawOverhead 阶段绘制
			canopyInfo := p.umbrellaShapedCanopyDrawPool.Get()
			canopyInfo.Image = tileInfo.SlicedHorizontal.CanopyImage
			canopyInfo.ScreenX = screenX
			canopyInfo.ScreenY = screenY // 树冠起始Y = 整体Y
			p.umbrellaShapedCanopyDrawList = append(p.umbrellaShapedCanopyDrawList, canopyInfo)
		} else {
			// 普通 tile，整体参与 Y-Sorting
			sortInfo := p.tileSortPool.Get()
			sortInfo.Image = tileInfo.Image
			sortInfo.ScreenX = screenX
			sortInfo.ScreenY = screenY
			sortInfo.WY = wy
			result = append(result, sortInfo)
		}
	}
	return result
}

// RecycleTileSortInfoSlice 回收 TileSortInfo 对象到池中（绘制完成后调用）
func (p *Map) RecycleTileSortInfoSlice(slice []commonrenderable.IRenderable) {
	for _, item := range slice {
		if sortInfo, ok := item.(*TileSortInfo); ok {
			p.tileSortPool.Put(sortInfo)
		}
	}
}

package cfg

import (
	"saClient/src/common"
	commonct "saClient/src/common/coordinatetransform"
)

// TiledMap 属性标签
const TiledMapTag_BgmFilePath = "backgroundMusicFilePath" // 背景音乐文件路径-标签
const TiledMapTag_LogicalGrid = "logicalGrid"             // 逻辑网格-标签 (格式: "宽,高", 如 "8,8")

// TiledMap Tiled 地图资源
type TiledMap struct {
	ID                      common.AssetID  // 地图ID
	TileCountW              int             // tile 数量 - 宽度
	TileCountH              int             // tile 数量 - 高度
	TileWPixel              int             // tile 宽度 - 像素
	TileHPixel              int             // tile 高度 - 像素
	Layers                  []*TiledLayer   // 图层列表
	Tiles                   []*TiledTileset // tileset 列表
	BackgroundMusicFilePath string          // 背景音乐文件路径

	WPixel int // 地图宽度 - 像素
	HPixel int // 地图高度 - 像素

	Boundary    *TiledMapBoundary    // 地图边界 (World 坐标)
	TileBlocked *TiledMapTileBlocked // 地图阻挡对象

	IsometricCT *commonct.Isometric // 坐标转换器

	LogicalGrid *TiledMapLogicalGridMgr // 逻辑网格 (用于 A* 寻路), nil 表示未启用
}

func NewTiledMap() *TiledMap {
	tiledMap := &TiledMap{}
	tiledMap.Boundary = NewTiledMapBoundary(tiledMap)
	tiledMap.TileBlocked = NewTiledMapTileBlocked(tiledMap)
	return tiledMap
}

// IsBlocked 综合检测位置是否被阻挡 (地图边界 + Tile阻挡 + Object阻挡)
// wx, wy: World 坐标
// 返回: clampedWX, clampedWY 限制在边界内的坐标, blocked 是否被阻挡
func (p *TiledMap) IsBlocked(wx, wy float32) (clampedWX, clampedWY float32, blocked bool) {
	// 限制在地图边界内
	clampedWX, clampedWY, blocked = p.Boundary.ClampWithW(wx, wy)
	if blocked {
		return clampedWX, clampedWY, true
	}

	// 检查 Tile 阻挡
	tx, ty := p.IsometricCT.W2T(clampedWX, clampedWY)
	if p.TileBlocked.IsBlockedWithTF(tx, ty) {
		return clampedWX, clampedWY, true
	}

	// 检查 Object 阻挡
	if _, objBlocked := p.FindBlockedByObject(clampedWX, clampedWY); objBlocked {
		return clampedWX, clampedWY, true
	}

	return clampedWX, clampedWY, false
}

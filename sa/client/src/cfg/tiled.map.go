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

	Boundary TiledMapBoundary // 地图边界 (World 坐标)

	IsometricCT *commonct.Isometric // 坐标转换器

	LogicalGrid *TiledMapLogicalGridMgr // 逻辑网格 (用于 A* 寻路), nil 表示未启用

	TileBlocked [][]bool // [废弃][参见 README.md] 阻挡2维数组 [W][H] true 表示阻挡 (由多个图层合成)
}

// IsBlockedByTileWithT 检查是否被阻挡-图块-指定 Tile 坐标
// tileX, tileY: Tile 坐标 (整数)
// 返回: true 表示该 tile 阻挡角色，false 表示可通行
func (p *TiledMap) IsBlockedByTileWithT(tileX, tileY int) bool {
	// 越界视为阻挡
	if tileX < 0 || p.TileCountW <= tileX || tileY < 0 || p.TileCountH <= tileY {
		return true
	}
	return p.TileBlocked[tileX][tileY]
}

// IsBlockedByTileWithTF 检查是否被阻挡-图块-指定 Tile 坐标 (float32 版本)
// tx, ty: Tile 坐标 (浮点数)
// 返回: true 表示该 tile 阻挡角色，false 表示可通行
func (p *TiledMap) IsBlockedByTileWithTF(tx, ty float32) bool {
	// 四舍五入到最近的 Tile
	tileX := int(tx + 0.5)
	tileY := int(ty + 0.5)
	return p.IsBlockedByTileWithT(tileX, tileY)
}

// IsBlocked 综合检测位置是否被阻挡 (地图边界 + Tile阻挡 + Object阻挡)
// wx, wy: World 坐标
// 返回: clampedWX, clampedWY 限制在边界内的坐标, blocked 是否被阻挡
func (p *TiledMap) IsBlocked(wx, wy float32) (clampedWX, clampedWY float32, blocked bool) {
	// 限制在地图边界内
	clampedWX, clampedWY = p.ClampMapBoundaryWithW(wx, wy)

	// 检查 Tile 阻挡
	tx, ty := p.IsometricCT.W2T(clampedWX, clampedWY)
	if p.IsBlockedByTileWithTF(tx, ty) {
		return clampedWX, clampedWY, true
	}

	// 检查 Object 阻挡
	if _, objBlocked := p.FindBlockedByObject(clampedWX, clampedWY); objBlocked {
		return clampedWX, clampedWY, true
	}

	return clampedWX, clampedWY, false
}

// IsInsideDiamond 检查 World 坐标是否在菱形地图边界内
func (p *TiledMap) IsInsideDiamond(wx, wy float32) bool {
	tx, ty := p.IsometricCT.W2T(wx, wy)
	return tx >= -0.5 && tx <= float32(p.TileCountW)-0.5 &&
		ty >= -0.5 && ty <= float32(p.TileCountH)-0.5
}

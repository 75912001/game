package cfg

import (
	"saClient/src/common"
	commonct "saClient/src/common/coordinatetransform"
)

// TiledMap 属性标签
const TiledMapTag_BgmFilePath = "backgroundMusicFilePath" // 背景音乐文件路径-标签

// TiledMap Tiled 地图资源
type TiledMap struct {
	ID                      common.AssetID  // 地图ID
	Width                   int             // 地图宽度（tile 数）
	Height                  int             // 地图高度（tile 数）
	TileWidth               int             // tile 宽度（像素）
	TileHeight              int             // tile 高度（像素）
	Layers                  []*TiledLayer   // 图层列表
	Tilesets                []*TiledTileset // tileset 列表
	BackgroundMusicFilePath string          // 背景音乐文件路径

	PixelW             int                 // 地图像素宽度
	PixelH             int                 // 地图像素高度
	TopWX, TopWY       float32             // 地图边界-菱形四角坐标-顶 (World 坐标)
	RightWX, RightWY   float32             // 地图边界-菱形四角坐标-右 (World 坐标)
	BottomWX, BottomWY float32             // 地图边界-菱形四角坐标-底 (World 坐标)
	LeftWX, LeftWY     float32             // 地图边界-菱形四角坐标-左 (World 坐标)
	IsometricCT        *commonct.Isometric // 坐标转换器
	TileBlocked        [][]bool            // 阻挡2维数组 [Width][Height] true 表示阻挡 (由多个图层合成)
}

// 限制在菱形地图边界内-叉积版 (World 坐标)
func (p *TiledMap) ClampMapBoundary(wx, wy float32) (clampedWX float32, clampedWY float32) {
	// cross 叉积，< 0 表示点在边外侧 (屏幕坐标系Y向下)
	cross := func(ax, ay, bx, by, wx, py float32) float32 {
		return (bx-ax)*(py-ay) - (by-ay)*(wx-ax)
	}
	// projectToEdge 将点投影到线段上
	projectToEdge := func(ax, ay, bx, by, wx, py float32) (float32, float32) {
		abx, aby := bx-ax, by-ay
		apx, apy := wx-ax, py-ay
		t := (apx*abx + apy*aby) / (abx*abx + aby*aby)
		if t < 0 {
			t = 0
		} else if t > 1 {
			t = 1
		}
		return ax + t*abx, ay + t*aby
	}
	// Top-Right 边
	if cross(p.TopWX, p.TopWY, p.RightWX, p.RightWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.TopWX, p.TopWY, p.RightWX, p.RightWY, wx, wy)
	}
	// Right-Bottom 边
	if cross(p.RightWX, p.RightWY, p.BottomWX, p.BottomWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.RightWX, p.RightWY, p.BottomWX, p.BottomWY, wx, wy)
	}
	// Bottom-Left 边
	if cross(p.BottomWX, p.BottomWY, p.LeftWX, p.LeftWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.BottomWX, p.BottomWY, p.LeftWX, p.LeftWY, wx, wy)
	}
	// Left-Top 边
	if cross(p.LeftWX, p.LeftWY, p.TopWX, p.TopWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.LeftWX, p.LeftWY, p.TopWX, p.TopWY, wx, wy)
	}
	return wx, wy
}

// ClampMapBoundaryWithW 限制在菱形地图边界内-Tile版 (高效版)
func (p *TiledMap) ClampMapBoundaryWithW(wx, wy float32) (float32, float32) {
	tx, ty := p.IsometricCT.W2T(wx, wy)

	// 菱形边界对应 Tile 坐标 [-0.5, Size-0.5]
	minT := float32(-0.5)
	maxTX := float32(p.Width) - 0.5
	maxTY := float32(p.Height) - 0.5

	clamped := false
	if tx < minT {
		tx = minT
		clamped = true
	} else if tx > maxTX {
		tx = maxTX
		clamped = true
	}
	if ty < minT {
		ty = minT
		clamped = true
	} else if ty > maxTY {
		ty = maxTY
		clamped = true
	}

	if clamped {
		return p.IsometricCT.T2W(tx, ty)
	}
	return wx, wy
}

// IsBlockedByTileWithT 检查是否被阻挡-图块-指定 Tile 坐标
// tileX, tileY: Tile 坐标 (整数)
// 返回: true 表示该 tile 阻挡角色，false 表示可通行
func (p *TiledMap) IsBlockedByTileWithT(tileX, tileY int) bool {
	// 越界视为阻挡
	if tileX < 0 || p.Width <= tileX || tileY < 0 || p.Height <= tileY {
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

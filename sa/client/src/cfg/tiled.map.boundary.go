package cfg

import xutil "github.com/75912001/xlib/util"

// TiledMapBoundary Tiled 地图边界 (菱形四角坐标 World 坐标)
type TiledMapBoundary struct {
	tiledMap           *TiledMap
	TopWX, TopWY       float32 // 地图边界-菱形四角坐标-顶 (World 坐标)
	RightWX, RightWY   float32 // 地图边界-菱形四角坐标-右 (World 坐标)
	BottomWX, BottomWY float32 // 地图边界-菱形四角坐标-底 (World 坐标)
	LeftWX, LeftWY     float32 // 地图边界-菱形四角坐标-左 (World 坐标)
}

func NewTiledMapBoundary(tiledMap *TiledMap) *TiledMapBoundary {
	tiledMapBoundary := &TiledMapBoundary{
		tiledMap: tiledMap,
	}
	return tiledMapBoundary
}

// ClampWithW 限制在菱形地图边界内
func (p *TiledMapBoundary) ClampWithW(wx, wy float32) (clampedWX float32, clampedWY float32, blocked bool) {
	tx, ty := p.tiledMap.IsometricCT.W2T(wx, wy)

	// 菱形边界对应 Tile 坐标 [-0.5, Size-0.5]
	minT := float32(-0.5)
	maxTX := float32(p.tiledMap.TileCountW) - 0.5
	maxTY := float32(p.tiledMap.TileCountH) - 0.5

	// 直接进行边界检查和 clamp，避免重复检查
	if xutil.Float32Less(tx, minT) {
		tx = minT
		blocked = true
	} else if xutil.Float32Less(maxTX, tx) {
		tx = maxTX
		blocked = true
	}
	if xutil.Float32Less(ty, minT) {
		ty = minT
		blocked = true
	} else if xutil.Float32Less(maxTY, ty) {
		ty = maxTY
		blocked = true
	}

	if blocked {
		clampedWX, clampedWY = p.tiledMap.IsometricCT.T2W(tx, ty)
		return clampedWX, clampedWY, true
	}
	return wx, wy, false
}

// IsInside 检查 World 坐标是否在菱形地图边界内
func (p *TiledMapBoundary) IsInside(wx, wy float32) bool {
	tx, ty := p.tiledMap.IsometricCT.W2T(wx, wy)

	minT := float32(-0.5)
	maxTX := float32(p.tiledMap.TileCountW) - 0.5
	maxTY := float32(p.tiledMap.TileCountH) - 0.5

	// 使用 xutil.Float32Less 进行边界判断
	// tx >= minT 等价于 !Float32Less(tx, minT)
	// tx <= maxTX 等价于 !Float32Less(maxTX, tx)
	return !xutil.Float32Less(tx, minT) && !xutil.Float32Less(maxTX, tx) &&
		!xutil.Float32Less(ty, minT) && !xutil.Float32Less(maxTY, ty)
}

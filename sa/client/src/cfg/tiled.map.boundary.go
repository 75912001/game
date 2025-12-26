package cfg

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

	if tx < minT {
		tx = minT
		blocked = true
	} else if maxTX < tx {
		tx = maxTX
		blocked = true
	}
	if ty < minT {
		ty = minT
		blocked = true
	} else if maxTY < ty {
		ty = maxTY
		blocked = true
	}

	if blocked {
		clampedWX, clampedWY = p.tiledMap.IsometricCT.T2W(tx, ty)
		return
	}
	clampedWX = wx
	clampedWY = wy
	return
}

// IsInside 检查 World 坐标是否在菱形地图边界内
func (p *TiledMapBoundary) IsInside(wx, wy float32) bool {
	tx, ty := p.tiledMap.IsometricCT.W2T(wx, wy)
	return tx >= -0.5 && tx <= float32(p.tiledMap.TileCountW)-0.5 &&
		ty >= -0.5 && ty <= float32(p.tiledMap.TileCountH)-0.5
}

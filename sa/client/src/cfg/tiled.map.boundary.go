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
	if cross(p.Boundary.TopWX, p.Boundary.TopWY, p.Boundary.RightWX, p.Boundary.RightWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.Boundary.TopWX, p.Boundary.TopWY, p.Boundary.RightWX, p.Boundary.RightWY, wx, wy)
	}
	// Right-Bottom 边
	if cross(p.Boundary.RightWX, p.Boundary.RightWY, p.Boundary.BottomWX, p.Boundary.BottomWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.Boundary.RightWX, p.Boundary.RightWY, p.Boundary.BottomWX, p.Boundary.BottomWY, wx, wy)
	}
	// Bottom-Left 边
	if cross(p.Boundary.BottomWX, p.Boundary.BottomWY, p.Boundary.LeftWX, p.Boundary.LeftWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.Boundary.BottomWX, p.Boundary.BottomWY, p.Boundary.LeftWX, p.Boundary.LeftWY, wx, wy)
	}
	// Left-Top 边
	if cross(p.Boundary.LeftWX, p.Boundary.LeftWY, p.Boundary.TopWX, p.Boundary.TopWY, wx, wy) < 0 {
		wx, wy = projectToEdge(p.Boundary.LeftWX, p.Boundary.LeftWY, p.Boundary.TopWX, p.Boundary.TopWY, wx, wy)
	}
	return wx, wy
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

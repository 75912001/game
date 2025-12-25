package cfg

// TiledMapLogicalGrid 逻辑网格 (用于 A* 寻路)
type TiledMapLogicalGrid struct {
	CellW int      // 单元格宽度 (像素)
	CellH int      // 单元格高度 (像素)
	GridW int      // 网格列数
	GridH int      // 网格行数
	Grid  [][]bool // 网格数据 [x][y], true=阻挡, false=可通行
}

// Grid2W 逻辑网格坐标 → World 坐标 (单元格中心)
func (p *TiledMapLogicalGrid) Grid2W(gx, gy int) (wx, wy float32) {
	wx = float32(gx*p.CellW + p.CellW/2)
	wy = float32(gy*p.CellH + p.CellH/2)
	return
}

// IsBlocked 检查逻辑网格坐标是否阻挡
// 返回: true=阻挡, false=可通行
func (p *TiledMapLogicalGrid) IsBlocked(gx, gy int) bool {
	if gx < 0 || gx >= p.GridW || gy < 0 || gy >= p.GridH {
		return true // 越界视为阻挡
	}
	return p.Grid[gx][gy]
}

// IsWalkable 检查逻辑网格坐标是否可通行 (IsBlocked 的反向)
// 返回: true=可通行, false=阻挡
func (p *TiledMapLogicalGrid) IsWalkable(gx, gy int) bool {
	return !p.IsBlocked(gx, gy)
}

// build 生成逻辑网格 (在 Assemble 阶段调用)
func (p *TiledMapLogicalGrid) build(tiledMap *TiledMap) {
	// 计算网格尺寸
	p.GridW = (tiledMap.PixelW + p.CellW - 1) / p.CellW // 向上取整
	p.GridH = (tiledMap.PixelH + p.CellH - 1) / p.CellH

	// 初始化网格数组
	p.Grid = make([][]bool, p.GridW)
	for gx := 0; gx < p.GridW; gx++ {
		p.Grid[gx] = make([]bool, p.GridH)
	}

	// 遍历所有格子，检测阻挡
	for gx := 0; gx < p.GridW; gx++ {
		for gy := 0; gy < p.GridH; gy++ {
			// 计算格子中心的 World 坐标
			wx := float32(gx*p.CellW + p.CellW/2)
			wy := float32(gy*p.CellH + p.CellH/2)

			// 检查是否在菱形边界内
			if !tiledMap.IsInsideDiamond(wx, wy) {
				p.Grid[gx][gy] = true // 菱形外 = 阻挡
				continue
			}

			// 调用现有阻挡检测
			_, _, blocked := tiledMap.IsBlocked(wx, wy)
			p.Grid[gx][gy] = blocked
		}
	}
}

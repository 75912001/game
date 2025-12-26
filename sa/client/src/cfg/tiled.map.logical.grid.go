package cfg

// TiledMapLogicalGridMgr 逻辑网格 (用于 A* 寻路)
type TiledMapLogicalGridMgr struct {
	CellPixelW int                      // 单元格宽度 (像素)
	CellPixelH int                      // 单元格高度 (像素)
	GridCountW int                      // 网格列数
	GridCountH int                      // 网格行数
	Grids      [][]*TiledMapLogicalGrid // 网格数据 [gridX][gridY]
}

func NewTiledMapLogicalGridMgr(cellW, cellH int) *TiledMapLogicalGridMgr {
	return &TiledMapLogicalGridMgr{
		CellPixelW: cellW,
		CellPixelH: cellH,
	}
}

type TiledMapLogicalGrid struct {
	Blocked bool    // 是否阻挡 true=阻挡, false=可通行
	WX      float32 // 单元格中心 World 坐标
	WY      float32 // 单元格中心 World 坐标
}

func NewTiledMapLogicalGrid() *TiledMapLogicalGrid {
	return &TiledMapLogicalGrid{}
}

func (p *TiledMapLogicalGridMgr) GetGrid(gx, gy int) *TiledMapLogicalGrid {
	return p.Grids[gx][gy]
}

// Grid2W 逻辑网格坐标 → World 坐标 (单元格中心)
func (p *TiledMapLogicalGridMgr) Grid2W(gx, gy int) (wx, wy float32) {
	wx = float32(gx*p.CellPixelW + p.CellPixelW/2)
	wy = float32(gy*p.CellPixelH + p.CellPixelH/2)
	return
}

// W2Grid World 坐标 → 逻辑网格坐标
func (p *TiledMapLogicalGridMgr) W2Grid(wx, wy float32) (gx, gy int) {
	gx = int(wx) / p.CellPixelW
	gy = int(wy) / p.CellPixelH
	// 边界限制
	if gx < 0 {
		gx = 0
	} else if p.GridCountW <= gx {
		gx = p.GridCountW - 1
	}
	if gy < 0 {
		gy = 0
	} else if p.GridCountH <= gy {
		gy = p.GridCountH - 1
	}
	return
}

// IsBlocked 检查逻辑网格坐标是否阻挡
// 返回: true=阻挡, false=可通行
func (p *TiledMapLogicalGridMgr) IsBlocked(gx, gy int) bool {
	if gx < 0 || gx >= p.GridCountW || gy < 0 || gy >= p.GridCountH { // 越界视为阻挡
		return true
	}
	return p.Grids[gx][gy].Blocked
}

// build 生成逻辑网格 (在 Assemble 阶段调用)
func (p *TiledMapLogicalGridMgr) build(tiledMap *TiledMap) {
	// 计算网格尺寸
	p.GridCountW = (tiledMap.WPixel + p.CellPixelW - 1) / p.CellPixelW // 向上取整
	p.GridCountH = (tiledMap.HPixel + p.CellPixelH - 1) / p.CellPixelH

	// 初始化网格数组
	p.Grids = make([][]*TiledMapLogicalGrid, p.GridCountW)
	for gx := 0; gx < p.GridCountW; gx++ {
		p.Grids[gx] = make([]*TiledMapLogicalGrid, p.GridCountH)
	}

	// 采样点偏移 (相对于格子左上角)
	// 检测: 中心 + 四角 + 四边中点 = 9 个采样点
	halfW := float32(p.CellPixelW) / 2
	halfH := float32(p.CellPixelH) / 2
	cellW := float32(p.CellPixelW)
	cellH := float32(p.CellPixelH)
	sampleOffsets := [][2]float32{
		{halfW, halfH}, // 中心
		{0, 0},         // 左上
		{cellW, 0},     // 右上
		{0, cellH},     // 左下
		{cellW, cellH}, // 右下
		{halfW, 0},     // 上中
		{halfW, cellH}, // 下中
		{0, halfH},     // 左中
		{cellW, halfH}, // 右中
	}

	// 遍历所有格子，检测阻挡
	for gx := 0; gx < p.GridCountW; gx++ {
		for gy := 0; gy < p.GridCountH; gy++ {
			// 格子左上角 World 坐标
			baseX := float32(gx * p.CellPixelW)
			baseY := float32(gy * p.CellPixelH)

			blocked := false
			for _, offset := range sampleOffsets {
				wx := baseX + offset[0]
				wy := baseY + offset[1]

				// 检查是否被阻挡
				_, _, isBlocked := tiledMap.IsBlocked(wx, wy)
				if isBlocked {
					blocked = true
					break
				}
			}
			logicalGrid := NewTiledMapLogicalGrid()
			logicalGrid.Blocked = blocked
			logicalGrid.WX, logicalGrid.WY = p.Grid2W(gx, gy)
			p.Grids[gx][gy] = logicalGrid
		}
	}
}

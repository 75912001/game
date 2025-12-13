package coordinatetransform

// ============================================================================
// 地图尺寸
// ============================================================================

// GetDiamondCorners 获取整个菱形地图的四个**几何尖端** (World 坐标)
// 适用于: 设定摄像机移动范围 (Camera Bounds)
// 返回: Top(x,y), Right(x,y), Bottom(x,y), Left(x,y)
func (p *Isometric) GetDiamondCorners() (topX, topY, rightX, rightY, bottomX, bottomY, leftX, leftY float32) {
	w := float32(p.Width)
	h := float32(p.Height)

	// Top Tip: 位于 (0,0) 瓦片的顶部中心
	// ImageX of (0,0) is offsetX.
	// Tip X = offsetX + halfTW. Tip Y = 0.
	topX = p.offsetX + p.halfTW
	topY = 0

	// Right Tip: 位于 (w-1, 0) 瓦片的右侧尖端
	// 但如果是计算整个菱形包围盒，我们可以用数学公式直接推：
	rightX = (w + h) * p.halfTW // 地图总像素宽度
	rightY = w * p.halfTH

	// Bottom Tip
	bottomX = topX               // 菱形上下对称
	bottomY = (w + h) * p.halfTH // 地图总像素高度

	// Left Tip
	leftX = 0
	leftY = h * p.halfTH

	return
}

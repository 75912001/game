package coordinatetransform

// 坐标系统说明:
//
// 1. Tile 坐标 (Tiled coordinates)
//    - 以 tile 为单位的网格坐标
//    - 原点 (0,0) 在地图左上角
//    - 范围: [0, Width-1] x [0, Height-1]
//
// 2. World 坐标 (World coordinates)
//    - 以像素为单位的世界坐标
//    - 原点在地图左上角
//    - 基于 tile 菱形中心
//
// 3. Screen 坐标 (Screen coordinates)
//    - 以像素为单位的屏幕坐标
//    - 原点在屏幕左上角
//    - Screen = World - Camera

var GCT *CT

// CT 坐标转换器 (等距地图)
type CT struct {
	Width      int // 地图宽度(tile 数)
	Height     int // 地图高度(tile 数)
	TileWidth  int // 瓦片宽度(像素)
	TileHeight int // 瓦片高度(像素)
}

// NewCT 创建坐标转换器
func NewCT(width, height, tileWidth, tileHeight int) *CT {
	return &CT{
		Width:      width,
		Height:     height,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
	}
}

// ============================================================================
// Tile <-> World 转换
// ============================================================================

// TileToWorld Tile坐标 -> World坐标 (tile 菱形中心)
func (p *CT) TileToWorld(tileX, tileY float64) (worldX, worldY float64) {
	halfTW := float64(p.TileWidth) / 2
	halfTH := float64(p.TileHeight) / 2
	offsetX := float64(p.Height-1) * halfTW

	// 先计算 tile 图像左上角位置
	imageX := (tileX-tileY)*halfTW + offsetX
	imageY := (tileX + tileY) * halfTH

	// 加上菱形中心偏移得到 World 坐标
	worldX = imageX + halfTW
	worldY = imageY + halfTH
	return
}

// WorldToTile World坐标 -> Tile坐标
func (p *CT) WorldToTile(worldX, worldY float64) (tileX, tileY float64) {
	halfTW := float64(p.TileWidth) / 2
	halfTH := float64(p.TileHeight) / 2
	offsetX := float64(p.Height-1) * halfTW

	// 将 World 坐标转换为 tile 图像左上角坐标
	imageX := worldX - halfTW
	imageY := worldY - halfTH

	// 逆变换公式
	sx := imageX - offsetX
	tileX = (sx/halfTW + imageY/halfTH) / 2
	tileY = (imageY/halfTH - sx/halfTW) / 2
	return
}

// ============================================================================
// World <-> Screen 转换
// ============================================================================

// WorldToScreen World坐标 -> Screen坐标
func (p *CT) WorldToScreen(worldX, worldY, cameraX, cameraY float64) (screenX, screenY float64) {
	screenX = worldX - cameraX
	screenY = worldY - cameraY
	return
}

// ScreenToWorld Screen坐标 -> World坐标
func (p *CT) ScreenToWorld(screenX, screenY, cameraX, cameraY float64) (worldX, worldY float64) {
	worldX = screenX + cameraX
	worldY = screenY + cameraY
	return
}

// ============================================================================
// Tile 图像位置 (用于渲染)
// ============================================================================

// TileToImagePos 获取 Tile 图像左上角的 World 坐标 (用于渲染)
// todo menglc 可以预先计算出每个 tile 的位置，避免每次绘制都计算一次
func (p *CT) TileToImagePos(tileX, tileY int) (imageX, imageY int) {
	halfTW := p.TileWidth / 2
	halfTH := p.TileHeight / 2
	offsetX := (p.Height - 1) * halfTW

	imageX = (tileX-tileY)*halfTW + offsetX
	imageY = (tileX + tileY) * halfTH
	return
}

// TileToImageScreenPos 获取 Tile 图像左上角的 Screen 坐标 (用于渲染)
func (p *CT) TileToImageScreenPos(tileX, tileY, cameraX, cameraY int) (screenX, screenY int) {
	imageX, imageY := p.TileToImagePos(tileX, tileY)
	screenX = imageX - cameraX
	screenY = imageY - cameraY
	return
}

// ============================================================================
// 边界检测和限制
// ============================================================================

// IsInBounds 检查 Tile 坐标是否在地图边界内
func (p *CT) IsInBounds(tileX, tileY float64) bool {
	return tileX >= 0 && tileX < float64(p.Width) && tileY >= 0 && tileY < float64(p.Height)
}

// ClampTile 将 Tile 坐标限制在地图边界内
func (p *CT) ClampTile(tileX, tileY float64) (clampedTX, clampedTY float64) {
	maxTX := float64(p.Width) - 0.01
	maxTY := float64(p.Height) - 0.01

	clampedTX, clampedTY = tileX, tileY

	if clampedTX < 0 {
		clampedTX = 0
	} else if clampedTX > maxTX {
		clampedTX = maxTX
	}

	if clampedTY < 0 {
		clampedTY = 0
	} else if clampedTY > maxTY {
		clampedTY = maxTY
	}
	return
}

// ============================================================================
// 地图尺寸
// ============================================================================

// GetDiamondCorners 获取菱形边界的四个角点 (World 坐标)
// 返回: top(顶), right(右), bottom(底), left(左)
func (p *CT) GetDiamondCorners() (topX, topY, rightX, rightY, bottomX, bottomY, leftX, leftY float64) {
	halfTW := float64(p.TileWidth) / 2
	halfTH := float64(p.TileHeight) / 2
	w := float64(p.Width)
	h := float64(p.Height)

	topX = h * halfTW
	topY = 0
	rightX = (w + h) * halfTW
	rightY = w * halfTH
	bottomX = w * halfTW
	bottomY = (w + h) * halfTH
	leftX = 0
	leftY = h * halfTH
	return
}

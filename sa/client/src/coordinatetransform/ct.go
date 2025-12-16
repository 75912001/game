package coordinatetransform

// 坐标系统说明:
//
// 1. Tile 坐标 (Tiled coordinates)
//    - 以 tile 为单位的网格坐标
//    - 原点 (0,0) 在菱形地图左上角
//
// 2. World 坐标 (World coordinates)
//    - 以像素为单位的世界坐标
//    - 原点在矩形地图左上角
//
// 3. Screen 坐标 (Screen coordinates)
//    - 以像素为单位的屏幕坐标
//    - 原点在屏幕左上角
//
// 核心转换函数 (最小完备集):
//
//         Tile ←───── S2T ─────── Screen
//           │                        ↑
//          T2W                      W2S
//           ↓                        │
//         World ─────────────────────┘
//
// 派生转换:
//   - W2T = S2T(world, camera=0)
//   - T2S = W2S(T2W(tile))
//   - S2W = screen + camera

// Isometric 等距地图坐标转换器
type Isometric struct {
	Width      int // 地图宽度(tile 数)
	Height     int // 地图高度(tile 数)
	TileWidth  int // 瓦片宽度(像素)
	TileHeight int // 瓦片高度(像素)

	// 缓存值 (避免重复计算)
	halfTW  float32 // TileWidth / 2
	halfTH  float32 // TileHeight / 2
	offsetX float32 // (Height - 1) * halfTW

	// int 版本 (用于渲染)
	halfTWi  int // TileWidth / 2
	halfTHi  int // TileHeight / 2
	offsetXi int // (Height - 1) * halfTWi
}

// NewIsometric 创建坐标转换器
func NewIsometric(width, height, tileWidth, tileHeight int) *Isometric {
	halfTW := float32(tileWidth) / 2
	halfTH := float32(tileHeight) / 2
	halfTWi := tileWidth / 2
	halfTHi := tileHeight / 2

	return &Isometric{
		Width:      width,
		Height:     height,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,

		halfTW:  halfTW,
		halfTH:  halfTH,
		offsetX: float32(height-1) * halfTW,

		halfTWi:  halfTWi,
		halfTHi:  halfTHi,
		offsetXi: (height - 1) * halfTWi,
	}
}

// ============================================================================
// 核心转换函数 (最小完备集)
// ============================================================================

// T2W Tile坐标 → World坐标 (菱形中心)
func (p *Isometric) T2W(tileX, tileY float32) (worldX, worldY float32) {
	worldX = (tileX-tileY)*p.halfTW + p.offsetX + p.halfTW
	worldY = (tileX+tileY)*p.halfTH + p.halfTH
	return
}

// W2S World坐标 → Screen坐标
func (p *Isometric) W2S(worldX, worldY, cameraX, cameraY float32) (screenX, screenY float32) {
	return worldX - cameraX, worldY - cameraY
}

// S2T Screen坐标 → Tile坐标
func (p *Isometric) S2T(screenX, screenY, cameraX, cameraY float32) (tileX, tileY float32) {
	// Screen → World
	worldX := screenX + cameraX
	worldY := screenY + cameraY

	// World → Tile (逆变换)
	sx := worldX - p.halfTW - p.offsetX
	sy := worldY - p.halfTH
	tileX = (sx/p.halfTW + sy/p.halfTH) / 2
	tileY = (sy/p.halfTH - sx/p.halfTW) / 2
	return
}

// ============================================================================
// 派生转换函数
// ============================================================================

// W2T World坐标 → Tile坐标
func (p *Isometric) W2T(worldX, worldY float32) (tileX, tileY float32) {
	return p.S2T(worldX, worldY, 0, 0)
}

// T2S Tile坐标 → Screen坐标
func (p *Isometric) T2S(tileX, tileY, cameraX, cameraY float32) (screenX, screenY float32) {
	worldX, worldY := p.T2W(tileX, tileY)
	return p.W2S(worldX, worldY, cameraX, cameraY)
}

// S2W Screen坐标 → World坐标
func (p *Isometric) S2W(screenX, screenY, cameraX, cameraY float32) (worldX, worldY float32) {
	return screenX + cameraX, screenY + cameraY
}

// ============================================================================
// Tile 图像位置 (用于渲染)
// ============================================================================

// TileImagePos 获取 Tile 图像左上角的 World 坐标
func (p *Isometric) TileImagePos(tileX, tileY int) (imageX, imageY int) {
	imageX = (tileX-tileY)*p.halfTWi + p.offsetXi
	imageY = (tileX + tileY) * p.halfTHi
	return
}

// TileImageScreenPos 获取 Tile 图像左上角的 Screen 坐标
func (p *Isometric) TileImageScreenPos(tileX, tileY, cameraX, cameraY int) (screenX, screenY int) {
	imageX, imageY := p.TileImagePos(tileX, tileY)
	return imageX - cameraX, imageY - cameraY
}

package cfg

// TiledMapTileBlocked 地图-tile 阻挡集合
//
//	[废弃][参见 README.md]
type TiledMapTileBlocked struct {
	tiledMap *TiledMap
	Blocked  [][]bool // 阻挡2维数组 [W][H] true 表示阻挡 (由多个图层合成)
}

func NewTiledMapTileBlocked(tiledMap *TiledMap) *TiledMapTileBlocked {
	return &TiledMapTileBlocked{
		tiledMap: tiledMap,
	}
}

// IsBlockedWithT 检查是否被阻挡-图块-指定 Tile 坐标
// tileX, tileY: Tile 坐标 (整数)
// 返回: true 表示该 tile 阻挡角色，false 表示可通行
func (p *TiledMapTileBlocked) IsBlockedWithT(tx, ty int) bool {
	if tx < 0 || p.tiledMap.TileCountW <= tx || ty < 0 || p.tiledMap.TileCountH <= ty { // 越界视为阻挡
		return true
	}
	return p.Blocked[tx][ty]
}

// IsBlockedWithTF 检查是否被阻挡-图块-指定 Tile 坐标 (float32 版本)
// tx, ty: Tile 坐标 (浮点数)
// 返回: true 表示该 tile 阻挡角色，false 表示可通行
func (p *TiledMapTileBlocked) IsBlockedWithTF(tx, ty float32) bool {
	// 四舍五入到最近的 Tile
	tileX := int(tx + 0.5)
	tileY := int(ty + 0.5)
	return p.IsBlockedWithT(tileX, tileY)
}

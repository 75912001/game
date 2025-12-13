package cfg

import (
	"image"
	ct "saClient/src/coordinatetransform"

	"saClient/src/common"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

type TiledLayerType uint32 // Tiled图层类型

const (
	TiledLayerType_Unknown     TiledLayerType = 0 // 未知图层
	TiledLayerType_TileLayer   TiledLayerType = 1 // 瓦片图层
	TiledLayerType_ObjectLayer TiledLayerType = 2 // 对象图层
)

const TiledMapBgmFilePathTag = "backgroundMusicFilePath" // 背景音乐文件路径-标签
const TiledObjectCollisionTag = "collision"              // 碰撞对象-标签
const TiledObjectTargetPortal = "targetPortal"           // 传送目标-标签 传送ID

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

	PixelW      int           // 地图像素宽度
	PixelH      int           // 地图像素高度
	IsometricCT *ct.Isometric // 坐标转换器
}

// TiledLayer Tiled 图层
type TiledLayer struct {
	ID      int            // 图层ID
	Type    TiledLayerType // 图层类型
	Visible bool           // 是否可见
	Opacity float32        // 透明度
	X       int            // 图层X偏移
	Y       int            // 图层Y偏移
	Width   int            // 图层宽度
	Height  int            // 图层高度
	Data    []int          // tile 数据（用于有限地图）
	Objects []*TiledObject // 对象列表（用于 objectgroup）
}

// TiledObject Tiled 对象（用于碰撞等）
type TiledObject struct {
	ID           int             // 对象ID
	Type         string          // 对象类型
	X            float32         // 对象X坐标
	Y            float32         // 对象Y坐标
	Width        float32         // 对象宽度
	Height       float32         // 对象高度
	Rotation     float32         // 旋转角度
	Visible      bool            // 是否可见
	Polygon      []*TiledPoint   // 多边形顶点（相对坐标）
	Collision    bool            // 是否碰撞
	TargetPortal common.PortalID // 目标传送点ID
}

// TiledPoint Tiled 点
type TiledPoint struct {
	X float32
	Y float32
}

// TiledTileset Tiled-瓦片集
type TiledTileset struct {
	FirstGID    int             // 起始 GID
	Name        string          // tileset 名称
	TileWidth   int             // tile 宽度
	TileHeight  int             // tile 高度
	TileCount   int             // tile 总数
	Columns     int             // 列数
	Image       *ebitenv2.Image // tileset 图片
	ImageWidth  int             // 图片宽度
	ImageHeight int             // 图片高度
}

// GetTileImage 从 tileset 获取指定 tile 的子图
func (t *TiledTileset) GetTileImage(localID int) *ebitenv2.Image {
	if t.Image == nil || localID < 0 || localID >= t.TileCount {
		return nil
	}
	col := localID % t.Columns
	row := localID / t.Columns
	x := col * t.TileWidth
	y := row * t.TileHeight
	return t.Image.SubImage(image.Rect(x, y, x+t.TileWidth, y+t.TileHeight)).(*ebitenv2.Image)
}

// ============================================================================
// 碰撞检测
// ============================================================================

// ContainsWorldPoint 检查 World 坐标点是否在对象内部
// worldX, worldY: World 坐标 (像素)
// tileHeight: 用于将 Tiled 像素坐标转换为 Tile 坐标
// ct: 坐标转换器
func (obj *TiledObject) ContainsWorldPoint(worldX, worldY float32, tileHeight int, coordTransform *ct.Isometric) bool {
	if len(obj.Polygon) > 0 {
		// 多边形检测: 顶点坐标是相对于 obj.X, obj.Y 的 World 坐标
		return obj.containsPointInPolygon(worldX, worldY)
	}
	// 矩形检测: 需要转换坐标系
	return obj.containsPointInRect(worldX, worldY, tileHeight, coordTransform)
}

// containsPointInPolygon 检查点是否在多边形内部 (射线法)
// 多边形顶点是 World 坐标 (obj.X + point.X, obj.Y + point.Y)
func (obj *TiledObject) containsPointInPolygon(worldX, worldY float32) bool {
	n := len(obj.Polygon)
	if n < 3 {
		return false
	}

	inside := false
	j := n - 1

	for i := 0; i < n; i++ {
		// 多边形顶点的 World 坐标
		xi := obj.X + obj.Polygon[i].X
		yi := obj.Y + obj.Polygon[i].Y
		xj := obj.X + obj.Polygon[j].X
		yj := obj.Y + obj.Polygon[j].Y

		// 射线法判断点是否在多边形内
		if ((yi > worldY) != (yj > worldY)) &&
			(worldX < (xj-xi)*(worldY-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// containsPointInRect 检查点是否在矩形内部
// Tiled 等距地图中，矩形对象的坐标需要转换
func (obj *TiledObject) containsPointInRect(worldX, worldY float32, tileHeight int, coordTransform *ct.Isometric) bool {
	// 将 Tiled 像素坐标转换为 Tile 坐标 (与 drawCollision 相同的转换逻辑)
	th := float32(tileHeight)
	tileX := obj.X/th + 0.5
	tileY := obj.Y/th + 0.5
	tileW := obj.Width / th
	tileH := obj.Height / th

	// 将待检测点从 World 转换为 Tile 坐标
	pointTileX, pointTileY := coordTransform.W2T(worldX, worldY)

	// 检查点是否在 Tile 矩形内
	return pointTileX >= tileX && pointTileX <= tileX+tileW &&
		pointTileY >= tileY && pointTileY <= tileY+tileH
}

// CheckCollision 检查 World 坐标点是否与任何碰撞对象相交
// 返回: 是否碰撞
func (m *TiledMap) CheckCollision(worldX, worldY float32) bool {
	for _, layer := range m.Layers {
		if layer.Type != TiledLayerType_ObjectLayer {
			continue
		}
		for _, obj := range layer.Objects {
			if !obj.Collision {
				continue
			}
			if obj.ContainsWorldPoint(worldX, worldY, m.TileHeight, m.IsometricCT) {
				return true
			}
		}
	}
	return false
}

// GetPortalAt 获取 World 坐标点所在的传送门
// 返回: 传送门对象 (如果存在), nil (如果不存在)
func (m *TiledMap) GetPortalAt(worldX, worldY float32) *TiledObject {
	for _, layer := range m.Layers {
		if layer.Type != TiledLayerType_ObjectLayer {
			continue
		}
		for _, obj := range layer.Objects {
			if obj.TargetPortal == 0 {
				continue
			}
			if obj.ContainsWorldPoint(worldX, worldY, m.TileHeight, m.IsometricCT) {
				return obj
			}
		}
	}
	return nil
}

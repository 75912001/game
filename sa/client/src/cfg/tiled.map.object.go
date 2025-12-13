package cfg

import (
	ct "saClient/src/coordinatetransform"

	"saClient/src/common"
)

// TiledObject 属性标签
const TiledObjectTag_Collision = "collision"       // 碰撞对象-标签
const TiledObjectTag_TargetPortal = "targetPortal" // 传送目标-标签 传送ID

type TiledObjectType string

const (
	TiledObjectType_Portal TiledObjectType = "portal" // 传送点
)

// TiledObject Tiled 对象(用于碰撞等)
type TiledObject struct {
	ID           int             // 对象ID
	Type         TiledObjectType // 对象类型
	X            float32         // 对象X坐标
	Y            float32         // 对象Y坐标
	Width        float32         // 对象宽度
	Height       float32         // 对象高度
	Visible      bool            // 是否可见
	Collision    bool            // 是否碰撞
	TargetPortal common.PortalID // 目标传送点ID
}

// ============================================================================
// 碰撞检测
// ============================================================================

// ContainsWorldPoint 检查 World 坐标点是否在对象内部
// worldX, worldY: World 坐标 (像素)
// ct: 坐标转换器
func (obj *TiledObject) ContainsWorldPoint(worldX, worldY float32, coordTransform *ct.Isometric) bool {
	return obj.containsPointInRect(worldX, worldY, coordTransform)
}

// containsPointInRect 检查点是否在矩形内部
// Tiled 等距地图中，矩形对象的坐标需要转换
func (obj *TiledObject) containsPointInRect(worldX, worldY float32, coordTransform *ct.Isometric) bool {
	// 将 Tiled 像素坐标转换为 Tile 坐标 (与 drawCollision 相同的转换逻辑)
	th := float32(coordTransform.TileHeight)
	tileX := obj.X/th - 0.5
	tileY := obj.Y/th - 0.5
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
			if obj.ContainsWorldPoint(worldX, worldY, m.IsometricCT) {
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
			if obj.ContainsWorldPoint(worldX, worldY, m.IsometricCT) {
				return obj
			}
		}
	}
	return nil
}

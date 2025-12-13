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
	TiledObjectType_Portal        TiledObjectType = "portal"        // 传送点
	TiledObjectType_ArrivalPortal TiledObjectType = "arrivalPortal" // 到达传送点
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
	PortalCfg    *PortalPoint    // 传送点配置
}

// ============================================================================
// 碰撞检测
// ============================================================================

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

// FindCollisionObject 查找 World 坐标点所在的碰撞对象
func (m *TiledMap) FindCollisionObject(worldX, worldY float32) (*TiledObject, bool) {
	for _, layer := range m.Layers {
		if layer.Type != TiledLayerType_ObjectLayer { // 只检查对象图层
			continue
		}
		for _, obj := range layer.Objects { // 遍历对象
			if !obj.Collision { // 只检查碰撞对象
				continue
			}
			if obj.containsPointInRect(worldX, worldY, m.IsometricCT) {
				return obj, true
			}
		}
	}
	return nil, false
}

package cfg

import (
	"saClient/src/common"
	commonct "saClient/src/common/coordinatetransform"
)

// TiledObject 属性标签

const TiledObjectProperty_Blocked = "blocked"             // 对象-变量-阻挡
const TiledObjectProperty_ArrivalPortal = "arrivalPortal" // 对象-变量-到达传送点
const TiledObjectProperty_TargetPortal = "targetPortal"   // 对象-变量-传送目标 传送ID
const TiledObjectProperty_EnemyGroupID = "enemyGroupID"   // 对象-变量-刷怪点敌人组ID
const TiledObjectProperty_PatrolRadius = "patrolRadius"   // 对象-变量-刷怪点巡逻半径(world坐标单位)

type TiledObjectType string

// TiledObject Tiled 对象(用于碰撞等)
type TiledObject struct {
	ID                      int             // 对象ID
	Type                    TiledObjectType // 对象类型
	X                       float32         // 对象X坐标
	Y                       float32         // 对象Y坐标
	Width                   float32         // 对象宽度
	Height                  float32         // 对象高度
	Visible                 bool            // 是否可见
	Blocked                 bool            // 阻挡
	IsWorldCoordinateSystem bool            // true: X/Y/Width/Height 是 World 坐标系; false: Tiled 等距像素坐标
	// portal 传送点
	TargetPortal common.PortalID // 目标传送点ID
	PortalCfg    *PortalPoint    // 传送点配置
	// property 刷怪
	EnemyGroupID uint32      // 刷怪点-敌人组ID (0表示非刷怪点)
	EnemyGroup   *EnemyGroup // 运行时-敌人组配置引用
	PatrolRadius float32     // 刷怪点-巡逻半径 (世界坐标单位)
}

// ============================================================================
// 碰撞检测
// ============================================================================

// containsPointInRect 检查点是否在矩形内部
func (p *TiledObject) containsPointInRect(worldX, worldY float32, coordTransform *commonct.Isometric) bool {
	if p.IsWorldCoordinateSystem {
		// World 坐标系：直接在 World 坐标中检测
		return worldX >= p.X && worldX <= p.X+p.Width &&
			worldY >= p.Y && worldY <= p.Y+p.Height
	}

	// Tiled 等距像素坐标：转换为 Tile 坐标检测
	th := float32(coordTransform.TileHeight)
	tileX := p.X/th - 0.5
	tileY := p.Y/th - 0.5
	tileW := p.Width / th
	tileH := p.Height / th

	// 将待检测点从 World 转换为 Tile 坐标
	pointTileX, pointTileY := coordTransform.W2T(worldX, worldY)

	// 检查点是否在 Tile 矩形内
	return pointTileX >= tileX && pointTileX <= tileX+tileW &&
		pointTileY >= tileY && pointTileY <= tileY+tileH
}

// FindBlockedByObject 查找 World 坐标点所在的阻挡对象
func (p *TiledMap) FindBlockedByObject(worldX, worldY float32) (*TiledObject, bool) {
	for _, layer := range p.Layers {
		if layer.LayerType != TiledLayerType_Collision { // 非碰撞图层
			continue
		}
		for _, obj := range layer.Objects { // 遍历对象
			if !obj.Blocked { // 非 阻挡
				continue
			}
			if obj.containsPointInRect(worldX, worldY, p.IsometricCT) {
				return obj, true
			}
		}
	}
	return nil, false
}

// FindPortalByObject 查找 World 坐标点所在的传送点对象
func (p *TiledMap) FindPortalByObject(worldX, worldY float32) (*TiledObject, bool) {
	for _, layer := range p.Layers {
		if layer.LayerType != TiledLayerType_Collision { // 非碰撞图层
			continue
		}
		for _, obj := range layer.Objects { // 遍历对象
			if obj.PortalCfg == nil { // 非 传送点
				continue
			}
			if obj.containsPointInRect(worldX, worldY, p.IsometricCT) {
				return obj, true
			}
		}
	}
	return nil, false
}

// FindSpawnEnemyGroupPoints 查找所有刷怪点对象
func (p *TiledMap) FindSpawnEnemyGroupPoints() []*TiledObject {
	var result []*TiledObject
	for _, layer := range p.Layers {
		if layer.LayerType != TiledLayerType_Collision { // 非碰撞图层
			continue
		}
		for _, obj := range layer.Objects { // 遍历对象
			if obj.EnemyGroupID == 0 { // 非刷怪点
				continue
			}
			result = append(result, obj)
		}
	}
	return result
}

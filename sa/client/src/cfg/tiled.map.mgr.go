package cfg

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	commonct "saClient/src/common/coordinatetransform"
	"strconv"
	"strings"

	"saClient/src/common"
	"saClient/src/proto"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pkg/errors"
)

var GTiledMapMgr = newTiledMapMgr()

// TiledMapMgr Tiled 地图资源管理器
type TiledMapMgr struct {
	Maps *xmap.MapMgr[common.AssetID, *TiledMap]
}

func newTiledMapMgr() *TiledMapMgr {
	return &TiledMapMgr{
		Maps: xmap.NewMapMgr[common.AssetID, *TiledMap](),
	}
}

// Load 加载 Tiled 地图资源
func (p *TiledMapMgr) Load() error {
	tiledMapDir := common.AppResTiledMapDir
	// 检查目录是否存在
	if _, err := os.Stat(tiledMapDir); os.IsNotExist(err) {
		return errors.WithMessagef(err, "tiledMapDir %s does not exist", tiledMapDir)
	}

	mapFiles, err := os.ReadDir(tiledMapDir)
	if err != nil {
		return errors.WithMessagef(err, "读取 Tiled 地图目录失败: %v %v", tiledMapDir, xruntime.Location())
	}

	for _, mapFile := range mapFiles {
		if mapFile.IsDir() {
			continue
		}
		fileName := mapFile.Name()
		ext := filepath.Ext(fileName)
		if ext != ".tmx" { // 不是 .tmx 格式
			continue
		}
		// 解析文件名获取地图ID: map.${mapID}.tmx
		baseName := fileName[:len(fileName)-len(ext)]
		if len(baseName) < 5 || baseName[:4] != "map." { // 不符合命名规范 使用 map. 前缀
			continue
		}
		mapIDStr := baseName[4:]
		mapID64, err := strconv.ParseUint(mapIDStr, 10, 32)
		if err != nil {
			return errors.WithMessagef(err, "解析 Tiled 地图文件名 %s 为ID失败: %v", fileName, xruntime.Location())
		}
		mapID := common.AssetID(mapID64)
		// 检查ID范围
		if mapID < common.AssetID(proto.AssetIDRange_AssetIDRange_Map_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Map_End) < mapID {
			return fmt.Errorf("Tiled 地图ID超出范围: %d %v", mapID, xruntime.Location())
		}
		// 加载地图
		tiledMap, err := p.loadTiledMap(mapID, filepath.Join(tiledMapDir, fileName))
		if err != nil {
			return errors.WithMessagef(err, "加载 Tiled 地图失败: %v %v", fileName, xruntime.Location())
		}
		ok := p.Maps.AddIfNotExist(mapID, tiledMap)
		if !ok {
			return fmt.Errorf("添加 Tiled 地图资源失败,地图已存在: %v %v", mapID, xruntime.Location())
		}
	}
	return nil
}

// loadTiledMap 加载单个 Tiled 地图 (TMX 格式)
func (p *TiledMapMgr) loadTiledMap(mapID common.AssetID, tmxPath string) (*TiledMap, error) {
	data, err := os.ReadFile(tmxPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "读取 TMX 文件失败: %v %v", tmxPath, xruntime.Location())
	}

	var mapXML tmxMap
	if err := xml.Unmarshal(data, &mapXML); err != nil {
		return nil, errors.WithMessagef(err, "解析 TMX XML 失败: %v %v", tmxPath, xruntime.Location())
	}
	if mapXML.Orientation != "isometric" {
		return nil, fmt.Errorf("仅支持等距(isometric)地图,当前地图 %v 类型: %s %v", mapID, mapXML.Orientation, xruntime.Location())
	}
	if mapXML.RenderOrder != "right-down" {
		return nil, fmt.Errorf("仅支持右下(right-down)渲染顺序,当前地图 %v 渲染顺序: %s %v", mapID, mapXML.RenderOrder, xruntime.Location())
	}
	if mapXML.Infinite != nil && *mapXML.Infinite != 0 {
		return nil, fmt.Errorf("不支持无限(infinite)地图,当前地图 %v infinite 属性: %v %v", mapID, mapXML.Infinite, xruntime.Location())
	}
	if mapXML.TileWidth/2 != mapXML.TileHeight {
		return nil, fmt.Errorf("仅支持等宽高比 2:1 的等距(isometric)地图,当前地图 %v 瓦片宽高: %d x %d %v", mapID, mapXML.TileWidth, mapXML.TileHeight, xruntime.Location())
	}

	tiledMap := &TiledMap{
		ID:         mapID,
		Width:      mapXML.Width,
		Height:     mapXML.Height,
		TileWidth:  mapXML.TileWidth,
		TileHeight: mapXML.TileHeight,
	}

	if mapXML.Properties != nil { // 解析地图属性
		for _, prop := range mapXML.Properties.Properties {
			switch prop.Name {
			case TiledMapTag_BgmFilePath: // 背景音乐文件路径
				tiledMap.BackgroundMusicFilePath = prop.Value // todo menglc 需要检查合法性
			default: // 未知属性 忽略
			}
		}
	}

	tmxDir := filepath.Dir(tmxPath)
	for _, tsRef := range mapXML.Tilesets { // 加载 tilesets
		tileset, err := p.loadTileset(tmxDir, tsRef)
		if err != nil {
			return nil, errors.WithMessagef(err, "加载地图 %v 的 tileset 失败: %v %v", mapID, tmxPath, xruntime.Location())
		}
		tiledMap.Tilesets = append(tiledMap.Tilesets, tileset) // todo menglc 这里不同的 mapID, 相同的 tileset 会重复加载,可以考虑缓存优化
	}

	for _, layerXML := range mapXML.Layers { // 加载 tile layers
		layer := &TiledLayer{
			ID:      layerXML.ID,
			Visible: layerXML.Visible == nil || *layerXML.Visible != 0,
			Opacity: layerXML.Opacity,
			Width:   layerXML.Width,
			Height:  layerXML.Height,
		}
		if layer.Opacity == 0 {
			layer.Opacity = 1.0
		}
		if layerType, ok := TiledLayerClassNameMap.Find(layerXML.Class); !ok {
			return nil, fmt.Errorf("不支持的图层类型, 当前地图 %v 图层 %d 类型: %s %v", mapID, layer.ID, layerXML.Class, xruntime.Location())
		} else {
			layer.LayerType = layerType
		}

		if layerXML.Data.Encoding == "csv" || layerXML.Data.Encoding == "" { // 解析有限地图的 data
			layer.Data, err = parseCSVData(layerXML.Data.Content)
			if err != nil {
				return nil, errors.WithMessagef(err, "解析 %v 图层 %d 的 CSV 数据失败: %v %v", mapID, layer.ID, tmxPath, xruntime.Location())
			}
		}
		if layerXML.Properties != nil { // 解析属性
			for _, prop := range layerXML.Properties.Properties {
				switch prop.Name {
				case TileLayerProperty_Blocked:
					layer.BlockedLayer = prop.Value == "true"
				}
			}
		}
		tiledMap.Layers = append(tiledMap.Layers, layer)
	}

	for _, ogXML := range mapXML.ObjectGroups { // 加载 object groups
		layer := &TiledLayer{
			ID:      ogXML.ID,
			Visible: ogXML.Visible == nil || *ogXML.Visible != 0,
			Opacity: 1.0,
		}
		if layerType, ok := TiledLayerClassNameMap.Find(ogXML.Class); !ok {
			return nil, fmt.Errorf("不支持的图层类型, 当前地图 %v 图层 %d 类型: %s %v", mapID, layer.ID, ogXML.Class, xruntime.Location())
		} else {
			layer.LayerType = layerType
		}

		for _, objXML := range ogXML.Objects {
			if !xutil.Float32Equal(objXML.Rotation, 0) {
				return nil, fmt.Errorf("为了提高效能, 不支持旋转对象, 当前地图 %v 对象 %d 旋转角度: %f %v", mapID, objXML.ID, objXML.Rotation, xruntime.Location())
			}
			if objXML.Polygon != nil { // 多边形对象
				return nil, fmt.Errorf("为了提高效能, 不支持多边形对象, 当前地图 %v 对象 %d 类型: polygon %v", mapID, objXML.ID, xruntime.Location())
			}
			obj := &TiledObject{
				ID:      objXML.ID,
				Type:    objXML.Type,
				X:       objXML.X,
				Y:       objXML.Y,
				Width:   objXML.Width,
				Height:  objXML.Height,
				Visible: objXML.Visible == nil || *objXML.Visible != 0,
			}
			if objXML.Properties != nil {
				for _, prop := range objXML.Properties.Properties {
					switch prop.Name {
					case TiledObjectProperty_Blocked:
						obj.Blocked = prop.Value == "true"
					case TiledObjectProperty_ArrivalPortal:
					case TiledObjectProperty_TargetPortal:
						v, err := strconv.ParseUint(prop.Value, 10, 32)
						if err != nil {
							return nil, errors.WithMessagef(err, "解析对象 %d 的目标传送点失败: %v %v", obj.ID, tmxPath, xruntime.Location())
						}
						obj.TargetPortal = common.PortalID(v)
					case TiledObjectProperty_EnemyGroupID:
						v, err := strconv.ParseUint(prop.Value, 10, 32)
						if err != nil {
							return nil, errors.WithMessagef(err, "解析对象 %d 的敌人组ID失败: %v %v", obj.ID, tmxPath, xruntime.Location())
						}
						obj.EnemyGroupID = uint32(v)
					case TiledObjectProperty_PatrolRadius:
						v, err := strconv.ParseFloat(prop.Value, 32)
						if err != nil {
							return nil, errors.WithMessagef(err, "解析对象 %d 的巡逻半径失败: %v %v", obj.ID, tmxPath, xruntime.Location())
						}
						obj.PatrolRadius = float32(v)
					case TiledObjectProperty_RespawnSecond:
						v, err := strconv.ParseInt(prop.Value, 10, 32)
						if err != nil {
							return nil, errors.WithMessagef(err, "解析对象 %d 的刷新间隔失败: %v %v", obj.ID, tmxPath, xruntime.Location())
						}
						obj.RespawnSecond = int(v)
					case TiledObjectProperty_SpawnRadius:
						v, err := strconv.ParseFloat(prop.Value, 32)
						if err != nil {
							return nil, errors.WithMessagef(err, "解析对象 %d 的刷新半径失败: %v %v", obj.ID, tmxPath, xruntime.Location())
						}
						obj.SpawnRadius = float32(v)
					default:
						// 忽略未知属性
					}
				}
			}
			layer.Objects = append(layer.Objects, obj)
		}
		tiledMap.Layers = append(tiledMap.Layers, layer)
	}
	return tiledMap, nil
}

// loadTileset 加载 tileset (TSX 格式)
func (p *TiledMapMgr) loadTileset(tmxDir string, tsRef tmxTilesetRef) (*TiledTileset, error) {
	tsxPath := filepath.Join(tmxDir, tsRef.Source)
	tsxPath = filepath.Clean(tsxPath)

	data, err := os.ReadFile(tsxPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "读取 TSX 文件失败: %v %v", tsxPath, xruntime.Location())
	}

	var tsxXML tsxTileset
	if err := xml.Unmarshal(data, &tsxXML); err != nil {
		return nil, errors.WithMessagef(err, "解析 TSX XML 失败: %v %v", tsxPath, xruntime.Location())
	}

	tileset := &TiledTileset{
		FirstGID:    tsRef.FirstGID,
		TileWidth:   tsxXML.TileWidth,
		TileHeight:  tsxXML.TileHeight,
		TileCount:   tsxXML.TileCount,
		Columns:     tsxXML.Columns,
		ImageWidth:  tsxXML.Image.Width,
		ImageHeight: tsxXML.Image.Height,
		TileBlocked: make(map[int][]*TiledTileBlocked),
	}

	// 解析 tileset 属性
	if tsxXML.Properties != nil {
		for _, prop := range tsxXML.Properties.Properties {
			switch prop.Name {
			case TiledLayerProperty_HorizontalSlicePixel: // overheadRatio
				if prop.Type != "int" {
					return nil, fmt.Errorf("解析 tileset 属性失败, 必须是 int 类型:%v %v", prop.Name, xruntime.Location())
				}
				slicePixel, err := strconv.Atoi(prop.Value)
				if err != nil {
					return nil, errors.WithMessagef(err, "解析 tileset 属性失败:%v %v", prop.Name, xruntime.Location())
				}
				if slicePixel <= 0 || tileset.TileHeight <= slicePixel {
					return nil, fmt.Errorf("解析 tileset 属性失败, ratio 必须在 (0, 1) 范围内:%v %v", prop.Name, xruntime.Location())
				}
				tileset.HorizontalSlicePixel = slicePixel
			case TiledLayerProperty_VerticalSlicePixel: // verticalSlicePixel
				if prop.Type != "int" {
					return nil, fmt.Errorf("解析 tileset 属性失败, 必须是 int 类型:%v %v", prop.Name, xruntime.Location())
				}
				slicePixel, err := strconv.Atoi(prop.Value)
				if err != nil {
					return nil, errors.WithMessagef(err, "解析 tileset 属性失败:%v %v", prop.Name, xruntime.Location())
				}
				if slicePixel <= 0 {
					return nil, fmt.Errorf("解析 tileset 属性失败, verticalSlicePixel 必须大于 0:%v %v", prop.Name, xruntime.Location())
				}
				if tileset.TileWidth <= slicePixel {
					return nil, fmt.Errorf("解析 tileset 属性失败, verticalSlicePixel 必须小于 TileWidth:%v %v", prop.Name, xruntime.Location())
				}
				tileset.VerticalSlicePixel = slicePixel
			}
		}
	}

	// 解析 tile 阻挡体 (缓存到 TileBlocked)
	for _, tile := range tsxXML.Tiles {
		if tile.ObjectGroup == nil || len(tile.ObjectGroup.Objects) == 0 {
			continue
		}
		var blockedSlice []*TiledTileBlocked
		for _, obj := range tile.ObjectGroup.Objects {
			blockedSlice = append(blockedSlice, &TiledTileBlocked{
				X:      obj.X,
				Y:      obj.Y,
				Width:  obj.Width,
				Height: obj.Height,
			})
		}
		tileset.TileBlocked[tile.ID] = blockedSlice
	}

	// 加载图片
	tsxDir := filepath.Dir(tsxPath)
	imgPath := filepath.Join(tsxDir, tsxXML.Image.Source)
	imgPath = filepath.Clean(imgPath)

	img, _, err := ebitenv2ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "加载 tileset 图片失败: %v %v", imgPath, xruntime.Location())
	}
	tileset.Image = img

	return tileset, nil
}

// parseCSVData 解析 CSV 格式的 tile 数据
func parseCSVData(content string) ([]int, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, nil
	}

	var result []int
	parts := strings.Split(content, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, errors.WithMessagef(err, "解析 CSV tile 数据失败: %v %v", err, xruntime.Location())
		}
		result = append(result, val)
	}
	return result, nil
}

// Check 检查 Tiled 地图资源
func (p *TiledMapMgr) Check() error {
	var err error
	p.Maps.Foreach(
		func(mapID common.AssetID, tiledMap *TiledMap) (isContinue bool) {
			exist := GMapMgr.Maps.IsExist(mapID)
			if !exist { // 检查地图是否合法
				err = fmt.Errorf("Tiled 地图资源 %v 未定义 %v", mapID, xruntime.Location())
				return false
			}
			for _, layer := range tiledMap.Layers {
				switch layer.LayerType {
				case TiledLayerType_Collision: // 碰撞图层
					for _, obj := range layer.Objects {
						if obj.TargetPortal != 0 {
							exist = GPortalMgr.Points.IsExist(obj.TargetPortal)
							if !exist { // 检查目标传送点是否存在
								err = fmt.Errorf("Tiled 地图资源 %v 中对象 %d 的目标传送点 %v 不存在 %v", mapID, obj.ID, obj.TargetPortal, xruntime.Location())
								return false
							}
						}
						if obj.EnemyGroupID != 0 {
							exist = GEnemyGroupMgr.EnemyGroups.IsExist(obj.EnemyGroupID)
							if !exist { // 检查敌人组是否存在
								err = fmt.Errorf("Tiled 地图资源 %v 中对象 %d 的敌人组 %v 不存在 %v", mapID, obj.ID, obj.EnemyGroupID, xruntime.Location())
								return false
							}
						}
					}
				}
			}
			return true
		},
	)
	return err
}

// Assemble 装配 Tiled 地图资源
func (p *TiledMapMgr) Assemble() error {
	p.Maps.Foreach(
		func(mapID common.AssetID, tiledMap *TiledMap) (isContinue bool) {
			// 等距地图像素尺寸 (包含完整的菱形内容区域)
			// +TileHeight 是为了包含第一行 tile 上方和最后一行 tile 下方的菱形边缘
			tiledMap.PixelW = (tiledMap.Width + tiledMap.Height) * (tiledMap.TileWidth / 2)
			tiledMap.PixelH = (tiledMap.Width + tiledMap.Height) * (tiledMap.TileHeight / 2)
			tiledMap.IsometricCT = commonct.NewIsometric(tiledMap.Width, tiledMap.Height, tiledMap.TileWidth, tiledMap.TileHeight)
			tiledMap.TopWX, tiledMap.TopWY, tiledMap.RightWX, tiledMap.RightWY, tiledMap.BottomWX, tiledMap.BottomWY, tiledMap.LeftWX, tiledMap.LeftWY = tiledMap.IsometricCT.GetDiamondCorners()
			// 初始化碰撞网格
			gridSize := tiledMap.Width * tiledMap.Height
			// 创建二维布尔切片，按 [width][height]
			tiledMap.TileBlocked = make([][]bool, tiledMap.Width)
			for x := 0; x < tiledMap.Width; x++ {
				tiledMap.TileBlocked[x] = make([]bool, tiledMap.Height)
			}

			// 查找或创建 Collision 图层 (用于存储 tile 碰撞体)
			var collisionLayer *TiledLayer
			for _, layer := range tiledMap.Layers {
				if layer.LayerType == TiledLayerType_Collision {
					collisionLayer = layer
					break // 找到任一个碰撞层, 将其作为 tile 碰撞体的缓存层
				}
			}
			if collisionLayer == nil { // 不存在则创建
				collisionLayer = &TiledLayer{
					LayerType: TiledLayerType_Collision,
					Visible:   true,
					Opacity:   1.0,
				}
				tiledMap.Layers = append(tiledMap.Layers, collisionLayer)
			}

			objIDCounter := 1000000
			for _, obj := range collisionLayer.Objects { // 找出现有碰撞体的最大 ID
				if objIDCounter <= obj.ID {
					objIDCounter = obj.ID + 1
				}
			}
			// 生成 tile 碰撞体 (缓存到 Collision 图层)
			for _, layer := range tiledMap.Layers {
				if layer.LayerType == TiledLayerType_Collision || len(layer.Data) == 0 {
					continue
				}
				for i, gid := range layer.Data {
					if gid == 0 {
						continue
					}
					// 查找 tileset
					tileset := tiledMap.findTilesetByGID(gid)
					if tileset == nil || len(tileset.TileBlocked) == 0 { // 该 tile 没有阻挡体
						continue
					}
					// 计算 tile 本地 ID
					localID := gid - tileset.FirstGID
					blockedSlice, ok := tileset.TileBlocked[localID]
					if !ok || len(blockedSlice) == 0 {
						continue
					}
					// 计算 tile 在地图中的位置
					tileX := i % tiledMap.Width
					tileY := i / tiledMap.Width
					// 计算 tile 的 World 中心坐标
					tileWorldX, tileWorldY := tiledMap.IsometricCT.T2W(float32(tileX), float32(tileY))
					// tile 图像左上角的 World 坐标
					// 参考 TileImagePos: imageY = (tileX+tileY)*halfTH (比 T2W 少 halfTH)
					// 加上渲染调整: screenY -= tileInfo.Height - 地图TileHeight
					// 所以: imgTopLeftY = worldY - halfTH - (tileset.TileHeight - 地图TileHeight)
					halfTW := float32(tiledMap.TileWidth) / 2
					halfTH := float32(tiledMap.TileHeight) / 2
					imgTopLeftWorldX := tileWorldX - halfTW
					imgTopLeftWorldY := tileWorldY - halfTH - (float32(tileset.TileHeight) - float32(tiledMap.TileHeight))

					for _, coll := range blockedSlice {
						obj := &TiledObject{
							ID:                      objIDCounter,
							X:                       imgTopLeftWorldX + coll.X,
							Y:                       imgTopLeftWorldY + coll.Y,
							Width:                   coll.Width,
							Height:                  coll.Height,
							Visible:                 true,
							Blocked:                 true, // tile 碰撞体默认阻挡
							IsWorldCoordinateSystem: true, // 使用 World 坐标系
						}
						collisionLayer.Objects = append(collisionLayer.Objects, obj)
						objIDCounter++
					}
				}
			}

			for _, layer := range tiledMap.Layers {
				switch layer.LayerType {
				case TiledLayerType_Collision: // 碰撞图层
					for _, obj := range layer.Objects {
						if obj.TargetPortal != 0 { // 传送点对象
							obj.PortalCfg = GPortalMgr.Points.Get(obj.TargetPortal)
						}
						if obj.EnemyGroupID != 0 { // 刷怪点对象
							obj.EnemyGroup = GEnemyGroupMgr.EnemyGroups.Get(obj.EnemyGroupID)
						}
					}
				default: // 其他图层
					// 如果图层有阻挡属性，该图层所有非空 tile 都阻挡
					if layer.BlockedLayer && len(layer.Data) > 0 {
						for i, gid := range layer.Data {
							if gid != 0 && i < gridSize {
								x := i % tiledMap.Width
								y := i / tiledMap.Width
								tiledMap.TileBlocked[x][y] = true
							}
						}
					}
				}
			}
			return true
		},
	)
	return nil
}

// findTilesetByGID 根据 GID 查找对应的 tileset
func (p *TiledMap) findTilesetByGID(gid int) *TiledTileset {
	var result *TiledTileset
	for _, ts := range p.Tilesets {
		if ts.FirstGID <= gid {
			if result == nil || ts.FirstGID > result.FirstGID {
				result = ts
			}
		}
	}
	return result
}

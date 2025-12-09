package res

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pkg/errors"
	"saClient/src/common"
	"saClient/src/proto"
)

var GMapMgr = newMapMgr()

// MapMgr 地图资源管理器
type MapMgr struct {
	Maps *xmap.MapMgr[common.AssetID, *Map]
}

func newMapMgr() *MapMgr {
	return &MapMgr{
		Maps: xmap.NewMapMgr[common.AssetID, *Map](),
	}
}

// Load 加载地图资源
func (p *MapMgr) Load() error {
	// 读取地图资源目录下所有文件
	mapFiles, err := os.ReadDir(common.AppResMapDir)
	if err != nil {
		return errors.WithMessagef(err, "读取地图资源目录失败: %v %v", common.AppResMapDir, xruntime.Location())
	}
	for _, mapFile := range mapFiles {
		if mapFile.IsDir() { // 跳过目录
			continue
		}
		fileName := mapFile.Name()
		// 匹配 map.${mapID}.jpeg 或 map.${mapID}.png 格式的文件
		ext := filepath.Ext(fileName)
		if ext != ".jpeg" && ext != ".jpg" && ext != ".png" {
			continue
		}
		// 解析文件名获取地图ID: map.${mapID}.ext
		baseName := fileName[:len(fileName)-len(ext)]
		// 期望格式: map.2000001
		if len(baseName) < 5 || baseName[:4] != "map." {
			continue
		}
		mapIDStr := baseName[4:]
		mapID64, err := strconv.ParseUint(mapIDStr, 10, 32)
		if err != nil {
			return errors.WithMessagef(err, "解析地图文件名 %s 为ID失败: %v", fileName, xruntime.Location())
		}
		mapID := common.AssetID(mapID64)
		// 检查ID范围
		if mapID < common.AssetID(proto.AssetIDRange_AssetIDRange_Map_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Map_End) < mapID {
			return fmt.Errorf("地图ID超出范围: %d %v", mapID, xruntime.Location())
		}
		// 加载地图图片
		imageFilePath := filepath.Join(common.AppResMapDir, fileName)
		img, _, err := ebitenv2ebitenutil.NewImageFromFile(imageFilePath)
		if err != nil {
			return errors.WithMessagef(err, "加载地图图片失败: %v %v", imageFilePath, xruntime.Location())
		}
		// 创建地图资源
		mapRes := &Map{
			ID:    mapID,
			Image: img,
		}
		ok := p.Maps.AddIfNotExist(mapID, mapRes)
		if !ok {
			return fmt.Errorf("添加地图资源失败,地图已存在: %v %v", mapID, xruntime.Location())
		}
	}
	return nil
}

// Check 检查地图资源
func (p *MapMgr) Check() error {
	return nil
}

// Assemble 装配地图资源
func (p *MapMgr) Assemble() error {
	return nil
}

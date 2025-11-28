package config

import (
	"airplaneClient/src/common"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var GPlaneMgr = newPlaneMgr()

// PlaneMgr 飞机配置管理器
type PlaneMgr struct {
	PlanesData map[common.PlaneKey]*Plane // key: common.PlaneKey
}

func newPlaneMgr() *PlaneMgr {
	return &PlaneMgr{
		PlanesData: make(map[common.PlaneKey]*Plane),
	}
}

// Load 加载飞机配置文件
func (p *PlaneMgr) Load() error {
	// 读取目录下所有文件
	files, err := os.ReadDir(common.AppConfDir)
	if err != nil {
		return fmt.Errorf("读取飞机配置目录失败: %v", err)
	}
	// 匹配 plane.001.001.json 格式的文件
	pattern := regexp.MustCompile(`^plane\.(\d{3})\.(\d{3})\.json$`)
	for _, file := range files {
		if file.IsDir() { // 跳过目录
			continue
		}
		fileName := file.Name()
		matches := pattern.FindStringSubmatch(fileName)
		if matches == nil { // 不匹配格式，跳过
			continue
		}
		// 加载配置文件
		filePath := filepath.Join(common.AppConfDir, fileName)
		config, err := loadConfig(filePath)
		if err != nil {
			return fmt.Errorf("加载配置文件 %s 失败: %v", fileName, err)
		}
		planeKey := common.PlaneKey{}
		{
			id64, err := strconv.ParseUint(matches[1], 10, 32)
			if err != nil {
				return fmt.Errorf("解析文件名中的 id 失败: %v", err)
			}
			level64, err := strconv.ParseUint(matches[2], 10, 32)
			if err != nil {
				return fmt.Errorf("解析文件名中的 level 失败: %v", err)
			}
			planeKey.ID = uint32(id64)
			planeKey.Level = uint32(level64)
		}
		err = loadPlane(planeKey, config)
		if err != nil {
			return fmt.Errorf("加载飞机 %s 失败: %v", fileName, err)
		}
	}
	return nil
}

func (p *PlaneMgr) GetPlane(planeKey common.PlaneKey) *Plane {
	return p.PlanesData[planeKey]
}

package config

import (
	"airplaneClient/src/common"
	"airplaneClient/src/resources"
	"encoding/json"
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ConfigPlane TexturePacker 导出的飞机配置
type ConfigPlane struct {
	Frames map[string]*ConfigPlanePart `json:"frames"` // key: 部件文件名称 "plane.001.001.body.000" val: 部件信息
}

// ConfigPlanePart 飞机-部件-信息
type ConfigPlanePart struct {
	Frame *Rect `json:"frame"` // 在大图中的位置和尺寸
}

// Rect 矩形
type Rect struct {
	X      int `json:"x"` // x 坐标
	Y      int `json:"y"` // y 坐标
	Width  int `json:"w"` // 宽度
	Height int `json:"h"` // 高度
}

// 加载单个配置文件
func loadConfig(filePath string) (*ConfigPlane, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config ConfigPlane
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// 加载单个飞机配置
func loadPlane(planeKey common.PlaneKey, configPlane *ConfigPlane) (err error) {
	var partMap [common.PlanePartTypeMax]map[uint32]*ConfigPlanePart // 分为不同部件的帧映射 key: 部件索引 [partIdx]
	for i := range common.PlanePartTypeMax {
		partMap[i] = make(map[uint32]*ConfigPlanePart)
	}
	for partFileName, partData := range configPlane.Frames {
		parts := strings.Split(partFileName, ".")
		if len(parts) != 5 {
			return fmt.Errorf("帧名称格式错误: %s", partFileName)
		}
		partName := parts[3]
		partIdx, err := strconv.ParseUint(parts[4], 10, 32)
		if err != nil {
			return fmt.Errorf("解析文件名中的部件帧索引失败: %v", err)
		}
		switch partName {
		case "nose":
			_, ok := partMap[common.PlanePartTypeNose][uint32(partIdx)]
			if ok { // partIdx 重复
				return fmt.Errorf("帧名称重复: %s", partFileName)
			}
			partMap[common.PlanePartTypeNose][uint32(partIdx)] = partData
		case "body":
			_, ok := partMap[common.PlanePartTypeBody][uint32(partIdx)]
			if ok { // partIdx 重复
				return fmt.Errorf("帧名称重复: %s", partFileName)
			}
			partMap[common.PlanePartTypeBody][uint32(partIdx)] = partData
		case "leftwing":
			_, ok := partMap[common.PlanePartTypeLeftWing][uint32(partIdx)]
			if ok { // partIdx 重复
				return fmt.Errorf("帧名称重复: %s", partFileName)
			}
			partMap[common.PlanePartTypeLeftWing][uint32(partIdx)] = partData
		default:
			return fmt.Errorf("帧名称包含未知部件: %s", partFileName)
		}
	}
	for _, part := range partMap {
		var maxPartIdx uint32 // 计算最大部件索引
		for partIdx, _ := range part {
			if maxPartIdx < partIdx {
				maxPartIdx = partIdx
			}
		}
		if len(part) != int(maxPartIdx)+1 { // 部件帧索引-不连续
			return fmt.Errorf("配置文件中的部件帧索引不连续")
		}
	}
	var configPlanePart [common.PlanePartTypeMax][]*ConfigPlanePart // 分为不同部件的帧切片，按索引有序 [0...]
	for part, partData := range partMap {
		if len(configPlanePart[part]) == 0 {
			configPlanePart[part] = make([]*ConfigPlanePart, len(partData), len(partData))
		}
		for partIdx, frame := range partData {
			configPlanePart[part][partIdx] = frame
		}
	}

	partsFrames, err := cutFrame(planeKey.ID, planeKey.Level, configPlanePart)
	if err != nil {
		return fmt.Errorf("裁剪飞机帧失败: %v", err)
	}

	GPlaneMgr.Plane[planeKey] = &Plane{
		PlaneKey:    planeKey,
		PartsFrames: partsFrames,
	}
	return nil
}

// 裁剪帧
func cutFrame(id uint32, level uint32, configPlanePartData [common.PlanePartTypeMax][]*ConfigPlanePart) (
	partsFrames [common.PlanePartTypeMax][]*ebitenv2.Image, err error) {
	// 构建图片路径
	var planeName string
	planeName = resources.GenPlaneName(id, level)
	planePath := filepath.Join(common.AppResourcesDir, "planes", planeName)

	// 加载图片
	img, _, err := ebitenv2ebitenutil.NewImageFromFile(planePath)
	if err != nil {
		log.Printf("加载飞机图片失败:%v %v", planePath, err)
		return partsFrames, err
	}
	// 计算并切分每个部件的帧
	for partType, partFrames := range configPlanePartData {
		if len(partFrames) == 0 { // 没有该部件，跳过
			continue
		}
		for _, configPlanePart := range partFrames {
			subImg := img.SubImage(image.Rect(
				configPlanePart.Frame.X,
				configPlanePart.Frame.Y,
				configPlanePart.Frame.Width,
				configPlanePart.Frame.Height,
			)).(*ebitenv2.Image)
			partsFrames[partType] = append(partsFrames[partType], subImg)
		}
	}
	return partsFrames, nil
}

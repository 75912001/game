package res

import (
	"encoding/json"
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"os"
	"path/filepath"
	"saClient/src/common"
	"strconv"
	"strings"
)

// CfgRole TexturePacker 导出的 角色 配置
type CfgRole struct {
	Frames map[string]*CfgRoleSprite `json:"frames"` // key: 文件名称 "role.${roleID}.${action}.${arg}" val: 信息
}

// CfgRoleSprite 角色-精灵-信息
type CfgRoleSprite struct {
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
func loadConfig(filePath string) (*CfgRole, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config CfgRole
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// 加载单个角色配置
func loadRole(roleID common.AssetID, cfgRole *CfgRole) (err error) {
	for roleFileName, roleData := range cfgRole.Frames {
		parts := strings.Split(roleFileName, ".")
		if len(parts) != 5 {
			return fmt.Errorf("帧名称格式错误: %s", roleFileName)
		}
		action := parts[2] // 动作
		switch action {
		case "move": // 动作-移动
			_ = roleData
			roleDirection := parts[3] // 方向
			switch roleDirection {
			case "up": // 上
			case "down": // 下
			case "downleft": // 左下
			case "upleft": // 左上
			case "left": // 左
			//case "downright": // 右下 不用加载, 通过翻转获得
			//case "upright": // 右上 不用加载, 通过翻转获得
			//case "right": // 右 不用加载, 通过翻转获得
			default:
				return fmt.Errorf("帧名称包含未知方向: %s", roleDirection)
			}
			partsFrames, err := cutFrame(planeKey.ID, planeKey.Level, configPlanePart)
			if err != nil {
				return fmt.Errorf("裁剪飞机帧失败: %v", err)
			}

			GRoleMgr.Roles[roleID] = &Role{
				ID: roleID,
				Move: &RoleMove{
					Frames:    partsFrames,
					FrameInfo: configPlanePart,
				},
			}
		default:
			return fmt.Errorf("帧名称包含未知行为: %s", action)
		}
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

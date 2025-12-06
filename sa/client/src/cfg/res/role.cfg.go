package res

import (
	"encoding/json"
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"os"
	"saClient/src/common"
	"saClient/src/proto"
	"strconv"
	"strings"
)

// RoleJson TexturePacker 导出的 角色 配置
type RoleJson struct {
	Frames map[string]*RoleImageSprite `json:"frames"` // key: 文件名称 "role.${roleID}.${action}.${arg}.${frameIndex}" val: 信息
}

// RoleImageSprite 角色-精灵
type RoleImageSprite struct {
	Frame *Rect `json:"frame"` // 在大图中的位置和尺寸
}

// 加载单个配置文件
func loadRoleJson(filePath string) (*RoleJson, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var roleJson RoleJson
	err = json.Unmarshal(data, &roleJson)
	if err != nil {
		return nil, err
	}

	return &roleJson, nil
}

// loadRoleImage 加载角色单个配置
func loadRoleImage(roleID common.AssetID, imageFilePath string, roleJson *RoleJson) (err error) {
	role, exists := GRoleMgr.Roles[roleID]
	if !exists { // 不存在则创建
		role = &Role{
			ID:   roleID,
			Move: NewRoleMove(),
		}
		GRoleMgr.Roles[roleID] = role
	}

	// 按方向分组帧信息
	directionFrames := make(map[proto.RoleDirection][]*roleFrameData)

	for frameFileName, spriteData := range roleJson.Frames {
		// 解析帧名称: role.${roleID}.${action}.${arg}.${frameIndex}
		parts := strings.Split(frameFileName, ".")
		if len(parts) < 5 {
			return fmt.Errorf("帧名称格式错误: %s", frameFileName)
		}
		action := parts[2]     // 动作
		direction := parts[3]  // 方向
		frameIndex := parts[4] // 帧索引
		roleAction := GetRoleActionByName(action)
		if roleAction == proto.RoleAction_RoleAction_Unknow {
			return fmt.Errorf("帧名称包含未知动作: %s", action)
		}
		switch roleAction {
		case proto.RoleAction_RoleAction_Move:
			roleDirection := GetRoleDirectionByName(direction)
			if roleDirection == proto.RoleDirection_RoleDirection_Unknow {
				return fmt.Errorf("帧名称包含未知方向: %s", direction)
			}
			idx, err := strconv.Atoi(frameIndex)
			if err != nil {
				return fmt.Errorf("解析帧索引失败: %s, %v", frameFileName, err)
			}
			directionFrames[roleDirection] = append(directionFrames[roleDirection], &roleFrameData{
				index:      idx,
				roleAction: roleAction,
				spriteInfo: spriteData,
			})
		default:
			return fmt.Errorf("帧名称包含未实现动作: %s", action)
		}
	}
	for _, frames := range directionFrames { // 按帧索引排序
		sortFrames(frames)
	}
	// 如果有帧数据，需要加载图片并裁剪
	if len(directionFrames) > 0 {
		err = roleCutFrames(role, imageFilePath, directionFrames)
		if err != nil {
			return fmt.Errorf("裁剪角色帧失败: %v", err)
		}
	}
	return nil
}

// roleFrameData 角色-帧数据(用于排序)
type roleFrameData struct {
	index      int
	roleAction proto.RoleAction
	spriteInfo *RoleImageSprite
}

// roleCutFrames 角色-裁剪-帧
func roleCutFrames(role *Role, imageFilePath string, directionFrames map[proto.RoleDirection][]*roleFrameData) error {
	// 加载图片
	img, _, err := ebitenv2ebitenutil.NewImageFromFile(imageFilePath)
	if err != nil {
		log.Printf("加载角色图片失败: %v %v", imageFilePath, err)
		return err
	}
	// 裁剪每个方向的帧
	for dir, frames := range directionFrames {
		// 裁剪并存储帧
		for _, fd := range frames {
			subImg := img.SubImage(image.Rect(
				fd.spriteInfo.Frame.X,
				fd.spriteInfo.Frame.Y,
				fd.spriteInfo.Frame.X+fd.spriteInfo.Frame.Width,
				fd.spriteInfo.Frame.Y+fd.spriteInfo.Frame.Height,
			)).(*ebitenv2.Image)
			switch fd.roleAction {
			case proto.RoleAction_RoleAction_Move:
				role.Move.Frames[dir] = subImg
				role.Move.FrameInfo[dir] = append(role.Move.FrameInfo[dir], fd.spriteInfo)
			default:
				return fmt.Errorf("未实现的角色动作裁剪: %v", fd.roleAction)
			}
		}
	}

	return nil
}

// sortFrames 按帧索引排序
func sortFrames(frames []*roleFrameData) {
	for i := 0; i < len(frames)-1; i++ {
		for j := i + 1; j < len(frames); j++ {
			if frames[i].index > frames[j].index {
				frames[i], frames[j] = frames[j], frames[i]
			}
		}
	}
}

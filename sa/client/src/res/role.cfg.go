package res

import (
	"encoding/json"
	"fmt"
	xruntime "github.com/75912001/xlib/runtime"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pkg/errors"
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
	Frames map[string]*RoleImageSprite `json:"frames"` // key: 文件名称 "role.${roleID}.${arg}.${data}.${frameIndex}" val: 信息
}

// RoleImageSprite 角色-精灵
type RoleImageSprite struct {
	Frame *Rect `json:"frame"` // 在大图中的位置和尺寸
}

// 加载单个配置文件
func loadRoleJson(filePath string) (*RoleJson, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "读取角色配置文件失败: %s %v", filePath, xruntime.Location())
	}
	var roleJson RoleJson
	err = json.Unmarshal(data, &roleJson)
	if err != nil {
		return nil, errors.WithMessagef(err, "解析角色配置文件失败: %s %v", filePath, xruntime.Location())
	}
	return &roleJson, nil
}

// loadRoleImage 加载角色单个配置
func loadRoleImage(roleID common.AssetID, imageFilePath string, roleJson *RoleJson) (err error) {
	role, exists := GRoleMgr.Roles.Find(roleID)
	if !exists { // 不存在则创建
		role = &Role{
			ID:   roleID,
			Move: NewRoleMove(),
		}
		ok := GRoleMgr.Roles.AddIfNotExist(roleID, role)
		if !ok {
			return fmt.Errorf("添加角色失败,角色已存在: %v %v", roleID, xruntime.Location())
		}
	}
	// 按方向分组帧信息
	orientationFrames := make(map[proto.AssetOrientation][]*roleFrameData)
	for frameFileName, spriteData := range roleJson.Frames { // 解析帧名称: role.${roleID}.${arg}.${data}.${frameIndex}
		args := strings.Split(frameFileName, ".")
		if len(args) < 5 {
			return fmt.Errorf("帧名称格式错误: %v %s %v", imageFilePath, frameFileName, xruntime.Location())
		}
		action := args[2]      // 动作
		orientation := args[3] // 方向
		frameIndex := args[4]  // 帧索引
		roleAction := GetRoleActionByName(action)
		if roleAction == proto.RoleAction_RoleAction_Unknow {
			return fmt.Errorf("帧名称包含未知动作: %v %s %v", imageFilePath, action, xruntime.Location())
		}
		switch roleAction {
		case proto.RoleAction_RoleAction_Move:
			roleOrientation := GetAssetOrientationByName(orientation)
			if roleOrientation == proto.AssetOrientation_AssetOrientation_Unknow {
				return fmt.Errorf("帧名称包含未知方向: %v %s %v", imageFilePath, orientation, xruntime.Location())
			}
			idx, err := strconv.Atoi(frameIndex)
			if err != nil {
				return errors.WithMessagef(err, "解析帧索引失败: %v %s %v", imageFilePath, frameFileName, xruntime.Location())
			}
			orientationFrames[roleOrientation] = append(orientationFrames[roleOrientation], &roleFrameData{
				index:      idx,
				roleAction: roleAction,
				spriteInfo: spriteData,
			})
		default:
			return fmt.Errorf("帧名称包含未实现动作: %v %s %v", imageFilePath, action, xruntime.Location())
		}
	}
	for _, frames := range orientationFrames { // 按帧索引排序
		sortFrames(frames)
	}
	// 如果有帧数据，需要加载图片并裁剪
	if 0 < len(orientationFrames) {
		err = roleCutFrames(role, imageFilePath, orientationFrames)
		if err != nil {
			return errors.WithMessagef(err, "裁剪角色帧失败: %v %v", imageFilePath, xruntime.Location())
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

// 角色-裁剪-帧
func roleCutFrames(role *Role, imageFilePath string, orientationFrames map[proto.AssetOrientation][]*roleFrameData) error {
	// 加载图片
	img, _, err := ebitenv2ebitenutil.NewImageFromFile(imageFilePath)
	if err != nil {
		log.Printf("加载角色图片失败: %v %v", imageFilePath, err)
		return err
	}
	// 裁剪每个方向的帧
	for orientation, frames := range orientationFrames {
		// 裁剪并存储帧
		for _, fd := range frames {
			subImg := img.SubImage(image.Rect(
				fd.spriteInfo.Frame.X,
				fd.spriteInfo.Frame.Y,
				fd.spriteInfo.Frame.X+fd.spriteInfo.Frame.W,
				fd.spriteInfo.Frame.Y+fd.spriteInfo.Frame.H,
			)).(*ebitenv2.Image)
			switch fd.roleAction {
			case proto.RoleAction_RoleAction_Move:
				role.Move.Frames[orientation] = append(role.Move.Frames[orientation], subImg)
				role.Move.FrameInfo[orientation] = append(role.Move.FrameInfo[orientation], fd.spriteInfo)
			default:
				return fmt.Errorf("未实现的角色动作裁剪: %v", fd.roleAction)
			}
		}
		// 生成镜像帧
		roleGenerateMirrorFrames(role, orientation)
	}
	return nil
}

// 生成镜像 (左->右)
func roleGenerateMirrorFrames(role *Role, orientation proto.AssetOrientation) {
	var mirrorOrientation proto.AssetOrientation // 镜像方向
	switch orientation {
	case proto.AssetOrientation_AssetOrientation_Left: // 左 -> 右
		mirrorOrientation = proto.AssetOrientation_AssetOrientation_Right
	case proto.AssetOrientation_AssetOrientation_DownLeft: // 左下 -> 右下
		mirrorOrientation = proto.AssetOrientation_AssetOrientation_DownRight
	case proto.AssetOrientation_AssetOrientation_UpLeft: // 左上 -> 右上
		mirrorOrientation = proto.AssetOrientation_AssetOrientation_UpRight
	default:
		return // 不需要镜像
	}
	for i, srcImg := range role.Move.Frames[orientation] {
		bounds := srcImg.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		// 创建镜像图片
		mirrorImg := ebitenv2.NewImage(w, h)
		op := &ebitenv2.DrawImageOptions{}
		op.GeoM.Scale(-1, 1)             // 水平翻转
		op.GeoM.Translate(float64(w), 0) // 平移回原位置
		mirrorImg.DrawImage(srcImg, op)
		role.Move.Frames[mirrorOrientation] = append(role.Move.Frames[mirrorOrientation], mirrorImg)
		// 复制帧信息(镜像帧信息与原帧相同)
		role.Move.FrameInfo[mirrorOrientation] = append(role.Move.FrameInfo[mirrorOrientation], role.Move.FrameInfo[orientation][i])
	}
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

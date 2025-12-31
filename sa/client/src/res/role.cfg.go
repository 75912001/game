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
	Frame    *Rect `json:"frame"`              // 在大图中的位置和尺寸
	HitFrame bool  `json:"hitFrame,omitempty"` // 是否 命中帧
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
// 从json配置文件中解析动作和方向信息, 统一处理所有动作类型
// 参数:
//   - roleID: 角色ID
//   - imageFilePath: 图片文件路径 (用于加载和裁剪)
//   - roleJson: 解析后的json配置数据
//
// 返回: 加载失败时返回错误
func loadRoleImage(roleID common.AssetID, imageFilePath string, roleJson *RoleJson) (err error) {
	role, exists := GRoleMgr.Roles.Find(roleID)
	if !exists { // 不存在则创建
		role = &Role{
			ID:      roleID,
			Actions: make(map[proto.RoleAction]*RoleActionData), // 初始化动作映射表
		}
		ok := GRoleMgr.Roles.AddIfNotExist(roleID, role)
		if !ok {
			return fmt.Errorf("添加角色失败,角色已存在: %v %v", roleID, xruntime.Location())
		}
	}
	// 按动作和方向分组帧信息
	// 外层key: 动作类型 (Move, AttackAxe等)
	// 内层key: 方向 (Up, Down, Left等)
	// value: 该动作-方向组合的所有帧数据
	actionOrientationFrames := make(map[proto.RoleAction]map[proto.AssetOrientation][]*roleFrameData)
	for frameFileName, spriteData := range roleJson.Frames { // 解析帧名称: role.${roleID}.${arg}.${data}.${frameIndex}
		args := strings.Split(frameFileName, ".")
		if len(args) < 5 {
			return fmt.Errorf("帧名称格式错误: %v %s %v", imageFilePath, frameFileName, xruntime.Location())
		}
		action := args[2]      // 动作名称 (move, attackAxe等)
		orientation := args[3] // 方向名称 (up, down, left等)
		frameIndex := args[4]  // 帧索引

		// 解析动作类型
		roleAction := GetRoleActionByName(action)
		if roleAction == proto.RoleAction_RoleAction_Unknow {
			return fmt.Errorf("帧名称包含未知动作: %v %s %v", imageFilePath, action, xruntime.Location())
		}

		// 解析方向
		roleOrientation := GetAssetOrientationByName(orientation)
		if roleOrientation == proto.AssetOrientation_AssetOrientation_Unknow {
			return fmt.Errorf("帧名称包含未知方向: %v %s %v", imageFilePath, orientation, xruntime.Location())
		}

		// 解析帧索引
		idx, err := strconv.Atoi(frameIndex)
		if err != nil {
			return errors.WithMessagef(err, "解析帧索引失败: %v %s %v", imageFilePath, frameFileName, xruntime.Location())
		}

		// 添加到分组映射表 (统一处理所有动作, 无需switch-case)
		// 确保动作存在于映射表中
		if actionOrientationFrames[roleAction] == nil {
			actionOrientationFrames[roleAction] = make(map[proto.AssetOrientation][]*roleFrameData)
		}
		actionOrientationFrames[roleAction][roleOrientation] = append(actionOrientationFrames[roleAction][roleOrientation], &roleFrameData{
			index:      idx,
			roleAction: roleAction,
			spriteInfo: spriteData,
		})
	}
	// 按帧索引排序 (遍历所有动作的所有方向)
	for _, orientationFrames := range actionOrientationFrames {
		for _, frames := range orientationFrames {
			sortFrames(frames)
		}
	}
	// 如果有帧数据，需要加载图片并裁剪
	if 0 < len(actionOrientationFrames) {
		err = roleCutFrames(role, imageFilePath, actionOrientationFrames)
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

// roleCutFrames 角色-裁剪-帧
// 从原始大图中裁剪出所有动作的各个方向的动画帧
// 参数:
//   - role: 角色资源对象
//   - imageFilePath: 原始大图路径
//   - actionOrientationFrames: 动作->方向->帧数据 的三层映射
//
// 返回: 裁剪失败时返回错误
func roleCutFrames(role *Role, imageFilePath string, actionOrientationFrames map[proto.RoleAction]map[proto.AssetOrientation][]*roleFrameData) error {
	// 加载图片
	img, _, err := ebitenv2ebitenutil.NewImageFromFile(imageFilePath)
	if err != nil {
		log.Printf("加载角色图片失败: %v %v", imageFilePath, err)
		return err
	}
	// 遍历所有动作类型
	for roleAction, orientationFrames := range actionOrientationFrames {
		// 确保该动作的数据结构存在
		if role.Actions[roleAction] == nil {
			role.Actions[roleAction] = NewRoleActionData()
		}
		actionData := role.Actions[roleAction]

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

				// 存储到对应动作的对应方向
				actionData.Frames[orientation] = append(actionData.Frames[orientation], subImg)
				actionData.FrameInfo[orientation] = append(actionData.FrameInfo[orientation], fd.spriteInfo)
			}
			// 生成镜像帧 (左->右, 左下->右下, 左上->右上)
			roleGenerateMirrorFrames(actionData, orientation)
		}
	}
	return nil
}

// roleGenerateMirrorFrames 生成镜像帧
// 为左侧方向的帧自动生成对应的右侧镜像帧
// 参数:
//   - actionData: 动作数据对象
//   - orientation: 原始方向 (仅处理左/左下/左上)
//
// 镜像规则:
//   - 左 -> 右
//   - 左下 -> 右下
//   - 左上 -> 右上
//   - 其他方向不处理
func roleGenerateMirrorFrames(actionData *RoleActionData, orientation proto.AssetOrientation) {
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
	for i, srcImg := range actionData.Frames[orientation] {
		bounds := srcImg.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		// 创建镜像图片
		mirrorImg := ebitenv2.NewImage(w, h)
		op := &ebitenv2.DrawImageOptions{}
		op.GeoM.Scale(-1, 1)             // 水平翻转
		op.GeoM.Translate(float64(w), 0) // 平移回原位置
		mirrorImg.DrawImage(srcImg, op)
		actionData.Frames[mirrorOrientation] = append(actionData.Frames[mirrorOrientation], mirrorImg)
		// 复制帧信息(镜像帧信息与原帧相同)
		actionData.FrameInfo[mirrorOrientation] = append(actionData.FrameInfo[mirrorOrientation], actionData.FrameInfo[orientation][i])
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

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

// PetJson TexturePacker 导出的 宠物 配置
type PetJson struct {
	Frames map[string]*PetImageSprite `json:"frames"` // key: 文件名称 "pet.${petID}.${arg}.${data}.${frameIndex}" val: 信息
}

// PetImageSprite 宠物-精灵
type PetImageSprite struct {
	Frame    *Rect `json:"frame"`              // 在大图中的位置和尺寸
	HitFrame bool  `json:"hitFrame,omitempty"` // 是否命中帧
}

// 加载单个配置文件
func loadPetJson(filePath string) (*PetJson, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "读取宠物配置文件失败: %s %v", filePath, xruntime.Location())
	}
	var petJson PetJson
	err = json.Unmarshal(data, &petJson)
	if err != nil {
		return nil, errors.WithMessagef(err, "解析宠物配置文件失败: %s %v", filePath, xruntime.Location())
	}
	return &petJson, nil
}

// loadPetImage 加载宠物单个配置
func loadPetImage(petID common.AssetID, imageFilePath string, petJson *PetJson) (err error) {
	pet, exists := GPetMgr.Pets.Find(petID)
	if !exists { // 不存在则创建
		pet = &Pet{
			ID:     petID,
			Move:   NewPetMove(),
			Attack: NewPetAttack(),
		}
		ok := GPetMgr.Pets.AddIfNotExist(petID, pet)
		if !ok {
			return fmt.Errorf("添加宠物失败,宠物已存在: %v %v", petID, xruntime.Location())
		}
	}
	// 按动作和方向分组帧信息
	actionOrientationFrames := make(map[proto.PetAction]map[proto.AssetOrientation][]*petFrameData)
	for frameFileName, spriteData := range petJson.Frames { // 解析帧名称: pet.${petID}.${arg}.${data}.${frameIndex}
		args := strings.Split(frameFileName, ".")
		if len(args) < 5 {
			return fmt.Errorf("帧名称格式错误: %v %s %v", imageFilePath, frameFileName, xruntime.Location())
		}
		action := args[2]      // 动作
		orientation := args[3] // 方向
		frameIndex := args[4]  // 帧索引
		petAction := GetPetActionByName(action)
		if petAction == proto.PetAction_PetAction_Unknow {
			return fmt.Errorf("帧名称包含未知动作: %v %s %v", imageFilePath, action, xruntime.Location())
		}
		petOrientation := GetAssetOrientationByName(orientation)
		if petOrientation == proto.AssetOrientation_AssetOrientation_Unknow {
			return fmt.Errorf("帧名称包含未知方向: %v %s %v", imageFilePath, orientation, xruntime.Location())
		}
		idx, err := strconv.Atoi(frameIndex)
		if err != nil {
			return errors.WithMessagef(err, "解析帧索引失败: %v %s %v", imageFilePath, frameFileName, xruntime.Location())
		}
		// 添加到分组映射表
		if actionOrientationFrames[petAction] == nil {
			actionOrientationFrames[petAction] = make(map[proto.AssetOrientation][]*petFrameData)
		}
		actionOrientationFrames[petAction][petOrientation] = append(actionOrientationFrames[petAction][petOrientation], &petFrameData{
			index:      idx,
			petAction:  petAction,
			spriteInfo: spriteData,
		})
	}
	// 按帧索引排序
	for _, orientationFrames := range actionOrientationFrames {
		for _, frames := range orientationFrames {
			sortPetFrames(frames)
		}
	}
	// 如果有帧数据，需要加载图片并裁剪
	if len(actionOrientationFrames) > 0 {
		err = petCutFrames(pet, imageFilePath, actionOrientationFrames)
		if err != nil {
			return errors.WithMessagef(err, "裁剪宠物帧失败: %v %v", imageFilePath, xruntime.Location())
		}
	}
	return nil
}

// petFrameData 宠物-帧数据(用于排序)
type petFrameData struct {
	index      int
	petAction  proto.PetAction
	spriteInfo *PetImageSprite
}

// 宠物-裁剪-帧
func petCutFrames(pet *Pet, imageFilePath string, actionOrientationFrames map[proto.PetAction]map[proto.AssetOrientation][]*petFrameData) error {
	// 加载图片
	img, _, err := ebitenv2ebitenutil.NewImageFromFile(imageFilePath)
	if err != nil {
		log.Printf("加载宠物图片失败: %v %v", imageFilePath, err)
		return err
	}
	// 遍历所有动作类型
	for petAction, orientationFrames := range actionOrientationFrames {
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
				switch petAction {
				case proto.PetAction_PetAction_Move:
					pet.Move.Frames[orientation] = append(pet.Move.Frames[orientation], subImg)
					pet.Move.FrameInfo[orientation] = append(pet.Move.FrameInfo[orientation], fd.spriteInfo)
				case proto.PetAction_PetAction_Attack:
					pet.Attack.Frames[orientation] = append(pet.Attack.Frames[orientation], subImg)
					pet.Attack.FrameInfo[orientation] = append(pet.Attack.FrameInfo[orientation], fd.spriteInfo)
				default:
					return fmt.Errorf("未实现的宠物动作裁剪: %v", petAction)
				}
			}
			// 生成镜像帧
			petGenerateMirrorFrames(pet, petAction, orientation)
		}
	}
	return nil
}

// 生成镜像 (左->右)
func petGenerateMirrorFrames(pet *Pet, petAction proto.PetAction, orientation proto.AssetOrientation) {
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

	// 根据动作类型选择帧数据
	var srcFrames []*ebitenv2.Image
	var srcFrameInfos []*PetImageSprite
	switch petAction {
	case proto.PetAction_PetAction_Move:
		srcFrames = pet.Move.Frames[orientation]
		srcFrameInfos = pet.Move.FrameInfo[orientation]
	case proto.PetAction_PetAction_Attack:
		srcFrames = pet.Attack.Frames[orientation]
		srcFrameInfos = pet.Attack.FrameInfo[orientation]
	default:
		return
	}

	for i, srcImg := range srcFrames {
		bounds := srcImg.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		// 创建镜像图片
		mirrorImg := ebitenv2.NewImage(w, h)
		op := &ebitenv2.DrawImageOptions{}
		op.GeoM.Scale(-1, 1)             // 水平翻转
		op.GeoM.Translate(float64(w), 0) // 平移回原位置
		mirrorImg.DrawImage(srcImg, op)

		// 存储到对应动作的镜像方向
		switch petAction {
		case proto.PetAction_PetAction_Move:
			pet.Move.Frames[mirrorOrientation] = append(pet.Move.Frames[mirrorOrientation], mirrorImg)
			pet.Move.FrameInfo[mirrorOrientation] = append(pet.Move.FrameInfo[mirrorOrientation], srcFrameInfos[i])
		case proto.PetAction_PetAction_Attack:
			pet.Attack.Frames[mirrorOrientation] = append(pet.Attack.Frames[mirrorOrientation], mirrorImg)
			pet.Attack.FrameInfo[mirrorOrientation] = append(pet.Attack.FrameInfo[mirrorOrientation], srcFrameInfos[i])
		default:
			return
		}
	}
}

// sortPetFrames 按帧索引排序
func sortPetFrames(frames []*petFrameData) {
	for i := 0; i < len(frames)-1; i++ {
		for j := i + 1; j < len(frames); j++ {
			if frames[i].index > frames[j].index {
				frames[i], frames[j] = frames[j], frames[i]
			}
		}
	}
}

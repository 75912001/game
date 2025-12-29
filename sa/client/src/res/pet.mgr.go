package res

import (
	"fmt"
	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"regexp"
	"saClient/src/common"
	"saClient/src/proto"
	"strconv"
)

var GPetMgr = newPetMgr()

// PetMgr 宠物配置管理器
type PetMgr struct {
	Pets *xmap.MapMgr[common.AssetID, *Pet]
}

func newPetMgr() *PetMgr {
	return &PetMgr{
		Pets: xmap.NewMapMgr[common.AssetID, *Pet](),
	}
}

// Load 加载宠物res配置文件
func (p *PetMgr) Load() error {
	// 读取目录下所有文件-pet
	petDirectories, err := os.ReadDir(common.AppResPetDir)
	if err != nil {
		return errors.WithMessagef(err, "读取宠物配置目录失败: %v %v", common.AppResPetDir, xruntime.Location())
	}
	for _, petDir := range petDirectories {
		if !petDir.IsDir() { // 非目录
			continue
		}
		petDirName := petDir.Name() // 目录名称-宠物ID
		var petID common.AssetID    // 宠物ID
		{
			id64, err := strconv.ParseUint(petDirName, 10, 32)
			if err != nil {
				return errors.WithMessagef(err, "解析宠物目录名称 %s 为 ID 失败: %v", petDirName, xruntime.Location())
			}
			petID = common.AssetID(id64)
		}
		// 读取 dirName 目录下的所有文件
		petDirPath := filepath.Join(common.AppResPetDir, petDirName)
		petFiles, err := os.ReadDir(petDirPath)
		if err != nil {
			return errors.WithMessagef(err, "读取宠物目录 %s 失败: %v", petDirPath, xruntime.Location())
		}
		// 匹配 pet.${petID}.${arg}.${data}.json 格式的文件 例如 `pet.4000003.move.up.json`
		petPattern := regexp.MustCompile(fmt.Sprintf(`^%v\.%v\.([^.]+)\.([^.]+)\.json$`, common.GetNameByAssetType(proto.AssetType_AssetType_Pet), petID))
		for _, petFile := range petFiles {
			if petFile.IsDir() { // 目录
				continue
			}
			jsonFileName := petFile.Name()
			if !petPattern.MatchString(jsonFileName) { // 不匹配格式
				continue
			}
			imageFilePath := filepath.Join(petDirPath, jsonFileName[:len(jsonFileName)-len(filepath.Ext(jsonFileName))]+".png")
			// 加载配置文件
			jsonFilePath := filepath.Join(petDirPath, jsonFileName)
			petJson, err := loadPetJson(jsonFilePath)
			if err != nil {
				return errors.WithMessagef(err, "加载配置文件失败 %s %v", jsonFileName, xruntime.Location())
			}
			err = loadPetImage(petID, imageFilePath, petJson)
			if err != nil {
				return errors.WithMessagef(err, "加载宠物失败 %v %v", petID, xruntime.Location())
			}
		}
	}
	return nil
}

// Check 验证宠物资源完整性
func (p *PetMgr) Check() error {
	// 验证每个宠物的8个方向是否完整
	var checkErr error
	p.Pets.Foreach(func(key common.AssetID, pet *Pet) bool {
		if err := checkPet8Directions(pet); err != nil {
			checkErr = err
			return false // 停止遍历
		}
		return true // 继续遍历
	})
	return checkErr
}

func (p *PetMgr) Assemble() error {
	return nil
}

// checkPet8Directions 验证宠物8个方向是否完整
func checkPet8Directions(pet *Pet) error {
	for _, value := range proto.AssetOrientation_value {
		orientation := proto.AssetOrientation(value)
		if orientation == proto.AssetOrientation_AssetOrientation_Unknow {
			continue
		}
		if orientation == proto.AssetOrientation_AssetOrientation_Max {
			continue
		}
		// 验证 Move 动画的8个方向
		if len(pet.Move.Frames[orientation]) == 0 {
			return fmt.Errorf("宠物 %v 移动动画缺少方向 %v %v", pet.ID, GetNameByAssetOrientation(orientation), xruntime.Location())
		}
	}

	return nil
}

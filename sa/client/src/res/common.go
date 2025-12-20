package res

import (
	"saClient/src/proto"
	"strings"
)

// Rect 矩形
type Rect struct {
	X      int `json:"x"` // x 坐标
	Y      int `json:"y"` // y 坐标
	Width  int `json:"w"` // 宽度
	Height int `json:"h"` // 高度
}

// GetNameByAssetOrientation 根据资产方向获取名称
// 例如: 枚举 AssetOrientation_AssetOrientation_Up 的字符串为 "AssetOrientation_Up" 返回 "up"
func GetNameByAssetOrientation(assetOrientation proto.AssetOrientation) string {
	str, ok := proto.AssetOrientation_name[int32(assetOrientation)]
	if ok {
		if i := strings.Index(str, "_"); i >= 0 && i+1 < len(str) {
			return strings.ToLower(str[i+1:])
		}
		name := strings.ToLower(str)
		if name == "max" { // 排除枚举中的最大值,该值不表示具体方向
			return "unknown"
		}
		return name
	}
	return "unknown"
}

// GetAssetOrientationByName 根据名称获取资产方向枚举值
func GetAssetOrientationByName(name string) proto.AssetOrientation {
	for k, v := range proto.AssetOrientation_name {
		str := v
		if i := strings.Index(str, "_"); i >= 0 && i+1 < len(str) {
			str = strings.ToLower(str[i+1:])
		} else {
			str = strings.ToLower(str)
		}
		if str == name {
			return proto.AssetOrientation(k)
		}
	}
	return proto.AssetOrientation_AssetOrientation_Unknow
}

// GetNameByRoleAction 根据角色动作获取名称
func GetNameByRoleAction(roleAction proto.RoleAction) string {
	str, ok := proto.RoleAction_name[int32(roleAction)]
	if ok {
		if i := strings.Index(str, "_"); i >= 0 && i+1 < len(str) {
			return strings.ToLower(str[i+1:])
		}
		name := strings.ToLower(str)
		if name == "max" { // 排除枚举中的最大值,该值不表示具体动作
			return "unknown"
		}
		return name
	}
	return "unknown"
}

// GetRoleActionByName 根据名称获取角色动作枚举值
func GetRoleActionByName(name string) proto.RoleAction {
	for k, v := range proto.RoleAction_name {
		str := v
		if i := strings.Index(str, "_"); i >= 0 && i+1 < len(str) {
			str = strings.ToLower(str[i+1:])
		} else {
			str = strings.ToLower(str)
		}
		if str == name {
			return proto.RoleAction(k)
		}
	}
	return proto.RoleAction_RoleAction_Unknow
}

// GetPetActionByName 根据名称获取宠物动作枚举值
func GetPetActionByName(name string) proto.PetAction {
	for k, v := range proto.PetAction_name {
		str := v
		if i := strings.Index(str, "_"); i >= 0 && i+1 < len(str) {
			str = strings.ToLower(str[i+1:])
		} else {
			str = strings.ToLower(str)
		}
		if str == name {
			return proto.PetAction(k)
		}
	}
	return proto.PetAction_PetAction_Unknow
}

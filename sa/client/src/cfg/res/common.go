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

// GetNameByAssetType 根据资产类型获取资源名称
// 例如: 枚举 AssetType_AssetType_Role 的字符串为 "AssetType_Role" 返回 "role"
func GetNameByAssetType(assetType proto.AssetType) string {
	value, ok := proto.AssetType_name[int32(assetType)]
	if ok {
		if i := strings.Index(value, "_"); i >= 0 && i+1 < len(value) {
			return strings.ToLower(value[i+1:])
		}
		name := strings.ToLower(value)
		if name == "max" { // 排除枚举中的最大值,该值不表示具体资源类型
			return "unknow"
		}
		return name
	}
	return "unknow"
}

// GetNameByRoleDirection 根据角色方向获取名称
// 例如: 枚举 RoleDirection_RoleDirection_Up 的字符串为 "RoleDirection_Up" 返回 "up"
func GetNameByRoleDirection(roleDirection proto.RoleDirection) string {
	str, ok := proto.RoleDirection_name[int32(roleDirection)]
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

// GetRoleDirectionByName 根据名称获取角色方向枚举值
func GetRoleDirectionByName(name string) proto.RoleDirection {
	for k, v := range proto.RoleDirection_name {
		str := v
		if i := strings.Index(str, "_"); i >= 0 && i+1 < len(str) {
			str = strings.ToLower(str[i+1:])
		} else {
			str = strings.ToLower(str)
		}
		if str == name {
			return proto.RoleDirection(k)
		}
	}
	return proto.RoleDirection_RoleDirection_Unknow
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

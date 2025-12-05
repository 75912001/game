package res

import (
	"saClient/src/proto"
	"strings"
)

// GetNameByType 根据资源类型获取资源名称
// 例如: 枚举 AssetType_AssetType_Role 的字符串为 "AssetType_Role" 返回 "role"
func GetNameByType(assetType proto.AssetType) string {
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

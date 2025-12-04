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
		return strings.ToLower(value)
	}
	return "unknow"
}

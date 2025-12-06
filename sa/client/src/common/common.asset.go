package common

import (
	"saClient/src/proto"
	"strings"
)

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

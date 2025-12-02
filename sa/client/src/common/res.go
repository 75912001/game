package common

import (
	"saClient/src/proto"
)

func GetResNameByType(resType proto.ResType) string {
	switch resType {
	case proto.ResType_ResType_Role:
		return "role"
	case proto.ResType_ResType_Map:
		return "map"
	case proto.ResType_ResType_Item:
		return "item"
	default:
		return "unknow"
	}
}

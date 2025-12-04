package cfg

import (
	"saClient/src/common"
)

type Role struct {
	ID   common.AssetID // 角色ID
	Move *RoleMove
}

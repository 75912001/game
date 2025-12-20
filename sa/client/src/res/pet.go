package res

import (
	"saClient/src/common"
)

type Pet struct {
	ID   common.AssetID // 宠物ID
	Move *PetMove
}

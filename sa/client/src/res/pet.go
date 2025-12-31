package res

import (
	"saClient/src/common"
)

type Pet struct {
	ID     common.AssetID // 宠物ID
	Move   *PetMove       // 移动动画
	Attack *PetAttack     // 攻击动画
}

package elemental

import "saClient/src/proto"

// Elemental 元素属性对象，总共10点可分配
// 只能分配给相邻的两个属性或单个属性
// 相邻关系:土-水,水-火,火-风,风-土(循环)

func Get(assetRoleBaseMap map[uint32]uint64) (earth, water, fire, wind uint32) {
	return GetEarth(assetRoleBaseMap), GetWater(assetRoleBaseMap), GetFire(assetRoleBaseMap), GetWind(assetRoleBaseMap)

}
func GetEarth(assetRoleBaseMap map[uint32]uint64) uint32 {
	value, ok := assetRoleBaseMap[uint32(proto.AssetIDRoleBase_AssetIDRoleBase_ElementalEarth)]
	if ok {
		return uint32(value)
	}
	return 0
}

func GetWater(assetRoleBaseMap map[uint32]uint64) uint32 {
	value, ok := assetRoleBaseMap[uint32(proto.AssetIDRoleBase_AssetIDRoleBase_ElementalWater)]
	if ok {
		return uint32(value)
	}
	return 0
}

func GetFire(assetRoleBaseMap map[uint32]uint64) uint32 {
	value, ok := assetRoleBaseMap[uint32(proto.AssetIDRoleBase_AssetIDRoleBase_ElementalFire)]
	if ok {
		return uint32(value)
	}
	return 0
}

func GetWind(assetRoleBaseMap map[uint32]uint64) uint32 {
	value, ok := assetRoleBaseMap[uint32(proto.AssetIDRoleBase_AssetIDRoleBase_ElementalWind)]
	if ok {
		return uint32(value)
	}
	return 0
}

// Validate 验证元素分配是否合法
// 规则：
// 1. 总和必须为10点
// 2. 只能分配给相邻的两个属性或单个属性
func Validate(earth, water, fire, wind uint32) error {
	total := earth + water + fire + wind

	// 验证总和是否为10
	if total != 10 {
		return &Error{Message: "元素总和必须为10点"}
	}

	// 统计有值的属性数量
	nonZeroCount := 0
	if earth > 0 {
		nonZeroCount++
	}
	if water > 0 {
		nonZeroCount++
	}
	if fire > 0 {
		nonZeroCount++
	}
	if wind > 0 {
		nonZeroCount++
	}

	// 只能有1个或2个属性有值
	if nonZeroCount > 2 {
		return &Error{Message: "只能分配给1个或2个相邻的属性"}
	}

	// 如果是2个属性，验证是否相邻
	if nonZeroCount == 2 {
		if !isAdjacent() {
			return &Error{Message: "两个属性必须相邻（土-水、水-火、火-风、风-土）"}
		}
	}

	return nil
}

// isAdjacent 检查两个非零属性是否相邻
func isAdjacent(earth, water, fire, wind uint32) bool {
	// 相邻关系：土(1)-水(2)、水(2)-火(3)、火(3)-风(4)、风(4)-土(1)
	if earth > 0 && water > 0 {
		return true // 土-水
	}
	if water > 0 && fire > 0 {
		return true // 水-火
	}
	if fire > 0 && wind > 0 {
		return true // 火-风
	}
	if wind > 0 && earth > 0 {
		return true // 风-土
	}
	return false
}

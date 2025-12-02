package common

type ElementalType uint32 // 元素类型

const (
	ElementalTypeUnknow ElementalType = iota // 无
	ElementalTypeEarth                       // 土
	ElementalTypeWater                       // 水
	ElementalTypeFire                        // 火
	ElementalTypeWind                        // 风
)

// Elemental 元素属性对象，总共10点可分配
// 只能分配给相邻的两个属性或单个属性
// 相邻关系：土-水、水-火、火-风、风-土（循环）
type Elemental struct {
	earth uint32 // 土-数值
	water uint32 // 水-数值
	fire  uint32 // 火-数值
	wind  uint32 // 风-数值
}

// NewElemental 创建新的元素属性对象
// 参数：各元素类型的点数
// 返回：元素对象和错误信息
func NewElemental(earth, water, fire, wind uint32) (*Elemental, error) {
	e := &Elemental{
		earth: earth,
		water: water,
		fire:  fire,
		wind:  wind,
	}

	if err := e.Validate(); err != nil {
		return nil, err
	}

	return e, nil
}

// Set 设置元素属性值
func (e *Elemental) Set(earth, water, fire, wind uint32) error {
	e.earth = earth
	e.water = water
	e.fire = fire
	e.wind = wind

	return e.Validate()
}

// Get 获取所有元素属性值
func (e *Elemental) Get() (earth, water, fire, wind uint32) {
	return e.earth, e.water, e.fire, e.wind
}

// GetByType 按元素类型获取属性值
func (e *Elemental) GetByType(elementalType ElementalType) uint32 {
	switch elementalType {
	case ElementalTypeEarth:
		return e.earth
	case ElementalTypeWater:
		return e.water
	case ElementalTypeFire:
		return e.fire
	case ElementalTypeWind:
		return e.wind
	default:
	}
	return 0
}

// Validate 验证元素分配是否合法
// 规则：
// 1. 总和必须为10点
// 2. 只能分配给相邻的两个属性或单个属性
func (e *Elemental) Validate() error {
	total := e.earth + e.water + e.fire + e.wind

	// 验证总和是否为10
	if total != 10 {
		return &ElementalError{Message: "元素总和必须为10点"}
	}

	// 统计有值的属性数量
	nonZeroCount := 0
	if e.earth > 0 {
		nonZeroCount++
	}
	if e.water > 0 {
		nonZeroCount++
	}
	if e.fire > 0 {
		nonZeroCount++
	}
	if e.wind > 0 {
		nonZeroCount++
	}

	// 只能有1个或2个属性有值
	if nonZeroCount > 2 {
		return &ElementalError{Message: "只能分配给1个或2个相邻的属性"}
	}

	// 如果是2个属性，验证是否相邻
	if nonZeroCount == 2 {
		if !e.isAdjacent() {
			return &ElementalError{Message: "两个属性必须相邻（土-水、水-火、火-风、风-土）"}
		}
	}

	return nil
}

// isAdjacent 检查两个非零属性是否相邻
func (e *Elemental) isAdjacent() bool {
	// 相邻关系：土(1)-水(2)、水(2)-火(3)、火(3)-风(4)、风(4)-土(1)
	if e.earth > 0 && e.water > 0 {
		return true // 土-水
	}
	if e.water > 0 && e.fire > 0 {
		return true // 水-火
	}
	if e.fire > 0 && e.wind > 0 {
		return true // 火-风
	}
	if e.wind > 0 && e.earth > 0 {
		return true // 风-土
	}
	return false
}

// ElementalError 元素错误类型
type ElementalError struct {
	Message string
}

func (e *ElementalError) Error() string {
	return e.Message
}

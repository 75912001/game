package common

type PlaneFrameType uint32 // 飞行姿态

const (
	PlaneFrameTypeStraight    PlaneFrameType = iota // 直飞
	PlaneFrameTypeSlightRight                       // 微右倾
	PlaneFrameTypeRight                             // 右倾
	PlaneFrameTypeSharpRight                        // 大右倾
	PlaneFrameTypeMax                               // 最大值
)

func (p *PlaneFrameType) GetFrameTypeCount() uint32 {
	return uint32(PlaneFrameTypeMax)
}

const PlaneFrameDelay = 5 // 每x帧切换一次动画

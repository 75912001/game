package common

type PlaneKey struct {
	ID    uint32 // id
	Level uint32 // 等级
}

type PlanePartType uint32

const (
	PlanePartTypeNose      PlanePartType = iota // 机头
	PlanePartTypeBody                           // 机身
	PlanePartTypeLeftWing                       // 机翼-左
	PlanePartTypeRightWing                      // 机翼-右
	PlanePartTypeMax                            // 最大值
)

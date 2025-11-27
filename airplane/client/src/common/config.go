package common

// ConfigPlaneFrame 飞机-帧信息-在大图中的位置和尺寸
type ConfigPlaneFrame struct {
	Frame *Rect `json:"frame"` // 在大图中的位置和尺寸
}

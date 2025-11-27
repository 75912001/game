package config

import "airplaneClient/src/common"

// PlaneConfig TexturePacker 导出的飞机配置
type PlaneConfig struct {
	Frames map[string]*common.ConfigPlaneFrame `json:"frames"` // key: 部件文件名称 "plane.001.001.body.000" val: 部件信息
}

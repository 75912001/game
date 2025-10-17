package resources

import "fmt"

// GenPlaneName 生成飞机资源名称
func GenPlaneName(id uint32, level uint32) string {
	return fmt.Sprintf("%03d.%03d.png", id, level)
}

package resources

import "fmt"

// GenMapName 生成-地图-资源名称
func GenMapName(id uint32, level uint32) string {
	return fmt.Sprintf("%03d.%03d.jpg", id, level)
}

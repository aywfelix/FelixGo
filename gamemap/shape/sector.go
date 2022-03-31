// 扇形
package shape

import (
	"math"
)

type Sector struct {
	Center   Point   // 扇形中心点
	Theta    float64 // 配置的角度
	Distance float64 // 配置的距离
}

func (s Sector) toPolarCoordinate(p Point) int {
	angle := math.Atan2(p.Y, p.X) * 180.0 / math.Pi
	return (int(angle) + 360) % 360
}

func (s Sector) IsInSector(p Point) bool {
	// 先判断距离，然后判断角度即可
	if IsFarThanDistance(s.Center, p, s.Distance) {
		return false
	}
	rePoint := s.Center.ToRelative(p)
	reAngle := s.toPolarCoordinate(rePoint)
	cangle := s.toPolarCoordinate(p)
	if math.Abs(float64(cangle-reAngle)) < (s.Theta / 2) {
		return true
	}
	return false
}

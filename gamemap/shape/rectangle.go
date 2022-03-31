// 矩形
package shape

import (
	"math"
)

type Rectangle struct {
	AP Point
	BP Point
}

func (r Rectangle) MinXY() (float64, float64) {
	return math.Min(r.AP.X, r.BP.X), math.Min(r.AP.Y, r.BP.Y)
}

func (r Rectangle) MaxXY() (float64, float64) {
	return math.Max(r.AP.X, r.BP.X), math.Max(r.AP.Y, r.BP.Y)
}

func (r Rectangle) IsInRect(p Point) bool {
	minX, minY := r.MinXY()
	maxX, maxY := r.MaxXY()
	if p.X >= minX && p.Y >= minY && p.X <= maxX && p.Y <= maxY {
		return true
	}
	return false
}

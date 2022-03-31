// 直线
package shape

import (
	"math"
)

type Line struct {
	Start Point
	End   Point
}

const eps float64 = 1e-4

// 判定一个点在一个直线上
func (l Line) IsInLine(c Point) bool {
	if math.Abs(c.Sub(l.Start).Cross(c.Sub(l.End))) <= eps {
		minX, minY := math.Min(l.Start.X, l.End.X), math.Min(l.Start.Y, l.End.Y)
		maxX, maxY := math.Max(l.Start.X, l.End.X), math.Max(l.Start.Y, l.End.Y)
		if c.X >= minX && c.Y >= minY && c.X <= maxX && c.Y <= maxY {
			return true
		}
	}
	return false
}

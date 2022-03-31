// 圆形
package shape

type Circular struct {
	Center Point   // 圆心
	R      float64 // 半径
}

func (c Circular) IsInCircular(p Point) bool {
	rePoint := c.Center.ToRelative(p)
	if (rePoint.X*rePoint.X + rePoint.Y*rePoint.Y) <= c.R*c.R {
		return true
	}
	return false
}

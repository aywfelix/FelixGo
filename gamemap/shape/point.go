// 某个坐标点
package shape

type Point struct {
	X float64
	Y float64
}

func (a Point) Distance(b Point) float64 {
	disX := b.X - a.X
	disY := b.Y - a.Y
	return disX*disX + disY*disY
}

func (a Point) Sub(b Point) Point {
	return Point{X: a.X - b.X, Y: a.Y - b.Y}
}

func (a Point) Cross(b Point) float64 {
	return a.X*b.Y - a.Y*b.X
}

func (a Point) ToRelative(b Point) Point {
	return Point{X: a.X - b.X, Y: a.Y - b.Y}
}

func IsFarThanDistance(a, b Point, distance float64) bool {
	if a.Distance(b) > distance*distance {
		return true
	}
	return false
}

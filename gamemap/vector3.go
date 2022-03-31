package gamemap

import (
	"math"
	"math/rand"
	"time"
)

// Vector3 代码位置的3D矢量
type Vector3 struct {
	X float64
	Y float64
	Z float64
}

// NewVector3 创建一个新的矢量
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{
		x,
		y,
		z,
	}
}

// Vector3_Zero 返回零值
func Vector3Zero() Vector3 {
	return Vector3{
		0,
		0,
		0,
	}
}

// IsEqual 相等
func (v Vector3) IsEqual(r Vector3) bool {
	if v.X-r.X > math.SmallestNonzeroFloat64 ||
		v.X-r.X < -math.SmallestNonzeroFloat64 ||
		v.Y-r.Y > math.SmallestNonzeroFloat64 ||
		v.Y-r.Y < -math.SmallestNonzeroFloat64 ||
		v.Z-r.Z > math.SmallestNonzeroFloat64 ||
		v.Z-r.Z < -math.SmallestNonzeroFloat64 {
		return false
	}

	return true
}

// Add 加
func (v Vector3) Add(o Vector3) Vector3 {
	return Vector3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

// AddS 加到自己身上
func (v *Vector3) AddS(o Vector3) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
}

// Sub 减
func (v Vector3) Sub(o Vector3) Vector3 {
	return Vector3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

// SubS 自已身上减
func (v *Vector3) SubS(o Vector3) {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
}

// Mul 乘
func (v Vector3) Mul(o float64) Vector3 {
	return Vector3{v.X * o, v.Y * o, v.Z * o}
}

// MulS 自己乘
func (v *Vector3) MulS(o float64) {
	v.X *= o
	v.Y *= o
	v.Z *= o
}

// Cross 叉乘
func (v Vector3) Cross(o Vector3) Vector3 {
	return Vector3{v.Y*o.Z - v.Z*o.Y, v.Z*o.X - v.X*o.Z, v.X*o.Y - v.Y*o.X}
}

// Dot 点乘
func (v Vector3) Dot(o Vector3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// Len 获取长度
func (v Vector3) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v *Vector3) Normalize() {
	len := v.Len()

	if len < math.SmallestNonzeroFloat64 {
		return
	}

	v.X = v.X / len
	v.Y = v.Y / len
	v.Z = v.Z / len
}

// RandXZ 在XZ平面上半径为r的圆内选取一个随机点
func RandXZ(v Vector3, r float32) Vector3 {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))

	tarR := randSeed.Float64() * float64(r)
	angle := randSeed.Float64() * 2 * math.Pi

	pos := Vector3{}
	pos.Y = 0

	pos.X = math.Cos(angle) * tarR
	pos.Z = math.Sin(angle) * tarR

	return v.Add(pos)
}

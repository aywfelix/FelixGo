package gamemap

import (
	"math"
)

type Vector2 struct {
	X, Y float64
}

func NewVector2(x, y float64) Vector2 {
	return Vector2{x, y}
}

func Vector2Zero() Vector2 {
	return Vector2{0, 0}
}

func Vector2Invalid() Vector2 {
	return Vector2{math.MaxFloat64, math.MaxFloat64}
}

func (v Vector2) IsEqual(o Vector2) bool {
	return v.X == o.X && v.Y == o.Y
}

// add
func (v Vector2) Add(o Vector2) {
	v.X += o.X
	v.Y += o.Y
}

// sub
func (v Vector2) Sub(o Vector2) Vector2 {
	return Vector2{v.X - o.X, v.Y - o.Y}
}

// subs
func (v Vector2) Subs(o Vector2) {
	v.X -= o.X
	v.Y -= o.Y
}

// mul
func (v Vector2) Mul(o Vector2) Vector2 {
	return Vector2{v.X * o.X, v.Y * o.Y}
}

// Dot 点乘
func (v Vector2) Dot(o Vector2) float64 {
	return v.X*o.X + o.Y*o.Y
}

// Len 获取长度
func (v Vector2) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

// Cross 叉乘
func (v Vector2) Cross(o Vector2) float64 {
	return v.X*o.Y - v.Y*o.X
}

// normalize
func (v Vector2) Normalize() {
	len := v.Len()
	if len < math.SmallestNonzeroFloat64 {
		return
	}

	v.X = v.X / len
	v.Y = v.Y / len
}

func (v Vector2) IsFarThanDistance(o Vector2, distance float64) bool {
	disX := v.X - o.X
	disY := v.Y - o.Y
	if (disX*disX + disY*disY) > distance*distance {
		return true
	}
	return false
}

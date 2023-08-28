package main

import "math"

type Vector2D struct {
	x, y float64
}

func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{x: v1.x + v2.x, y: v1.y + v2.y}
}

func (v1 Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{x: v1.x - v2.x, y: v1.y - v2.y}
}

func (v1 Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{x: v1.x * v2.x, y: v1.y * v2.y}
}

func (v1 Vector2D) AddVal(d float64) Vector2D {
	return Vector2D{x: v1.x + d, y: v1.y + d}
}

func (v1 Vector2D) SubtractVal(d float64) Vector2D {
	return Vector2D{x: v1.x - d, y: v1.y - d}
}

func (v1 Vector2D) MultiplyByVal(d float64) Vector2D {
	return Vector2D{x: v1.x * d, y: v1.y * d}
}

func (v1 Vector2D) DivideByVal(d float64) Vector2D {
	return Vector2D{x: v1.x / d, y: v1.y / d}
}

func (v1 Vector2D) Limit(lower, upper float64) Vector2D {
	return Vector2D{
		x: math.Min(math.Max(v1.x, lower), upper),
		y: math.Min(math.Max(v1.y, lower), upper),
	}
}

func (v1 Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v1.x-v2.x, 2) + math.Pow(v1.y-v2.y, 2))
}

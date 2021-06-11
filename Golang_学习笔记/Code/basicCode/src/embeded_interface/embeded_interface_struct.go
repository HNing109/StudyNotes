package embeded_interface

import "math"

type Square struct {
	Side float32
}

type Circle struct {
	Radius float32
}

func (sq *Square) Area() float32 {
	return sq.Side * sq.Side
}

func (ci *Circle) Area() float32 {
	return ci.Radius * ci.Radius * math.Pi
}

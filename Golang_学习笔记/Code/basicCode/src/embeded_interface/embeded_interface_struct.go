package embeded_interface

import "math"

type Square struct {
	Side float32
	Name string
}

type Circle struct {
	Radius float32
}

func (sq *Square) Area() float32 {
	return sq.Side * sq.Side
}

func (sq *Square) Set(name ...string)string{
	sq.Name = name[0]
	return "set sq.Name = " + name[0] + " success"
}


func (ci *Circle) Area() float32 {
	return ci.Radius * ci.Radius * math.Pi
}

package embeded_interface

type Shaper interface {
	Area() float32
}

//嵌套接口
type AllShaper interface {
	Shaper
	Color() string
}

package main

import "fmt"

//接口
type Abser interface {
	Abs() float64
}

//结构体
type Vertex struct {
	x float64
	y float64
}

//实现结构体的方法
func (v *Vertex) Abs() float64{
	var res = v.x - v.y
	if res >= 0{
		return res
	} else{
		return -res
	}
}

func main(){
	var a Abser
	v := Vertex{3,4}
	//*Vertex 实现了 Abser的方法Abs()
	a = &v
	//调用方法
	fmt.Println(a.Abs())
}

package main

import "fmt"

/**
接口的隐式实现
 */

//接口
type I interface{
	Say()
}

//结构体
type T struct{
	str string
}

//实现接口中的方法
func (t *T) Say(){
	fmt.Println(t.str)
}

func main()  {
	var impl I = &T{str: "xxxx"}
	impl.Say()
}
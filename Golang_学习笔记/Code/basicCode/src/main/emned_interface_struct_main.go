package main

import (
	em "embeded_interface"
	"fmt"
)

//断言：用于判断结构体，是否继承至某个接口
func main() {
	var shaper em.Shaper
	sq := &em.Square{Side: 18}

	shaper = sq

	fmt.Println(shaper.Set("chris"))

	//断言： 接口变量.(*结构体名)
	if val, ok := shaper.(*em.Square); ok {
		fmt.Println("ok = , val = ", ok, val)
	}
}

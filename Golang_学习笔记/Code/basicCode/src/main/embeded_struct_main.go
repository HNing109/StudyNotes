package main

import (
	"embeded_struct"
	"fmt"
)

func main() {
	var outer = embeded_struct.NewOuter("chris", 18)
	//输出：{chris 18 { }}
	fmt.Println(*outer)
	//访问内部类
	outer.Inner.Name = "Fyj"
	outer.Inner.Sex = "women"
	//输出：{chris 18 {Fyj women}}
	fmt.Println(*outer)
}

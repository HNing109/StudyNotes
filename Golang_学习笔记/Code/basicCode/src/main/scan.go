package main

import (
	"fmt"
)

var(
	first_name string
	address string
	f float32
	j int
	s string
	input = "56.12 / 5212 / Go"
	format_0 = "%f %d %s"
	format_1 = "%f / %d / %s"
)

func main() {
	fmt.Println("input name and address:")
	//遇见换行符，终止输入
	fmt.Scanln(&first_name, &address)
	fmt.Println(first_name, address)

	fmt.Println("input data ：")
	//按照特定格式，输入数据
	fmt.Scanf(format_0, &f, &j, &s)
	fmt.Println( f, j, s)

	//按照format格式，读取input字符串，并分发给对应的参数
	fmt.Sscanf(input, format_1, &f, &j, &s)
	fmt.Println("From the string we read: ", f, j, s)

}
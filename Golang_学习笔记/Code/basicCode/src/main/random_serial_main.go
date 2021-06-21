package main

import (
	"fmt"
	"random"
)

/**
随机生成0、1序列
 */
func main() {
	var instance = new(random.RandomSerial)
	ch := make(chan int)

	for index := 0; index < 100; index++{
		res := instance.GetRandomSerial(ch)
		fmt.Print(res)
	}

}

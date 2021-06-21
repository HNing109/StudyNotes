package main

import (
	"fmt"
	"natural"
)

func main() {
	var instance = new(natural.GenerateNatural)
	ch := instance.GetData() 	// 自然数序列: 2, 3, 4, ...
	for i := 0; i < 100; i++ {
		prime := <-ch // 新出现的素数
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = instance.PrimeFilter(ch, prime) // 基于新素数构造的过滤器
	}
}

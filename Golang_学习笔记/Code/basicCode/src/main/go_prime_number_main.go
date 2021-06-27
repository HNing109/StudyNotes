package main

import (
	"fmt"
	gpn "go_prime_number"
)

/**
 * Created by Chris on 2021/6/27.
 */

func main(){
	var goPrimeNumber = new(gpn.GoPrimeNumber)
	//设置启动的协程数
	goPrimeNumber.GoroutineNum = 8
	//设置计算素数的范围：0开始
	goPrimeNumber.PrimeNumber = 80000
	//计算素数
	resChan, totalTime := goPrimeNumber.Count()
	//遍历结果
	for val := range resChan{
		fmt.Println(val)
	}
	fmt.Println("total time = ", totalTime)
}
 

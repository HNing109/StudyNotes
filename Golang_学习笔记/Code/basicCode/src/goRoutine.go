package main

import "fmt"

func goForSum(arr []int, ch chan int) {
	res := 0
	for _, val := range arr{
		res += val
	}
	//结果存入信道
	ch <- res
}

func main(){
	var arr []int = []int{1,2,3,4,5,6,7,8,9}
	//创建信道：缓冲区为2，即：信道中最多可存储2个数据
	ch := make(chan int, 3)
	//开启两个协程：计算求和
	go goForSum(arr[ : len(arr) / 2], ch)
	go goForSum(arr[len(arr) / 2 :], ch)
	//从信道中取出结果
	res1 := <- ch
	res2 := <- ch
	fmt.Println(res1, res2, res1 + res2)

	var ch2 = make(chan int, 1)
	ch2 <- 2
	fmt.Println(<- ch2)
}
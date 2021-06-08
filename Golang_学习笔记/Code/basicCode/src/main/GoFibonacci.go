package main

import "fmt"

func fibonacci(num int, ch chan int){
	pre := 0
	next := 1
	temp := 0
	for index := 0; index < num; index++{
		ch <- pre
		temp = pre
		pre = next
		next = pre + temp
	}
	//结束数据输入后：关闭信道
	close(ch)
}

func main(){
	ch := make(chan int, 10)
	go fibonacci(cap(ch), ch)
	for val := range ch{
		fmt.Println(val)
	}
}
package main

import (
	"fmt"
	"time"
)

/**
使用Go协程计算斐波那契数列
 */

func dup3(in <-chan int64) (<-chan int64, <-chan int64, <-chan int64) {
	a, b, c := make(chan int64, 2), make(chan int64, 2), make(chan int64, 2)
	go func() {
		for {
			x := <-in
			a <- x
			b <- x
			c <- x
		}
	}()
	return a, b, c
}

func fib() <-chan int64 {
	x := make(chan int64, 2)
	a, b, out := dup3(x)
	go func() {
		//初始化
		x <- 0
		x <- 1
		//丢弃一个数据：0
		<-a
		for {
			//自动堵塞，若a、b有数就会自动取出
			x <- <-a + <-b
		}
	}()
	<- out
	return out
}

func main() {
	start := time.Now()
	x := fib()
	for i := 0; i < 80; i++ {
		fmt.Println(<-x)
	}
	end := time.Now()
	delta := end.Sub(start)
	fmt.Printf("longCalculation took this amount of time: %s\n", delta)
}

package main

import (
	"fmt"
	"time"
)

func selectFibonacci(ch chan int, quitCh chan int){
	pre := 0
	next := 1
	for{
		//实现go协程的等待多个通信操作：
		//select 会阻塞到某个分支可以继续执行为止，这时就会执行该分支。当多个分支都准备好时会随机选择一个执行
		select{
		case ch <- pre:
		pre, next = next , pre + next
		case <- quitCh :
			fmt.Println("quit")
			return
		//若没有准备好的协程，就执行default
		default:
			fmt.Println("execute default")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main(){
	ch := make(chan int, 10)
	quit := make(chan int, 1)
	//使用go协程打印数据
	go func() {
		for index := 0; index < cap(ch); index++{
			fmt.Println(<- ch)
		}
		quit <- 0
	}()
	//执行计算
	selectFibonacci(ch, quit)
}

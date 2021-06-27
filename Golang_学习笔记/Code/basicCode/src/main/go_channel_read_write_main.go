package main

import (
	"fmt"
	"time"
)

func main() {
	//数据管道
	var chData = make(chan int, 10)
	//退出管道
	var chExit = make(chan int)
	var totalNum = 20

	//写入数据的协程
	go func(chData chan int) {
		for index := 0; index < totalNum; index++{
			time.Sleep(500 * time.Millisecond)
			chData <- index
			fmt.Printf("th %d: input data = %d, len = %d\n", index, index, len(chData))
		}
		//写入完成后，关闭管道
		close(chData)
	}(chData)


	//读取数据的协程
	go func(chData chan int, chExit chan int){
		//方式1：for-range，遍历结束前一定要关闭管道，否则会出现死锁
		for val := range chData{
			fmt.Printf("get data = %d\n", val)
		}

		//方式2：while，循环判断是否关闭管道
		//for {
		//	if val, open := <- chData; !open{
		//		break
		//	} else{
		//		fmt.Printf("get data = %d\n", val)
		//	}
		//}

		//关闭退出管道
		close(chExit)
	}(chData, chExit)


	//检测到退出管道关闭后，程序向下执行（即：退出main主函数）
	for {
		if _, open := <- chExit; !open{
			break
		}
	}
}
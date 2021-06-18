package csp

import (
	"fmt"
	"sync"
)

type GoWaitWaitGroup struct{

}

func(g *GoWaitWaitGroup) Test() {
	var wg sync.WaitGroup

	// 开N个后台打印线程
	for i := 0; i < 10; i++ {
		//增加一个等待时间数
		wg.Add(1)
		go func() {
			fmt.Println("你好, 世界")
			//完成一个等待时间
			wg.Done()
		}()
	}

	// 等待N个后台线程完成
	wg.Wait()
	fmt.Println("finish all")
	// 等待N个后台线程完成
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

package go_prime_number

import (
	"fmt"
	"time"
)

/**
 * Created by Chris on 2021/6/27.
 */

type GoPrimeNumber struct{
	//goroutine的数量
	GoroutineNum int
	//求解素数的范围
	PrimeNumber int
}

/*
存入数据
*/
func(g GoPrimeNumber) inputNum(dataChan chan int){
	for i := 0; i < g.PrimeNumber; i++{
		dataChan <- i
	}
	close(dataChan)
}

/**
计算素数
*/
func(g GoPrimeNumber) countPrimeNumber(dataChan chan int, resChan chan int, exitChan chan bool){
	var flag bool
	for{
		num, open := <- dataChan
		if !open{
			break
		}
		//假设是素数
		flag = true
		//判断素数
		for i := 2; i < num; i++{
			//该数不是素数
			if num % i == 0{
				flag = false
				break
			}
		}
		//存入素数
		if flag{
			resChan <- num
		}
	}
	fmt.Println("有一个Goroutine执行结束")
	exitChan <- true
}

/**
计算素数
 */
func(g GoPrimeNumber) Count() (chan int, time.Duration){
	dataChan := make(chan int, g.PrimeNumber)
	resChan := make(chan int, g.PrimeNumber / 2)
	exitChan := make(chan bool, g.GoroutineNum)

	//存入数据
	go g.inputNum(dataChan)

	//计时开始
	startTime := time.Now()

	//计算素数：开启多个协程，计算素数
	for i := 0; i < g.GoroutineNum; i++{
		g.countPrimeNumber(dataChan, resChan, exitChan)
	}

	//等待：计算完毕
	go func() {
		for i := 0; i < g.GoroutineNum; i++{
			<- exitChan
		}
		//关闭管道
		close(resChan)
	}()

	//计时结束
	endTime := time.Now()

	return resChan, endTime.Sub(startTime)
}



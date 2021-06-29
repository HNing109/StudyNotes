package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

/**
运行代码后，执行go tool trace trace.out
 * Created by Chris on 2021/6/29.
 */

func main() {
	//创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	//main
	fmt.Println("Hello World")
}
 
 
 
 
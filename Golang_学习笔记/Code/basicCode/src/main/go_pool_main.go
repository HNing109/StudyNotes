package main

import "go_pool"

/**
测试协程池
 * Created by Chris on 2021/6/28.
 */

func main() {
	goPool := new(go_pool.GoPool)
	goPool.Start()
}
 
 
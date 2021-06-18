package main

import (
	"fmt"
	singleton "gof/singleton"
)
/**
创建两个对象，获得的是同一个实例（单例模式）
 */
func main(){
	var instance1 = singleton.GetSingletonInstance("chris")
	fmt.Println(instance1)
	var instance2 = singleton.GetSingletonInstance("Fyj")
	fmt.Println(instance2)
}

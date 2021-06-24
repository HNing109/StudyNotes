package main

import (
	"fmt"
	"math/rand"
	pi "polymorphic_interface"
	"sort"
	"strconv"
)

/**
实现接口方法，调用sort.Sort()方法。
 */
func main() {
	var ps pi.PersonSlice
	//初始化：随机值
	for index := 0; index < 10; index++{
		var p pi.Person
		p.SetName("hero id : " + strconv.Itoa(rand.Intn(100)))
		p.SetAge(rand.Intn(100))
		ps = append(ps, p)
	}
	fmt.Println("sort before: ")
	for index, val := range ps{
		fmt.Println("th : " + strconv.Itoa(index), val)
	}

	//排序（对年龄进行排序-升序）
	fmt.Println("sorted after: ")
	sort.Sort(ps)
	for index, val := range ps{
		fmt.Println("th : " + strconv.Itoa(index), val)
	}
}

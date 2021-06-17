package main

import "fmt"

func main() {
	// defer 和 return之间的顺序是先返回值, i=0，后defer
	fmt.Println("return:", test1())
	//i = 2
	fmt.Println("return:", test2())
}

//返回值未命名（此函数的defer不改变return值）
func test1() int {
	var i int
	defer func() {
		i++
		fmt.Println("defer1", i) //作为闭包引用，i=2
	}()
	defer func() {
		i++
		fmt.Println("defer2", i) //作为闭包引用，i=1
	}()
	return i
}


//返回值有命名（此函数的defer可改变return值）
func test2() (i int) { //返回值命名i
	defer func() {
		i++
		fmt.Println("defer1", i)
	}()
	defer func() {
		i++
		fmt.Println("defer2", i)
	}()
	return i
}
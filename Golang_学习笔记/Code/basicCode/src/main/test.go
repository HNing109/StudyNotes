package main

import (
	"bytes"
	"embeded_interface"
	"fmt"
	"sort"
	//"time"
)
import "longparam"

const NUM int = 41

var redis [NUM]int64

func main() {

	var index, min = longparam.GetData(1, 2, 3, 4, 0)
	fmt.Printf("min number :index = %d, val = %d\n", index, min)

	sum := func(i int) int {
		s := 0
		for index := 0; index < i; index++ {
			s += i
		}
		return s
	}(10)
	fmt.Println(sum)

	//var start = time.Now()
	//for index := 0; index < NUM; index++{
	//	fmt.Println(fibonacci(int64(index)))
	//}
	//var end = time.Now()
	//fmt.Println("time = ",end.Sub(start))
	//
	//var start1 = time.Now()
	//for index := 0; index < NUM; index++{
	//	fmt.Println(fibo(int64(index)))
	//}
	//var end1 = time.Now()
	//fmt.Println("time = ",end1.Sub(start1))

	var arr1 [5]int
	var arr2 = new([10]int)
	arr3 := [3]int{1, 2, 3}
	fmt.Printf("arr1 = %T, arr2 = %T, arr3 = %T\n", arr1, arr2, arr3)

	//arr2[1] = 2
	for val := range arr2 {
		fmt.Print(val)
	}
	for val := range arr1 {
		fmt.Print(val)
	}

	silce := new([100]int)[0:100]
	ss := make([]int, 4, 10)
	for val := range silce {
		fmt.Println(val)
	}

	for val := range ss {
		fmt.Println(val)
	}

	var strs = []string{"xx", "sss", "wwww"}
	var buffer bytes.Buffer
	for _, val := range strs {
		buffer.WriteString(string(val))
	}
	fmt.Println(buffer.String())

	qq := make([]int, 10)
	fmt.Println(len(qq), cap(qq), len(qq[10:10]))

	str := "chris"
	ch := []byte(str)
	ch[0] = 'F'
	str1 := string(ch)
	fmt.Println(str, str1)

	fmt.Println(bytes.Compare([]byte(str1), []byte(str1)))

	arr9 := make([]int, 10)
	for index := 0; index < len(arr9)-2; index++ {
		arr9[index] = index
	}
	arr9[3] = 1
	isSorted := sort.IntsAreSorted(arr9)
	fmt.Println(isSorted)
	sort.Ints(arr9)
	fmt.Println(arr9)
	fmt.Println(sort.SearchInts(arr9, 1))
}

func init() {
	fmt.Println("init package: main")
}

func fibonacci(num int64) int64 {
	if redis[num] != 0 {
		return redis[num]
	}
	temp := int64(0)
	if num <= 1 {
		return 1
	} else {
		temp = fibonacci(num-1) + fibonacci(num-2)
	}
	redis[num] = temp
	return temp
}

func fibo(num int64) int64 {
	if num <= 1 {
		return 1
	} else {
		return fibo(num-1) + fibo(num-2)
	}
}

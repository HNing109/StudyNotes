package main

import (
	"bytes"
	"myreflect"
	"reflect"

	//"embeded_interface"
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


	wwwwwww()
	qqq()

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



var i = 5
var str = "ABC"

type Person struct {
	name string
	age  int
}

type Any interface{}

func wwwwwww() {
	var val Any
	val = 5
	fmt.Printf("val has the value: %v\n", val)
	val = str
	fmt.Printf("val has the value: %v\n", val)
	pers1 := new(Person)
	pers1.name = "Rob Pike"
	pers1.age = 55
	val = pers1
	fmt.Printf("val has the value: %v\n", val)
	switch t := val.(type) {
	case int:
		fmt.Printf("Type int %T\n", t)
	case string:
		fmt.Printf("Type string %T\n", t)
	case bool:
		fmt.Printf("Type boolean %T\n", t)
	case *Person:
		fmt.Printf("Type pointer to Person %T\n", t)
	default:
		fmt.Printf("Unexpected type %T", t)
	}

	lambda := func(any interface{}) string{
		switch val := any.(type){
		case bool:
			return "this param type is bool"
		case string:
			return "this param " + val +" type is string"
		default:
			return "unknow param type"
		}
	}

	var str = "chris"
	res := lambda(str)
	fmt.Println(res)

}

func qqq() {
	ref := myreflect.ReflectTest{Name:"chris"}
	value := reflect.ValueOf(ref)
	typeOf := reflect.TypeOf(ref)
	fmt.Println(value, typeOf)

	//获取对象中的属性值
	for index := 0; index < value.NumField(); index++ {
		fmt.Println("file ", index , " : ", value.Field(index))
	}

	//获取对象中的方法
	for index := 0; index < value.NumMethod(); index++ {
		fmt.Println("method ", index, " : ", value.Method(index))
	}
}
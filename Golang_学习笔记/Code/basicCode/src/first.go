package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

//定义全局变量
var globalNum1 int = 10
//可自动根据赋值数据，获取数据类型
var globalStr1 = "chris"

var (
	name string = "FYJ"
	age = 18
	isWoman = true
)

//声明常量：不能使用:=
const Pi = 3.14

func main(){
	fmt.Println("hello chris")
	fmt.Println("this time is : ", time.Now())
	fmt.Println("get rand number: ", rand.Intn(10))
	fmt.Println("Pi = ", math.Pi)
	fmt.Println("add() : ", add(10, 20))
	fmt.Println("sub(): ", subtraction(20, 10))
	var str1  = "xxxxx"
	var str2  = "wwwwww"
	str1, str2 = returnMulString(str1, str2)
	fmt.Println("returnMulString() : ", str1, "----", str2)

	//使用:= 定义变量（替代var），只能在函数内使用
	res1, res2 := split(10, 20)
	fmt.Println("res1 = ", res1, " ;  res2 = ", res2)

	//输出数据的类型、数值
	fmt.Printf("Type = %T, Value = %v \n", name, name)
	fmt.Printf("Type = %T, Value = %v\n", age, age)
	fmt.Printf("Type = %T, Value = %v\n", isWoman, isWoman)

	fmt.Println("sum(): ", getAllData(10))

	fmt.Println("ifTest(): ", ifTest(30,2))

	fmt.Println("switchTest(): ", switchTest("chris"))

	deferTest()
	point()
	structTest()
	arrTest()
	slice()
	mapTest()
	funcValTest()

	closureTest()

	wayTest()
}

func add(num1 int, num2 int) int{
	return num1 + num2
}

func subtraction(num1, num2 int) int {
	return int(math.Abs(float64(num1 - num2)))
}

func returnMulString(str1, str2 string) (string, string){
	return str2, str1
}

func split(num1 int , num2 int) (res1 float32, res2 float32){
	res1 = float32(num1 + num2)
	res2 = float32(num1 - num2)
	return
}

func getAllData(num int) int{
	var sum = 0
	for index := 0; index < num; index++ {
		sum += index
	}
	return sum
}

func ifTest(num1 int , num2 int) int{
	//在if中定义局部变量v，仅可在if-else中使用
	if v := num1 - num2; v < num2{
		return num2
	} else{
		fmt.Println("return V=：", v)
	}
	//这里不能使用v
	return num1 - num2
}

/**
Go 自动提供了每个 case 后面所需的 break 语句
 */
func switchTest(str string) string{
	switch sw := str; sw{
	case "chris":
		return "select " + str
	case "fyj":{
		return  "select " + str
	}
	default:
		return "null"
	}
}

/**
defer推迟调用函数：仅当外层函数执行完后，才执行defer
推迟的函数调用会被压入一个栈中。当外层函数返回时，被推迟的函数会按照后进先出的顺序调用
 */
func deferTest(){
	for index := 0; index < 10; index++{
		defer fmt.Print(index)
	}
}

/**
指针
指针的零值是 nil
 */
func point(){
	var val int = 2
	var p *int
	//间接引用：通过指针p，间接引用cal的值
	p = &val
	fmt.Printf("*p = %v, val = %v \n", *p, val)
	//重定向：通过指针p改变val的值
	*p = 3
	fmt.Printf("*p = %v, val = %v \n", *p, val)
}


/**
结构体
 */
type myStruct struct{
	name string
	age int
}

var(
	//1.创建结构体
	p1 = myStruct{"aa", 2}
	//2.使用默认值创建结构体：默认初始值为0
	p2 = myStruct{name: "bb"}
	p3 = myStruct{}
	//3.创建结构体指针：该指针指向结构体
	p4 = &myStruct{name: "cc", age: 1}
)
func structTest(){
	fmt.Println("p1: ", p1)
	fmt.Println("p2.name: " + p2.name + "  p2.age: " + string(p2.age))
	fmt.Println("p3: ", p3)
	fmt.Println("p4.name: " + p4.name + " p4.age: " + strconv.Itoa(p4.age))
}

/**
数组
 */
func arrTest(){
	//方式1
	var arr1 [10]int
	for index := 0; index < len(arr1); index++{
		arr1[index] = index
	}
	for index := 0; index < len(arr1); index++{
		fmt.Print(" ", arr1[index])
	}
	fmt.Println()

	//方式2：输出[1 2 3 4 5]
	arr2 := [5]int{1,2,3,4,5}
	fmt.Println(arr2)
	fmt.Println()
}

/**
slice切片：提供动态容量的数组
切片的零值是 nil
 */
func slice(){
	//默认未赋值的位置，数据=0
	//数组的创建方式
	//arr := [5]int{1,2,3,4}
	//1.切片的创建方式（效果同上）
	arr := []int{1,2,3,4}

	//2.make创建切片：
	//len=5，cap=5,元素均为0
	arr1 := make([]int, 5)
	fmt.Printf("%d, len = %d, cap = %d\n", arr1, len(arr1), cap(arr1))
	//len=0，cap=5,没有元素的切片
	arr2 := make([]int, 0, 5)
	fmt.Printf("%d, len = %d, cap = %d\n", arr2, len(arr2), cap(arr2))

	//切片: 获取arr索引值为：1~2的数据
	//var sl []int = arr[1: 3]
	var sl = arr[1: 3]

	fmt.Println(sl, arr)
	//更改切片中的数据，会修改对应引用的数据
	sl[0] = 1000
	fmt.Println(sl, arr)

	//遍历切片
	for index, val := range arr{
		fmt.Printf("index = %d, val = %d\n", index, val)
	}
}

/**
映射：类似于java中的hashmap：存放的是键值对
 */
type websites struct{
	name string
	ip string
}

var m = map[string]websites{
	//websites可以省略
	"baidu": websites{
		"baidu",
		"1,2,3,4",
	},
	"google": websites{
		"google",
		"8.7.9.7",
	},
}
func mapTest(){
	fmt.Println(m)
	//获取元素
	fmt.Println(m["baidu"])

	//修改元素
	m["baidu"] = websites{m["baidu"].name, "2.2.2.2"}
	fmt.Println(m["baidu"])

	//查找元素是否存在
	if val, ok := m["baidu"]; ok{
		//删除元素
		delete(m, "baidu")
		fmt.Println("delete val = ", val)
	}
	fmt.Println(m)
}

/*
函数值：
将函数作为参数传入另一个函数中，在该函数中可以之接使用传入的函数，进行运算
 */
//f为传入的函数值
func funcVal(f func(float64, float64) float64, x float64, y float64) float64{
	return f(x, y)
}

func funcValTest(){
	mySqrt := func(x float64, y float64) float64{
		sum := 1.0
		for index := 1; index <= int(y); index++{
			sum = sum * x
		}
		return sum
	}
	//两者等效
	fmt.Println(funcVal(mySqrt, 3, 4))
	fmt.Println(funcVal(math.Pow, 3, 4))
}

/**
闭包
本质：函数返回另一个函数的返回值，并且函数中的局部变量可被缓存、重复使用
 */
func closure() func(int) int{
	//该值可被叠加使用
	sum := 0
	//将x值循环叠加
	return func(x int) int{
		sum += x
		return sum
	}
}

func closureTest(){
	//f1、 f2对应一个闭包：闭包中的数值会一直存在，可以循环叠加
	f1, f2 := closure(), closure()
	for index := 0; index < 10; index++{
		//index即为传入的x值
		fmt.Println(f1(index), f2(-2 * index))
	}
}

/**
方法：即定义结构体的方法
（Go中没有类，使用该方式给结构体增加方法，等同于java中类的方法）
是一类带特殊的 接收者 参数的函数。（接收者的类型定义和方法声明必须在同一包内）
方法接收者在它自己的参数列表内，位于 func 关键字和方法名之间。
Abs 方法拥有一个名为 v，类型为 way 的接收者
 */
type way struct{
	x float64
	y float64
}

/**
使用值接收者（少用）：myWay way，该方式仅仅修改way结构体中数据的副本（退出该函数后，不影响原有的数据）
 */
func (myWay way) ABS() float64{
	if myWay.x - myWay.y >= 0{
		return myWay.x - myWay.y
	} else{
		return myWay.y - myWay.x
	}
}

/**
使用指针接收者（常用）：myWay *way，可以直接改变way结构体中的数据
 */
func (myWay *way) Scale(num float64){
	myWay.x = myWay.x * num
	myWay.y = myWay.y * num
}
//将上述的  方法  重写为  函数
func ScaleFunc(myWay *way, num float64){
	myWay.x = myWay.x * num
	myWay.y = myWay.y * num
}

func wayTest(){
	w := way{
		x: 3,
		y: 4,
	}
	w1 := &way{
		x: 5,
		y: 6,
	}
	//调用 ： 方法
	w.Scale(10)
	//调用 ： 函数   (两者等效)
	//ScaleFunc(&w, 10)

	fmt.Printf("w = %v, w1 = %v \n", w, w1)
	fmt.Println(w.ABS())
}



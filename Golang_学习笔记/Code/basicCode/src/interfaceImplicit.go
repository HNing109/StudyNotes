package main

import (
	"fmt"
)

/**
接口的隐式实现
 */

//接口
type I interface{
	Say()
}

//结构体
type T struct{
	name string
	age int
}

//实现接口中的方法
func (t *T) Say(){
	fmt.Println(t.name)
}

//重写fmt中的String()方法
func (t *T) String() string{
	return fmt.Sprintf("%v :(%d year)", t.name, t.age)
}

//重写Error()方法
func (t *T) Error() string{
	if t.name == ""  || t.age <= 0 {
		return fmt.Sprintf("name or age is error")
	}
	return ""
}

//创建新的对象
func CreateNew(name string, age int) error{
	return &T{name, age}
}

func main()  {
	var impl I = &T{name: "chris", age: 18}
	impl.Say()
	//重写String(), 使用自定义的输出格式
	fmt.Println(impl)
	//重写error()
	if err := CreateNew("xx", -1); err != nil {
		fmt.Println(err)
	} else{
		fmt.Println("success")
	}

}


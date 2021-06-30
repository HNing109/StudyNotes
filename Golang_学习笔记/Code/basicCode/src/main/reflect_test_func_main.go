package main

import "my_reflect"

/**
 * Created by Chris on 2021/6/30.
 */

func main() {
	var instance = new(my_reflect.ReflectFuncTest)
	//反射变量
	instance.ReflectVar(1)
	//反射结构体：使用new、&{}创建的对象都是ptr指针类型
	var student = new(my_reflect.Student)
	student.Name = "chris"
	student.Age = 18
	instance.ReflectStruct(student)
}

package myreflect

import (
	"fmt"
	"reflect"
)
type ReflectTest struct{
	Name string
	age int
}

func (this *ReflectTest)ReflectTypeValue(){
	var param float64 = 3.14
	val := reflect.ValueOf(param)
	typ := reflect.TypeOf(param)

	fmt.Println(val, typ)
	//打印类型
	fmt.Println("type = ", val.Type())
	fmt.Println("kind = ", val.Kind())
	//reflect.ValueOf(param).Float()：打印的就是param的数值
	fmt.Println("value = ", val.Float())

	//通过反射修改对象中的数据
	//需要使用指针，获取对象中的数据
	val_1 := reflect.ValueOf(&param)
	//
	val_1 = val_1.Elem()
	//该参数是否可以设置
	fmt.Println("can set = ", val_1.CanSet())
	val_1.SetFloat(22)
	fmt.Println("set value = ",val_1.Interface() )
}



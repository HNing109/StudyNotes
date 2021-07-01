package my_reflect

import (
	"fmt"
	"reflect"
)

/**
测试reflect包
 * Created by Chris on 2021/6/30.
 */

type ReflectFuncTest struct{

}

/**
对常量进行反射
 */
func(r *ReflectFuncTest) ReflectVar(v interface{}){
	//获取reflect.TypeOf
	rType := reflect.TypeOf(v)
	fmt.Printf("rTpye= %v\n", rType)

	//获取reflect.ValueOf
	rValue := reflect.ValueOf(v)
	//获取reflect.ValueOf的值、类型
	fmt.Printf("rValue = %v, type = %T \n", rValue, rValue)
	//获取参数值
	fmt.Println("rValue.Int() = ", rValue.Int())

	//将reflect.ValueOf 转 interface{}
	rVal2Interface := rValue.Interface()
	fmt.Println("rVal2Interface = ", rVal2Interface)

	//interface{} 转 变量的实际类型  (断言)
	value := rVal2Interface.(int)
	fmt.Println("value = ", value)
}



/**
对结构体进行反射
 */
type Student struct{
	//添加标签
	Name string	`json:"name"`
	Age int		`json:"age"`
}

func(s Student) Say(param string){
	fmt.Println("Student.Say(): ", param)
}

func(r *ReflectFuncTest) ReflectStruct(s interface{}){
	//获取reflect.TypeOf
	rType := reflect.TypeOf(s)
	fmt.Println("rType = ", rType)

	//获取reflect.ValueOf
	rValue := reflect.ValueOf(s)

	//改变数据：通过Elem() 获得指针值指向的元素值对象（即：传入的对象s），因此调用函数时必须传入指针类型的s
	//且对应的结构体属性，首字母必须大写，否则无法获取到此属性
	rValueChange := rValue.Elem().FieldByName("Age")
	//需要进行类型强转：Set只能接收Value类型的数据
	rValueChange.Set(reflect.ValueOf(20))

	//获取Kind
	rTypeKind := rType.Kind()
	RvalueKind := rValue.Kind()
	fmt.Printf("rTypeKind = %v, rValueKind = %v\n", rTypeKind, RvalueKind)

	//将reflect.ValueOf 转 interface{}
	rValue2Interface := rValue.Interface()
	fmt.Printf("value = %v, type = %T\n", rValue2Interface, rValue2Interface)

	//interface{} 转 变量的实际类型   (断言)
	student , ok := rValue2Interface.(*Student)
	if ok {
		fmt.Println(student)
	}
}

/**
获取结构体的属性值、标签，通过反射调用结构体的方法
 */
func(r *ReflectFuncTest) ReflectGetFieldAndMethod(s interface{}){
	rValue := reflect.ValueOf(s)
	rType := reflect.TypeOf(s)
	/******************** 获取结构体的属性值、标签 *************************/
	numVal := rValue.NumField()
	for i := 0; i < numVal; i++{
		//获取属性值
		fmt.Printf("Filed %d: val = %v\n", i, rValue.Field(i))
		//获取标签：为json的字段值
		tagVal := rType.Field(i).Tag.Get("json")
		if tagVal != ""{
			fmt.Printf("Field %d: tag = %v\n", i, tagVal)
		}
	}

	/******************** 获取结构体的方法、并调用 *************************/
	fmt.Println("s method number = ", rValue.NumMethod())
	var params []reflect.Value
	params = append(params, reflect.ValueOf("say reflect get struct function"))
	rValue.Method(0).Call(params)
}

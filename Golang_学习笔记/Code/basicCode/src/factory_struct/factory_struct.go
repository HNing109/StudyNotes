package factory_struct

//将结构体首字母小写，强制用户只能使用NewFactoryStruct方法创建对象
type factoryStruct struct {
	name string
	age int
	address string
}

//等同于java中的构造函数
func NewFactoryStruct(name string, age int, address string) *factoryStruct{
	object := new(factoryStruct)
	object.name = name
	object.age = age
	object.address = address
	return object
}

func (this *factoryStruct) SetName(name string){
	this.name = name
}

func (this *factoryStruct) GetName() string{
	return this.name
}



package embeded_struct

type outer struct{
	Name string
	age int
	//内嵌结构体
	Inner inner
}

type inner struct{
	Name string
	Sex string
	tel string
}

//构造器s
func NewOuter(name string, age int) *outer{
	return &outer{Name: name, age: age}
}

func (this *outer) GetInnerInfo() string {
	return this.Inner.tel
}



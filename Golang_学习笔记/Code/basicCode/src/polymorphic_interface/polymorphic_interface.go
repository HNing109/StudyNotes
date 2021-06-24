package polymorphic_interface

/**
实现接口方法，调用sort.Sort()方法。
*/

type Person struct{
	name string
	age int
}

func(p *Person) SetName(name string){
	p.name = name
}

func(p *Person) SetAge(age int){
	p.age = age
}


//接口的切片（在此切片类型中，实现Interface{}接口，即可调用sort.Sort()）
type PersonSlice []Person

/******************实现Interface接口，用于使用sort.Sort()方法*********************/
func(p PersonSlice) Len()int{
	return len(p)
}

func(p PersonSlice) Less(i, j int) bool{
	return p[i].age < p[j].age
}

func(p PersonSlice) Swap(i, j int){
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}






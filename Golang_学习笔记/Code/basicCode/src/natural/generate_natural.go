package natural

/**
生成素数
 */

type GenerateNatural struct{

}

func(g *GenerateNatural) GetData()chan int{
	ch := make(chan int)
	go func(){
		for index := 2; ; index++{
			ch <- index
			//fmt.Println("GetData ,index = ", index)
		}
	}()
	return ch
}


func(g *GenerateNatural) PrimeFilter(input chan int, prime int)chan int{
	out := make(chan int)
	go func() {
		for{
			if index := <- input; index % prime != 0{
				out <- index
			}
			//fmt.Print("PrimeFilter ,index = ", index)
		}
	}()
	return out
}
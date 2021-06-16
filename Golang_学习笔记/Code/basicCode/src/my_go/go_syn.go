package my_go

import (
	"fmt"
	"time"
)

type Go_syn struct{

}

func(g *Go_syn) SendData(ch chan string){
	arr := []string{"111", "222", "333"}
	for index := 0; index < len(arr); index++{
		time.Sleep(1 * time.Second)
		ch <- arr[index]
	}
	defer close(ch)
}


func (g *Go_syn) GetData(ch chan string){
	//for val := range ch{
	//	fmt.Println(val)
	//}
	for {
		input, open := <-ch
		if !open {
			break
		}
		fmt.Printf("%s ", input)
	}
}
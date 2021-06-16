package my_go

import "fmt"

type Go_select struct{

}

/**
存入数据
 */
func(g *Go_select) Pump1(ch chan interface{}){
	for index := 0; ; index++{
		ch <- index
	}
}

func(g *Go_select) Pump2(ch chan interface{}){
	for index := 0; ; index ++{
		ch <- index
	}
}


/**
取出数据
 */
func(g *Go_select) GetData(ch1, ch2 chan interface{}){
	for{
		select{
		case val := <-ch1:
			fmt.Println("<- ch1 : ",val)
		case val:= <- ch2:
			fmt.Println("<- ch2 : ", val)
		}
	}

}

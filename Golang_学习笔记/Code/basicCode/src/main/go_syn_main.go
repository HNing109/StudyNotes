package main

import "my_go"

func main() {
	var g = new(my_go.Go_syn)
	ch := make(chan string, 10)
	g.SendData(ch)
	g.GetData(ch)
}

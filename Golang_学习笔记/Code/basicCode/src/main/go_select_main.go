package main

import (
	"my_go"
	"time"
)

func main() {
	var g = new(my_go.Go_select)
	var ch1 = make(chan interface{})
	var ch2 = make(chan interface{})
	go g.Pump1(ch1)
	go g.Pump2(ch2)
	go g.GetData(ch1, ch2)
	time.Sleep(1 * time.Second)
}

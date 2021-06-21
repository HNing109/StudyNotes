package main

import "fmt"

var done = make(chan bool)
var msg string

func aGoroutine() {
	msg = "xxx"
	close(done)
}

func main() {
	ch := func() <-chan int {
		ch := make(chan int)
		go func() {
			for i := 0; ; i++ {
				ch <- i
			}
		} ()
		return ch
	}()

	for v := range ch {
		fmt.Println(v)
		if v == 5 {
			break
		}
	}
}
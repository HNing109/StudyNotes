package main

var done = make(chan bool)
var msg string

func aGoroutine() {
	msg = "xxx"
	close(done)
}

func main() {
	go aGoroutine()
	<-done
	println(msg)
}
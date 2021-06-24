package main

import "fmt"

var done = make(chan bool)
var msg string

func aGoroutine() {
	msg = "xxx"
	close(done)
}

type Q struct{
	name string
	age int
}

func(q *Q) aa(n string){
	q.name = n
}

func (q Q) bb(n string){
	q.name = n
}

type human struct{
	name string
}

func(h *human) say(){
	fmt.Println("xxxx")
}

type man struct{
	human
	age int
}

func main() {
	var q = new(Q)
	q.name = "hhh"
	fmt.Println(q)
	q.aa("xxx")
	fmt.Println(q)
	q.bb("dddd")
	fmt.Println(q)

	num := 10
	fmt.Printf("%x\n",num)
	fmt.Printf("%p",&num)

	var p = new(man)
	p.name = "chris"
	fmt.Println(p)
	p.say()
}
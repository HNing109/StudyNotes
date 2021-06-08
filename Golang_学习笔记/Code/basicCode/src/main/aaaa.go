package main

import "fmt"

var a = "G"

func main() {
LABEL1:
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5; j++ {
			if j == 4 {
				continue LABEL1
			}
			fmt.Printf("i is: %d, and j is: %d\n", i, j)
		}
	}
}

func init() {
	fmt.Println("init package: main")
}
func n() { print(a) }

func mm() {
	a = "O"
	print(a)
	aa := 2222.45678765
	fmt.Printf("%x", aa)
}
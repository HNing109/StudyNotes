package main

import (
	"errors"
	"fmt"
	"math"
	"os"
)

//自定义错误
var errorNoNegative error = errors.New("Number is negative")

func MySqrt(number float64)(float64, error){
	if number < 0 {
		return -1, errorNoNegative
	}
	return math.Sqrt(number), nil
}

func main(){
	number := -2.0
	//错误处理
	if val, err := MySqrt(number); err != nil{
		fmt.Println(err)
	} else{
		fmt.Printf("sqrt(%f) = %f",number, val)
	}

	user := os.Getenv("USER")
	if user == ""{
		panic("Unkown user")
	} else{
		fmt.Println(user)
	}

}
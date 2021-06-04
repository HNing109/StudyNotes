package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64  //创建一个新的类型

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number:  %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	//返回错误：小于0时，无法进行开根号操作
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	return math.Sqrt(x), nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
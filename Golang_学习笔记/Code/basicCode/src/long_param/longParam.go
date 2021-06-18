package long_param

import "fmt"

/*
...  ：表示slices是一个不定长度的参数
 */
func GetData(arr ...int) (index int, val int){
	//代码追踪
	defer un(trace("GetData() function"))
	if len(arr) == 0{
		return -1, -1
	}
	var min = arr[0]
	var minIndex = 0
	for index, val := range arr {
		if val < min{
			min = val
			minIndex = index
		}
	}
	return minIndex, min

}

/**
代码追踪
 */
func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}

func un(s string) {
	fmt.Println("leaving:", s)
}

func main() {
	var index, min = GetData(1,2,3,4,0)
	fmt.Printf("min number :index = %d, val = %d", index, min)
}

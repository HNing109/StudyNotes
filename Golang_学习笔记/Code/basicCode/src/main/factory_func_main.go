package main

import (
	factory "factory_function"
	"fmt"
)

func main() {
	//使用工厂函数：创建新的函数
	addBmp := factory.AddSuffix(".bmp")
	addJpg := factory.AddSuffix(".jpg")


	//动态添加后缀
	bmpFile := addBmp("picture_1")
	jpgFile := addJpg("picture_2")
	fmt.Println(bmpFile, jpgFile)


}
package main

import (
	"flag"
	"fmt"
)

//定义变量所对应的命令行参数：仅当命令行有输入参数 -n  时，NewLing = true
var NewLine = flag.Bool("n", false, "print newline")

const (
	Space   = " "
	Newline = "\n"
)

//测试输入： -u chris -p 22
func main() {
	//定义变量，接收命令行输入的参数
	var user string
	var pwd string
	var port int

	//指定变量对应的命令行输入参数、默认值、说明
	flag.StringVar(&user, "u", "admin", "用户名，默认为admin")
	flag.StringVar(&pwd, "pwd", "ADMIN", "密码，默认为ADMIN")
	flag.IntVar(&port, "p", 88, "端口号， 默认为88")

	//获取命令行输入的参数：从os.Args[1:]中解析注册的flag，必须在所有flag都注册好时且未访问值时执行。
	flag.Parse()

	//打印参数
	fmt.Printf("user = %s, pwd = %s, port = %d", user, pwd, port)
}

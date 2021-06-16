package http_pkg

import (
	"fmt"
	"log"
	"net/http"
)

type HttpPkg struct{

}
/**
http.ResponseWriter：通过此对象进行数据输出
http.Request：网页服务器发送的响应对象
 */
func (h *HttpPkg) HelloWordServer(w http.ResponseWriter, req *http.Request){
	fmt.Println("Hello world server")
	//[1:] ： 从1开始，是为了滤除根目录/
	fmt.Fprintf(w, "hello, " + req.URL.Path[1:])
}

func (h * HttpPkg) Test(){
	//访问的URL、对应的处理函数
	http.HandleFunc("/", h.HelloWordServer)
	//监听本地端口8080
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil{
		log.Fatal("ListenAndServer: ", err.Error())
	}
}

package rpc_pkg

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

/**
RPC服务器：开启RPC的功能
 */
type RpcServer struct{

}
func(r *RpcServer) StartServer(){
	log.Println("Starting Rc-server......")
	//创建rpc调用对象
	rpcObject := new(Args)
	//注册RPC调用的对象
	rpc.Register(rpcObject)
	rpc.HandleHTTP()
	//以tcp方式监听本地端口
	listen, err := net.Listen("tcp", "localhost:8888")
	if err != nil{
		log.Fatal("Starting Rpc-server Error :", err)
	}
	//使用协程处理监听到的http请求
	go http.Serve(listen, nil)
	time.Sleep(1000 * time.Second)
}

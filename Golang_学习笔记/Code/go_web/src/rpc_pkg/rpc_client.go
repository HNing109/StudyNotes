package rpc_pkg

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
)
/**
RCP客户端：模拟远程调用RCP
 */
type RpcClient struct{

}
//远程RPC服务器的地址
const RpcServerHost = "localhost:"
const RpcPort = 8888

func (r *RpcClient) StartClient(){
	log.Println("Starting Rc-client......")
	//以TCP方式拨号，连接RPC服务器
	client, errServer := rpc.DialHTTP("tcp", RpcServerHost + strconv.Itoa(RpcPort))
	if errServer != nil{
		log.Fatal("Starting Rpc-client error: ", errServer)
	}
	//创建相应的RPC调用的对象
	args := NewRpcObject(3,4)
	var reply int
	//调用远程的Args.Multiply()方法，传入args对象作为参数，响应结果存入reply
	errClient := client.Call("Args.Multiply", args, &reply)
	if errClient != nil{
		log.Fatal("Args Error: ", errClient)
	}
	fmt.Printf("call rpc success: args.n = %d, args.M = %d, reply = %d", args.N, args.M, reply)
}

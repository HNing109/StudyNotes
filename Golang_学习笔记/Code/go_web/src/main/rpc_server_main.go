package main

import "rpc_pkg"

func main(){
	var rpcServer = new(rpc_pkg.RpcServer)
	rpcServer.StartServer()
}

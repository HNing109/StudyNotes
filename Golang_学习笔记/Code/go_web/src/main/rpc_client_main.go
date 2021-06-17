package main

import (
	"rpc_pkg"
)

func main(){
	var rpcClient = new(rpc_pkg.RpcClient)
	rpcClient.StartClient()
}

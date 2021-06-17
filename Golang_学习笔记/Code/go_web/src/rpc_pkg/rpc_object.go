package rpc_pkg

/**
Rpc：调用的对象
 */
type Args struct{
	N, M int
}

func(a *Args) Multiply(args * Args, reply *int) error{
	*reply = args.N * args.M
	return nil
}

func NewRpcObject(N, M int) *Args{
	obj := new(Args)
	obj.N = N
	obj.M = M
	return obj
}


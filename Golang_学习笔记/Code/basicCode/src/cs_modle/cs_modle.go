package cs_modle

import "fmt"

/**
模拟c/s请求处理
 */

type Request struct {
	a, b      int
	replyc    chan int // reply channel inside the Request
}

type binOp func(a, b int) int

func run(op binOp, req *Request) {
	req.replyc <- op(req.a, req.b)
}

func server(op binOp, service chan *Request) {
	for {
		req := <-service // requests arrive here
		// start goroutine for request:
		go run(op, req)  // don’t wait for op to complete
	}
}

func startServer(op binOp) chan *Request {
	reqChan := make(chan *Request)
	go server(op, reqChan)
	return reqChan
}

func CSTestFunc() {
	adder := startServer(func(a, b int) int { return a + b })
	const N = 100
	var reqs [N]Request
	for i := 0; i < N; i++ {
		req := &reqs[i]
		req.a = i
		req.b = i + N
		req.replyc = make(chan int)
		adder <- req  // adder is a channel of requests
	}
	// checks:
	for i := N - 1; i >= 0; i-- {
		// doesn’t matter what order
		if <-reqs[i].replyc != N+2*i {
			fmt.Println("fail at", i)
		} else {
			fmt.Println("Request ", i, "is ok!")
		}
	}
	fmt.Println("done")
}
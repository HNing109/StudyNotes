package random

/**
使用select生成随机0、1序列
 */
type RandomSerial struct{

}

func(r *RandomSerial) GetRandomSerial(ch chan int) int{
	go func(){
		select{
			case ch <- 0:
			case ch <- 1:
		}
	}()
	return <-ch
}


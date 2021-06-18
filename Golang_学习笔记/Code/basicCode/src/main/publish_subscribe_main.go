package main

import (
	"csp"
	"fmt"
	"strings"
	"time"
)

func main() {
	p := csp.NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	subscribe_all_topic := p.SubscribeAllTopic()
	topic := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "topic")
		}
		return false
	})
	//发布消息
	p.Publish("hello, world!")
	p.Publish("hello, topic!")

	go func() {
		for  msg := range subscribe_all_topic {
			fmt.Println("subscribe_all_topic:", msg)
		}
	} ()

	go func() {
		for  msg := range topic {
			fmt.Println("topic:", msg)
		}
	} ()

	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}

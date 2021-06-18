package csp

import (
	"sync"
	"time"
)

/*
发布者-订阅者模式
 */

type(
	//订阅者：是一个管道
	subscriber chan interface{}
	//主题：是一个过滤器
	topicFunc func(val interface{}) bool
)

type Publisher struct{
	//锁
	mu sync.RWMutex
	//订阅队列缓存大小
	buffer int
	//超时时间
	timeout time.Duration
	//订阅者信息
	subscribes map[subscriber]topicFunc

}

//新建发布者
func NewPublisher(publisherTimeout time.Duration, buffer int) *Publisher{
	return &Publisher{
		buffer: buffer,
		timeout:publisherTimeout,
		subscribes: make(map[subscriber]topicFunc),
	}
}

//添加一个订阅者，使用过滤器订阅某个主题
func(p *Publisher) SubscribeTopic(topic topicFunc) subscriber{
	ch := make(chan interface{}, p.buffer)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribes[ch] = topic
	return ch
}

//添加订阅者：订阅所有主题
func(p *Publisher) SubscribeAllTopic() subscriber{
	return p.SubscribeTopic(nil)
}

//退订
func (p Publisher) Evict(sub chan interface{}){
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subscribes, sub)
	close(sub)
}

// 发布一个主题
func (p *Publisher) Publish(val interface{}) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribes {
		wg.Add(1)
		go p.sendTopic(sub, topic, val, &wg)
	}
	wg.Wait()
}

// 关闭发布者对象，同时关闭所有的订阅者管道。
func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for sub := range p.subscribes {
		delete(p.subscribes, sub)
		close(sub)
	}
}


// 发送主题，可以容忍一定的超时
func (p *Publisher) sendTopic(
								sub subscriber,
								topic topicFunc,
								val interface{},
								wg *sync.WaitGroup, ) {
	defer wg.Done()
	if topic != nil && !topic(val) {
		return
	}

	select {
	case sub <- val:
	case <-time.After(p.timeout):
	}
}


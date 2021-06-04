package main

import (
	"fmt"
	"sync"
	"time"
)

/**
测试互斥锁
 */
type mutex struct{
	//map：存放键值对
	myMap map[string]int
	//互斥锁
	mux sync.Mutex
}

/**
增加key对应的val
 */
func (m *mutex) IncVal(key string){
	//上锁
	m.mux.Lock()
	if val, ok := m.myMap[key]; ok{
		m.myMap[key] = val + 1
	}
	//释放锁
	m.mux.Unlock()
}

/**
获取key对应的value
 */
func (m *mutex) getValue(key string) int{
	m.mux.Lock()
	var res int
	if val, ok := m.myMap[key]; ok{
		res = val
	} else{
		res = -1
	}
	m.mux.Unlock()
	return res
}

func main() {
	exam := mutex{
		myMap: make(map[string]int),
		mux:   sync.Mutex{},
	}
	exam.myMap["chris"] = 0
	for index := 0; index < 10; index++{
		go exam.IncVal("chris")
	}
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("key = chris , value = ", exam.myMap["chris"])
}
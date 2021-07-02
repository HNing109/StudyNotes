package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)
/**
 * Created by Chris on 2021/7/1.
 */

//redis连接池
var pool *redis.Pool

//启动时，自动初始化连接池
func init(){
	pool = &redis.Pool{
		// 初始化链接的代码， 链接哪个ip的redis
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", "192.168.154.19:4000")
		},
		TestOnBorrow:    nil,   //提供一个方法，用来诊断一个连接的健康状态
		MaxIdle:         8,		//最大空闲链接数
		MaxActive:       0,		//表示和数据库的最大链接数， 0 表示没有限制
		IdleTimeout:     100,	//最大空闲时间
		Wait:            false, //如果Wait被设置成true，则Get()方法将会阻塞
		MaxConnLifetime: 0,		//关闭早于此时间的连接。如果为0，则连接池不会根据年龄关闭连接。
	}
}

func main(){
	//获取连接
	conn := pool.Get()
	//退出程序时，关闭连接
	defer conn.Close()

	//存入数据：string
	_, err := conn.Do("Set", "name", "chris_zhang")
	if err != nil{
		fmt.Println("set error = ", err)
		return
	}

	//取出数据：string
	getData, err := redis.String(conn.Do("Get", "name"))
	if err != nil{
		fmt.Println("get error = ", err)
		return
	}
	fmt.Println("getData = ", getData)

	//写入数据：key-val数据： [string] int/string
	_, err = conn.Do("HMSet", "user01", "name", "FJY", "age", 19)
	if err != nil {
		fmt.Println("HMSet  err=", err)
		return
	}

	//读取数据：key-val数据： [string] int/string
	r, err := redis.Strings(conn.Do("HMGet","user01", "name", "age"))
	if err != nil {
		fmt.Println("HMGet  err=", err)
		return
	}
	for i, v := range r {
		fmt.Printf("r[%d]=%s\n", i, v)
	}
}
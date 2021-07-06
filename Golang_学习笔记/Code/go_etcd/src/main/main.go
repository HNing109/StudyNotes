package main

/**
 * Created by Chris on 2021/7/6.
 */

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {

	//连接
	cli, err := clientv3.New(clientv3.Config{
		//此处填写--listen-client-urls监听地址（各个ETCD节点监听地址）
		Endpoints:   []string{
			"192.168.83.130:2379",
			"192.168.83.130:22379",
			"192.168.83.130:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}
	fmt.Println("connect success")
	defer cli.Close()


	ctx, _ := context.WithCancel(context.TODO())

	//存数据
	startTime :=time.Now()
	_, err = cli.Put(ctx, "name", "chris")
	//操作完毕，取消连接
	// cancel()
	endTime :=time.Now()
	fmt.Println("put耗时", endTime.Sub(startTime))
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	//设置超时为5秒
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	//取数据
	startTime = time.Now()
	resp, err := cli.Get(ctx, "name")
	fmt.Println("get 耗时:",time.Now().Sub(startTime))
	// 	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}

	//打印数据
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

}

 
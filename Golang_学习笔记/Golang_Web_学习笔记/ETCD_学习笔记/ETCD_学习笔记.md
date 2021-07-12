# 1、磁盘性能评估

- 工具：fio

- 下载地址：https://github.com/axboe/fio/releases/tag/fio-3.27

- 安装命令：

  - ./configure
  - make
  - make install

- 测试命令：

  该命令运行在测试数据存储文件夹所在目录

  - fio --rw=write --ioengine=sync --fdatasync=1 --directory=test-data-fio --size=22m --bs=2300 --name=mytest
  - 参数：
    - --directory=测试数据存储的位置

- 性能指标： [Prometheus](https://prometheus.io/) 指标中的wal_fsync_duration_seconds，该指标的**第 99 个百分位数应小于 10 毫秒**，即可认为存储设备足够快。

  eg：

  ```shell
  #测试结果
  mytest: (g=0): rw=write, bs=(R) 2300B-2300B, (W) 2300B-2300B, (T) 2300B-2300B, ioengine=sync, iodepth=1
  fio-3.27
  Starting 1 process
  Jobs: 1 (f=1)
  mytest: (groupid=0, jobs=1): err= 0: pid=12041: Wed Jun 30 15:01:07 2021
    write: IOPS=8240, BW=18.1MiB/s (19.0MB/s)(22.0MiB/1217msec); 0 zone resets
      clat (usec): min=2, max=3090, avg=44.71, stdev=48.88
       lat (usec): min=2, max=3090, avg=44.87, stdev=48.88
      clat percentiles (usec):
       |  1.00th=[    3],  5.00th=[    4], 10.00th=[    4], 20.00th=[    4],
       | 30.00th=[    4], 40.00th=[    4], 50.00th=[   71], 60.00th=[   72],
       | 70.00th=[   73], 80.00th=[   74], 90.00th=[   85], 95.00th=[   88],
       | 99.00th=[  115], 99.50th=[  124], 99.90th=[  221], 99.95th=[  375],
       | 99.99th=[  486]
     bw (  KiB/s): min=18355, max=18736, per=100.00%, avg=18545.50, stdev=269.41, samples=2
     iops        : min= 8172, max= 8342, avg=8257.00, stdev=120.21, samples=2
    lat (usec)   : 4=41.72%, 10=1.93%, 20=0.16%, 50=0.03%, 100=54.25%
    lat (usec)   : 250=1.81%, 500=0.08%
    lat (msec)   : 4=0.01%
    
    #################### （查看第99次sync percentiles）此处的时间单位：us #################### 
    fsync/fdatasync/sync_file_range:
      sync (usec): min=52, max=552, avg=73.71, stdev=14.16
      sync percentiles (usec):
       |  1.00th=[   67],  5.00th=[   69], 10.00th=[   70], 20.00th=[   71],
       | 30.00th=[   71], 40.00th=[   72], 50.00th=[   72], 60.00th=[   73],
       | 70.00th=[   74], 80.00th=[   75], 90.00th=[   77], 95.00th=[   84],
       | 99.00th=[  116], 99.50th=[  130], 99.90th=[  343], 99.95th=[  367],
       | 99.99th=[  453]
    cpu          : usr=0.00%, sys=63.40%, ctx=21292, majf=0, minf=16
    IO depths    : 1=200.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
       submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
       complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
       issued rwts: total=0,10029,0,0 short=10029,0,0,0 dropped=0,0,0,0
       latency   : target=0, window=0, percentile=100.00%, depth=1
  
  Run status group 0 (all jobs):
    WRITE: bw=18.1MiB/s (19.0MB/s), 18.1MiB/s-18.1MiB/s (19.0MB/s-19.0MB/s), io=22.0MiB (23.1MB), run=1217-1217msec
  
  Disk stats (read/write):
    sda: ios=5211/9278, merge=0/0, ticks=303/549, in_queue=0, util=91.86%
  ```

  

# 2、ETCD基准测试

- 工具：etcd/tools/benchmark
- 下载：go get go.etcd.io/etcd/v3/tools/benchmark



# 3、ETCD入门

## 3.1、ETCD部署要求

### 3.1.1、系统要求

- 系统

  amd64-linux

- 存储设备：

  80GB以上的SDD。ETCD数据需要写入磁盘中，因此为增加数据读写速度，ETCD数据存储设备需要使用SSD，且需要验证存储设备的性能。可使用fio工具进行检测。

- 处理器

  双核以上

- RAM

  最大8GB，最小2GB。因为ETCD需要强制设置一个默认的RAM大小（2GB）
  
  - 设置占用RAM的大小：--quota-backend-bytes



### 3.1.2、部署原则

- ETCD集群的数量必须为奇数

  因为ETCD主机出现宕机时，需要从其他节点中选取leaer。若集群节点数量为偶数，则可能会出现某个节点所得票数相等，无法选举出leader。

- ETCD集群最大数量 ≤ 7

  虽然增加集群的数量，可以增加ETCD的容错性，但是集群数量越大数据的写入性能就会下降，因为需要将数据复制至更多的节点中。推荐集群数量 = 5，可容忍2个节点故障

- ETCD可用于跨区域、跨数据中心部署

  但是由于集群中的节点分布在不同的网段中，会增加数据请求的延时。并且，节点之间的数据复制，将会占用带宽。

- 删除节点的操作步骤：

  此顺序不能颠倒。

  - 首先删除节点
  - 然后添加新的节点





## 3.2、安装ETCD

### 3.2.1、使用二进制文件安装

- release文件下载：https://github.com/etcd-io/etcd/releases/
- 解压之后，将etcd二进制文件存入/usr/local/bin中
- 测试是否安装完成（在任意路径下）：etcd --version



### 3.2.2、使用源码安装

- 前提：Go1.16以上的版本

- 安装命令：

  - 下载：git clone -b v3.5.0 https://github.com/etcd-io/etcd.git

  - 进入目录etcd，执行构建脚本：./build.sh

  - 构建后生成的二进制文件位于`bin`目录下，将`bin`目录的完整路径添加到系统环境变量中：

    export PATH="$PATH:pwd/bin"

  - 测试是否安装成功：etcd --version



## 3.3、启动ETCD

### 3.3.1、启动单个ETCD

- 启动： 

  - 命令：etcd， 或者./etcd

- 查看启动的详情：

  - 命令：etcdctl member list

  - 结果：默认情况下：client address  =  http://127.0.0.1:2379

  | ID               | STATUS  | NAME   | PEER ADDRS             | CLIENT ADDRS          | IS LEARNER |
  | ---------------- | ------- | ------ | ---------------------- | --------------------- | ---------- |
  | 8211f1d0f64f3269 | started | infra1 | http://127.0.0.1:12380 | http://127.0.0.1:2379 | false      |



### 3.3.2、启动ETCD集群

#### 3.3.2.1、部署在Ubuntu系统

- **准备：**

  - **安装goreman**

    - 命令：go get github.com/mattn/goreman 

    - goreman 程序来控制基于 Procfile 的配置文件的ETCD程序。通过goreman + Procfile配置文件，可以快速启动ETCD集群。

      

  - **Procfile配置文件**

    - 先准备Procfile配置文件，在该文件中写入需要启动节点的信息。

      

    - **启动参数说明：**

      1. --name：

         etcd 集群中的节点名，不可重复

      2. --data-dir：

         存放数据目录，节点ID，集群ID，Snapshot文件，集群初始化配置，WAL 文件

      3. **--listen-peer-urls：**

         本member使用，用于监听其他member发送信息的地址。ip为全0代表监听本member侧所有接口。(应用场景：如选举，数据同步等)

      4. **--initial-advertise-peer-urls** ：

         其他member通过该地址与本member交互信息。该参数的value一定要同时配置到--initial-cluster参数中。

      5. **<font color='red'>--listen-client-urls</font>** ：

         监听外部客户端的请求（即：**etcd节点对外提供服务的地址**，**可添加多个，使用,分隔**）。

         - **注意事项：**

           <font color='red'>--listen-client-urls，必须配置一个url = **http://本机网卡的IP+端口**，远程客户端才能访问本机的ETCD数据库。</font>

           否则会出现错误：

           ```shell
           grpc: Conn.resetTransport failed to create client transport: connection error: desc = "transport: dial tcp 192.168.83.130:2379: connectex: No connection could be made because the target machine actively refused it.";
           ```

           

         **eg：**etcdctl工具、curl请求、go编写的访问etcd服务代码

         ```shell
         #启动etcd节点时，建议的配置：
         #保证本地etcdctl工具、远程客户端均可请求访问etcd服务端
         --listen-client-urls http://127.0.0.1:2379,http://本地网卡IP:2379 --advertise-client-urls http://127.0.0.1:2379
         ```

         

      6. **<font color='red'>--advertise-client-urls</font>** ：

         该节点的客户端（广播客户端），用于此url与其他etcd节点通信（**可添加多个，使用,分隔**）（**注意：此url不是用于远程客户端访问该节点的**）。

         **此处，不能配置为  空、localhost。**

      7. --initial-cluster-token：

         集群的token值，设置该值后集群将生成唯一id,并为每个节点也生成唯一id,当使用相同配置文件再启动一个集群时，只要该token值不一样，ETCD集群就不会相互影响

      8. **--initial-cluster** ：

         是集群中所有--initial-advertise-peer-urls 的合集，即：**此处配置的url必须和--initial-advertise-peer-urls配置的一致** 

      9. --initial-cluster-state new：

         新建集群的标志

         

    - **参考官方配置文件：**https://github.com/etcd-io/etcd/blob/main/Procfile
      
      - 客户端通信地址为：
        - 各个节点访问其他节点的客户端地址：http://127.0.0.1:2379 , http://127.0.0.1:22379 ,  http://127.0.0.1:32379
        - **远程客户端访问ETCD节点的地址**：http://192.168.83.130:2379，http://192.168.83.130:22379，http://192.168.83.130:32379
      
      ```shell
      # Use goreman to run `go get github.com/mattn/goreman`
      # Change the path of bin/etcd if etcd is located elsewhere
      
      # ETCD1
      etcd1: etcd --name infra1 --listen-client-urls http://127.0.0.1:2379,http://192.168.83.130:2379 --advertise-client-urls http://127.0.0.1:2379 --listen-peer-urls http://127.0.0.1:12380 --initial-advertise-peer-urls http://127.0.0.1:12380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
      
      # ETCD2
      etcd2: etcd --name infra2 --listen-client-urls http://127.0.0.1:22379,http://192.168.83.130:22379 --advertise-client-urls http://127.0.0.1:22379 --listen-peer-urls http://127.0.0.1:22380 --initial-advertise-peer-urls http://127.0.0.1:22380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
      
      # ETCD3
      etcd3: etcd --name infra3 --listen-client-urls http://127.0.0.1:32379,http://192.168.83.130:32379 --advertise-client-urls http://127.0.0.1:32379 --listen-peer-urls http://127.0.0.1:32380 --initial-advertise-peer-urls http://127.0.0.1:32380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
      #proxy: etcd grpc-proxy start --endpoints=127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379 --listen-addr=127.0.0.1:23790 --advertise-client-url=127.0.0.1:23790 --enable-pprof
      
      # A learner node can be started using Procfile.learner
      ```

  

- **启动集群：**

  - 命令：goreman -f **Procfile** start 

  - Procfile为之前路径中编写的配置文件（因此，**该命令必须在Procfile文件所在路径中使用**）

    

- **使用etcdctl客户端工具**

  - 配置所使用的etcdctl工具命令版本

    - 命令：export ETCDCTL_API=3 
    - 默认使用 v2 API 来和 etcd 数据库通信，这是为了向后兼容 etcdctl。若需要使用v3 的etcdctl API，就必须配置环境变量：ETCDCTL_API 设置为版本3。

    

  - 查看集群中的信息

    etcdctl --write-out=table --endpoints=localhost:2379 member list

    | ID               | STATUS  | NAME   | PEER ADDRS             | CLIENT ADDRS           | IS LEARNER |
    | ---------------- | ------- | ------ | ---------------------- | ---------------------- | ---------- |
    | 8211f1d0f64f3269 | started | infra1 | http://127.0.0.1:12380 | http://127.0.0.1:2379  | false      |
    | 91bc3c398fb3c146 | started | infra2 | http://127.0.0.1:12380 | http://127.0.0.1:22379 | false      |
    | fd422379fda50e48 | started | infra3 | http://127.0.0.1:12380 | http://127.0.0.1:32379 | false      |




#### 3.3.2.2、部署在Docker中





## 3.4、ETCD工具

### 3.4.1、etcdctl

这是命令行客户端工具，可直接对etcd数据库进行操作。

#### 3.4.1.1、<font color='red'>命令的附加选项</font>

- 指定节点启动的地址、端口

  --endpoints=localhost:2379  

- 指定范围

  key1 keyn

  eg：etcdctl get foo1 foo4    ：注意，获取区间间 **[key1, key4)**，第4个键值无法被获取到

- 指定前缀（可实现模糊匹配）

  --prefix=xxx

- 限制数量

  --limit=n

- 大于等于某个key对应的value

  --from-key key1

  eg：etcdctl get --from-key foo

- 查询历史版本

  --rev=n  ：查找版本为n的键值对

- 数据显式格式

  - -w=json

    以json的格式显示数据（无论数据是否存在，均可以显示）

    eg：etcdctl get foo -w=json 。



#### 3.4.1.2、节点相关命令

- **查看节点信息**

  - 方式1：

    etcdctl member list

    查看成员信息（精简版）

  - 方式2：

    etcdctl --write-out=table --endpoints=localhost:2379 member list 

    查看节点localhost:2379中的成员信息



- **停止节点**

  goreman run stop **节点名**

  - **eg：**goreman run stop etcd2 

  - **注意**：这个命令无法停止etcd，最后还是用ps命令找出pid，然后kill 

    - ps -ef | grep etcd | grep endpoints的地址 

    - kill PID号

      

- **重启节点**

  goreman run restart **节点名**

  **eg：**goreman run restart etcd2



#### 3.4.1.3、数据相关命令

- **存数据**

  - 命令：

    - 方式1：etcdctl  put **键 值**

    - 方式2：etcdctl --endpoints=localhost:2379 put **键 值**

      

- **取数据**

  - 方式1：etcdctl  get **键 值**
  - 方式2：etcdctl --endpoints=localhost:2379 get **键**

  

- **删除键值对**

  etcdctl del 键

  

- **<font color='red'>监控数据</font>**

  etcdctl watch key

  实时监控key所对应值的变化

  - **范围监控**

    etcdctl watch key1 keyn

  - **监控指定的多个key-val**

    etcdctl watch -i

    watch key1

    watch key3

  

- **压缩数据**

  etcdctl compact 版本号

   压缩之后，版本号之前的数据均无法访问，即：不可以获通过etcdctl get --rev=版本号-1 ，获取历史版本为（版本号-1）的数据

  

- **<font color='red'>存活时间（租约）</font>**

  - **设置存活时间**

    - 新建一个存活时间（单位：秒）

      etcdctl lease grant 时间

      - eg：etcdctl lease grant 10

      - 结果：lease 32695410dcc0ca06 granted with TTL(10s)    

        **存活时间配置id：**32695410dcc0ca06 ，需要使用该id取设置key的存活时间

        

    - 给指定key设置存活时间

      etcdctl put --lease=存活时间配置id key value

      

  - **撤销存活时间**

    etcdctl lease revoke 存活时间配置id

    注意：当存活时间撤销之后，配置了该存活时间的key-val将被删除。

    

  - **维持存活时间**

    - 新建一个存活时间配置

      etcdctl lease grant 时间

    - 维持旧的存活时间

      etcdctl lease keep-alive 新建的存活时间配置id

    

  - **查看存活时间的信息**

    - 查看存活时间配置

      etcdctl lease timetolive 存活时间配置id

    - 查看使用了存活时间配置的key-value

      etcdctl lease timetolive --keys 存活时间配置id

- 



### 3.4.2、grpc-gateway

网关接收 ETCD的 protocol buffer 消息定义的 JSON mapping 。key 和 value 字段被定义为 byte 数组，因此必须在 JSON 中以 base64 编码。任何HTTP/JSON 客户端都可使用 curl访问ETCD数据库。

**需要使用到v3alpha/kv/中的组件。**

- 存key-value

  curl -L http://localhost:2379/v3alpha/kv/put \ 

  -X POST -d '{"key": "Zm9v", "value": "YmFy"}' 

  

- 读取key-value

  curl -L http://localhost:2379/v3alpha/kv/range \ 

  -X POST -d '{"key": "Zm9v"}' 

  

- 监控数据

  curl http://localhost:2379/v3alpha/watch \ 

  -X POST -d '{"create_request": {"key":"Zm9v"} }' & 

  

- 事务

  curl -L http://localhost:2379/v3alpha/kv/txn \ 

  -X POST \ 

  -d '{"compare":[{"target":"CREATE","key":"Zm9v","createRevision":"2"}],"success":[ 

  {"requestPut":{"key":"Zm9v","value":"YmFy"}}]}' 

  

- 搭建认证服务

  - 创建 root 用户 

    curl -L http://localhost:2379/v3alpha/auth/user/add \ 

    -X POST -d '{"name": "root", "password": "pass"}' 

  - 创建 root 角色 

    curl -L http://localhost:2379/v3alpha/auth/role/add \ 

    -X POST -d '{"name": "root"}' 

  - 授予 root 角色 

    curl -L http://localhost:2379/v3alpha/auth/user/grant \ 

    -X POST -d '{"user": "root", "role": "root"}' 

  - 开启认证 

    curl -L http://localhost:2379/v3alpha/auth/enable -X POST -d '{}' 

  - root 用户 获取认证 token 

    curl -L http://localhost:2379/v3alpha/auth/authenticate \ 

    -X POST -d '{"name": "root", "password": "pass"}' 

  - 为认证 token 设置 Authorization header，以便在获取键时使用认证证书

    curl -L http://localhost:2379/v3alpha/kv/put \ 

    -H 'Authorization : sssvIpwfnLAcWAQH.9' \ 

    -X POST -d '{"key": "Zm9v", "value": "YmFy"}' 



## 3.5、go代码远程访问ETCD

- **启动ETCD集群**
  - 前提：注意启动节点时，所配置的<font color='red'>--listen-client-urls</font>，此处必须有ETCD所在机器的IP（网卡地址），否则无法远程访问ETCD
  - 命令：goreman -f Procfile start     （Procfile为配置文件：包含各个节点的配置详情）
  - 如果还是出现refuse it的错误：
    - 方式1：ETCD的启动配置中，增加--listen-client-urls的监听地址 
    - 方式2：
      - 查看已经开放的端口：netstat -tlnp
      - 开放远程访问端口：ufw allow 端口号



- **代码**

  ```go
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
  		//此处IP为：--listen-client-urls监听地址（各个ETCD节点监听地址）
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
  ```

  
























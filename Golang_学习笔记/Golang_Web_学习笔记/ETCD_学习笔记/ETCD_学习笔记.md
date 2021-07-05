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

- **准备：**

  - 安装goreman

    - 命令：go get github.com/mattn/goreman 

    - goreman 程序来控制基于 Procfile 的配置文件的ETCD程序。通过goreman + Procfile配置文件，可以快速启动ETCD集群。

  - Procfile配置文件：

    - 先准备Procfile配置文件，在该文件中写入需要启动节点的信息。

    - 参考官方配置文件：https://github.com/etcd-io/etcd/blob/main/Procfile
      - 启动节点为：localhost:2379 , localhost:22379 , 和 localhost:32379

      ```go
      # Use goreman to run `go get github.com/mattn/goreman`
      # Change the path of bin/etcd if etcd is located elsewhere
      
      #ETCD节点1
      etcd1: etcd --name infra1 --listen-client-urls http://127.0.0.1:2379 --advertise-client-urls http://127.0.0.1:2379 --listen-peer-urls http://127.0.0.1:12380 --initial-advertise-peer-urls http://127.0.0.1:12380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
      
      #ETCD节点2
      etcd2: etcd --name infra2 --listen-client-urls http://127.0.0.1:22379 --advertise-client-urls http://127.0.0.1:22379 --listen-peer-urls http://127.0.0.1:22380 --initial-advertise-peer-urls http://127.0.0.1:22380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
      
      #ETCD节点3
      etcd3: etcd --name infra3 --listen-client-urls http://127.0.0.1:32379 --advertise-client-urls http://127.0.0.1:32379 --listen-peer-urls http://127.0.0.1:32380 --initial-advertise-peer-urls http://127.0.0.1:32380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
      
      #proxy: etcd grpc-proxy start --endpoints=127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379 --listen-addr=127.0.0.1:23790 --advertise-client-url=127.0.0.1:23790 --enable-pprof
      
      # A learner node can be started using Procfile.learner
      ```

  

- **启动集群：**

  - 命令：goreman -f **Procfile** start 

  - Procfile为之前路径中编写的配置文件（因此，**该命令必须在Procfile文件所在路径中使用**）

    

- **使用etcdctl客户端工具**

  - 配置使用命令的版本

    - 命令：export ETCDCTL_API=3 
    - 默认使用 v2 API 来和 etcd 数据库通信，这是为了向后兼容 etcdctl。若需要使用v3 的etcdctl API，就必须配置环境变量：ETCDCTL_API 设置为版本3。

    

  - 查看集群中的信息

    etcdctl --write-out=table --endpoints=localhost:12379 member list

    | ID               | STATUS  | NAME   | PEER ADDRS             | CLIENT ADDRS           | IS LEARNER |
    | ---------------- | ------- | ------ | ---------------------- | ---------------------- | ---------- |
    | 8211f1d0f64f3269 | started | infra1 | http://127.0.0.1:12380 | http://127.0.0.1:2379  | false      |
    | 91bc3c398fb3c146 | started | infra2 | http://127.0.0.1:12380 | http://127.0.0.1:22379 | false      |
    | fd422379fda50e48 | started | infra3 | http://127.0.0.1:12380 | http://127.0.0.1:32379 | false      |

    

## 3.4、ETCD工具

### 3.4.1、etcdctl

这是命令行客户端工具，可直接对etcd数据库进行操作。

#### 3.4.1.1、<font color='red'>基本命令的附加选项</font>

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

    etcdctl --write-out=table --endpoints=localhost:12379 member list 

    查看节点localhost:12379中的成员信息



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






















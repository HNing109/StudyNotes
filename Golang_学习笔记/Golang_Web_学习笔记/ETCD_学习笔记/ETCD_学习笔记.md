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

## 3.1、安装ETCD

### 3.1.1、使用二进制文件安装

- release文件下载：https://github.com/etcd-io/etcd/releases/
- 解压之后，将etcd二进制文件存入/usr/local/bin中
- 测试是否安装完成：etcd --version



### 3.1.2、使用源码安装

- 前提：Go1.16以上的版本

- 安装命令：

  - 下载：git clone -b v3.5.0 https://github.com/etcd-io/etcd.git

  - 进入目录etcd，执行构建脚本：./build.sh

  - 构建后生成的二进制文件位于`bin`目录下，将`bin`目录的完整路径添加到系统环境变量中：

    export PATH="$PATH:pwd/bin"

  - 测试是否安装成功：etcd --version



## 3.2、ETCD部署要求

### 3.2.1、系统要求

- 系统

  amd64-linux

- 存储设备：

  80GB以上的SDD。ETCD数据需要写入磁盘中，因此为增加数据读写速度，ETCD数据存储设备需要使用SSD，且需要验证存储设备的性能。可使用fio工具进行检测。

- 处理器

  双核以上

- RAM

  最大8GB，最小2GB。因为ETCD需要强制设置一个默认的RAM大小（2GB）



### 3.2.2、部署原则

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



- 








**<font color='red'>注意：笔者选择的开发环境为：VMware、Ubuntu18、Goland</font>**

 

# 1、配置host文件

安装switchhosts软件，新建配置文件，添加以下配置（否则gostack中的私有包无法下载）

```shell
# 屏蔽的域名
#0.0.0.0 bs.studycoder.com

# My hosts
# from liwei
172.28.8.27 git.ctyun.cn
172.28.8.27 rancher.ctyun.cn
172.28.8.27 nexus.ctyun.cn
172.28.8.27 mockbin.ctyun.cn
172.28.8.27 nginx.ctyun.cn
172.28.8.27 gen.ctyun.cn
172.28.8.241 docker.ctyun.cn
172.28.8.27 rancher.ctyun.cn
172.28.8.27 swagger.ctyun.cn
172.28.8.27 moho-manual
172.28.8.27 moho-manual.ctyun.cn
172.28.8.27 doc.ctyun.cn
172.28.8.27 sonar.ctyun.cn
172.28.8.27 ntp.ctyun.cn
172.28.8.27 host.ctyun.cn
10.1.34.173 gitlab.engineering.ctyun.cn
10.1.34.176 intranet.engineering.ctyun.cn
172.28.8.27 gostack.git.ctyun.cn

172.28.8.248 mariadb.cty.os
172.28.8.248 keystone-public.cty.os
172.28.8.248 keystone-internal.cty.os
172.28.8.248 keystone-admin.cty.os 
172.28.8.248 glance-api.cty.os 
172.28.8.248 nova-api.cty.os 
172.28.8.248 nova-api-metadata.cty.os 
172.28.8.248 nova-placement-api.cty.os 
172.28.8.248 neutron-server.cty.os 
172.28.8.248 cinder-api.cty.os 
172.28.8.248 gostack-compute.cty.os
172.28.8.248 gostack-cron.cty.os
```



# 2、配置sshuttle（Linux）

- 安装软件：

  apt install sshuttle

- 配置：

  需要使用本地机器上的ssh_pub（ssh公钥）

- 连接远程服务器：

  sshuttle --dns -r root@172.28.8.248 10.114.194.0/24 -e 'ssh -p 10000'




# 2、配置ProxyCap5.36（Windows）

- **<font color='red'>安装OpenSSH</font>**    **（不安装此软件也行，可直接使用git配置时生成的SSH）**

  此密钥不同于git配置时生成的ssh，ProxyCap需要使用open-ssh密钥。ProxyCap推荐安装7.4

  - 下载：

    下载地址：https://www.mls-software.com/opensshd.html

    ![image-20210629103426442](GoStack_工程配置.assets/image-20210629103426442.png)

  - 安装步骤：

    - 安装教程：https://jingyan.baidu.com/article/9158e0002c159ea254122821.html
    - 安装路径必须在默认路径中（C:\Program Files\OpenSSH）

    - 测试是否安装成功：CMD中执行ssh

  - 生成open SSH步骤：

    - 命令：

      ssh-keygen -t rsa

      **注意：**ssh密钥保存路径可以放在其他位置，默认在C:\Users\Lenovo\ .ssh。会和git生成的ssh密钥存放路径冲突。

      即：输入上述命令后，修改.ssh默认存放位置

      Enter file in which to save the key (C:\Users\Lenovo\ .ssh): C:\Users\Lenovo\ .ssh-openssh
    
      
    
    - 在目标服务器上配置公钥：
    
      将生成的id_rsa.pub发给威哥
      
      

- **<font color='red'>安装ProxyCap5.36（注意5.26不可使用）</font>**

  - 注意事项：

    安装软件时，直接使用默认路径安装。一定要用破解软件包里面的软件进行安装，否则会出现无法访问代理服务器的情况。

  - proxycap下载地址：
    - 链接：https://pan.baidu.com/s/1kX82793TktdQBIPRjzTtxQ 
    - 提取码：w84e

  - 破解步骤：

    - 关闭杀毒软件：

      关闭windows defender，并且在运行proxycap破解软件path时，需要允许操作，否则无法执行path软件。

      <img src="GoStack_工程配置.assets/image-20210629110522137.png" alt="image-20210629110522137" style="zoom:67%;" />

    - 关闭系统后台运行的proxycap进程：

      打开任务管理器，关闭proxycap相关进程

      - ProxyCap UI
      - ProxyCap service

      

    - 运行破解软件proxycap.5.xx.64bit-patch.exe：

      图中不对，正确结果，应该是file path，OK。

      <img src="GoStack_工程配置.assets/image-20210629110358414.png" alt="image-20210629110358414" style="zoom:67%;" />

    - 查看破解结果：

      **运行破解软件后，重启计算机，可查看软件的破解结果。**
      
      <img src="GoStack_工程配置.assets/image-20210629141618328.png" alt="image-20210629141618328" style="zoom:80%;" />

  - 配置ProxyCap

    - 配置代理

      <img src="GoStack_工程配置.assets/image-20210629134057703.png" alt="image-20210629134057703" style="zoom:80%;" />

    - 配置规则

      ![image-20210806112342588](GoStack_工程配置.assets/image-20210806112342588.png)

    - 测试是否可以使用：

      **<font color='red'>必须先开启OpenVPN，连接河北VPN，然后再使用ProxyCap进行代理</font>**

      - 测试连接代理服务器：

        <img src="GoStack_工程配置.assets/image-20210629135647274.png" alt="image-20210629135647274" style="zoom:80%;" />

      - 判断ProxyCap是否已经建立代理连接

        通过查看ProxyCap的Status and Log，可查看已经建立的代理连接。 

        <img src="GoStack_工程配置.assets/image-20210806112034059.png" alt="image-20210806112034059" style="zoom: 80%;" />

      - 测试连接保定02测试环境ETCD数据库

        - 浏览器访问地址：

          http://10.114.194.115:12000/etcdkeeper/

        - 服务器地址：10.114.194.115:12379

          用户名：root 

          密码:：CTyun2020!
      
        - 成功连接：
      
          **此处填写ETCD的IP地址**
        
          <img src="GoStack_工程配置.assets/image-20210629135912990.png" alt="image-20210629135912990" style="zoom:80%;" />
      
      

- **开发环境ETCD地址: （本地调试）**

  - 10.114.194.115:12379 
  - 10.114.194.115:22379 
  - 10.114.194.115:32379 

- **测试环境ETCD地址:** 

  - 10.114.194.116:12379 
  - 10.114.194.116:22379 
  - 10.114.194.116:32379

-  **注意：**若ProxyCap软件无法正常运行，直接重装软件，即可解决问题。

  

# 3、安装OpenVPN

查看说明手册



# 4、安装libvirt

- Linux环境

  sudo apt-get install libvirt-dev libvirt-daemon libvirt-clients

  （不安装libvirt，在gostack/cmd/gostack中执行go build时，将导致参数无法读取，无法完成build操作）

- Windows环境

  

# 5、go env配置

- go编译器版本：

  15.12

  

- 配置命令： 

  go env -w GO111MODULE=on
  go env -w GOPROXY=https://goproxy.io,direct

  go env -w GOINSECURE=git.ctyun.cn
  go env -w GOPRIVATE=git.ctyun.cn

   

  ```shell
  查看是否生效，打开以下路径的文件：
  windows：C:\Users\Lenovo\AppData\Roaming\go\env
  ubuntu:/home/chris/.config/go
  
  GO111MODULE=on
  GOINSECURE=git.ctyun.cn
  GOPRIVATE=git.ctyun.cn
  GOPROXY=https://goproxy.io,direct
  ```

  

  使用go env命令（查看配置）：

  ```shell
  C:\Code\gostack>go env
  set GO111MODULE=on
  set GOARCH=amd64
  set GOBIN=
  set GOCACHE=C:\Users\Lenovo\AppData\Local\go-build
  set GOENV=C:\Users\Lenovo\AppData\Roaming\go\env
  set GOEXE=.exe
  set GOFLAGS= -mod=
  set GOHOSTARCH=amd64
  set GOHOSTOS=windows
  set GOINSECURE=git.ctyun.cn
  set GOMODCACHE=C:\Users\Lenovo\go\pkg\mod
  set GONOPROXY=git.ctyun.cn
  set GONOSUMDB=git.ctyun.cn
  set GOOS=windows
  set GOPATH=C:\Users\Lenovo\go
  set GOPRIVATE=git.ctyun.cn
  set GOPROXY=https://goproxy.io,direct
  set GOROOT=C:\Program Files\Go
  set GOSUMDB=sum.golang.org
  set GOTMPDIR=
  set GOTOOLDIR=C:\Program Files\Go\pkg\tool\windows_amd64
  set GCCGO=gccgo
  set AR=ar
  set CC=gcc
  set CXX=g++
  set CGO_ENABLED=1
  set GOMOD=C:\Code\gostack\go.mod
  set CGO_CFLAGS=-g -O2
  set CGO_CPPFLAGS=
  set CGO_CXXFLAGS=-g -O2
  set CGO_FFLAGS=-g -O2
  set CGO_LDFLAGS=-g -O2
  set PKG_CONFIG=pkg-config
  set GOGCCFLAGS=-m64 -mthreads -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=C:\Users\Lenovo\AppData\Local\Temp\go-build594281098=/tmp/go-build -gno-re
  cord-gcc-switches
  ```

  

- 获取go.mod中配置的包：

  命令：go mod download

  

  Q：若上述命令执行后未出现错误，但Gostack工程依然不能import外部包

  A：方式一：重启Goland

  ​      方式二：删除Gostack工程，重新git clone工程，打开Goland编译器，即可自动导入外部包。
  
  

- **<font color='red'>在获取GitLab中的包时，需要填写自己的Gitlab账号、密码</font>** 

  

# 6、Goland配置

- 无需配置Goland，否则会出现包无法使用的情况

  ![image-20210622101124589](GoStack_工程配置.assets/image-20210622101124589.png)



# 7、启动代码

**注意：**gostack中的所有模块，都需要逐个手动启动，无法一次性全部启动

**启动代码的前提：**本机必须完成上述配置，并且 ***开启OpenVPN、ProxyCap，确保本机能够访问ETCD数据库***

- **方式一：**

  该方式比较适合Linux环境使用

  - 打包gostack整个工程：

    命令：

    - cd cmd/gostack/

    - go build

    （打包完成后，该目录下会生成一个gostack可执行文件）

  - 启动需要执行的模块

    在cmd/gostack/目录下，使用命令：./gostack xxx

    eg：启动scheduler，  ./gostack scheduler
    
    

- **方式二：**

  自己新建一个main函数，调用RunXxx，并配置对应的etc/app.yml文件路径

  eg：windows下，启动scheduler模块

  ![img](GoStack_工程配置.assets/PUODS6_Z91R3{FEOZX%%2.png)

  
  
  ```go
  //图中的代码如下
  package main
  import "git.ctyun.cn/gostack/gostack/scheduler"
  
  func main() {
  	path := "./scheduler/etc/app.yml"
  	scheduler.RunScheduler(path)
  }
  ```

  
  
  

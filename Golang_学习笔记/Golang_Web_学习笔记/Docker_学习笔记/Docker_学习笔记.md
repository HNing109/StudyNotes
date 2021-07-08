Docker-Hub：https://www.docker.com/products/docker-hub

# 1、Docker入门

## 1.1、虚拟技术、容器技术

- **虚拟技术**

  例如Vmware，直接在操作系统中通过软件的方式**虚拟化一个完整的操作系统**，该系统和真实系统没有区别。需要加载内核、Lib相关库，然后其余的APP应用可在该系统中运行。

  - **缺点**

    1、 资源占用十分多：需要加载很多Lib组件

    2、 冗余步骤多：

    3、 启动很慢：加载的系统kernel内核庞大

    4、当Lib组件崩溃时，容易导致系统无法运行。

  ![image-20210708091223300](Docker_学习笔记.assets/image-20210708091223300.png)

- **容器技术**

  容器化并没有模拟一个完整的操作系统，只需要加载系统核心部分的内核。将Lib组件和APP都封装起来放进一个容器中，每个容器都是相互隔离的，不会相互影响。即使某个容器崩溃，其余容器可照常运行。

  ![image-20210708091502212](Docker_学习笔记.assets/image-20210708091502212.png)

  - **为什么Docker比VMware快**
    1、docker有着比虚拟机更少的抽象层。由于docker不需要Hypervisor实现硬件资源虚拟化,运行在docker容器上的程序直接使用的都是实际物理机的硬件资源。因此在CPU、内存利用率上docker将会在效率上有明显优势。
    2、docker利用的是宿主机的内核,而不需要Guest OS。
  
    

## 1.2、基本概念

- **主机：**用于运行docker环境的系统（实体主机、虚拟机都可以）；

- **客户端：**使用命令来操作主机中的docker环境；

- **镜像：**指的是一个个软件，用户可以自己打包自己开发环境的某个软件，生成对应的镜像——例如：tomcat、mysql；

- **容器：**镜像在dokcer中运行后就变成容器，容器具有启动、删除、停止等操作；

- **仓库：（eg：dokcer hub）**所有的镜像都可以保存在远程仓库，等需要使用的时候直接下载即可，不再需要对系统进行软件安装。

  ![image-20210701234524627](Docker_学习笔记.assets/image-20210701234524627.png)



## 1.3、Ubuntu安装Docker

官网教程：https://docs.docker.com/engine/install/ubuntu/

- 更新安装软件列表

  ```shell
  sudo apt-get update
  ```

- 安装必要的组件

  ```shell
  sudo apt-get install \
      apt-transport-https \
      ca-certificates \
      curl \
      gnupg \
      lsb-release
  ```

- 卸载旧版本的Docker

  ```shell
  sudo apt-get remove docker \
      docker-client \
      docker-client-latest \
      docker-common \
      docker-latest \
      docker-latest-logrotate \
      docker-logrotate \
      docker-engine
  ```

- （可选）修改ubuntu的镜像仓库

  ```shell
  #地址：/etc/apt/sources.list ， 
  
  #阿里源
  deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
  deb http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
  deb http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
  deb http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
  deb http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
  deb-src http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
  deb-src http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
  deb-src http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
  deb-src http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
  deb-src http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
  
  
  
  #将原有的镜像替换为163镜像
  deb http://mirrors.163.com/ubuntu/ bionic-backports main restricted universe multiverse
  deb http://mirrors.163.com/ubuntu/ bionic-security main restricted universe multiverse
  deb http://mirrors.163.com/ubuntu/ bionic-backports main restricted universe multiverse
  deb http://mirrors.163.com/ubuntu/ bionic-security main restricted universe multiverse
  
  ```
  
- 添加GPG key

  ```shell
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
  
  echo \
    "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
  ```
  
  
  
- 安装Docker（推荐使用最新版Docker）

  ```shell
  #方式1
  sudo apt-get update
  sudo apt-get upgrade
  sudo apt-get install docker-ce docker-ce-cli containerd.io
  
  #方式2
  #列出Docker版本列表
  sudo apt-cache madison docker-ce
  #指定使用的Dokcer版本
  sudo apt-get install docker-ce=<VERSION_STRING> docker-ce-cli=<VERSION_STRING> containerd.io
  ```

- 测试Docker是否安装成功

  ```shell
  #需要等待一会儿，因为本地仓库中没有hello-world镜像，需要临时下载
  sudo docker run hello-world
  ```

- **修改Docker的软件下载镜像仓库**

  使用阿里云官网提供的：https://cr.console.aliyun.com/cn-hangzhou/instances/mirrors

  通过修改daemon配置文件/etc/docker/daemon.json来使用加速器

  ```shell
  sudo mkdir -p /etc/docker
  sudo tee /etc/docker/daemon.json <<-'EOF'
  {
    "registry-mirrors": ["https://w0vftfrg.mirror.aliyuncs.com"]
  }
  EOF
  sudo systemctl daemon-reload
  sudo systemctl restart docker
  ```

  



## 1.4、基本使用

### 1.4.1、常用命令

#### 1.4.1.1、镜像命令

- docker images：

  查看当前系统已经pull安装的镜像；

  ```shell
  REPOSITORY    TAG       IMAGE ID       CREATED        SIZE
  hello-world   latest    d1165f221234   4 months ago   13.3kB
  ```

  

- docker search 镜像名字

  搜索镜像

  

- docker pull 镜像名字：  （等同于docker image pull）

  从docker仓库**下载镜像**；默认下载latest最新版

  **下载指定版本**：docker pull mysql:5.7

  ```shell
  root@chris:~# docker pull mysql:5.7
  5.7: Pulling from library/mysql		#如果不写tag，默认就是latest
  b4d181a07f80: Already exists 		#分层下载： docker image 的核心，UFS联合文件系统（防止下载重复的镜像）
  a462b60610f5: Pull complete 		#从第二层开始现在，上一层已经存在，故不重复下载
  578fafb77ab8: Pull complete 
  524046006037: Pull complete 
  d0cbe54c8855: Pull complete 
  aa18e05cc46d: Pull complete 
  32ca814c833f: Pull complete 
  52645b4af634: Pull complete 
  bca6a5b14385: Pull complete 
  309f36297c75: Pull complete 
  7d75cacde0f8: Pull complete 
  Digest: sha256:1a2f9cd257e75cc80e9118b303d1648366bc2049101449bf2c8d82b022ea86b7	# 签名 防伪
  Status: Downloaded newer image for mysql:5.7
  docker.io/library/mysql:5.7	#真实地址
  ```

  

- docker rmi 镜像名字：   （等同于docker image rm）

  删除某个**下载的镜像**；（即：删除docker pull下载的镜像文件）

  ```shell
  docker rmi -f 镜像id					 	#删除指定的镜像
  docker rmi -f 镜像id 镜像id 镜像id 镜像id	#删除指定的镜像
  docker rmi -f $(docker images -aq) 		 #删除全部的镜像
  ```
  
  

#### 1.4.1.2、容器命令

- **docker run 镜像名**

  新建、启动容器；

  ```shell
  docker run 可选参数 image名 （或者，docker container run [可选参数] image名）
  #可选参数说明
  --name="Name"		容器名字 tomcat01 tomcat02 用来区分容器
  -d					后台方式运行
  -it 				使用交互方式运行，进入容器查看内容
  -p					指定容器的端口 
  		-p ip:主机端口:容器端口
  		-p 主机端口:容器端口 (常用)，eg：-p 8080(宿主机):8080(容器)
  		-p 容器端口 （默认只有容器端口）
  		
  --rm image名        一般是用来测试，用完就删除容器
  -P(大写) 			  随机指定端口
  ```

  - **eg：**

    - 启动、进入容器中的centos

      docker run **-it** centos /bin/bash

    - 测试tomcat是否能启动

      docker run -it **--rm** tomcat:9.0

    

  - **常见问题：**

    直接使用命令docker run **-d** centos启动，使用docker ps发现centos 停止了。**这是常见的坑**，docker容器使用后台运行，就必须要有要一个前台进程，docker发现没有应用，就会自动停止。

  

- **docker start 容器id**

  启动指定id的容器（启动之前已经docker run运行的容器）

  

- **docker restart 容器id**

  重启之前运行的镜像。

  **eg：**当VMware重启之后，使用无法访问之前开启的mysql，使用docker ps无法看到mysql正在运行，使用docker ps –a可以看到，此时只需要使用该命令，就可以重启这个镜像，外部可以正常访问mysql数据库；

  

- **docker stop 容器id**

  停止运行某个镜像；

  

- **docker kill 容器id**

  强制停止当前容器

  

- **进入正在运行的容器中**

  ```shell
  docker 可选参数 -it 容器id 路径
  #可选参数说明
  	exec   # 进入当前容器后开启一个新的终端，可以在里面操作。（常用，必须是docker ps，显示的正在运行的容器）
  	attach # 进入容器正在执行的终端
  ```

  

- **退出当前进入的容器**

   ```shell
   #方式1（命令）
   exit 		   				#直接退出、停止容器
   #方式2（按键）
   （ctrl +P） + （ctrl +Q）     #退出容器，但不停止容器
   ```

  

- **docker rm 容器Id**：

  删除**正在运行**的某个镜像；（即：停止docker ps -a，查找出正在运行的镜像）

  ```shell
  docker rm 容器id   				#删除指定的容器，不能删除正在运行的容器，如果要强制删除 rm -rf
  docker rm -f $(docker ps -aq)    #删除指定的容器
  docker ps -a -q|xargs docker rm  #删除所有的容器
  ```

  

- **docker ps：**

  显示所有运行的镜像；

  ```shell
  docker ps 可选参数
  #（默认）docker ps：显示 有端口映射、且正在运行的镜像；（使用-p启动）
  #可选参数说明
    -a, --all             显示所有容器
    -n, --last int        显示指定个数的容器
    -q, --quiet           只显示容器di
  ```



- **<font color='red'>从容器中拷贝数据至宿主机</font>**

  docker cp 容器id:容器内路径  宿主机目的路径

  

- **<font color='red'>从宿主机拷贝数据至容器中</font>**

  docker cp 宿主机目的路径 容器id:容器内路径 

  

- **<font color='red'>挂载宿主机文件夹至容器中</font>**

  docker run -it -v 宿主机绝对路径**:**容器绝对路径 容器名 bashshell

  若容器中不存在文件夹，会自动新建该文件夹。注意，**必须在启动新容器时，才能进行-v挂载操作**。

  eg：

  ```shell
  docker run -it -v /home/chris/test:/soft centos /bin/bash
  ```

  

- 查看docker的日志

  ```shell
  docker logs -t --tail n 容器id #查看指定容器的n行日志
  docker logs -ft 容器id         #查看指定容器的所有日志
  ```

  

- docker top 容器id

  查看容器的任务管理器

  

- docker inspect 容器id

  查看镜像的元数据

  

- docker stats 

  查看docker容器使用的cpu、内存情况

  

- **Docker可视化工具**

  - portainer（不常用）

    安装命令：

    ```shell
    docker run -d -p 8080:9000 \
    --restart=always -v /var/run/docker.sock:/var/run/docker.sock --privileged=true portainer/portainer
    ```

    浏览器访问：http://192.168.83.136:8080/ ，创建账号后可以访问。（这边的IP为Docker宿主机的网卡地址）

    ![image-20210708135705937](Docker_学习笔记.assets/image-20210708135705937.png)

  - Rancher（CI/CD）

    



### 1.4.2、Docker使用MySQL

- **安装**

  - docker search mysql

  - docker pull mysql （默认安装latest版本）

    （注意：安装的latest的mysql，使用较低版本的navicat是无法登陆，需要更新navicat版本）

- **启动Mysql**

  MySQL默认端口：3306

  - **命令：**

    docker run --name mysql04 -e MYSQL_ROOT_PASSWORD=123456 -d -p 5001:3306 mysql 

  - **解释：**

    - --name：为这次启动的mysql起别名；
    - -e：设置MySQL启动时的环境变量（默认参数），启动时必须配置密码，否则无法连接MySQL数据库
    - MYSQLROOT_PASSWORD：mysql登陆密码；
    - -d：后台运行（daemon）；
    - **-p 5001:3306**  ：将docker容器的端口映射到Linux系统的端口上。mysql的映射端口：**linux为5001，docker为3306**。因此，使用Navicat连接MySQL时，连接端口为Linux的5001端口（而不是Docker的3306端口）
    - mysql：启动的docker镜像（真实的应用名称）



### 1.4.3、Docker安装Redis

安装、启动Redis的方式，参见安装MySQL步骤.

Redis默认端口：6379

- 启动

  - 命令：docker run --name redis01 -d -p 4000:6379 redis

    



### 1.4.4、Docker安装RabbitMQ

- **安装**

  - rabbitmq**（带web界面）**：docker pull **rabbitmq:management** 

  - rabbitmq**（没有web界面）**：docker pull **rabbitmq**

    

- **启动**

  - **命令：**

    docker run -dit --name Myrabbitmq -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin -p 15672:15672 -p 5672:5672 **rabbitmq:management**

  - **解释：**

    - Myrabbitmq：设定此次启动容器的别名

    - RABBITMQ_DEFAULT_USER=admin：设定用户名

    - RABBITMQ_DEFAULT_PASS=admin：设定密码

    - -p 15672:15672：设定浏览器登陆web端口

    - -p 5672:5672：设定客户端端口（springboot中.yml配置的端口）

    - **rabbitmq:management**：此次启动的docker镜像（真实的应用名称）

      

  - **检测是否成功启动：**

    使用浏览器的web界面，IP + 端口15672

    ![image-20210702000735222](Docker_学习笔记.assets/image-20210702000735222.png)




# 2、Docker进阶

## 2.1、容器数据卷





## 2.2、数据卷容器



## 2.3、DockerFile





## 2.4、制作自己的镜像

### 2.4.1、制作流程



### 2.4.2、发布镜像至DockerHub



# 3、Docker高级

## 3.1、Docker网络





## 3.2、Docker Compose（yml）







## 3.3、Docker Swarm（kubernetes）





## 3.4、CI/CD之Jenkins






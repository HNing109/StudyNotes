本教程安装QueenS版本的OpenStack，采用all in one的方式安装（即：所有组件均安装在同一个机器上）

OpenStack官网安装手册：https://docs.openstack.org/install-guide/



其余的参考手册

1. 官网安装指南：https://docs.openstack.org/install-guide/
2. 华为鲲鹏的帮助文档（供参考）：https://support.huaweicloud.com/dpmg-kunpengcpfs/kunpengcpfs_04_0002.html
3. 在openEuler（类似centos）上安装openstack queens版本（供参考）：https://docs.openeuler.org/zh/docs/20.03_LTS_SP2/docs/thirdparty_migration/OpenStack-queens.html 



# 1、搭建基础环境

## 1.1、安装VMWare

使用VMWare16 pro版本，下载地址：https://www.vmware.com/cn/products/workstation-pro/workstation-pro-evaluation.html



## 1.2、下载CentoOS7

使用CentOS7系统，下载地址：https://www.vmware.com/cn/products/workstation-pro/workstation-pro-evaluation.html



## 1.3、安装虚拟机

- 安装教程：https://blog.csdn.net/babyxue/article/details/80970526

- 虚拟机配置如下：

  RAM：大于10GB，否则OpenStack安装完成后，虚拟机运行会卡顿。

  ROM：100GB，分两块磁盘，sda用于存放CentOS系统，sdb用于存放LVM（cinder的数据卷）。

  网卡：enss33用于虚拟机访问网络，enss37用于OpenStack组件之间的通信（all in one的方式，实际上不需要用到enss37网卡）

  ![image-20210724154502009](OpenStack_安装笔记.assets/image-20210724154502009.png)



# 2、准备工作

## 2.1、禁用SELinux

作用：SELinux致使很多服务端口默认是关闭，可能会导致Apache服务无法启动。

错误如下："Devstack fail to start apache2 -"Address already in use":"coild not bind to address"

解决方式：配置/etc/selinux/config，关闭SELINUX

```shell
#配置
[root@controller /]# cat /etc/selinux/config

# This file controls the state of SELinux on the system.
# SELINUX= can take one of these three values:
#     enforcing - SELinux security policy is enforced.
#     permissive - SELinux prints warnings instead of enforcing.
#     disabled - No SELinux policy is loaded.
SELINUX=disabled
# SELINUXTYPE= can take one of three values:
#     targeted - Targeted processes are protected,
#     minimum - Modification of targeted policy. Only selected processes are protected. 
#     mls - Multi Level Security protection.
SELINUXTYPE=targeted 

#查看配置结果
[root@controller openstack]# sestatus -v
SELinux status:                 disabled
```



## 2.2、关闭防火墙

作用：免去配置OpenStack各个节点时，需要不断手动开启端口的麻烦。但是在实际生产环境中绝对不允许这样做，这会引起安全问题。

```shell
[root@controller ~]# systemctl stop firewalld 
[root@controller ~]# systemctl disable firewalld
```



## 2.3、配置网卡

CentOS系统默认使用*DHCP*（动态主机配置协议）分配网卡IP。此处，建议使用static（静态）网卡IP，因为在VMWare的虚拟机重启之后，虚拟机的网卡IP可能会变化，这就会导致OpenStack其他节点配置的网卡IP无法使用。

- 查看当前网卡IP

  ```shell
  [root@controller /]# ifconfig
  ens33: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
          inet 192.168.83.139  netmask 255.255.255.0  broadcast 192.168.83.255
          inet6 fe80::c62f:bd78:ba20:b81  prefixlen 64  scopeid 0x20<link>
          ether 00:0c:29:5b:0e:1c  txqueuelen 1000  (Ethernet)
          RX packets 15733  bytes 10884363 (10.3 MiB)
          RX errors 0  dropped 0  overruns 0  frame 0
          TX packets 11840  bytes 6242032 (5.9 MiB)
          TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
  ```

  

- 配置网卡IP

  ```shell
  [root@controller /]# cat /etc/sysconfig/network-scripts/ifcfg-ens33
  TYPE=Ethernet
  PROXY_METHOD=none
  BROWSER_ONLY=no
  DEFROUTE=yes
  IPV4_FAILURE_FATAL=no
  IPV6INIT=yes
  IPV6_AUTOCONF=yes
  IPV6_DEFROUTE=yes
  IPV6_FAILURE_FATAL=no
  IPV6_ADDR_GEN_MODE=stable-privacy
  NAME=ens33
  UUID=038e74ab-c77d-4f83-b67b-dabde3a03c1d
  DEVICE=ens33
  #更改的部分
  ONBOOT=yes
  BOOTPROTO=static
  IPADDR=192.168.83.139
  GATEWAY=192.168.83.2
  NETMASK=255.255.255.0
  DNS1=114.114.114.114
  ```

  若需要配置第二张网卡enss37，则将上面的NAME、DEVICE更改后，删除UUID，并配置新的IPADDR，即可。



## 2.4、配置主机名

```shell
#设置主机名
[root@controller ~]# hostnamectl set-hostname controller 

#查看设置
[root@controller ~]# hostname 
controller 

#修改hosts文件：controller域名，解析时指向本机的网卡IP
[root@controller ~]# cat /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6

192.168.83.139 controller
```



## 2.5、设置时区

```shell
[root@controller ~]# cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime cp: overwrite ‘/etc/localtime’? y 
[root@controller ~]# date Tue Jul 20 16:13:44 CST 2021
```



# 3、安装软件

**注意：安装软件时，如果出现安装失败的情况，可能是因为网络原因导致的。只需要再次执行安装命令即可。**

## 3.1、安装NTP

用于同步各个节点的时间。需要安装Chrony，通常情况下安装在控制器节点，其他节点则同步控制节点的时间，来达到各个节点时间一致的效果。

```shell
#安装chrony
[root@controller /]# yum install chrony -y

#配置（可选）
[root@controller /]# vim /etc/chrony.conf
server 0.centos.pool.ntp.org iburst

#启动软件、设置开机启动
[root@controller /]# systemctl start chronyd.service
[root@controller /]# systemctl enable chronyd.service

```



## 3.2、添加OpenStack软件源

网络上发布的OpenStack存在多个不同版本，eg：Victoria、Queens、Train等。本次使用的时Queens版本。

```shell
#添加Queens版本源
[root@controller ~]# yum install centos-release-openstack-queens -y 

#升级软件包
[root@controller ~]# yum upgrade -y

#安装OpenStack Client
[root@controller ~]# yum install python-openstackclient -y
```



## 3.3、安装SQL数据库

通常情况下，数据库安装在控制节点上，用于存储OpenStack的数据。

```shell
#安装SQL
[root@controller ~]# yum install mariadb mariadb-server python2-PyMySQL -y

#配置SQL
[root@controller /]# vim /etc/my.cnf.d/openstack.cnf
[mysqld]
#修改为本机的网卡IP
bind-address = 192.168.83.139

default-storage-engine = innodb
innodb_file_per_table = on
max_connections = 4096
collation-server = utf8_general_ci
character-set-server = utf8

#设置开机自启、启动SQL
[root@controller /]# systemctl enable mariadb.service
[root@controller /]# systemctl start mariadb.service

#设置SQL数据库的root账号、密码
[root@controller /]# mysql_secure_installation
```



## 3.4、安装消息队列

消息队列服务通常安装在控制器节点上。用于协调服务之间的操作和状态信息。

```shell
#安装MQ
[root@controller /]# yum install rabbitmq-server

#设置开机自启、启动MQ
[root@controller /]# systemctl enable rabbitmq-server.service
[root@controller /]# systemctl start rabbitmq-server.service

#添加MQ用户，设置密码
[root@controller /]# rabbitmqctl add_user openstack RABBIT_PASS

#允许openstack用户进行配置、写入和读取访问
[root@controller /]# rabbitmqctl set_permissions openstack ".*" ".*" ".*"
```



## 3.5、安装Memcached

Memcached 服务通常安装在控制器节点上。服务的身份服务认证机制使用 Memcached 来缓存token令牌。

```shell
#安装
[root@controller /]# yum install memcached python-memcached

#配置
[root@controller /]# vim /etc/sysconfig/memcached
PORT="11211"
USER="memcached"
MAXCONN="1024"
CACHESIZE="64"
#新增controller域名，其他节点可以通过该域名直接访问
OPTIONS="-l 127.0.0.1,::1,controller"

#设置开机自启、启动软件
[root@controller /]# systemctl enable memcached.service
[root@controller /]# systemctl start memcached.service
```



## 3.6、安装ETCD数据库

用于分布式密钥锁定、存储配置、跟踪服务实时性和其他场景

```shell
#安装etcd
[root@controller /]# yum install etcd

#配置etcd：
#ETCD_INITIAL_CLUSTER, ETCD_INITIAL_ADVERTISE_PEER_URLS, ETCD_ADVERTISE_CLIENT_URLS,ETCD_LISTEN_CLIENT_URLS设置为控制节点的IP（即：本机的网卡IP）
[root@controller /]# vim /etc/etcd/etcd.conf
[Member]
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="http://192.168.83.139:2380"
ETCD_LISTEN_CLIENT_URLS="http://localhost:2379,http://192.168.83.139:2379"
ETCD_NAME="controller"
[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="http://192.168.83.139:2380"
ETCD_ADVERTISE_CLIENT_URLS="http://192.168.83.139:2379"
ETCD_INITIAL_CLUSTER="controller=http://192.168.83.139:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster-01"
ETCD_INITIAL_CLUSTER_STATE="new"

#设置开机自启、启动软件
[root@controller /]# systemctl enable etcd
[root@controller /]# systemctl start etcd
```



# 4、安装OpenStack

注意：本次安装的是Queens版本

## 4.1、默认密码

本次安装OpenStack所有的组件，均使用默认密码。

| Password name                        | Description                                        |
| :----------------------------------- | :------------------------------------------------- |
| Database password (no variable used) | Root password for the database                     |
| `ADMIN_PASS`                         | Password of user `admin`                           |
| `CINDER_DBPASS`                      | Database password for the Block Storage service    |
| `CINDER_PASS`                        | Password of Block Storage service user `cinder`    |
| `DASH_DBPASS`                        | Database password for the Dashboard                |
| `DEMO_PASS`                          | Password of user `demo`                            |
| `GLANCE_DBPASS`                      | Database password for Image service                |
| `GLANCE_PASS`                        | Password of Image service user `glance`            |
| `KEYSTONE_DBPASS`                    | Database password of Identity service              |
| `METADATA_SECRET`                    | Secret for the metadata proxy                      |
| `NEUTRON_DBPASS`                     | Database password for the Networking service       |
| `NEUTRON_PASS`                       | Password of Networking service user `neutron`      |
| `NOVA_DBPASS`                        | Database password for Compute service              |
| `NOVA_PASS`                          | Password of Compute service user `nova`            |
| `PLACEMENT_PASS`                     | Password of the Placement service user `placement` |
| `RABBIT_PASS`                        | Password of RabbitMQ user `openstack`              |



## 4.2、安装Keystone

### 4.2.1、安装Keystone步骤

- 作用：用于管理身份验证、授权和服务目录。通过身份验证后，最终用户可以使用他们的身份访问其他 OpenStack 服务。

- 创建数据库

  ```shell
  #登录数据库
  [root@controller /]# mysql -u root -p
  
  #创建keystone数据库
  MariaDB [(none)]> CREATE DATABASE keystone;
  
  #用户登录授权
  #本地登录
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON keystone.* TO 'keystone'@'localhost' \
  IDENTIFIED BY 'KEYSTONE_DBPASS';
  #远程登录
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON keystone.* TO 'keystone'@'%' \
  IDENTIFIED BY 'KEYSTONE_DBPASS';
  ```

  

- 安装、配置keystone

  ```shell
  #安装软件
  [root@controller /]# yum install openstack-keystone httpd mod_wsgi
  
  #配置keystone.conf
  [root@controller /]# vim /etc/keystone/keystone.conf 
  
  [database]
  connection = mysql+pymysql://keystone:KEYSTONE_DBPASS@controller/keystone
  
  [token]
  provider = fernet
  
  #填充keystone数据库
  [root@controller /]# su -s /bin/sh -c "keystone-manage db_sync" keystone
  
  #初始化Fernet密钥库
  [root@controller /]# keystone-manage fernet_setup --keystone-user keystone --keystone-group keystone
  [root@controller /]# keystone-manage credential_setup --keystone-user keystone --keystone-group keystone
  
  #引导身份服务
  [root@controller /]# keystone-manage bootstrap --bootstrap-password ADMIN_PASS \
                        --bootstrap-admin-url http://controller:5000/v3/ \
                        --bootstrap-internal-url http://controller:5000/v3/ \
                        --bootstrap-public-url http://controller:5000/v3/ \
                        --bootstrap-region-id RegionOne
                        
  ```

   

- 配置Apache HTTP服务

  ```shell
  #配置httpd.conf
  [root@controller /]# vim /etc/httpd/conf/httpd.conf
  ServerName controller
  
  #创建软链接
  [root@controller /]# /usr/share/keystone/wsgi-keystone.conf
  ```

  

- 配置系统环境

  ```shell
  [root@controller /]# systemctl enable httpd.service
  [root@controller /]# systemctl start httpd.service
  
  #配置管理账户
  [root@controller /]# export OS_USERNAME=admin \
  					 export OS_PASSWORD=ADMIN_PASS \
  					 export OS_PROJECT_NAME=admin \
  					 export OS_USER_DOMAIN_NAME=Default \
  					 export OS_PROJECT_DOMAIN_NAME=Default \
  					 export OS_AUTH_URL=http://controller:5000/v3 \
  					 export OS_IDENTITY_API_VERSION=3
  ```



- 创建域、项目、用户、角色信息

  - 此处创建一个示例域：example

    ```shell
    #创建域
    [root@controller /]# openstack domain create --description "An Example Domain" example
    +-------------+----------------------------------+
    | Field       | Value                            |
    +-------------+----------------------------------+
    | description | An Example Domain                |
    | enabled     | True                             |
    | id          | 2f4f80574fd84fe6ba9067228ae0a50c |
    | name        | example                          |
    +-------------+----------------------------------+
    
    #创建项目
    [root@controller /]# openstack project create --domain default \
      					 --description "Service Project" service
    +-------------+----------------------------------+
    | Field       | Value                            |
    +-------------+----------------------------------+
    | description | Service Project                  |
    | domain_id   | default                          |
    | enabled     | True                             |
    | id          | 24ac7f19cd944f4cba1d77469b2a73ed |
    | is_domain   | False                            |
    | name        | service                          |
    | parent_id   | default                          |
    +-------------+----------------------------------+
    ```

    

  - 创建一个demo项目

    ```shell
    #创建demo项目
    [root@controller /]# openstack project create --domain default \
                         --description "Demo Project" demo
    +-------------+----------------------------------+
    | Field       | Value                            |
    +-------------+----------------------------------+
    | description | Demo Project                     |
    | domain_id   | default                          |
    | enabled     | True                             |
    | id          | 231ad6e7ebba47d6a1e57e1cc07ae446 |
    | is_domain   | False                            |
    | name        | demo                             |
    | parent_id   | default                          |
    +-------------+----------------------------------+
    
    #创建demo项目的用户
    [root@controller /]# openstack user create --domain default \
                         --password-prompt demo
    #填入默认密码：DEMO_PASS
    User Password:
    Repeat User Password:
    +---------------------+----------------------------------+
    | Field               | Value                            |
    +---------------------+----------------------------------+
    | domain_id           | default                          |
    | enabled             | True                             |
    | id                  | aeda23aa78f44e859900e22c24817832 |
    | name                | demo                             |
    | options             | {}                               |
    | password_expires_at | None                             |
    +---------------------+----------------------------------+
    
    #创建角色
    [root@controller /]# openstack role create user
    +-----------+----------------------------------+
    | Field     | Value                            |
    +-----------+----------------------------------+
    | domain_id | None                             |
    | id        | 997ce8d05fc143ac97d83fdfb5998552 |
    | name      | user                             |
    +-----------+----------------------------------+
    
    #为demo项目的用户：demoe，赋予角色：user
    [root@controller /]# openstack role add --project demo --user demo user
    ```

    

### 4.2.2、验证是否成功安装

```shell
#取消设置临时OS_AUTH_URL和OS_PASSWORD 环境变量
[root@controller /]# unset OS_AUTH_URL OS_PASSWORD

#
[root@controller /]# openstack --os-auth-url http://controller:35357/v3 \
  --os-project-domain-name Default --os-user-domain-name Default \
  --os-project-name admin --os-username admin token issue
#此处填入的密码为用户demo的密码，默认为DEMO_PASS
Password:
+------------+-----------------------------------------------------------------+
| Field      | Value                                                           |
+------------+-----------------------------------------------------------------+
| expires    | 2016-02-12T20:14:07.056119Z                                     |
| id         | gAAAAABWvi7_B8kKQD9wdXac8MoZiQldmjEO643d-e_j-XXq9AmIegIbA7UHGPv |
|            | atnN21qtOMjCFWX7BReJEQnVOAj3nclRQgAYRsfSU_MrsuWb4EDtnjU7HEpoBb4 |
|            | o6ozsA_NmFWEpLeKy0uNn_WeKbAhYygrsmQGA49dclHVnz-OMVLiyM9ws       |
| project_id | 343d245e850143a096806dfaefa9afdc                                |
| user_id    | ac3377633149401296f6c0d92d79dc16                                |
+------------+-----------------------------------------------------------------+
```



### 4.2.3、创建OpenStack客户端脚本

- 为什么需要使用脚本？

  由于之前，配置系统环境时，使用的是export命令，该命令用于临时配置环境变量，当系统重启或者执行export明令的终端关闭后，该环境变量就无法使用了，导致openstack命令无法被执行。

  

- 需要创建的脚本如下：

  为管理（admin）和演示（demo）项目以及用户创建客户端环境脚本。即：获取相应权限的CLI命令(Command Line Interface)，通过命令的方式进行交互。这两个脚本可以存放在同一个目录中。

  - **admin-openrc**

    ```shell
    export OS_PROJECT_DOMAIN_NAME=Default
    export OS_USER_DOMAIN_NAME=Default
    export OS_PROJECT_NAME=admin
    export OS_USERNAME=admin
    export OS_PASSWORD=ADMIN_PASS
    export OS_AUTH_URL=http://controller:5000/v3
    export OS_IDENTITY_API_VERSION=3
    export OS_IMAGE_API_VERSION=2
    ```

    

  - **demo-openrc**

    ```shell
    export OS_PROJECT_DOMAIN_NAME=Default
    export OS_USER_DOMAIN_NAME=Default
    export OS_PROJECT_NAME=demo
    export OS_USERNAME=demo
    export OS_PASSWORD=DEMO_PASS
    export OS_AUTH_URL=http://controller:5000/v3
    export OS_IDENTITY_API_VERSION=3
    export OS_IMAGE_API_VERSION=2
    ```

  

- 使用脚本的方式

  ```shell
  [root@controller /]# . admin-openrc
  
  #请求身份验证令牌
  [root@controller /]# openstack token issue
  +------------+-----------------------------------------------------------------+
  | Field      | Value                                                           |
  +------------+-----------------------------------------------------------------+
  | expires    | 2016-02-12T20:44:35.659723Z                                     |
  | id         | gAAAAABWvjYj-Zjfg8WXFaQnUd1DMYTBVrKw4h3fIagi5NoEmh21U72SrRv2trl |
  |            | JWFYhLi2_uPR31Igf6A8mH2Rw9kv_bxNo1jbLNPLGzW_u5FC7InFqx0yYtTwa1e |
  |            | eq2b0f6-18KZyQhs7F3teAta143kJEWuNEYET-y7u29y0be1_64KYkM7E       |
  | project_id | 343d245e850143a096806dfaefa9afdc                                |
  | user_id    | ac3377633149401296f6c0d92d79dc16                                |
  +------------+-----------------------------------------------------------------+
  ```

  

## 4.3、安装Glance

### 4.3.1、安装Glance步骤

- 作用：提供发现、注册和检索虚拟机镜像的功能。

- 创建数据库

  ```shell
  #登录数据库
  [root@controller /]# mysql -u root -p
  
  #创建glance数据库
  MariaDB [(none)]> CREATE DATABASE glance;
  
  #授权
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON glance.* TO 'glance'@'localhost' \
    IDENTIFIED BY 'GLANCE_DBPASS';
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON glance.* TO 'glance'@'%' \
    IDENTIFIED BY 'GLANCE_DBPASS';
  ```

  

- 获取管理员的CLI命令

  ```shell
  [root@controller /]# . admin-openrc
  ```

  

- 创建服务凭证：用户、角色信息，镜像服务端点

  ```shell
  #创建glance用户
  [root@controller /]# openstack user create --domain default --password-prompt glance
  #此处使用默认密码：GLANCE_PASS
  User Password:
  Repeat User Password:
  +---------------------+----------------------------------+
  | Field               | Value                            |
  +---------------------+----------------------------------+
  | domain_id           | default                          |
  | enabled             | True                             |
  | id                  | 3f4e777c4062483ab8d9edd7dff829df |
  | name                | glance                           |
  | options             | {}                               |
  | password_expires_at | None                             |
  +---------------------+----------------------------------+
  
  #为glance用户添加角色信息
  [root@controller /]# openstack role add --project service --user glance admin
  
  #创建glance服务实例
  [root@controller /]# openstack service create --name glance \
                       --description "OpenStack Image" image
  
  +-------------+----------------------------------+
  | Field       | Value                            |
  +-------------+----------------------------------+
  | description | OpenStack Image                  |
  | enabled     | True                             |
  | id          | 8c2c7f1b9b5049ea9e63757b5533e6d2 |
  | name        | glance                           |
  | type        | image                            |
  +-------------+----------------------------------+
  ```

  

- 创建镜像服务的API端点：一共需要创建3种

  ```shell
  #public类型
  [root@controller /]# openstack endpoint create --region RegionOne \
                       image public http://controller:9292
  
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 340be3625e9b4239a6415d034e98aace |
  | interface    | public                           |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | 8c2c7f1b9b5049ea9e63757b5533e6d2 |
  | service_name | glance                           |
  | service_type | image                            |
  | url          | http://controller:9292           |
  +--------------+----------------------------------+
  
  #admin类型
  [root@controller /]# openstack endpoint create --region RegionOne \
                       image admin http://controller:9292
  
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 0c37ed58103f4300a84ff125a539032d |
  | interface    | admin                            |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | 8c2c7f1b9b5049ea9e63757b5533e6d2 |
  | service_name | glance                           |
  | service_type | image                            |
  | url          | http://controller:9292           |
  +--------------+----------------------------------+
  
  #internal类型
  [root@controller /]# openstack endpoint create --region RegionOne \
                       image internal http://controller:9292
  
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | a6e4b153c2ae4c919eccfdbb7dceb5d2 |
  | interface    | internal                         |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | 8c2c7f1b9b5049ea9e63757b5533e6d2 |
  | service_name | glance                           |
  | service_type | image                            |
  | url          | http://controller:9292           |
  +--------------+----------------------------------+
  ```

  

- 安装、配置glance

   ```shell
   #安装glance
   [root@controller /]# yum install openstack-glance
   
   #配置glance-api.conf
   [root@controller /]# vim /etc/glance/glance-api.conf 
   [database]
   connection = mysql+pymysql://glance:GLANCE_DBPASS@controller/glance
   
   [keystone_authtoken]
   auth_uri = http://controller:5000
   auth_url = http://controller:5000
   memcached_servers = controller:11211
   auth_type = password
   project_domain_name = Default
   user_domain_name = Default
   project_name = service
   username = glance
   password = GLANCE_PASS
   
   [paste_deploy]
   flavor = keystone
   
   [glance_store]
   stores = file,http
   default_store = file
   filesystem_store_datadir = /var/lib/glance/images/
   
   #配置glance-registry.conf
   [root@controller /]# vim /etc/glance/glance-registry.conf
   [database]
   connection = mysql+pymysql://glance:GLANCE_DBPASS@controller/glance
   
   [keystone_authtoken]
   auth_uri = http://controller:5000
   auth_url = http://controller:5000
   memcached_servers = controller:11211
   auth_type = password
   project_domain_name = Default
   user_domain_name = Default
   project_name = service
   username = glance
   password = GLANCE_PASS
   
   [paste_deploy]
   flavor = keystone
   
   #填充lance数据库
   [root@controller /]# su -s /bin/sh -c "glance-manage db_sync" glance
   
   ```

  

- 配置环境

  开机自启动、启动glance

  ```shell
  [root@controller /]# systemctl enable openstack-glance-api.service \
                       openstack-glance-registry.service
  [root@controller /]# systemctl start openstack-glance-api.service \
                       openstack-glance-registry.service
  ```

  

### 4.3.2、验证是否成功安装

使用[CirrOS](http://launchpad.net/cirros)验证 Image 服务是否可用：使用OpenStack 部署一个小型 Linux 映像（eg：cirros-cloud）

```shell
#获取凭证
[root@controller /]# . admin-openrc

#下载镜像
[root@controller /]# wget http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img

#上传镜像
[root@controller /]# openstack image create "cirros" \
                     --file cirros-0.4.0-x86_64-disk.img \
                     --disk-format qcow2 --container-format bare \
                     --public
+------------------+------------------------------------------------------+
| Field            | Value                                                |
+------------------+------------------------------------------------------+
| checksum         | 133eae9fb1c98f45894a4e60d8736619                     |
| container_format | bare                                                 |
| created_at       | 2015-03-26T16:52:10Z                                 |
| disk_format      | qcow2                                                |
| file             | /v2/images/cc5c6982-4910-471e-b864-1098015901b5/file |
| id               | cc5c6982-4910-471e-b864-1098015901b5                 |
| min_disk         | 0                                                    |
| min_ram          | 0                                                    |
| name             | cirros                                               |
| owner            | ae7a98326b9c455588edd2656d723b9d                     |
| protected        | False                                                |
| schema           | /v2/schemas/image                                    |
| size             | 13200896                                             |
| status           | active                                               |
| tags             |                                                      |
| updated_at       | 2015-03-26T16:52:10Z                                 |
| virtual_size     | None                                                 |
| visibility       | public                                               |
+------------------+------------------------------------------------------+

[root@controller /]# openstack image list
+--------------------------------------+--------+--------+
| ID                                   | Name   | Status |
+--------------------------------------+--------+--------+
| 38047887-61a7-41ea-9b49-27987d5e8bb9 | cirros | active |
+--------------------------------------+--------+--------+
```



## 4.4、安装Nova

### 4.4.1、安装、配置控制节点

- 创建Nova数据库

  ```shell
  #进入数据库
  [root@controller /]# mysql -u root -p
  
  #创建数据库
  MariaDB [(none)]> CREATE DATABASE nova_api;
  MariaDB [(none)]> CREATE DATABASE nova;
  MariaDB [(none)]> CREATE DATABASE nova_cell0;
  
  #数据库授权
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON nova_api.* TO 'nova'@'localhost' \
    IDENTIFIED BY 'NOVA_DBPASS';
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON nova_api.* TO 'nova'@'%' \
    IDENTIFIED BY 'NOVA_DBPASS';
  
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON nova.* TO 'nova'@'localhost' \
    IDENTIFIED BY 'NOVA_DBPASS';
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON nova.* TO 'nova'@'%' \
    IDENTIFIED BY 'NOVA_DBPASS';
  
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON nova_cell0.* TO 'nova'@'localhost' \
    IDENTIFIED BY 'NOVA_DBPASS';
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON nova_cell0.* TO 'nova'@'%' \
    IDENTIFIED BY 'NOVA_DBPASS';
  ```

  

- 获取管理员的CLI命令

  ```shell
  [root@controller /]# . admin-openrc
  ```

  

- 创建计算服务凭证

  ```shell
  #创建nova用户
  [root@controller /]# openstack user create --domain default --password-prompt nova
  #使用默认密码：NOVA_PASS
  User Password:
  Repeat User Password:
  +---------------------+----------------------------------+
  | Field               | Value                            |
  +---------------------+----------------------------------+
  | domain_id           | default                          |
  | enabled             | True                             |
  | id                  | 8a7dbf5279404537b1c7b86c033620fe |
  | name                | nova                             |
  | options             | {}                               |
  | password_expires_at | None                             |
  +---------------------+----------------------------------+
  
  #为nova用户添加角色
  [root@controller /]# openstack role add --project service --user nova admin
  
  #创建nova服务实例
  [root@controller /]# openstack service create --name nova \
                       --description "OpenStack Compute" compute
  
  +-------------+----------------------------------+
  | Field       | Value                            |
  +-------------+----------------------------------+
  | description | OpenStack Compute                |
  | enabled     | True                             |
  | id          | 060d59eac51b4594815603d75a00aba2 |
  | name        | nova                             |
  | type        | compute                          |
  +-------------+----------------------------------+
  
  #创建计算API端点：3种
  [root@controller /]# openstack endpoint create --region RegionOne \
                       compute public http://controller:8774/v2.1
  +--------------+-------------------------------------------+
  | Field        | Value                                     |
  +--------------+-------------------------------------------+
  | enabled      | True                                      |
  | id           | 3c1caa473bfe4390a11e7177894bcc7b          |
  | interface    | public                                    |
  | region       | RegionOne                                 |
  | region_id    | RegionOne                                 |
  | service_id   | 060d59eac51b4594815603d75a00aba2          |
  | service_name | nova                                      |
  | service_type | compute                                   |
  | url          | http://controller:8774/v2.1               |
  +--------------+-------------------------------------------+
  
  [root@controller /]# openstack endpoint create --region RegionOne \
                       compute public http://controller:8774/v2.1
  +--------------+-------------------------------------------+
  | Field        | Value                                     |
  +--------------+-------------------------------------------+
  | enabled      | True                                      |
  | id           | 3c1caa473bfe4390a11e7177894bcc7b          |
  | interface    | public                                    |
  | region       | RegionOne                                 |
  | region_id    | RegionOne                                 |
  | service_id   | 060d59eac51b4594815603d75a00aba2          |
  | service_name | nova                                      |
  | service_type | compute                                   |
  | url          | http://controller:8774/v2.1               |
  +--------------+-------------------------------------------+
  
  [root@controller /]# openstack endpoint create --region RegionOne \
                       compute admin http://controller:8774/v2.1
  +--------------+-------------------------------------------+
  | Field        | Value                                     |
  +--------------+-------------------------------------------+
  | enabled      | True                                      |
  | id           | 38f7af91666a47cfb97b4dc790b94424          |
  | interface    | admin                                     |
  | region       | RegionOne                                 |
  | region_id    | RegionOne                                 |
  | service_id   | 060d59eac51b4594815603d75a00aba2          |
  | service_name | nova                                      |
  | service_type | compute                                   |
  | url          | http://controller:8774/v2.1               |
  +--------------+-------------------------------------------+
  
  
  #创建放置服务用户（placement）
  [root@controller /]# openstack user create --domain default --password-prompt placement
  #使用默认密码：PLACEMENT_PASS
  User Password:
  Repeat User Password:
  +---------------------+----------------------------------+
  | Field               | Value                            |
  +---------------------+----------------------------------+
  | domain_id           | default                          |
  | enabled             | True                             |
  | id                  | fa742015a6494a949f67629884fc7ec8 |
  | name                | placement                        |
  | options             | {}                               |
  | password_expires_at | None                             |
  +---------------------+----------------------------------+
  
  #为placement用户添加角色信息
  [root@controller /]# openstack role add --project service --user placement admin
  
  #在服务目录中创建Placement API条目
  [root@controller /]# openstack service create --name placement --description "Placement API" placement
  +-------------+----------------------------------+
  | Field       | Value                            |
  +-------------+----------------------------------+
  | description | Placement API                    |
  | enabled     | True                             |
  | id          | 2d1a27022e6e4185b86adac4444c495f |
  | name        | placement                        |
  | type        | placement                        |
  +-------------+----------------------------------+
  
  #创建Placement API服务端点：3种
  [root@controller /]#  openstack endpoint create --region RegionOne placement public http://controller:8780
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 2b1b2637908b4137a9c2e0470487cbc0 |
  | interface    | public                           |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | 2d1a27022e6e4185b86adac4444c495f |
  | service_name | placement                        |
  | service_type | placement                        |
  | url          | http://controller:8780           |
  +--------------+----------------------------------+
  
  [root@controller /]# openstack endpoint create --region RegionOne placement internal http://controller:8780
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 02bcda9a150a4bd7993ff4879df971ab |
  | interface    | internal                         |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | 2d1a27022e6e4185b86adac4444c495f |
  | service_name | placement                        |
  | service_type | placement                        |
  | url          | http://controller:8780           |
  +--------------+----------------------------------+
  
  [root@controller /]# openstack endpoint create --region RegionOne placement admin http://controller:8780
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 3d71177b9e0f406f98cbff198d74b182 |
  | interface    | admin                            |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | 2d1a27022e6e4185b86adac4444c495f |
  | service_name | placement                        |
  | service_type | placement                        |
  | url          | http://controller:8780           |
  +--------------+----------------------------------+
  ```

  

-  安装、配置Glance服务

  ```shell
  #安装Nova服务
  [root@controller /]# zypper install openstack-nova-api openstack-nova-scheduler \
                       openstack-nova-conductor openstack-nova-consoleauth \
                       openstack-nova-novncproxy openstack-nova-placement-api \
                       iptables
  
  #配置nova.conf
  [root@controller /]# vim /etc/nova/nova.conf 
  [DEFAULT]
  enabled_apis = osapi_compute,metadata
  
  [api_database]
  connection = mysql+pymysql://nova:NOVA_DBPASS@controller/nova_api
  
  [database]
  connection = mysql+pymysql://nova:NOVA_DBPASS@controller/nova
  
  [DEFAULT]
  transport_url = rabbit://openstack:RABBIT_PASS@controller
  
  [api]
  auth_strategy = keystone
  
  [keystone_authtoken]
  auth_url = http://controller:5000/v3
  memcached_servers = controller:11211
  auth_type = password
  project_domain_name = default
  user_domain_name = default
  project_name = service
  username = nova
  password = NOVA_PASS
  
  [DEFAULT]
  my_ip = 10.0.0.11
  use_neutron = True
  firewall_driver = nova.virt.firewall.NoopFirewallDriver
  
  [vnc]
  enabled = true
  server_listen = $my_ip
  server_proxyclient_address = $my_ip
  
  [glance]
  api_servers = http://controller:9292
  
  [oslo_concurrency]
  lock_path = /var/run/nova
  
  [placement]
  os_region_name = RegionOne
  project_domain_name = Default
  project_name = service
  auth_type = password
  user_domain_name = Default
  auth_url = http://controller:5000/v3
  username = placement
  password = PLACEMENT_PASS
  
  #填充Nova-api数据库
  [root@controller /]# su -s /bin/sh -c "nova-manage api_db sync" nova
  
  #注册cell0数据库
  [root@controller /]# su -s /bin/sh -c "nova-manage cell_v2 map_cell0" nova
  
  #创建cell1单元：此处会生成一个UUID
  [root@controller /]# su -s /bin/sh -c "nova-manage cell_v2 create_cell --name=cell1 --verbose" nova
  
  #填充nova数据库
  [root@controller /]# su -s /bin/sh -c "nova-manage db sync" nova
  
  #验证novecell0、cell1是否正确注册
  [root@controller /]# nova-manage cell_v2 list_cells
  +-------+--------------------------------------+
  | Name  | UUID                                 |
  +-------+--------------------------------------+
  | cell1 | 109e1d4b-536a-40d0-83c6-5f121b82b650 |
  | cell0 | 00000000-0000-0000-0000-000000000000 |
  +-------+--------------------------------------+
  ```

  

- 配置环境

  ```shell
  #重命名配置文件：将源配置文件拿过来使用
  [root@controller /]# mv /etc/apache2/vhosts.d/nova-placement-api.conf.sample /etc/apache2/vhosts.d/nova-placement-api.conf
  
  #重新加载apache服务
  [root@controller /]# systemctl reload apache2.service
  
  #开机自启动、启动nova服务
  [root@controller /]# systemctl enable openstack-nova-api.service \
                       openstack-nova-consoleauth.service openstack-nova-scheduler.service \
                       openstack-nova-conductor.service openstack-nova-novncproxy.service
  [root@controller /]# systemctl start openstack-nova-api.service \
                       openstack-nova-consoleauth.service openstack-nova-scheduler.service \
                       openstack-nova-conductor.service openstack-nova-novncproxy.service
  ```




### 4.4.2、安装、配置计算节点

- 安装、配置Nova-compute

  ```shell
  #安装软件
  [root@controller /]# yum install openstack-nova-compute
  
  #配置nova.conf
  [root@controller /]# vim /etc/nova/nova.conf
  [DEFAULT]
  enabled_apis = osapi_compute,metadata
  transport_url = rabbit://openstack:RABBIT_PASS@controller
  
  [api]
  auth_strategy = keystone
  
  [keystone_authtoken]
  auth_url = http://controller:5000/v3
  memcached_servers = controller:11211
  auth_type = password
  project_domain_name = default
  user_domain_name = default
  project_name = service
  username = nova
  password = NOVA_PASS
  
  [DEFAULT]
  #MANAGEMENT_INTERFACE_IP_ADDRESS：替换为本机的网卡IP
  my_ip = MANAGEMENT_INTERFACE_IP_ADDRESS
  use_neutron = True
  firewall_driver = nova.virt.firewall.NoopFirewallDriver
  
  [vnc]
  enabled = True
  server_listen = 0.0.0.0
  server_proxyclient_address = $my_ip
  novncproxy_base_url = http://controller:6080/vnc_auto.html
  
  [glance]
  api_servers = http://controller:9292
  
  [oslo_concurrency]
  lock_path = /var/lib/nova/tmp
  
  [placement]
  os_region_name = RegionOne
  project_domain_name = Default
  project_name = service
  auth_type = password
  user_domain_name = Default
  auth_url = http://controller:5000/v3
  username = placement
  password = PLACEMENT_PASS
  ```

  

- 配置环境

  ```shell
  #判断VMWware的虚拟机是否开启虚拟加速
  [root@controller /]# egrep -c '(vmx|svm)' /proc/cpuinfo
  
  #若上述命令的输出结果 = 0，则需要执行此步骤（说明未开启虚拟加速，需要配置nova.conf的virt_type ）
  [root@controller /]# vim  /etc/nova/nova.conf
  [libvirt]
  # 默认是KVM
  virt_type = qemu
  
  #设置开机自启动、启动nova-compute服务
  [root@controller /]# systemctl enable libvirtd.service openstack-nova-compute.service
  [root@controller /]# systemctl start libvirtd.service openstack-nova-compute.service
  ```

   

- 获取管理员的CLI命令

  ```shell
  [root@controller /]# . admin-openrc
  ```

  

- 将计算节点添加至cell数据库种

  ```shell
  #确认数据库中存在host主机
  [root@controller /]# openstack compute service list --service nova-compute
  +----+-------+--------------+------+-------+---------+----------------------------+
  | ID | Host  | Binary       | Zone | State | Status  | Updated At                 |
  +----+-------+--------------+------+-------+---------+----------------------------+
  | 1  | node1 | nova-compute | nova | up    | enabled | 2017-04-14T15:30:44.000000 |
  +----+-------+--------------+------+-------+---------+----------------------------+
  
  #发现计算主机
  [root@controller /]# su -s /bin/sh -c "nova-manage cell_v2 discover_hosts --verbose" nova
  Found 2 cell mappings.
  Skipping cell0 since it does not contain hosts.
  Getting compute nodes from cell 'cell1': ad5a5985-a719-4567-98d8-8d148aaae4bc
  Found 1 computes in cell: ad5a5985-a719-4567-98d8-8d148aaae4bc
  Checking host mapping for compute host 'compute': fe58ddc1-1d65-4f87-9456-bc040dc106b3
  Creating host mapping for compute host 'compute': fe58ddc1-1d65-4f87-9456-bc040dc106b3
  ```

  

## 4.5、安装Neutron

### 4.5.1、安装、配置控制节点

- 创建数据库

  ```shell
  #进入数据库
  [root@controller /]# mysql -u root -p
  
  #创建数据库
  MariaDB [(none)] CREATE DATABASE neutron;
  
  #授权
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON neutron.* TO 'neutron'@'localhost' \
    IDENTIFIED BY 'NEUTRON_DBPASS';
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON neutron.* TO 'neutron'@'%' \
    IDENTIFIED BY 'NEUTRON_DBPASS';
    
  ```

  

- 获取管理员的CLI命令

  ```shell
  [root@controller /]# . admin-openrc
  ```

  

- 创建服务凭证

  ```shell
  #创建Neutron用户
  [root@controller /]# openstack user create --domain default --password-prompt neutron
  #此处使用默认密码：NEUTRON_PASS
  User Password:
  Repeat User Password:
  +---------------------+----------------------------------+
  | Field               | Value                            |
  +---------------------+----------------------------------+
  | domain_id           | default                          |
  | enabled             | True                             |
  | id                  | fdb0f541e28141719b6a43c8944bf1fb |
  | name                | neutron                          |
  | options             | {}                               |
  | password_expires_at | None                             |
  +---------------------+----------------------------------+
  
  #添加Neutron用户的juese
  [root@controller /]# openstack role add --project service --user neutron admin
  
  #创建Neutron服务实例
  [root@controller /]# openstack service create --name neutron \
                       --description "OpenStack Networking" network
  +-------------+----------------------------------+
  | Field       | Value                            |
  +-------------+----------------------------------+
  | description | OpenStack Networking             |
  | enabled     | True                             |
  | id          | f71529314dab4a4d8eca427e701d209e |
  | name        | neutron                          |
  | type        | network                          |
  +-------------+----------------------------------+
  
  
  ```

   

- 创建网络服务API端点

  ```shell
  [root@controller /]# openstack endpoint create --region RegionOne \
                       network public http://controller:9696
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 85d80a6d02fc4b7683f611d7fc1493a3 |
  | interface    | public                           |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | f71529314dab4a4d8eca427e701d209e |
  | service_name | neutron                          |
  | service_type | network                          |
  | url          | http://controller:9696           |
  +--------------+----------------------------------+
  
  [root@controller /]# openstack endpoint create --region RegionOne \
                       network internal http://controller:9696
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 09753b537ac74422a68d2d791cf3714f |
  | interface    | internal                         |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | f71529314dab4a4d8eca427e701d209e |
  | service_name | neutron                          |
  | service_type | network                          |
  | url          | http://controller:9696           |
  +--------------+----------------------------------+
  
  [root@controller /]# openstack endpoint create --region RegionOne \
                       network admin http://controller:9696
  +--------------+----------------------------------+
  | Field        | Value                            |
  +--------------+----------------------------------+
  | enabled      | True                             |
  | id           | 1ee14289c9374dffb5db92a5c112fc4e |
  | interface    | admin                            |
  | region       | RegionOne                        |
  | region_id    | RegionOne                        |
  | service_id   | f71529314dab4a4d8eca427e701d209e |
  | service_name | neutron                          |
  | service_type | network                          |
  | url          | http://controller:9696           |
  +--------------+----------------------------------+
  ```

  

- 配置网络

  此处有两种网络可选：provider networks、self-service networks。本次使用self-service networks。

  - 安装组件

    ```shell
    [root@controller /]# yum install openstack-neutron openstack-neutron-ml2 \
                         openstack-neutron-linuxbridge ebtables
    ```

    

  - 配置neutron.conf

    ```shell
    [root@controller /]# vim /etc/neutron/neutron.conf
    [database]
    connection = mysql+pymysql://neutron:NEUTRON_DBPASS@controller/neutron
    
    [DEFAULT]
    core_plugin = ml2
    service_plugins = router
    allow_overlapping_ips = true
    
    [DEFAULT]
    transport_url = rabbit://openstack:RABBIT_PASS@controller
    auth_strategy = keystone
    
    [keystone_authtoken]
    auth_uri = http://controller:5000
    auth_url = http://controller:35357
    memcached_servers = controller:11211
    auth_type = password
    project_domain_name = default
    user_domain_name = default
    project_name = service
    username = neutron
    password = NEUTRON_PASS
    
    [DEFAULT]
    notify_nova_on_port_status_changes = true
    notify_nova_on_port_data_changes = true
    
    [nova]
    auth_url = http://controller:35357
    auth_type = password
    project_domain_name = default
    user_domain_name = default
    region_name = RegionOne
    project_name = service
    username = nova
    password = NOVA_PASS
    
    [oslo_concurrency]
    lock_path = /var/lib/neutron/tmp
    ```

    

  - 配置ml2_conf.ini

    配置模块化第 2 层 (ML2) 插件

    ```shell
    [root@controller /]# vim /etc/neutron/plugins/ml2/ml2_conf.ini
    [ml2] 
    type_drivers = flat,vlan,vxlan
    tenant_network_types = vxlan
    mechanism_drivers = linuxbridge,l2population
    extension_drivers = port_security
    
    [ml2_type_flat]
    flat_networks = provider
    
    [ml2_type_vxlan]
    vni_ranges = 1:1000
    
    [securitygroup]
    enable_ipset = true
    ```

    

  - 配置linuxbridge_agent.ini

    配置 Linux 网桥代理

    ```shell
    #配置文件
    [root@controller /]# vim /etc/neutron/plugins/ml2/linuxbridge_agent.ini
    [linux_bridge] 
    #PROVIDER_INTERFACE_NAME：替换为本机网卡名，eg：ens33
    physical_interface_mappings  =  provider:PROVIDER_INTERFACE_NAME
    
    [vxlan] 
    enable_vxlan  =  true 
    #OVERLAY_INTERFACE_IP_ADDRESS：替换为本机网卡IP
    local_ip  =  OVERLAY_INTERFACE_IP_ADDRESS 
    l2_population  =  true
    
    [securitygroup] 
    enable_security_group = true
    firewall_driver = neutron.agent.linux.iptables_firewall.IptablesFirewallDriver
    
    #验证linux系统内核是否支持网桥筛选器，下面的参数输出值均为1，说明支持。
    [root@controller /]# sysctl net.bridge.bridge-nf-call-iptables
    net.bridge.bridge-nf-call-iptables = 1
    [root@controller /]# sysctl net.bridge.bridge-nf-call-ip6tables
    net.bridge.bridge-nf-call-ip6tables = 1
    
    #（可选）若上述的参数不为1，则需要手动配置/etc/sysctl.conf
    [root@controller /]# vim /etc/sysctl.conf
    net.bridge.bridge-nf-call-iptables = 1
    net.bridge.bridge-nf-call-ip6tables = 1
    ```

    

  - 配置第三层代理

    为自助服务虚拟网络提供路由和 NAT 服务。

    ```shell
    [root@controller /]# vim /etc/neutron/l3_agent.ini
    [DEFAULT]
    interface_driver = linuxbridge
    ```

    

  - 配置DHCP代理

    为虚拟网络提供 DHCP 服务

    ```shell
    [root@controller /]# vim /etc/neutron/dhcp_agent.ini 
    [DEFAULT]
    interface_driver = linuxbridge
    dhcp_driver = neutron.agent.linux.dhcp.Dnsmasq
    enable_isolated_metadata = true
    ```

    

- 配置元数据代理

  ```shell
  [root@controller /]# vim /etc/neutron/metadata_agent.ini
  [DEFAULT]
  nova_metadata_host = controller
  metadata_proxy_shared_secret = METADATA_SECRET
  ```

  

- 配置Compute服务使用Networking服务

  即：配置Nova服务使用Neutron服务

  ```shell
  [root@controller /]# vim /etc/nova/nova.conf 
  [neutron]
  url = http://controller:9696
  auth_url = http://controller:35357
  auth_type = password
  project_domain_name = default
  user_domain_name = default
  region_name = RegionOne
  project_name = service
  username = neutron
  password = NEUTRON_PASS
  service_metadata_proxy = true
  metadata_proxy_shared_secret = METADATA_SECRET
  ```

  

- 配置环境

  ```shell
  #网络服务初始化脚本需要一个符号链接/etc/neutron/plugin.ini，
  #指向ML2插件配置文件/etc/neutron/plugins/ML2/ml2conf.ini
  [root@controller /]# ln -s /etc/neutron/plugins/ml2/ml2_conf.ini /etc/neutron/plugin.ini
  
  #填充数据库
  [root@controller /]# su -s /bin/sh -c "neutron-db-manage --config-file /etc/neutron/neutron.conf \
    --config-file /etc/neutron/plugins/ml2/ml2_conf.ini upgrade head" neutron
  
  #重启Nova-API服务
  [root@controller /]# systemctl restart openstack-nova-api.service
  
  #配置网络服务开机自启、启动网络服务
  [root@controller /]# systemctl enable neutron-server.service \
                       neutron-linuxbridge-agent.service neutron-dhcp-agent.service \
                       neutron-metadata-agent.service
  [root@controller /]# systemctl start neutron-server.service \
                       neutron-linuxbridge-agent.service neutron-dhcp-agent.service \
                       neutron-metadata-agent.service
                       
  #配置第三层服务开机自启、启动服务
  [root@controller /]# systemctl enable neutron-l3-agent.service
  [root@controller /]# systemctl start neutron-l3-agent.service
  ```

  

### 4.5.2、安装、配置计算节点

- 安装组件

  ```shell
  [root@controller /]# yum install openstack-neutron-linuxbridge ebtables ipset
  ```

  

- （该步骤可省略）配置通用组件

  该步骤可省略，因为all in one安装方式，计算节点和控制节点在同一台主机中，4.5.1中已经配置了neutron.conf文件。因此该步骤可省略。

  ```shell
  [root@controller /]# vim /etc/neutron/neutron.conf
  #（注意：使用all in one的安装方式，此处的数据库不可以注释，使用所有节点均在同一台主机中，注释掉会导致网络组件无法使用）
  #注释[database]的connection，此处不需要使用数据库
  [database]
  #connection=mysql+pymysql://neutron:NEUTRON_DBPASS@controller/neutron
  
  [DEFAULT]
  transport_url = rabbit://openstack:RABBIT_PASS@controller
  auth_strategy = keystone
  
  [keystone_authtoken]
  auth_uri = http://controller:5000
  auth_url = http://controller:35357
  memcached_servers = controller:11211
  auth_type = password
  project_domain_name = default
  user_domain_name = default
  project_name = service
  username = neutron
  password = NEUTRON_PASS
  
  [oslo_concurrency]
  lock_path = /var/lib/neutron/tmp
  ```

  

- （该步骤可省略）配置网络

  该步骤可省略，因为all in one安装方式，计算节点和控制节点在同一台主机中，4.5.1中已经配置了linuxbridge_agent.ini文件。因此该步骤可省略。

  ```shell
  #配置文件
  [root@controller /]# vim /etc/neutron/plugins/ml2/linuxbridge_agent.ini
  [linux_bridge] 
  #PROVIDER_INTERFACE_NAME：替换为本机网卡名，eg：ens33
  physical_interface_mappings  =  provider:PROVIDER_INTERFACE_NAME
  
  [vxlan] 
  enable_vxlan  =  true 
  #OVERLAY_INTERFACE_IP_ADDRESS：替换为本机网卡IP
  local_ip  =  OVERLAY_INTERFACE_IP_ADDRESS 
  l2_population  =  true
  
  [securitygroup] 
  enable_security_group = true
  firewall_driver = neutron.agent.linux.iptables_firewall.IptablesFirewallDriver
  
  #验证linux系统内核是否支持网桥筛选器，下面的参数输出值均为1，说明支持。
  [root@controller /]# sysctl net.bridge.bridge-nf-call-iptables
  net.bridge.bridge-nf-call-iptables = 1
  [root@controller /]# sysctl net.bridge.bridge-nf-call-ip6tables
  net.bridge.bridge-nf-call-ip6tables = 1
  
  #（可选）若上述的参数不为1，则需要手动配置/etc/sysctl.conf
  [root@controller /]# vim /etc/sysctl.conf
  net.bridge.bridge-nf-call-iptables = 1
  net.bridge.bridge-nf-call-ip6tables = 1
  
  ```

  

  验证网络是否配置好

  ```shell
  [root@controller /]# . admin-openrc
  
  #列出加载的扩展以验证neutron-server进程是否成功启动
  [root@controller /]# openstack extension list --network
  +---------------------------+---------------------------+----------------------------+
  | Name                      | Alias                     | Description                |
  +---------------------------+---------------------------+----------------------------+
  | Default Subnetpools       | default-subnetpools       | Provides ability to mark   |
  |                           |                           | and use a subnetpool as    |
  |                           |                           | the default                |
  | Availability Zone         | availability_zone         | The availability zone      |
  ......
  ......
  ```

  

- （该步骤可省略）配置Compute服务使用Networking服务

  即：配置Nova服务使用Neutron服务

  该步骤可省略，因为all in one安装方式，计算节点和控制节点在同一台主机中，4.5.1中已经配置了nova.conf 文件。因此该步骤可省略。

  ```shell
  [root@controller /]# vim /etc/nova/nova.conf 
  [neutron]
  url = http://controller:9696
  auth_url = http://controller:35357
  auth_type = password
  project_domain_name = default
  user_domain_name = default
  region_name = RegionOne
  project_name = service
  username = neutron
  password = NEUTRON_PASS
  ```

  

- 配置环境

  重启、开机自启动

  ```shell
   [root@controller /]# systemctl restart openstack-nova-compute.service
   [root@controller /]# systemctl enable neutron-linuxbridge-agent.service
   [root@controller /]# systemctl start neutron-linuxbridge-agent.service
  ```

  

## 4.6、安装Horizon

### 4.6.1、安装Horizon步骤

该组件选用手动安装的方式，不使用source的方式安装

- 安装dashboard软件

  ```shell
  [root@controller /]# yum install openstack-dashboard
  ```

  

- 配置local_settings

  ```shell
  [root@controller /]# vim /etc/openstack-dashboard/local_settings
  #配置主机名
  OPENSTACK_HOST = "controller"
  
  #'*'  允许所有主机访问
  ALLOWED_HOSTS = ['*', 'two.example.com']
  
  #新增SESSION_ENGINE
  SESSION_ENGINE = 'django.contrib.sessions.backends.cache'
  
  #修改原有的CACHES，为下面的部分
  CACHES = {
      'default': {
           'BACKEND': 'django.core.cache.backends.memcached.MemcachedCache',
           'LOCATION': 'controller:11211',
      }
  }
  
  OPENSTACK_KEYSTONE_URL = "http://%s:5000/v3" % OPENSTACK_HOST
  
  OPENSTACK_KEYSTONE_MULTIDOMAIN_SUPPORT = True
  
  #修改原有的OPENSTACK_API_VERSIONS为：
  OPENSTACK_API_VERSIONS = {
      "identity": 3,
      "image": 2,
      "volume": 2,
  }
  
  OPENSTACK_KEYSTONE_DEFAULT_DOMAIN = "Default"
  
  OPENSTACK_KEYSTONE_DEFAULT_ROLE = "user"
  
  #（不需要设置）只有在选择网络为provider networks时才需要设置该步骤
  OPENSTACK_NEUTRON_NETWORK = {
      ...
      'enable_router': False,
      'enable_quotas': False,
      'enable_distributed_router': False,
      'enable_ha_router': False,
      'enable_lb': False,
      'enable_firewall': False,
      'enable_vpn': False,
      'enable_fip_topology_check': False,
  }
  
  #（可选）设置时区
  TIME_ZONE = "TIME_ZONE"
  
  ```

  

- 配置openstack-dashboard.conf

  ```shell
  [root@controller /]# vim /etc/httpd/conf.d/openstack-dashboard.conf 
  #新增
  WSGIApplicationGroup %{GLOBAL}
  ```

  

- 配置环境

  ```shell
  #重启httpd、memcached服务
  [root@controller /]# systemctl restart httpd.service memcached.service
  ```



### 4.6.2、验证是否成功安装

- 在浏览器，登录下面的地址

  ```shell
  #方式1：
  http://controller/dashboard
  
  #方式2：
  http://本机网卡IP/dashboard
  ```

  

- 登录选项

  - domain：default

  - username：admin

  - password：ADMIN_PASS

    

- 结果 

  ![image-20210725162439337](OpenStack_安装笔记.assets/image-20210725162439337.png)

  

## 4.7、安装Cinder

### 4.7.1、安装、配置存储节点

- 安装LVM

  ```shell
  [root@controller /]# yum install lvm2 device-mapper-persistent-data
  
  #设置开机自启动
  [root@controller /]# systemctl enable lvm2-lvmetad.service
  [root@controller /]# systemctl start lvm2-lvmetad.service
  ```

  

  - 创建LVM物理卷

    注意：此处的sdb，表示的是第二块硬盘。（同理，sda表示第一块硬盘）

    ```shell
    [root@controller /]# pvcreate /dev/sdb
    Physical volume "/dev/sdb" successfully created
    ```

    

  - 创建 LVM 卷组`cinder-volumes`：

    ```shell
    [root@controller /]# vgcreate cinder-volumes /dev/sdb
    Volume group "cinder-volumes" successfully created
    ```

    

  - 配置lvm.conf

    ```shell
    [root@controller /]# vim /etc/lvm/lvm.conf
    #设置过滤器，a：接受。 r：拒绝
    #下面的filter，只需要选择一个进行配置即可。
    
    #存储节点使用LVM （本教程使用这个）
    filter  =  [ "a/sda/", "a/sdb/", "r/.*/"]
    
    #计算节点使用LVM
    filter  =  [ "a/sda/", "r/.*/"]
    ```

    

- 安装cinder组件

  ```shell
  [root@controller /]# yum install openstack-cinder targetcli python-keystone
  ```

  

  -  配置cinder.conf

    ```shell
    [root@controller /]# vim /etc/cinder/cinder.conf
    [database]
    connection = mysql+pymysql://cinder:CINDER_DBPASS@controller/cinder
    
    [DEFAULT]
    transport_url = rabbit://openstack:RABBIT_PASS@controller
    auth_strategy = keystone
    
    [keystone_authtoken]
    auth_uri = http://controller:5000
    auth_url = http://controller:5000
    memcached_servers = controller:11211
    auth_type = password
    project_domain_id = default
    user_domain_id = default
    project_name = service
    username = cinder
    password = CINDER_PASS
    
    [DEFAULT]
    #配置为本机的网卡IP
    my_ip = MANAGEMENT_INTERFACE_IP_ADDRESS
    
    [lvm]
    volume_driver = cinder.volume.drivers.lvm.LVMVolumeDriver
    volume_group = cinder-volumes
    iscsi_protocol = iscsi
    iscsi_helper = lioadm
    
    [DEFAULT]
    enabled_backends = lvm
    glance_api_servers = http://controller:9292
    
    [oslo_concurrency]
    lock_path = /var/lib/cinder/tmp
    ```

    

  - 配置环境

    ```shell
    [root@controller /]# systemctl enable openstack-cinder-volume.service target.service
    [root@controller /]# systemctl start openstack-cinder-volume.service target.service
    ```

    

### 4.7.2、安装、配置控制节点

在控制器节点上安装和配置块存储服务.

- 创建数据库

  ```shell
  #进入数据库
  [root@controller /]# mysql -u root -p
  
  #创建cinder数据库
  MariaDB [(none)]> CREATE DATABASE cinder;
  
  #数据库授权
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON cinder.* TO 'cinder'@'localhost' \
    IDENTIFIED BY 'CINDER_DBPASS';
  MariaDB [(none)]> GRANT ALL PRIVILEGES ON cinder.* TO 'cinder'@'%' \
    IDENTIFIED BY 'CINDER_DBPASS';
  ```

  

- 创建服务凭证

  ```shell
  #获取管理员的CLI命令
  [root@controller /]# . admin-openrc
  
  #创建cinde用户
  [root@controller /]#openstack user create --domain default --password-prompt cinder
  #此处使用默认密码：CINDER_PASS
  User Password:
  Repeat User Password:
  +---------------------+----------------------------------+
  | Field               | Value                            |
  +---------------------+----------------------------------+
  | domain_id           | default                          |
  | enabled             | True                             |
  | id                  | 9d7e33de3e1a498390353819bc7d245d |
  | name                | cinder                           |
  | options             | {}                               |
  | password_expires_at | None                             |
  +---------------------+----------------------------------+
  
  #为cinder用户添角色信息
  [root@controller /]# openstack role add --project service --user cinder admin
  
  #创建cinder服务实例
  #cinderv2
  [root@controller /]# openstack service create --name cinderv2 \
    --description "OpenStack Block Storage" volumev2
  +-------------+----------------------------------+
  | Field       | Value                            |
  +-------------+----------------------------------+
  | description | OpenStack Block Storage          |
  | enabled     | True                             |
  | id          | eb9fd245bdbc414695952e93f29fe3ac |
  | name        | cinderv2                         |
  | type        | volumev2                         |
  +-------------+----------------------------------+
  #cinderv3
  [root@controller /]# openstack service create --name cinderv3 \
    --description "OpenStack Block Storage" volumev3
  +-------------+----------------------------------+
  | Field       | Value                            |
  +-------------+----------------------------------+
  | description | OpenStack Block Storage          |
  | enabled     | True                             |
  | id          | ab3bbbef780845a1a283490d281e7fda |
  | name        | cinderv3                         |
  | type        | volumev3                         |
  +-------------+----------------------------------+
  
  
  ```

  

- 创建块存储服务API端点

  - 创建v2版本的端点

    ```shell
    [root@controller /]# openstack endpoint create --region RegionOne \
      volumev2 public http://controller:8776/v2/%\(project_id\)s
    +--------------+------------------------------------------+
    | Field        | Value                                    |
    +--------------+------------------------------------------+
    | enabled      | True                                     |
    | id           | 513e73819e14460fb904163f41ef3759         |
    | interface    | public                                   |
    | region       | RegionOne                                |
    | region_id    | RegionOne                                |
    | service_id   | eb9fd245bdbc414695952e93f29fe3ac         |
    | service_name | cinderv2                                 |
    | service_type | volumev2                                 |
    | url          | http://controller:8776/v2/%(project_id)s |
    +--------------+------------------------------------------+
    
    [root@controller /]# openstack endpoint create --region RegionOne \
      volumev2 internal http://controller:8776/v2/%\(project_id\)s
    +--------------+------------------------------------------+
    | Field        | Value                                    |
    +--------------+------------------------------------------+
    | enabled      | True                                     |
    | id           | 6436a8a23d014cfdb69c586eff146a32         |
    | interface    | internal                                 |
    | region       | RegionOne                                |
    | region_id    | RegionOne                                |
    | service_id   | eb9fd245bdbc414695952e93f29fe3ac         |
    | service_name | cinderv2                                 |
    | service_type | volumev2                                 |
    | url          | http://controller:8776/v2/%(project_id)s |
    +--------------+------------------------------------------+
    
    [root@controller /]# openstack endpoint create --region RegionOne \
      volumev2 admin http://controller:8776/v2/%\(project_id\)s
    +--------------+------------------------------------------+
    | Field        | Value                                    |
    +--------------+------------------------------------------+
    | enabled      | True                                     |
    | id           | e652cf84dd334f359ae9b045a2c91d96         |
    | interface    | admin                                    |
    | region       | RegionOne                                |
    | region_id    | RegionOne                                |
    | service_id   | eb9fd245bdbc414695952e93f29fe3ac         |
    | service_name | cinderv2                                 |
    | service_type | volumev2                                 |
    | url          | http://controller:8776/v2/%(project_id)s |
    +--------------+------------------------------------------+
    ```

    

  - 创建v3版本的端点

    ```shell
    [root@controller /]# openstack endpoint create --region RegionOne \
      volumev3 public http://controller:8776/v3/%\(project_id\)s
    +--------------+------------------------------------------+
    | Field        | Value                                    |
    +--------------+------------------------------------------+
    | enabled      | True                                     |
    | id           | 03fa2c90153546c295bf30ca86b1344b         |
    | interface    | public                                   |
    | region       | RegionOne                                |
    | region_id    | RegionOne                                |
    | service_id   | ab3bbbef780845a1a283490d281e7fda         |
    | service_name | cinderv3                                 |
    | service_type | volumev3                                 |
    | url          | http://controller:8776/v3/%(project_id)s |
    +--------------+------------------------------------------+
    
    [root@controller /]# openstack endpoint create --region RegionOne \
      volumev3 internal http://controller:8776/v3/%\(project_id\)s
    +--------------+------------------------------------------+
    | Field        | Value                                    |
    +--------------+------------------------------------------+
    | enabled      | True                                     |
    | id           | 94f684395d1b41068c70e4ecb11364b2         |
    | interface    | internal                                 |
    | region       | RegionOne                                |
    | region_id    | RegionOne                                |
    | service_id   | ab3bbbef780845a1a283490d281e7fda         |
    | service_name | cinderv3                                 |
    | service_type | volumev3                                 |
    | url          | http://controller:8776/v3/%(project_id)s |
    +--------------+------------------------------------------+
    
    [root@controller /]# openstack endpoint create --region RegionOne \
      volumev3 admin http://controller:8776/v3/%\(project_id\)s
    +--------------+------------------------------------------+
    | Field        | Value                                    |
    +--------------+------------------------------------------+
    | enabled      | True                                     |
    | id           | 4511c28a0f9840c78bacb25f10f62c98         |
    | interface    | admin                                    |
    | region       | RegionOne                                |
    | region_id    | RegionOne                                |
    | service_id   | ab3bbbef780845a1a283490d281e7fda         |
    | service_name | cinderv3                                 |
    | service_type | volumev3                                 |
    | url          | http://controller:8776/v3/%(project_id)s |
    +--------------+------------------------------------------+
    ```

    

- 安装cinder组件

  ```shell
  [root@controller /]# yum install openstack-cinder
  ```

  

  - 配置**cinder.conf**

    ```shell
    [root@controller /]# vim /etc/cinder/cinder.conf
    [database]
    connection = mysql+pymysql://cinder:CINDER_DBPASS@controller/cinder
    
    [DEFAULT]
    transport_url = rabbit://openstack:RABBIT_PASS@controller
    auth_strategy = keystone
    my_ip = 10.0.0.11
    
    [keystone_authtoken]
    auth_uri = http://controller:5000
    auth_url = http://controller:5000
    memcached_servers = controller:11211
    auth_type = password
    project_domain_id = default
    user_domain_id = default
    project_name = service
    username = cinder
    password = CINDER_PASS
    
    [oslo_concurrency]
    lock_path = /var/lib/cinder/tmp
    
    
    ```

    

  - 填充数据库

    ```shell
    [root@controller /]# su -s /bin/sh -c "cinder-manage db sync" cinder
    ```

    

- 配置Nova使用块存储

  ```shell
  [root@controller /]# vim /etc/nova/nova.conf 
  ```

  

- 配置环境

  ```shell
  [root@controller /]# systemctl restart openstack-nova-api.service
  [root@controller /]# systemctl enable openstack-cinder-api.service openstack-cinder-scheduler.service
  [root@controller /]# systemctl start openstack-cinder-api.service openstack-cinder-scheduler.service
  ```

  

### 4.7.3、（可选）安装、配置备份服务

本教程没有安装swift组件，因此可省略此步骤

- 安装cinder

  ```shell
  [root@controller /]# yum install openstack-cinder
  ```

  - 配置cinder.conf

    ```shell
    [root@controller /]# /etc/cinder/cinder.conf
    [DEFAULT]
    backup_driver = cinder.backup.drivers.swift
    #替换SWIFT_URL为；为对象存储服务的 URL
    #可通过命令查看具体的替换数值：openstack catalog show object-store
    backup_swift_url = SWIFT_URL
    ```

  

- 配置环境

  ```shell
  [root@controller /]# systemctl enable openstack-cinder-backup.service
  [root@controller /]# systemctl start openstack-cinder-backup.service
  ```



### 4.7.4、验证是否成功安装

```shell
[root@controller /]# . admin-openrc

[root@controller /]# openstack volume service list
#此处未安装swift组件，因此只有存储节点、控制节点
+------------------+------------+------+---------+-------+----------------------------+
| Binary           | Host       | Zone | Status  | State | Updated_at                 |
+------------------+------------+------+---------+-------+----------------------------+
| cinder-scheduler | controller | nova | enabled | up    | 2016-09-30T02:27:41.000000 |
| cinder-volume    | block@lvm  | nova | enabled | up    | 2016-09-30T02:27:46.000000 |
+------------------+------------+------+---------+-------+----------------------------+
```





# 5、安装ceph

参考连接：

- 官网：
  - 下载ceph：https://docs.ceph.com/en/mimic/start/quick-start-preflight/#rhel-centos
  - 本地安装ceph：https://docs.ceph.com/en/mimic/start/quick-ceph-deploy/
  - ceph对接enstack：https://docs.ceph.com/en/latest/rbd/rbd-openstack/
- 博客（ceph对接openstack）：https://blog.csdn.net/ygtlovezf/article/details/78983249

## 5.1、下载ceph工具

- 安装存储库软件包

  ```shell
  sudo yum install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
  ```

  

- 配置镜像源

  ```shell
  cat << EOM > /etc/yum.repos.d/ceph.repo
  [ceph-noarch]
  name=Ceph noarch packages
  baseurl=https://download.ceph.com/rpm-luminous/el7/noarch
  enabled=1
  gpgcheck=1
  type=rpm-md
  gpgkey=https://download.ceph.com/keys/release.asc
  EOM
  ```

  

- 更新软件包、安装ceph部署工具

  ```shell
  [root@controller /]# sudo yum update
  [root@controller /]# sudo yum install ceph-deploy
  ```

  

## 5.2、本地安装ceph

- 创建my-cluster文件夹：存放后面生成的密钥文件

  ```shell
  [root@controller /]# mkdir my-cluster
  [root@controller /]# cd my-cluster  
  ```

  

- 创建ceph集群

  ```shell
  [root@controller my-cluster]# ceph-deploy new controller
  ```

  

- 各个节点安装ceph

  本示例，只有一个controller节点（因为其他的节点（eg：compute等）都安装在controller节点中中）

  ```shell
  [root@controller my-cluster]# ceph-deploy install controller
  ```

  - 执行上述命令时，可能出现的错误

    ```shell
    #错误
    [node1][WARNIN] ensuring that /etc/yum.repos.d/ceph.repo contains a high priority
    [ceph_deploy][ERROR ] RuntimeError: NoSectionError: No section: 'ceph'
    
    #解决方式
    [root@controller my-cluster]# yum remove ceph-release
    [root@controller my-cluster]# ceph-deploy install controller
    ```

    

- 部署监控，并获取证书（密钥）

  该命令在当前路径下生成证书，因此必须在指定的目录中执行，否则后面的命令会找不到对用的证书

  ```shell
  [root@controller my-cluster]# ceph-deploy mon create-initial
  ceph.client.admin.keyring
  ceph.bootstrap-mgr.keyring
  ceph.bootstrap-osd.keyring
  ceph.bootstrap-mds.keyring
  ceph.bootstrap-rgw.keyring
  ceph.bootstrap-rbd.keyring
  ```

  

- 使用ceph deploy将配置文件和管理密钥复制到管理节点和ceph节点，以便可以使用ceph CLI，而无需在每次执行命令时指定监视器地址和ceph.client.admin.keyring

  ```shell
  [root@controller my-cluster]# ceph-deploy admin controller
  ```

  

- 设置开机自启动

  ```shell
  [root@controller my-cluster]# ceph-deploy mgr create controller
  ```

  

- 添加OSD：

  将挂载的硬盘添加为OSD，一共5块硬盘

  ```shell
  [root@controller my-cluster]# ceph-deploy osd create --data /dev/sdc controller
  [root@controller my-cluster]# ceph-deploy osd create --data /dev/sdd controller
  [root@controller my-cluster]# ceph-deploy osd create --data /dev/sde controller
  [root@controller my-cluster]# ceph-deploy osd create --data /dev/sdf controller
  [root@controller my-cluster]# ceph-deploy osd create --data /dev/sdg controller
  ```

  - 执行上述命令时，可能出现的错误

    ```shell
    #错误
    [controller][WARNIN]  stderr: purged osd.0
    [controller][WARNIN] -->  RuntimeError: command returned non-zero exit status: 5
    [controller][ERROR ] RuntimeError: command returned non-zero exit status: 1
    [ceph_deploy.osd][ERROR ] Failed to execute command: /usr/sbin/ceph-volume --cluster ceph lvm create --bluestore --data /dev/sdc
    [ceph_deploy][ERROR ] GenericError: Failed to create 1 OSDs
    
    #原因：没有添加空的磁盘
    
    #解决方式：在/etc/lvm/lvm.conf的filter添加其他磁盘：eg：添加c~g，5个磁盘
    #查看挂载的硬盘（作为OSD的硬盘）
    [root@controller my-cluster]# fdisk -l
    #修改配置
    [root@controller my-cluster]# vim /etc/lvm/lvm.conf
    filter = [ "a/sda/", "a/sdb/","a/sdc/","a/sdd/","a/sde/","a/sdf/","a/sdg/", "r/.*/"]
    ```

    

- 查看集群是否健康

  ```shell
  [root@controller my-cluster]# ssh controller ceph health
  root@controller's password: 
  HEALTH_OK
  ```

  

## 5.3、ceph接入openstack

- 配置默认的size副本大小

  安装ceph后，默认的size大小为3，但是本示例的副本数量为1，因为只有一个机器。

  ```shell
  [root@controller my-cluster]# vim /etc/ceph/ceph.conf 
  [global]
  fsid = 1cde6d0b-7e92-41fb-93a7-63763daefd76
  mon_initial_members = controller
  mon_host = 192.168.83.139
  auth_cluster_required = cephx
  auth_service_required = cephx
  auth_client_required = cephx
  
  #新增部分
  osd_pool_default_size = 1
  osd_pool_default_min_size = 1
  
  #重启系统，使配置生效（否则后面创建的池，无法使用）
  [root@controller my-cluster]# reboot
  ```

  

- 创建、初始化池

  创建池的时候，需要计算pg_num、pgp_num

  ```shell
  #创建池
  [root@controller my-cluster]# ceph osd pool create volumes 256
  [root@controller my-cluster]# ceph osd pool create images 256
  [root@controller my-cluster]# ceph osd pool create backups 256
  [root@controller my-cluster]# ceph osd pool create vms 256
  
  #初始化池
  [root@controller my-cluster]# rbd pool init volumes
  [root@controller my-cluster]# rbd pool init images
  [root@controller my-cluster]# rbd pool init backups
  [root@controller my-cluster]# rbd pool init vms
  
  #查看已经创建的池
  [root@controller my-cluster]# ceph osd lspools
  #查看size的大小
  [root@controller my-cluster]# ceph osd pool get volumes size
  #查看pg_num
  [root@controller my-cluster]# ceph osd pool get volumes pg_num
  #查看pgp_num
  [root@controller my-cluster]# ceph osd pool get volumes pgp_num
  ```

  

- 配置、安装openstack ceph客户端

  ```shell
  #配置：将my-cluster的ceph.conf数据，复制进/etc/ceph/ceph.conf
  [root@controller my-cluster]# ssh controller sudo tee /etc/ceph/ceph.conf < ceph.conf
  
  #安装
  [root@controller my-cluster]# sudo yum install python-rbd
  [root@controller my-cluster]# sudo yum install ceph-common
  ```

  

- 设置 CEPH 客户端身份验证（**重点**，若/etc/ceph/ceph.conf中不存在cephx值，则无需执行该步骤）

  - 在ceph中创建cinder、glance等用户，并进行权限控制

    （注意：官网教程中的这部分是不能使用的）

    ```shell
    [root@controller my-cluster]# ceph auth get-or-create client.cinder mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=volumes, allow rwx pool=vms, allow rx pool=images'
    
    [root@controller my-cluster]# ceph auth get-or-create client.glance mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=images'
    
    [root@controller my-cluster]# ceph auth get-or-create client.cinder-backup mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=backups'
    ```

    

  - 将上面生成的keyring文件，保存在相应的节点上，并修改为相应的权限

    ```shell
    [root@controller my-cluster]# ceph auth get-or-create client.glance | ssh controller sudo tee /etc/ceph/ceph.client.glance.keyring
    ssh controller sudo chown glance:glance /etc/ceph/ceph.client.glance.keyring
    
    [root@controller my-cluster]# ceph auth get-or-create client.cinder | ssh controller sudo tee /etc/ceph/ceph.client.cinder.keyring
    ssh controller sudo chown cinder:cinder /etc/ceph/ceph.client.cinder.keyring
    
    [root@controller my-cluster]# ceph auth get-or-create client.cinder-backup | ssh controller sudo tee /etc/ceph/ceph.client.cinder-backup.keyring
    ssh controller sudo chown cinder:cinder /etc/ceph/ceph.client.cinder-backup.keyring
    ```

    

  - 在nova-compute节点上保存和cinder-volume相同的keyring

     ```shell
     [root@controller my-cluster]# ceph auth get-or-create client.cinder | ssh controller sudo tee /etc/ceph/ceph.client.cinder.keyring
     ```

    

  -  在运行节点上，创建密钥的临时副本

    ```shell
    [root@controller my-cluster]# ceph auth get-key client.cinder | ssh controller tee client.cinder.key
    ```

    

  -  生成、添加密钥

    ```shell
    #生成密钥
    [root@controller my-cluster]# uuidgen
    457eb676-33da-42ec-9a8c-9293d545c337
    
    #修改secret.xml文件
    [root@controller my-cluster]# cat > secret.xml <<EOF
    <secret ephemeral='no' private='no'>
      <uuid>457eb676-33da-42ec-9a8c-9293d545c337</uuid>
      <usage type='ceph'>
        <name>client.cinder secret</name>
      </usage>
    </secret>
    EOF
    
    #查看修改后的密钥
    [root@controller my-cluster]# sudo virsh secret-define --file secret.xml
    Secret 457eb676-33da-42ec-9a8c-9293d545c337 created
    
    #将密钥添加到libvirt并删除密钥的临时副本
    #（注意：需要修改client.cinder.key、secret.xml文件的路径，否则会找不到）
    [root@controller my-cluster]# sudo virsh secret-set-value --secret 457eb676-33da-42ec-9a8c-9293d545c337 --base64 $(cat /root/client.cinder.key) && rm /root/client.cinder.key secret.xml
    ```

    

- 配置OpenStack使用ceph，主要是配置各个组件的配置文件

  - 配置glance

    ```shell
    [root@controller my-cluster]# vim /etc/glance/glance-api.conf
    
    [DEFAULT]
    show_image_direct_url = True
    
    [glance_store]
    stores = rbd,file,http
    default_store = rbd
    rbd_store_pool = images
    rbd_store_user = glance
    rbd_store_ceph_conf = /etc/ceph/ceph.conf
    rbd_store_chunk_size = 8
    
    #注释掉原来LVM配置的：stores、default_store
    #stores = file,http
    #default_store = file
    filesystem_store_datadir = /var/lib/glance/images/
    
    [paste_deploy]
    flavor = keystone
    ```

    

  - 配置cinder

    ```shell
    [root@controller my-cluster]# vim /etc/cinder/cinder.conf
    
    [DEFAULT]
    #添加ceph
    enabled_backends = lvm,ceph
    glance_api_version = 2
    
    
    [ceph]
    #配置ceph后端名称
    volume_backend_name = ceph
    #配置驱动
    volume_driver = cinder.volume.drivers.rbd.RBDDriver
    rbd_pool = volumes
    rbd_ceph_conf = /etc/ceph/ceph.conf
    rbd_flatten_volume_from_snapshot = false
    rbd_max_clone_depth = 5
    rbd_store_chunk_size = 4
    rados_connect_timeout = -1
    #这是之前创建的cinder用户，以及UUID
    rbd_user = cinder
    rbd_secret_uuid = 457eb676-33da-42ec-9a8c-9293d545c337
    
    [lvm]
    #（新增）配置lvm后端名称
    volume_backend_name = lvm
    volume_driver = cinder.volume.drivers.lvm.LVMVolumeDriver
    volume_group = cinder-volumes
    iscsi_protocol = iscsi
    iscsi_helper = lioadm
    
    ```

    

  - 配置nova

    ```shell
    [root@controller my-cluster]# vim /etc/nova/nova.conf
    [libvirt]
    rbd_user = cinder
    rbd_secret_uuid = 457eb676-33da-42ec-9a8c-9293d545c337
    #images_type = rbd
    #images_rbd_pool = vms
    #images_rbd_ceph_conf = /etc/ceph/ceph.conf
    #disk_cachemodes="network=writeback"
    
    ```

    

- 重启相关服务

  ```shell
  [root@controller my-cluster]# sudo glance-control api restart
  [root@controller my-cluster]# sudo service nova-compute restart
  [root@controller my-cluster]# sudo service cinder-volume restart
  [root@controller my-cluster]# sudo service cinder-backup restart
  ```



- 至此，openstack与ceph对接完毕，但是还是无法使用ceph，需要手动创建ceph类型的存储。（详见下面的步骤）



## 5.4、使用ceph

- 在cinder中创建ceph的类型的存储

  ```shell
  #创建存储类型
  [root@controller openstack]# cinder type-create ceph
  +--------------------------------------+------+-------------+-----------+
  | ID                                   | Name | Description | Is_Public |
  +--------------------------------------+------+-------------+-----------+
  | df9c985c-1bb7-4d31-b040-9c8b97d37957 | ceph | -           | True      |
  +--------------------------------------+------+-------------+-----------+
  
  #给存储类型，设置后端名称
  [root@controller openstack]# cinder type-key ceph set volume_backend_name=ceph
  ```

  

- 查看当前已有的cinder存储类型

  ```shell
  [root@controller openstack]# cinder extra-specs-list
  +--------------------------------------+------+---------------------------------+
  | ID                                   | Name | extra_specs                     |
  +--------------------------------------+------+---------------------------------+
  | a915e3ff-5b25-4113-a1c2-6d1f6c388e9c | lvm  | {'volume_backend_name': 'lvm'}  |
  | df9c985c-1bb7-4d31-b040-9c8b97d37957 | ceph | {'volume_backend_name': 'ceph'} |
  +--------------------------------------+------+---------------------------------+
  ```

  

  dashboard中显示的卷类型：lvm、ceph

  ![image-20210802175929626](OpenStack_安装笔记.assets/image-20210802175929626.png) 

  

  

- 查看cinder的服务

  ```shell
  [root@controller openstack]# cinder service-list
  +------------------+-----------------+------+---------+-------+----------------------------+-----------------+
  | Binary           | Host            | Zone | Status  | State | Updated_at                 | Disabled Reason |
  +------------------+-----------------+------+---------+-------+----------------------------+-----------------+
  | cinder-backup    | controller      | nova | enabled | up    | 2021-08-02T08:41:24.000000 | -               |
  | cinder-scheduler | controller      | nova | enabled | up    | 2021-08-02T08:41:22.000000 | -               |
  | cinder-volume    | controller@ceph | nova | enabled | up    | 2021-08-02T08:41:25.000000 | -               |
  | cinder-volume    | controller@lvm  | nova | enabled | up    | 2021-08-02T08:41:19.000000 | -               |
  +------------------+-----------------+------+---------+-------+----------------------------+-----------------+
  ```

  
  
- 查看ceph集群挂载的osd

  ```shell
  [root@controller openstack]# ceph osd tree
  ID CLASS WEIGHT  TYPE NAME           STATUS REWEIGHT PRI-AFF 
  -1       0.04898 root default                                
  -3       0.04898     host controller                         
   0   hdd 0.00980         osd.0           up  1.00000 1.00000 
   1   hdd 0.00980         osd.1           up  1.00000 1.00000 
   2   hdd 0.00980         osd.2           up  1.00000 1.00000 
   3   hdd 0.00980         osd.3           up  1.00000 1.00000 
   4   hdd 0.00980         osd.4           up  1.00000 1.00000 
  ```
  
  
  
- **<font color='red'>注意：</font>**openstack对接完ceph，创建实例镜像时，需要注意原来上传的镜像是否可用，不可用的话，就需要重新上传镜像。否则无法创建实例。

  
  
- 创建实例，并使用ceph中的卷

  创建好实例后，并在项目中新建一个ceph类型的卷，然后把卷挂载到实例中
  
  - 新建实例
  
    ![image-20210803103346575](OpenStack_安装笔记.assets/image-20210803103346575.png)
  
  - 挂载ceph卷
  
    ![image-20210803103447726](OpenStack_安装笔记.assets/image-20210803103447726.png)
  
  

# 6、Dashboard创建虚拟机流程

使用OpenStack的UI界面，基于镜像来创建虚拟机实例。仅当OpenStack可以创建虚拟机时，才算所有组件安装成功。

- **创建实例类型**

  这个实例类型是之后用来创建实例的模板，可根据自己的需求进行配置。

  ![image-20210729183434196](OpenStack_安装笔记.assets/image-20210729183434196.png)

  

- **上传镜像**

  - 需要本地下载好Centos7镜像

    ```shell
    本示例需要下载CentOS7镜像，下载地址：http://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-20140929_01.qcow2
    ```

  - 使用命令上传

    ```shell
    [root@controller Downloads]# openstack image create "centos7" \
                                 --file CentOS-7-x86_64-GenericCloud-20140929_01.qcow2 \
                                 --disk-format qcow2 --container-format bare \
                                 --public    
    ```

  - 查看上传结果

    - 使用命令查看

      ```shell
      [root@controller Downloads]# openstack image list
      +--------------------------------------+-----------+--------+
      | ID                                   | Name      | Status |
      +--------------------------------------+-----------+--------+
      | c3587da6-1e0b-44ef-8c84-f3e4aaab1790 | centos7   | active |
      +--------------------------------------+-----------+--------+
      ```

    - 使用dashboard

      ![image-20210729171440091](OpenStack_安装笔记.assets/image-20210729171440091.png)

  

- **创建外部网络信息** 

  外部网络：

  - 查看本机的网关IP

    ```shell
    #网关IP为：192.168.83.2
    [root@controller chris]# route -n 
    Kernel IP routing table
    Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
    0.0.0.0         192.168.83.2    0.0.0.0         UG    0      0        0 ens33
    169.254.0.0     0.0.0.0         255.255.0.0     U     1002   0        0 ens33
    169.254.0.0     0.0.0.0         255.255.0.0     U     1003   0        0 ens37
    192.168.83.0    0.0.0.0         255.255.255.0   U     0      0        0 ens33
    192.168.83.0    0.0.0.0         255.255.255.0   U     0      0        0 ens37
    192.168.122.0   0.0.0.0         255.255.255.0   U     0      0        0 virbr0
    
    ```

    

  - 创建网络

    ![image-20210729175033738](OpenStack_安装笔记.assets/image-20210729175033738.png)

    

  - 配置子网

    - 子网的网络IP：必须和外部网络在同个网段。即：和本机网卡IP在同一个网段。
    - 子网的网关：必须和本机的网关IP一致。

    ![image-20210729172645824](OpenStack_安装笔记.assets/image-20210729172645824.png)

    

  - 配置子网详情

    分配地址池为：浮动IP池。多个IP采用  `,`  分割 

    ![image-20210729175227272](OpenStack_安装笔记.assets/image-20210729175227272.png)

    

- **创建内部网络**

  内部网络：项目中的虚拟机使用的网络。

  - 创建网络

    ![image-20210729175631736](OpenStack_安装笔记.assets/image-20210729175631736.png)

  - 配置子网

    子网的网络地址可自定义，不需要限制网段。创建的虚拟机网卡地址将在包含于这个网络地址中。

    ![image-20210729175825235](OpenStack_安装笔记.assets/image-20210729175825235.png)

  - 子网详情

    配置浮动IP

    ![image-20210729180019076](OpenStack_安装笔记.assets/image-20210729180019076.png)

  

- **创建路由**

  - 新建路由

    该路由是项目中的路由，用于连接内部网络和外部网络的。

    ![image-20210729180248904](OpenStack_安装笔记.assets/image-20210729180248904.png)

  - 配置路由接口

    此处的IP地址配置为：内部网络的网关IP

    ![image-20210729180542431](OpenStack_安装笔记.assets/image-20210729180542431.png)

  - 创建好的网络，可查看到网络拓扑图

    <img src="OpenStack_安装笔记.assets/image-20210729180915810.png" alt="image-20210729180915810" style="zoom: 67%;" />

  

- **创建实例（创建虚拟机）**

  - 配置实例详情

    ![image-20210729181258800](OpenStack_安装笔记.assets/image-20210729181258800.png)

  - 配置镜像源

    选择之前上传的镜像：centos7

    ![image-20210729183156120](OpenStack_安装笔记.assets/image-20210729183156120.png)

  - 配置实例类型

    选择之前配置的实例类型（实例模板），用于配置所创建虚拟机的CPU、RAM、ROM等参数。

    ![image-20210729183216572](OpenStack_安装笔记.assets/image-20210729183216572.png)

  - 配置网络

    选择之前配置的内部网络。 

    ![image-20210729181557293](OpenStack_安装笔记.assets/image-20210729181557293.png)

  - 其余选项，均采用默认值即可。

  - 查看创建好的虚拟机

    ![image-20210729183007326](OpenStack_安装笔记.assets/image-20210729183007326.png)

  - 登录创建虚拟机的控制台

    <img src="OpenStack_安装笔记.assets/image-20210729182924043.png" alt="image-20210729182924043" style="zoom:80%;" />



- **删除所创建的虚机流程：**
  - 流程： 关闭虚拟机 =》删除实例 =》 删除实例类型 =》 删除路由 =》 删除网络（外部、内部网络） 
  - 若，不按照该流程，则可能会出现部分组件删除失败的情况。**错误：**删除失败，该组件正在使用中。



# 7、常见问题

**1、安装keystone后，执行命令：openstack domain create --description "An Example Domain" example**

**出现错误：**

```shell
Failed to discover available identity versions when contacting http://controller:5000/v3. Attempting to parse version from URL.
Unable to establish connection to http://controller:5000/v3/auth/tokens: HTTPConnectionPool(host='controller', port=5000): Max retries exceeded with url: /v3/auth/tokens (Caused by NewConnectionError('<urllib3.connection.HTTPConnection object at 0x7fa20bba02b0>: Failed to establish a new connection: [Errno -3] Temporary failure in name resolution',))
```

**解决方式：**

关闭SELinux，然后使用命令setnforce 0（或者重启系统）



**2、配置网卡时，需要注意的事项**

配置网络时，需要先使用dhcp，获取虚拟机网卡ip，然后再使用static的方式，配置网卡ip。如果直接使用dhcp的方式，可能会导致该ip无法使用，无法ping通。



**3.2、安装OpenStack各个组件时，使用命令： yum install xxxx -y。**

**出现错误：**

```shell
Error downloading packages:
  python-dogpile-cache-0.6.2-1.el7.noarch: [Errno 256] No more mirrors to try.
  python-keyring-5.7.1-1.el7.noarch: [Errno 256] No more mirrors to try.
  python-dogpile-core-0.4.1-2.el7.noarch: [Errno 256] No more mirrors to try.
  python-cmd2-0.6.8-8.el7.noarch: [Errno 256] No more mirrors to try.
```

**解决方式：**
（1）方式1：再次执行安装命令 yum install python-openstackclient -y
（2）方式2：yum update =》yum install python-openstackclient -y
若还是出现错误：Error downloading packages: lttng-ust-2.10.0-1.el7.x86_64: [Errno 256] No more mirrors to try.
使用yum clean cache、yum makecache，然后再次进行yum update



**4、horizon登录注意事项**

- 地址：http://controller/dashboard
  不行的话，换成: http://controller对用的IP地址/dashboard
- 登录填写的资料：
  domain：default
  username：admin
  password：ADMIN_PASS



**5、禁用 SELinux**

- **否则会出现错误：**the request you have made requires authentication(HTTPConnection 401)

- **解决方式：**

  ```shell
  [root@controller /]# vim /etc/selinux/config
  SELINUX=disabled 
  SELINUXTYPE=targeted
  ```

  然后setnforce 0（或者重启系统）



**6、VMWare挂载磁盘的注意事项** 
**查看挂载：**

```shell
[root@controller /]# fdisk -l 
```

**若新增的磁盘未挂载上，处理方式：**

```shell
#1、查看当前有的磁盘文件：
[root@controller /]# cd /sys/class/scsi_host/ 

#2、获取当前目录下的所有文件名
[root@controller /]# ls
host0 host1 host2

#3、逐个执行上面文件中出现的文件
[root@controller /]# echo "- - -" >  /sys/class/scsi_host/host0/scan
[root@controller /]# echo "- - -" >  /sys/class/scsi_host/host1/scan
[root@controller /]# echo "- - -" >  /sys/class/scsi_host/host2/scan
```



**7、重启虚拟机后，无法使用openstack命令**

- **原因：**

  因为使用export设置环境变量，因此该环境变量只限定于当前执行命令的终端使用。系统关机、开启新的终端均无法使用这个环境变量。

- **解决方式：**

  需要运行源文件命令，添加环境变量。

  ```shell
  [root@controller /]# . admin-openrc
  ```

  

  ```shell
  #admin-openrc文件中的内容
  [root@controller /]# vim admin-openrc
  export OS_PROJECT_DOMAIN_NAME=Default
  export OS_USER_DOMAIN_NAME=Default
  export OS_PROJECT_NAME=admin
  export OS_USERNAME=admin
  export OS_PASSWORD=ADMIN_PASS
  export OS_AUTH_URL=http://controller:5000/v3
  export OS_IDENTITY_API_VERSION=3
  export OS_IMAGE_API_VERSION=2
  ```



**8、执行openstack volume service list时，cinder-volume运行失败**

```shell
[root@controller /]# openstack volume service list

+------------------+------------+------+---------+-------+----------------------------+
| Binary           | Host       | Zone | Status  | State | Updated_at                 |
+------------------+------------+------+---------+-------+----------------------------+
| cinder-scheduler | controller | nova | enabled | up    | 2016-09-30T02:27:41.000000 |
| cinder-volume    | block@lvm  | nova | disabled| down  | 2016-09-30T02:27:46.000000 |
+------------------+------------+------+---------+-------+----------------------------+
```

- **原因：**
  若，存储节点在操作系统磁盘上使用LVM，则需要将关联的设备添加到过滤器中。

- **解决方式：**
  配置存储节点时，需要挂载相应的卷。

  ```shell
  #修改配置文件：/etc/lvm/lvm.conf
  [root@controller /]# vim /etc/lvm/lvm.conf
  devices {
  ...
  filter = [ "a/sda/", "a/sdb/", "r/.*/"]
  }
  ```

  

**9、Nova-comput无法启动**

```shell
#无法使用
[root@controller /]# openstack compute service list --service nova-compute

```

- **原因：**

  防火墙阻止访问5672端口

- **解决方式：**

  ```shell
  #查看nova-comput日志：发现端口未启用
  cat /var/log/nova/nova-compute.log
  
  #关闭防火墙
  systemctl stop firewalld
  systemctl disable firewalld
  
  #重启nova服务
  systemctl start libvirtd.service openstack-nova-compute.service
  ```

  








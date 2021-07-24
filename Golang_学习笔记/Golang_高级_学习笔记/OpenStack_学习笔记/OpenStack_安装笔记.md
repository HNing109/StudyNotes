本教程安装QueenS版本的OpenStack，采用all in one的方式安装（即：所有组件均安装在同一个机器上）

OpenStack官网安装手册：https://docs.openstack.org/install-guide/

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

作用：不禁用SELinux可能会导致openstack newton Apache无法启动。

错误如下："Devstack fail to start apache2 -"Address already in use":"coild not bind to address"

解决方式：配置/etc/selinux/config，关闭SELINUX

```shell
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
```



## 2.2、关闭防火墙

作用：免去配置OpenStack各个节点时，需要不断手动开启端口的麻烦。但是在实际生产环境中绝对不允许这样做，这会引起安全问题。

```shell
[root@controller ~]# systemctl stop NetworkManager 
[root@controller ~]# systemctl disable NetworkManager
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
#[Member]
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="http://192.168.83.139:2380"
ETCD_LISTEN_CLIENT_URLS="http://localhost:2379,http://192.168.83.139:2379"
ETCD_NAME="controller"
#[Clustering]
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

    

- 验证keystone安装是否成功

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

  

- 创建OpenStack客户端脚本

  - 为什么需要使用脚本？

    由于之前，配置系统环境时，使用的是export命令，该命令用于临时配置环境变量，当系统重启或者执行export明令的终端关闭后，该环境变量就无法使用了，导致openstack命令无法被执行。

    

  - 需要创建的脚本：

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

    

## 4.2、安装Glance

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

  

## 4.4、安装Nova

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

  

## 4.5、安装Neutron





## 4.6、安装Horizon





## 4.7、安装Cinder


























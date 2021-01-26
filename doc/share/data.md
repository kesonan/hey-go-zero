# Redis&Mysql安装

为了演示，因此将三个应用安装在同一台centOS虚拟机上。

> ## 注意 
> 本文档并非我原创，以下内容参考了他人博客，为了方便大家走流程，因此搬到这里来了。

# 准备工作
* centOS 7
* 设置固定ip

# 设置固定ip
虚拟机固定ip设置这里就不重复了，参考[centOS安装知设置静态ip](./centos_install.md)
本次演示设置固定ip为`172.16.100.135`

# Redis安装

原文: [CentOS7安装MySQL（完整版）](https://blog.csdn.net/qq_36582604/article/details/80526287)

## 安装gcc
**检测gcc**
```shell
$ gcc -v
```
```text
加载 "fastestmirror" 插件
Config time: 0.017
Yum version: 3.4.3
没有该命令：gcc。请使用 /usr/bin/yum --help
```
> 没有安装gcc，因此需要安装gcc

**安装gcc**
```shell
$ yum -y install gcc
$ yum -y install centos-release-scl
$ yum -y install devtoolset-9-gcc devtoolset-9-gcc-c++ devtoolset-9-binutils
$ scl enable devtoolset-9 bash
```

## 安装wget
```shell
$ yum -y install wget
```

## 安装redis
```shell
$ wget http://download.redis.io/releases/redis-6.0.6.tar.gz
$ tar xzf redis-6.0.6.tar.gz
$ cd redis-6.0.6
$ make
$ make install PREFIX=/usr/local/redis
```

## 设置redis.conf
**新建redis目录**
```shell
$ sudo mkdir /etc/redis
$ cp redis.conf /etc/redis/
```

**修改redis为后台启动**
编辑/etc/redis/redis.conf，修改`daemonize no`为`daemonize yes`
添加`bind ${IP}`供外部访问,${IP}为机器ip，如：`bind 172.16.100.135`

**设置开机启动**
```shell
$ /etc/systemd/system/redis.service
```
添加如下内容:
```text
[Unit]
Description=redis-server
After=network.target
[Service]
Type=forking
ExecStart=/usr/local/redis/bin/redis-server /etc/redis/redis.conf
PrivateTmp=true
[Install]
WantedBy=multi-user.target
```
```shell
$ systemctl daemon-reload
$ systemctl start redis.service
$ systemctl enable redis.service
```

## 创建软链接
```shell
$ ln -s /usr/local/redis/bin/redis-cli /usr/bin/redis-cli
```

## 测试redis是否安装成功
```shell
$ redis-cli -h 172.16.100.135 -p 6379
```
```text
172.16.100.135:6379> ping
PONG
```

## 常见命令
**启动/停止/重启redis**
```
systemctl start redis.service 
systemctl stop redis.service 
systemctl restart redis.service 
```
**开启启动/停止启动**
```shell
systemctl enable redis.service   
systemctl disable redis.service 
```
**查看redis状态**
```shell
systemctl status redis.service
```

# 安装mysql
原文：[CentOS7安装MySQL（完整版）](https://blog.csdn.net/qq_36582604/article/details/80526287)

## 新建mysql目录
```shell
$ mkdir /usr/local/mysql
```

## 下载并安装MySQL(依次执行)
```shell
$ cd /usr/local/mysql
$ wget -i -c http://mirrors.ustc.edu.cn/mysql-repo/mysql57-community-release-el7-10.noarch.rpm
```
```shell
$ rpm -ivh mysql57-community-release-el7-10.noarch.rpm
```
```shell
$ yum install mysql-server
```

## 启动mysql
```shell
$ systemctl start  mysqld.service
```

## 开机启动mysql
```shell
$ systemctl enable mysqld.service
```

## 首次登录mysql

**找到root用户临时密码**
```shell
$ grep "password" /var/log/mysqld.log
```
```shell
$ mysql -uroot -p
```

## 修改默认密码
```shell
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY '${new_password}';
mysql> flush privileges; 
```

> ${new_password}为新密码

## 开启mysql的远程访问
```shell
mysql> grant all privileges on *.* to 'root'@'${IP}' identified by '${password}' with grant option;
```

> ${IP}为机器IP，${password}为新密码

## 防火墙开放端口
```shell
$ firewall-cmd --zone=public --add-port=3306/tcp --permanent
$ firewall-cmd --reload
```

## 更改mysql的编码

**退出mysql**
```shell
mysql> exit;
```

**编辑配置文件**
```shell
$ vi /etc/my.cnf
```
新增几处内容
```text
default-server-set=utf8
character-server-set=utf8
collation-server=utf8_general_ci
```

最终内容
```text
 For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/5.7/en/server-configuration-defaults.html
[client]
default-character-set=utf8

[mysqld]
#
# Remove leading # and set to the amount of RAM for the most important data
# cache in MySQL. Start at 70% of total RAM for dedicated server, else 10%.
# innodb_buffer_pool_size = 128M
#
# Remove leading # to turn on a very important data integrity option: logging
# changes to the binary log between backups.
# log_bin
#
# Remove leading # to set options mainly useful for reporting servers.
# The server defaults are faster for transactions and fast SELECTs.
# Adjust sizes as needed, experiment to find the optimal values.
# join_buffer_size = 128M
# sort_buffer_size = 2M
# read_rnd_buffer_size = 2M
datadir=/var/lib/mysql
socket=/var/lib/mysql/mysql.sock
character-set-server=utf8
collation-server=utf8_general_ci

# Disabling symbolic-links is recommended to prevent assorted security risks
symbolic-links=0

log-error=/var/log/mysqld.log
pid-file=/var/run/mysqld/mysqld.pid
```
保存，重启`service mysqld restart`。

# nginx安装
```shell
$ yum install gcc-c++
$ yum install -y pcre pcre-devel
$ yum install -y zlib zlib-devel
$ yum install -y openssl openssl-devel
$ wget http://nginx.org/download/nginx-1.18.0.tar.gz
$ tar -zxvf nginx-1.18.0.tar.gz
$ cd nginx-1.18.0
$ ./configure
$ make
$ make install
```

## 启动nginx
```shell
$ cd /usr/local/nginx/sbin
$ ./nginx
```
## 常用命令
* ./nginx -s stop
* ./nginx -s quit 
* ./nginx -s reload

## 开机启动
```shell
$ cd /lib/systemd/system
$ vi nginx.service
```
添加如下内容
```text
[Unit]
Description=nginx 
After=network.target 
   
[Service] 
Type=forking 
ExecStart=/usr/local/nginx/sbin/nginx
ExecReload=/usr/local/nginx/sbin/nginx reload
ExecStop=/usr/local/nginx/sbin/nginx quit
PrivateTmp=true 
   
[Install] 
WantedBy=multi-user.target
```

设置开机启动
```shell
$ systemctl enable nginx.service
```

启动nginx
```shell
$ systemctl start nginx.service
```

停止nginx
```shell
$ systemctl stop nginx.service
```

重启nginx
```shell
$ systemctl restart nginx.service
```

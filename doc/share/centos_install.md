# centOS安装

> 作者：[松妹子](https://github.com/anqiansong)
>
> 日期：2020年01月14日

# 下载centOS

* 版本： centOS 7
* 下载地址：[centOS下载入口](https://www.centos.org/download/)

# 安装centOS

## 虚拟机准备

* 打开VMware Fusion，新建
  ![选择安装方法](../../resource/centos_01.png)
* `继续`，选择`从光盘或映像中安装`，找到已经下载好的centOS iso文件
  ![创建新的虚拟机](../../resource/centos_02.png)
* `继续`，后续默认即可，直至完成
  ![选择固件类型](../../resource/centos_03.png)
  ![完成](../../resource/centos_04.png)
* 完成保存为`master`

  
## 克隆虚拟机
> 至此，一台虚拟机已经准备好了，由于我们需要一台虚拟机来做master节点，2台虚拟机来做node节点，因此我们需要再克隆两台虚拟机出来。
> 在克隆前先关掉虚拟机

![虚拟机克隆](../../resource/vm_clone.png)
克隆出两台虚拟机名称分别为node1、node2

## 虚拟机设置

> 为了保证后续k8s稳定安装，需要设置一下虚拟机配置，设置处理器为2核+，内存为2048MB+

![虚拟机设置](../../resource/vm_setting.png)
![虚拟机设置](../../resource/vm_cpu.png)
![虚拟机设置](../../resource/vm_cpu_mem.png)

## 设置静态ip
* 分别在进入虚拟机设置，并对`网络适配器`进行设置，点击`高级选项`，生成MAC地址，要确保三个虚拟机的MAC地址不能一样，记住每个虚拟机名称对应的mac地址。
  ![网络适配器](../../resource/vm_net.png)
  ![网络适配器](../../resource/vm_mac.png)
  ![网络适配器](../../resource/vm_mac2.png)
  ![网络适配器](../../resource/vm_mac3.png)

* 编辑dhcpd.conf
  ``` shell
  $ sudo vi /Library/Preferences/VMware\ Fusion/vmnet8/dhcpd.conf
  ```
* 指定master的ip为172.16.100.131，node1的ip为172.16.100.132，node2的ip为172.16.100.133,在文末填充如下内容后，保存后`重新启动VMware Fusion`。
  ``` text
  host master{
          hardware ethernet 00:50:56:3F:71:76;
          fixed-address 172.16.100.131;
  }
  
  host node1{
          hardware ethernet 00:50:56:3B:B7:98;
          fixed-address 172.16.100.132;
  }
  
  host node2{
          hardware ethernet 00:50:56:2F:70:48;
          fixed-address 172.16.100.133;
  }
  ```
  
## 安装

* 分别启动`master`、`node1`、`node2`虚拟机，进入centOS启动页
  ![centOS启动页](../../resource/centos_home.png)
* 通过上下键移动，选择`Install CentOS 7`，回车
  ![centOS初始化](../../resource/centos_init.png)
* 选择语言`中文`后`继续`
  ![centOS选择语言](../../resource/centos_language.png)
* 安装位置，我们选择`自动分区`就行了
* 网络设置
  ![centOS网络设置](../../resource/centos_network.png)
* 选择`打开`网络，并分别修改`主机`名称为`master`、`node1`、`node2`，点击`应用`，最后点击`完成`
  ![centOS网络设置](../../resource/centos_network_on.png)
* 点击`开始安装`
  ![centOS开始安装](../../resource/centos_install_start.png)
* 点击`ROOT 密码`进入密码设置页面
  ![centOS设置密码](../../resource/centos_set_password.png)
* 设置好密码后，点击`完成`(密码自行记住，后续登录需要用，这里暂时不添加用户了)
  ![centOS设置密码](../../resource/centos_password.png)
* 等待安装完成后`重启`
  ![centOS安装完成](../../resource/centos_installed.png)
* 重启后登录系统验证,输入用户名`root`，密码为刚刚自行设置的密码
  ![centOS登录](../../resource/centos_login.png)
* 至此，centOS就安装完成了

## centOS设置

> 这里为了操作美感，在启动虚拟机后通过在iTerm2中利用ssh链接到centOS后进行后续操作。

## 登录ssh
> ## 温馨提示
> 为了快速安装，如果你也是在mac上使用iTerm2，可以开始3个tab用来分别连接到三台虚拟机上，然后通过command+shift+i来同时输入命令。
![多窗口命令同步](../../resource/muti_tab.png)

``` shell
$ ssh root@172.16.100.131 -p 22
```

``` text
The authenticity of host '172.16.100.131 (172.16.100.131)' can't be established.
ECDSA key fingerprint is SHA256:sv99JzCzhchM0zKNnS3RNMgJbqCbE0nLRDXXdEQiuBE.
Are you sure you want to continue connecting (yes/no)? 
```

在第一次登录时会询问你是否继续链接，输入`yes`

``` text
yes
```

``` text
Warning: Permanently added '172.16.100.135' (ECDSA) to the list of known hosts.
root@172.16.100.135's password:
Last login: Thu Jan 14 21:58:25 2021
```

## 安装vim
``` shell
$ yum -y install vim
```
![安装vim](../../resource/install_vim.png)

## 安装wget

``` shell
$ yum -y install wget
```

![安装wget](../../resource/wget_install.png)

## 备份repo

``` shell
$ mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.bak
```

## yum配置阿里云

* 下载阿里云repo

``` shell
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
```
![备份repo](../../resource/backup_repo.png)

* 清除和生成缓存

``` shell
$ yum clean all
```

``` shell
$ yum makecache
```
![yum元数据缓存](../../resource/yum_makecache.png)

## 关闭防火墙

``` shell
$ systemctl stop firewalld.service
$ systemctl disable firewalld.service
```
![关闭防火墙](../../resource/firewalld.png)

## 关闭selinux

* 临时关闭

``` shell
$ setenforce 0
```

* 永久关闭
G
``` shell
$ vim /etc/selinux/config
```

将SELINUX=enforcing改为`disabled`
![关闭selinux](../../resource/selinux_disabled.png)

> ## 温馨提示
> 关闭防火墙和selinux是为了后续顺利安装k8s做准备。

# 结尾
至此，centOS准备工作已经结束

# 猜你想看
* [k8s集群安装](./k8s_install.md)

## 参考链接

以上文档部分参考自[Centos7.6操作系统安装及优化全纪录](https://blog.51cto.com/3241766/2398136)


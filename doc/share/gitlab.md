# gitlab安装

# 准备工作
* centOS 7


# 安装curl
``` shell
$ yum -y install curl
```
``` text
yum -y install curl
已加载插件：fastestmirror
Loading mirror speeds from cached hostfile
 * base: mirrors.ustc.edu.cn
 * extras: mirrors.ustc.edu.cn
 * updates: mirrors.ustc.edu.cn
 ...
```

# 安装gitlab
## 添加yum仓库源
``` shell
$ curl -s https://packages.gitlab.com/install/repositories/gitlab/gitlab-ce/script.rpm.sh | sudo bash
```


## 安装
``` shell
$ yum  -y  install  gitlab-ce
```

## 修改访问地址和端口
```shell
$ vi /etc/gitlab/gitlab.rb
```
找到`external_url`，将其修改为`http://${IP}:8090`，比如我的机器ip是172.16.100.134，
则配置修改为`http://172.16.100.134:8090`，端口可以随意，只要不和已有端口冲突即可。

## nginx配置
修改监听端口为8090，server_name 为本机ip

``` shell
$ /var/opt/gitlab/nginx/conf/gitlab-http.conf
```
``` text 
server {
    listen:         *:8090;
    server_name:    172.16.100.134;
    ....
}
```

## 重新启动gitlab
```shell
$ gitlab-ctl reconfigure
$ gitlab-ctl restart
```

# 访问gitlab
访问 `http://${IP}:8090`即可进入到gitlab首页，其中`${IP}`为你的机器ip。

# 准备工作
* 环境安装
    * Go语言环境
    * Goctl工具
    * Protobuf
    * Protoc-gen-go
    * Etcd
    * Redis
    * Postman
    * Beanstalkd
    

# Go语言环境
* [下载](https://golang.org/dl/) Go安装包
* [安装](https://golang.org/doc/install) Go语言环境
* 查看Go版本

    ```shell script
    $ go version
    ```
    ```text
    go version go1.15.1 darwin/amd64
    ```
  
# 开启Go Module
在我们后续演示过程中，均已Go Module形式创建工程，这里不对Go Path工程做演示，如果对Go Path比较熟悉的同学可以
使用Go Path（但对于后续Go Path问题不做回答）。

查看当前go module状态

```shell script
$ go env GO111MODULE
```
```text
on
```

如果不是on的话可以通过如下方式开启

```shell script
$ go env -w GO111MODULE="on"
```

# 配置代理
查看当前go proxy

```shell script
$ go env GOPROXY
```
```text
https://goproxy.cn
```

如果当前go proxy不是`https://goproxy.cn`的话建议你设置为该值（中国地区）

```shell script
$ go env -w GOPROXY=https://goproxy.cn
```

# 配置环境变量path

```shell script
$ vi /etc/paths
```

添加执行路径（如:$GOPATH）到末尾,这里建议创建一个自己方便浏览的目录来管理一些可执行文件
插入后有如下内容

```text
/usr/local/bin
/usr/bin
/bin
/usr/sbin
/sbin
/Users/xxx/workspace/private/path [1]
```

> 说明: 在我的电脑是以(`$HOME/workspace/private/path`)来存放可执行文件。

> [1] xxx为用户名称

# Goctl工具安装

```shell script
$ go get -u github.com/tal-tech/go-zero/tools/goctl
```

由于通过`go get`获取到的goctl二进制文件在`$GOPATH/bin`目录下，我们需要将其移动到我们之前指定的path路径下，便于管理。

```shell script
$ mv $GOPATH/bin/goctl $HOME/workspace/private/path [2]
```

> [2] `$GOPATH`是一个变量值，在终端下，其具体值可通过`go env GOPATH`查看，随后将其拼接称完成命令即可，如

```shell script
$ mv /Users/xxx/go/bin/goctl $HOME/workspace/private/path
```

查看`goctl`版本

```shell script
$ goctl -v
```
```text
goctl version 20201125 darwin/amd64 [3]
```

> [3] `goctl version`为固定标志符，`20201125`为发版时间，`darwin/amd64`为操作系统和操作系统架构，开发人员
> 在后续遇到任何goctl问题，可指定某一个版本进行说明。

# Protobuf安装
Protobuf及Protoc-gen-go是用于后续生产rpc服务的工具依赖。

* 进入github选择自己操作系统对应的二进制文件[下载](https://github.com/protocolbuffers/protobuf/releases)
* 解压后将bin目录中的`protoc`存放到我们之前指定的path(`$HOME/workspace/private/path`)目录下即可
* 查看版本

    ```shell script
    $ protoc --version
    ```
    ```text
    libprotoc 3.14.0
    ```
  
> 我这里下载了 [protoc-3.14.0-osx-x86_64.zip](https://github-production-release-asset-2e65be.s3.amazonaws.com/23357588/42d3ec00-25c2-11eb-81d8-19b6fba46513?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201201%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20201201T142830Z&X-Amz-Expires=300&X-Amz-Signature=fc829e7700c6cd3f7e3c39b5038db842f2ab9f738262fe207693e04bfa4c381a&X-Amz-SignedHeaders=host&actor_id=10302073&key_id=0&repo_id=23357588&response-content-disposition=attachment%3B%20filename%3Dprotoc-3.14.0-osx-x86_64.zip&response-content-type=application%2Foctet-stream)

# Protoc-gen-go安装

```shell script
$  go get -u github.com/golang/protobuf/protoc-gen-go
```

将`$GOPATH/bin`目录中的`protoc-gen-go`移动到我们之前指定的path(`$HOME/workspace/private/path`)目录下即可

```shell script
$ mv $GOPATH/bin/protoc-gen-go $HOME/workspace/private/path
```

> 注意:
> 这里我们是从`github.com/golang/protobuf`编译，而不是`https://github.com/protocolbuffers` 项目中的`protoc-gen-go`，这里大家都从前面那个项目编译。后者项目目前Goctl没有做兼容。

# Etcd安装
在本地演示项目中，我们采用Etcd来作为服务发现，更多关于Etcd的介绍请跳转至[官网](https://etcd.io/)

### brew安装(仅类Unix操作系统)

```shell script
$ brew install etcd
```
```text
==> Downloading https://mirrors.ustc.edu.cn/homebrew-bottles/bottles/etcd-3.4.13
Already downloaded: /Users/xxx/Library/Caches/Homebrew/downloads/1e85ac78899a479fed7a4726ad381dc357eb1215dc3972fbb8b3a87087f90c93--etcd-3.4.13.mojave.bottle.tar.gz
==> Pouring etcd-3.4.13.mojave.bottle.tar.gz
==> Caveats
To have launchd start etcd now and restart at login:
  brew services start etcd
Or, if you don't want/need a background service you can just run:
  etcd
==> Summary
🍺  /usr/local/Cellar/etcd/3.4.13: 8 files, 38.7MB
==> `brew cleanup` has not been run in 30 days, running now...
```

### 查看etcd版本

```shell script
$ etcd --version
```

```text
etcd Version: 3.4.13
Git SHA: Not provided (use ./build instead of go build)
Go Version: go1.15
Go OS/Arch: darwin/amd64
```

### 查看etcdctl版本

```shell script
$ etcdctl version
```
```text
etcdctl version: 3.4.13 [4]
API version: 3.4
```

> [4] 为了方便记忆，这里可以将etcdctl起一个别名`etl`

```shell script
$ vi ~/.zshrc [5]
```

> [5] zsh安装可自行google，当然你也可以使用bash,这里可自行google去设置别名。

在末尾添加`alias etl=etcdctl`，然后`source ~/.zshrc`即可。

# Redis安装

```
$ brew install redis
```
```text
==> Downloading https://mirrors.ustc.edu.cn/homebrew-bottles/bottles/redis-6.0.6.mojave.bottle.tar.gz
######################################################################## 100.0%
==> Pouring redis-6.0.6.mojave.bottle.tar.gz
==> Caveats
To have launchd start redis now and restart at login:
  brew services start redis
Or, if you don't want/need a background service you can just run:
  redis-server /usr/local/etc/redis.conf
==> Summary
🍺  /usr/local/Cellar/redis/6.0.6: 13 files, 3.8MB
```

查看redis版本

```shell script
$ redis-cli -v
```
```text
redis-cli 6.0.6
```

>说明：windows安装教程请[下载安装](https://redis.io/download)

# Postman安装(可选)
为了方便接口测试，这里建议大家安装一下postman工具，方便后期api调试，当然你也可以使用其他工具如`curl`、Idea工具中的
`Http Client`等，选择一个你熟悉的工具即可。

# Beanstalk安装

```shell script
$ brew install beanstalkd
```
```text
==> Downloading https://mirrors.ustc.edu.cn/homebrew-bottles/bottles/beanstalkd-1.12.mojave.bottle.tar.gz
######################################################################## 100.0%
==> Pouring beanstalkd-1.12.mojave.bottle.tar.gz
==> Caveats
To have launchd start beanstalkd now and restart at login:
  brew services start beanstalkd
Or, if you don't want/need a background service you can just run:
  beanstalkd
==> Summary
🍺  /usr/local/Cellar/beanstalkd/1.12: 8 files, 65.4KB
```

```shell script
$ beanstalkd -v
```
```text
beanstalkd 1.12
```

>说明：windows安装教程请[下载安装](https://beanstalkd.github.io/download.html)

# End

上一篇 [《首页》](../index.md)

下一篇 [《Goctl介绍》](./goctl-intro.md)

# 猜你想

* [《目录说明》](../index.md)
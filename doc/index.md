# 前言
和传统开发项目的流程一样，本演示项目也从需求出发，到服务部署，所以，这是一个漫长的过程，
在这个过程中你需要保持耐心，一步一个脚印向前走，否则你会不知所措。在之前的go-zero文档中，
有不少开发者反馈我们的文档不全、难懂，又或者对新手来说不易入手，又或者因为版本迭代更新快，
而导致文档更新不及时而引起的阅读难以理解，文档前后跳转大，对于以上种种原因，我利用工作
之余的编写此演示项目，希望对大家有所帮助。

# 硬件环境
* 电脑

# 软件环境
* mac OS、windows、linux
* Go
* IDE(Goland、Atom、VSCode)
* Goctl
* Etcd
* Redis
* Protoc&Protoc-gen-go
* Postman

# 集成依赖
* go-zero（core）

# 本机环境
* mac OS(10.14.6)
* go version go1.15.1 darwin/amd64
* Goland 2020.2.3
* goctl version 20201125 darwin/amd64
* Etcd Version: 3.4.13
* Redis-cli 6.0.6
* Protoc: libprotoc 3.14.0
* Protoc-gen-go

# 常见概念介绍
`$` 代表一个shell命令/可执行文件开始，如
```shell script
$ echo hello
```
代表执行`echo hello`命令输出`hello`

`Goctl` go-zero自带的代码生产工具

# 目录
* [首页](../readme.md)
* [准备工作](./prepare)
    * [准备工作](./prepare/prepare.md)
    * [Goctl介绍](./prepare/goctl-intro.md)
    * [服务目录](./prepare/service-structure.md)
    * [Api语法介绍](./prepare/api-grammar.md)
    * [Proto使用说明](./prepare/proto-rule.md)
    * [创建工程](./prepare/project-create.md)
    * [数据库准备工作](./prepare/db-create.md)
    * [常见FAQ集合](./prepare/faq.md)
    
* [需求概况](./requirement)
    * [需求说明](./requirement/summary.md)


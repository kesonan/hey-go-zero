# 服务目录
在正式进入演示工程前，开发人员需要对服务的目录结构作一定的规划，一个好的目录结构不仅可以让
人好容易理解不同层级的业务，也便于开发人员的维护和管理。

> 这里介绍的服务目录仅为内部使用的结构目录，并不强制要求大家按照此目录结构来。当然如果你是新手，
> 或者使用`goctl`来生成服务，那么你必须得了解，否则后面你会不知所措。

# Project Structure

```text
greet                           [1]
├── common                      [2]
│   └── util
├── go.mod
├── go.sum
└── service                     [3]
    ├── app                     [4]
    │   ├── schedule      [5]
    │   ├── selection     [6]
    │   └── user          [7]
    │       ├── api       [8]
    │       ├── rpc       [9]
    │       ├── script    [10]
    │       └── sync      [11]
    │       └── model     [12]
    ├── hera                    [13]
    └── ops                     [14]
```

* [1] 工程层级
* [2] 工具层级，如：一些common工具，除此之外，开发人员可以根据自己的需要在相同层级创建其他package
* [3] 服务层级，存各个项目组的服务
* [4][13][14] 项目组层级[optional]，如果在你所在企业有多个项目组，而大家的代码都在一个大仓库下，建议预留一定的项目组目录进行区分，否则可能每个项目组之间都会有一个user服务，导致无法区分，当然，如果没有项目组的划分，那可省去此层级。如上述结构中分别对应为:
    * app c端项目组
    * hera 数据组
    * ops 运维组
* [5][6][7] 微服务层级，如user、selection、schedule分别为不同的微服务。
* [8][9][10][11] 服务类型层级，除上述列举的分类外，你还可以根据自己的场景添加其他分类，如消息队列服务`rmq`等。
    * api api服务，外部http访问入口
    * rpc rpc服务，微服务间通信
    * script 脚本服务，存放一些定时/临时脚本
    * sync 同步服务
* [12] 在和服务类型层级同级的还有一个model目录，这里是存放data access入口，如常见的数据库访问层将放在这里。

> 注意: 以上目录树仅供参考，开发人员可以根据实际场景建立适合自己公司业务的目录结构。

# Service Structure
在Project Structure中我们已经提到了在微服务层级中还会对每个微服务进行分类，如`api`、`rpc`、`script`等，在这一层级上，我们还需要在其下层创建目录层级，这里推荐以下目录结构，这里并非固定为
该结构，你可以适当根据需要添加或者增减。如api服务还会有`handler`、`svc`、`types`层

```text
.
├── etc                             [1]
│   └── user.yaml
├── internal                        [2]
│   ├── config                [3]
│   │   └── config.go
│   └── logic                 [4]
│       └── userlogic.go
└── user.go                         [5]
```

* [1] yaml/json参考配置文件，一般用于给运维进行配置参考
* [2] 服务内部逻辑，在这里定义的文件仅供该服务访问
* [3] yaml/json对应的配置声明
* [4] 业务逻辑层，在这里实现业务需求
* [5] main入口，一个服务的执行入口

# Api Structure (Goctl生成版本)

```text
api
├── etc                                 [1]
│   └── user-api.yaml
├── internal                            [2]
│   ├── config                    [3]
│   │   └── config.go
│   ├── handler                   [4]
│   │   ├── loginhandler.go [5]
│   │   └── routes.go       [6]
│   ├── logic                     [7]
│   │   └── loginlogic.go
│   ├── svc                       [8]
│   │   └── servicecontext.go
│   └── types                     [9]
│       └── types.go
├── user.api                            [10]
└── user.go                             [11]
```

> 上述tree为goctl生成的api服务目录结构，如果你是使用goctl来生成服务，我们建议在不熟悉的情况下切勿随意更改。

* [1] yaml/json参考配置文件存储目录，开发人员可在这里编写参考配置供运维参考[rw]
* [2] 服务内部逻辑，在这里定义的文件仅供该服务访问
* [3] yaml/json对应的配置声明 [rw]
* [4] handler层，即`HandlerFunc`适配器的实现
* [5] handler处理文件,一个协议对应一个handler.go文件[rw]
* [6] 路由定义[ro]
* [7] 业务层，在这里填充业务code[rw]
* [8] 资源依赖层，将service层所需要的资源依赖均存放在这里，便于资源的管理和控制。[rw]
* [9] service对外接收或响应的结构体定义[r]
* [10] api定义文件，通过.api定义我们所需要的服务
* [11] main入口，服务执行入口[r]

> `r` 可读，`w` 可写，`o`重新生成将覆盖原文件，因此切勿在标记`o`的文件上做修改，否则在代码重新生成时将全部覆盖。

# Rpc Structure (Goctl生成版本)

```text
rpc
├── etc                                         [1]
│   └── user.yaml
├── internal                                    [2]
│   ├── config                            [3]
│   │   └── config.go
│   ├── logic                             [4]
│   │   └── getuserinfologic.go
│   ├── server                            [5]
│   │   └── userserviceserver.go
│   └── svc                               [6]
│       └── servicecontext.go
├── user                                        [7]
│   └── user.pb.go
├── user.go                                     [8]
├── user.proto                                  [9]
└── userservice                                 [10]
    └── userservice.go
```

> 上述tree为goctl生成的rpc服务目录结构，如果你是使用goctl来生成服务，我们建议在不熟悉的情况下切勿随意更改。

* [1] yaml/json参考配置文件存储目录，开发人员可在这里编写参考配置供运维参考[rw]
* [2] 服务内部逻辑，在这里定义的文件仅供该服务访问
* [3] yaml/json对应的配置声明 [rw]
* [4] 业务层，在这里填充业务code[rw]
* [5] rpc server实现[ro]
* [6] 资源依赖层，将service层所需要的资源依赖均存放在这里，便于资源的管理和控制。[rw]
* [7] pb.go存放目录
* [8] main入口，服务执行入口[r]
* [9] proto文件[rw]
* [10] rpc client调用入口[ro]
> `r` 可读，`w` 可写，`o`重新生成将覆盖原文件，因此切勿在标记`o`的文件上做修改，否则在代码重新生成时将全部覆盖。

# Model Structure (Goctl生成版本)

```text
model
├── usermodel.go
└── vars.go
```

# End

上一篇 [《Goctl介绍》](./goctl-intro.md)

下一篇 [《Api语法介绍》](./api-grammar.md)
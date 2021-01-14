# Api语法
在上一章节[《服务目录》](./service-structure.md)我们提到过.api文件，该文件是对api服务的一个定义文件，像proto定义rpc服务一样，
我们可以通过.api文件声明服务名称，请求路由，请求方法，请求入参，请求出参，中间件声明，jwt声明等。利用该文件，我们可以通过goctl工具快速
生成api服务。

# 为什么需要Api文件
这里我们来假设，在没有.api文件的情况下，你会通过什么方式去省时省力的创建一个服务，我想没有几个人无不会用到工具其替代，如果说你的手速能达到这个效率，
那我们另当别论，请你跳过这一章节。我们需要.api文件可以带来哪些好处？
* 便于快速生成api服务
* 0学习成本
* 阅读性强
* 能生成各种语言的调用代码
* 替代api文档
* 扩展性强

# Api概览
api语法及其简单易懂，0学习成本，我们来看一下下列api定义

``` text
syntax = "v1"
import "common.api" [1]

type LoginReq {
    Username string `json:"username"`
    password string `json:"password"`
}

type LoginReply {
    Token string `json:"token"`
    Expire int64 `json:"expire"`
    RefreshToken string `json:"refreshToken"`
}

type UserInfoReply {
    Name string `json:"name"`
    Age int `json:"age"`
    Gender string `json:"gender"`
    Birthday string `json:"birthday"`
}

type UserPasswordEditReq {
    OldPassword string `json:"oldPassword"`
    NewPassword string `json:"newPassword"`
}

@server(
    jwt: Auth   [2]
    middleware: MetricMiddleware,LogMiddleware  [3]
    group: user [4]
)
service user-api {      [5]
    @handler healthCheck    [6]
    post /user/health/check     [7]

    @handler login
    post /user/login (LoginReq) returns (LoginReply)

    @handler userInfo
    get /user/info returns (UserInfoReply)

    @handler passwordEdit
    post /user/password/edit (UserPasswordEditReq)
}
```

在没了解api语法之前，我们可以从中简单的看出该文件定了4个路由，分别为

* `/user/health/check` 健康检查
* `/user/login` 登录
* `/user/info` 获取用户信息
* `/user/password/edit` 编辑用户密码

在service block head定义了user-api服务需要`jwt`授权，需要`MetricMiddleware`和`LogMiddleware`中间件和`user`分组。

# Api结构

```tetx
api
├── import block
├── info block
├── type block
└── service block
```

# Api语法
api语法详见[api语法描述](https://github.com/tal-tech/go-zero/blob/master/tools/goctl/api/parser/readme.md)

# 补充说明
上述api语法中"语法定义"模块的相关代码为antlr4语法，更多信息可参考[官方文档](https://www.antlr.org/)

为了快速编写api文件，我们开发了idea、vsCode插件，该插件提供了高亮、跳转、模板、格式化等功能，[点击这里](https://github.com/tal-tech/goctl-plugins) 查看更多详情

# End

上一篇 [《服务目录》](./service-structure.md)

下一篇 [《Proto使用说明》](./proto-rule.md)
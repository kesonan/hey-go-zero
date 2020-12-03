# user创建步骤
通过本文档介绍我在编写演示项目的每一步流程，这样不至于你在阅读的时候忽然发现，怎么到了这里，那里是怎么回事。

> 说明：本文档对新手比较适用，如果已经很熟悉go-zero、goctl的同学可以跳过本文档。

# 创建api目录
进入`user`模块下创建api目录得到

```text
service
    ├── course
    ├── schedule
    ├── selection
    └── user
        └── api
```

> 上述tree是以`service`作为root目录。

# 新建user.api
* 在`api`目录文件夹上`右键`->`New Api File`->`输入user`->`选择Empty file`->`回车`
* 修改user.api文件内容为

    ```text
    type (
        UserRegisterReq {
            Username string `json:"username"`
            Passowrd string `json:"passowrd"`
            // 定义用户角色，仅允许student|teacher两个枚举值。
            Role string `json:"role,options=student|teacher"`
        }
    
        UserLoginReq {
            Username string `json:"username"`
            Passowrd string `json:"passowrd"`
            Role string `json:"role,options=student|teacher"`
        }
    
        UserLoginReply {
            Id string `json:"id"`
            Token string `json:"token"`
            ExpireAt int64 `json:"expireAt"`
        }
    )
    
    type (
        UserInfoReply {
            Id string `json:"id"`
            Name string `json:"name"`
            Gender string `json:"gender"`
            Birthday string `json:"birthday"`
            Role string `json:"role"`
        }
    
        UserInfoReq {
            Name string `json:"name,optional"`
            Gender string `json:"gender,optional"`
            Birthday string `json:"birthday"`
        }
    )
    
    @server(
        group: noauth
    )
    service user-api {
        @handler register
        post /api/user/register (UserRegisterReq)
    
        @handler login
        post /api/user/login (UserLoginReq) returns (UserLoginReply)
    }
    
    @server(
        jwt: Auth
        group: auth
    )
    service user-api {
        @handler userInfo
        get /api/user/info/self returns (UserInfoReply)
    
        @handler userInfoEdit
        post /path (UserInfoReq)
    }
    ```

# 生成代码
在Goland中生成代码有三种方式（任意一种均可）
* project面板区文件右键生成
    * 选中`user.api`文件->`右键`->`New`->`Go Zero`->`Api Code`
    * `Api Code`回车后会弹出一个文件对话框询问你需要生成服务的目标目录，默认为`user.api`所在目录，我们这里选择默认，点击`OK`确认生成。
* api文件编辑区右键生成
    * 打开`user.api`文件->`编辑区`->`右键`->`Generate..`->`Api Code`
    * `Api Code`回车后会弹出一个文件对话框询问你需要生成服务的目标目录，默认为`user.api`所在目录，我们这里选择默认，点击`OK`确认生成。
* 终端生成（推荐）
    * 选中`user.api`文件->`右键`->`Open in Terminal`
    * 执行`goctl api go -api user.api -dir .`命令即可
    
        ```shell script
        $ goctl api go -api user.api -dir .
        ```
        ```text
        Done.
        ```
接下来我们看一下生成代码的目录树，在终端下进入`user/api`目录

```shell script
$ tree
```
```text
.
├── etc
│   └── user-api.yaml
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── auth
│   │   │   ├── userinfoedithandler.go
│   │   │   └── userinfohandler.go
│   │   ├── noauth
│   │   │   ├── loginhandler.go
│   │   │   └── registerhandler.go
│   │   └── routes.go
│   ├── logic
│   │   ├── auth
│   │   │   ├── userinfoeditlogic.go
│   │   │   └── userinfologic.go
│   │   └── noauth
│   │       ├── loginlogic.go
│   │       └── registerlogic.go
│   ├── svc
│   │   └── servicecontext.go
│   └── types
│       └── types.go
├── user.api
└── user.go

```

> 说明： 这个时候进入`user.go`文件查看，发现代码有多处地方报红
> 解决方案：在终端进入`user/api`执行
> ```
> $ go test -race ./...
> ```
> 为了方便，可将`go test -race ./...`设置一个别名为`gt`，后续我们的`go test`均用`gt`命令替代。

到这里，user api服务便创建好了。我们首先来尝试调用获取用户信息接口看看效果。

# 完善yaml配置文件
yaml配置文件需要配置什么配置项完全参考于`api/internal/config`下你定义的配置。由于我们在之前user.api文件中声明需要`jwt`鉴权

```
@server(
	jwt: Auth
	group: Auth
)
```
所以在生成代码时，配置项定义也生成好了，接下来看一下目前已经定义的配置：

```go
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
```

从上述内容可以看出，目前拥有两大块配置定义
* rest.RestConf

    该配置定义声明了一个rest api服务基础信息，通过源码你可以看到
    * 服务相关定义`service.ServiceConf`
    * 指定监听的`Host`和端口`Port`
    * 日志信息`logx.LogConf`
    * 超时时间`Timeout`等

> `rest.RestConf`配置说明见[rest api 服务基本配置说明](../../../doc/gozero/rest-api-conf.md)

* Auth

Auth配置包含`AccessSecret`和`AccessExpire`两个配置项，分别为jwt密钥和过期时间设置。更多jwt信息请参考[jwt官方说明文档](https://jwt.io/introduction/)

接下来我们编辑`api/etc/user-api.yaml`文件，添加配置上述配置项

```yaml
Name: user-api
Host: 0.0.0.0
Port: 8888
Auth:
  AccessSecret: 1e69481b-7405-4369-9ce3-9aaffdb56ce3
  AccessExpire: 3600
```

> 注意：`AccessSecret`这里只是一个示例，在真实环境中，请自行从实际场景出发去设置，切勿用示例值。

# 启动user api服务

```shell script
$ go run user.go
```
```text
Starting server at 0.0.0.0:8888...
```

# 尝试访问服务
这里我们先来访问一下获取用户信息的协议

```shell script
$ curl -i -X GET \
    http://localhost:8888/api/user/info/self
```
```text
HTTP/1.1 401 Unauthorized
Date: Thu, 03 Dec 2020 14:40:11 GMT
Content-Length: 0
```
不出所料，由于`api/user/info/self`协议需要`jwt`鉴权，通过curl可以看到，目前并没有任何jwt token 信息传递给http server，因此得到`401`的http状态响应。

> 注意：windows版本在终端用`curl`进行http请求，且请求体为`json`类型时，需要将json进行转义。

# 填充user api服务逻辑

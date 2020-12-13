# user-api创建步骤
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
    info(
    	title: "用户系统"
    	desc: "用户模块api描述文件，详细需求说明请见hey-go-zero/doc/requirement/user.md"
    	author: "songmeizi"
    	version: "1.0"
    )
    
    type (
    	UserRegisterReq {
    		Username string `json:"username"`
    		Passowrd string `json:"password"`
    		// 定义用户角色，仅允许student|teacher两个枚举值。
    		Role string `json:"role,options=student|teacher"`
    	}
    	
    	UserLoginReq {
    		Username string `json:"username"`
    		Passowrd string `json:"password"`
    	}
    	
    	UserLoginReply {
    		Id int64 `json:"id"`
    		Token string `json:"token"`
    		ExpireAt int64 `json:"expireAt"`
    	}
    )
    
    type (
    	UserInfoReply {
    		Id int64 `json:"id"`
    		Name string `json:"name"`
    		Gender string `json:"gender"`
    		Role string `json:"role"`
    	}
    	
    	UserInfoReq {
    		Name string `json:"name,optional"`
    		Gender string `json:"gender,optional"`
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
    	post /api/user/info/edit (UserInfoReq)
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
├── etc // yaml配置文件
│   └── user-api.yaml
├── internal // 仅user api服务可访问的内部文件
│   ├── config  // yaml配置文件对应的结构定义
│   │   └── config.go
│   ├── handler // http.HandlerFunc实现
│   │   ├── auth   // 文件分组1，来自user.api定义中的group值
│   │   │   ├── userinfoedithandler.go
│   │   │   └── userinfohandler.go
│   │   ├── noauth  // 文件分组2，来自user.api定义中的group值
│   │   │   ├── loginhandler.go
│   │   │   └── registerhandler.go
│   │   └── routes.go // 路由定义
│   ├── logic // 业务逻辑
│   │   ├── auth // 文件分组1，来自user.api定义中的group值
│   │   │   ├── userinfoeditlogic.go
│   │   │   └── userinfologic.go
│   │   └── noauth // 文件分组2，来自user.api定义中的group值
│   │       ├── loginlogic.go
│   │       └── registerlogic.go
│   ├── svc // 资源依赖
│   │   └── servicecontext.go
│   └── types
│       └── types.go
├── readme.md
├── user.api // api定义
└── user.go // main入口

```

> 说明：上述目录中的注释是为了大家能够快速知道该目录结构的用途，是后期我加入的，实际生成的tree不会带注释和readme.md文件。
> 另：这个时候进入`user.go`文件查看，发现代码有多处地方报红
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

# 创建user表

```mysql
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `username` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '登录用户名',
  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '登录用户密码',
  `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户姓名',
  `gender` tinyint(1) DEFAULT '0' COMMENT '用户性别 0-未知，1-男，2-女',
  `role` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户角色 student-学生,teacher-教师，manager-管理员',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

>说明：请将上述create ddl 复制后自行创建，这里就不过多演示了。

# 生成带redis缓存的usermodel代码
首先进入`service/user`目录，右键`user`文件夹进入终端

```shell script
$ goctl model mysql datasource -url="ugozero@tcp(127.0.0.1:3306)/heygozero" -table="user" -c -dir ./model
```
```text
Done.
```

生成完毕后会在`service/user`目录下会多一个`model`文件夹，其包含内容如下:

```text
model
├── usermodel.go
└── vars.go
```

# 添加regex.go
在`hey-go-zero`下添加一个`common/regex`和`common/codeerror`文件夹，

创建`regex.go`文件，填充代码:

```go
package regex

import "regexp"

const (
	Username = `(?m)[a-zA-Z_0-9]{6,20}`
	Password = `(?m)[a-zA-Z_0-9.-]{6,18}`
)

func Match(s, reg string) bool {
	r := regexp.MustCompile(reg)
	ret := r.FindString(s)
	return ret == s
}
```


# 添加`Mysql`和`CacheRedis`配置定义和yaml配置项
* 编打开`service/user/api/internal/config/config.go`，添加`Mysql`、`CacheRedis`定义

    ```go
    package config
    
    import (
    	"github.com/tal-tech/go-zero/core/stores/cache"
    	"github.com/tal-tech/go-zero/rest"
    )
    
    type Config struct {
    	rest.RestConf
    	Auth struct {
    		AccessSecret string
    		AccessExpire int64
    	}
    	Mysql struct {
    		DataSource string
    	}
    	CacheRedis cache.CacheConf
    }
    ```
  
* 打开`service/user/api/etc/user-api.yaml`文件，添加`Mysql`、`CacheRedis`配置项

    ```yaml
    Name: user-api
    Host: 0.0.0.0
    Port: 8888
    Auth:
      AccessSecret: 1e69481b-7405-4369-9ce3-9aaffdb56ce3
      AccessExpire: 3600
    Mysql:
      DataSource: ugozero@tcp(127.0.0.1:3306)/heygozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
    CacheRedis:
      - Host: 127.0.0.1:6379
      Type: node
    ```

    >说明： 我本地redis没有设置密码，因此没有配置`Password`配置项。

# ServiceContext增加`UserModel`资源
打开`service/user/api/internal/svc/servicecontext.go`，添加`UserModel`依赖。

```go
package svc

import (
	"hey-go-zero/service/user/api/internal/config"
	"hey-go-zero/service/user/model"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
```
# 填充user api服务逻辑

### 准备
在填充注册逻辑之前建议先阅读[自定义错误处理](../../../doc/gozero/http-error.md)，在下面的逻辑中我们将会用到。

### 添加`error.go`文件
在`service/user/api/internal/logic`下创建`error.go`文件，添加自定义错误类型

```go
var (
	InvalidUsername = errorx.NewInvalidParameterError("username")
	InvalidPassword = errorx.NewInvalidParameterError("password")
)
```

### 填充注册逻辑
打开`service/user/api/internal/logic/noauth/registerlogic.go`文件，编辑`Register`方法：

```go
if !regex.Match(req.Username, regex.Username) {
    return logic.InvalidUsername
}

if !regex.Match(req.Passowrd, regex.Password) {
    return logic.InvalidPassword
}

_, err := l.svcCtx.UserModel.FindOneByUsername(req.Username)
switch err {
case nil:
    return errorx.NewDescriptionError("用户名已存在")
case model.ErrNotFound:
    _, err = l.svcCtx.UserModel.Insert(model.User{
        Username: req.Username,
        Password: req.Passowrd,
        Role:     req.Role,
    })
    return err
default:
    return err
}
```

启动redis

```shell script
$ redis-server
```

启动user api服务，访问注册协议。

```shell script
$ go run user.go
```
```text
Starting server at 0.0.0.0:8888...
```

访问注册协议

```shell script
$ curl -i -X POST \
    http://localhost:8888/api/user/register \
    -H 'content-type: application/json' \
    -d '{
          "username":"songmeizi",
          "password":"111111",
          "role":"student"
  }'
```
```text
HTTP/1.1 200 OK
Date: Fri, 04 Dec 2020 09:46:58 GMT
Content-Length: 0
```
再次发起同样的请求你得到

```text
HTTP/1.1 406 Not Acceptable
Content-Type: application/json
Date: Fri, 04 Dec 2020 13:19:11 GMT
Content-Length: 39

{"code":-1,"desc":"用户名已存在"}
```

由于上述提示`用户名已存在`错误了，而且我们启用了redis缓存，如果不出意外的话，redis中已经有缓存了，分别为:
* 唯一索引`username`缓存的`主键id`值
* `主键id`缓存的用户行记录

我们访问redis查看一下。

```shell script
$ 127.0.0.1:6379> get cache#User#username#songmeizi
  "1"
  127.0.0.1:6379> get cache#User#id#1
  "{\"Username\":\"songmeizi\",\"Password\":\"111111\",\"Name\":\"\",\"Gender\":0,\"Role\":\"student\",\"CreateTime\":\"2020-12-04T17:46:58+08:00\",\"UpdateTime\":\"2020-12-04T17:46:58+08:00\",\"Id\":1}"
  127.0.0.1:6379>
```

> 说明：在`usermodel.go`中可查看到redis key prefix，具体拼接规则，你可以自行看一下`usermodel.go`中代码。
> 如：
> ```text
> cacheUserUsernamePrefix = "cache#User#username#"
> cacheUserIdPrefix       = "cache#User#id#"
> ```

> 恭喜！🎉🎉🎉 走到这里你已经成功的实现了第一条协议，你有没有发现你写得最多的代码是`Register`函数，填充注册逻辑，而持久层、缓存层及handler相关的代码你都没有编写，甚至你可能都不知道用到了这些代码。
用`go-zero`实现一个服务就是这么easy！接下来还有很长的路要走，不过大部分工作都像写`注册`代码一样，你只负责填充逻辑就行，其他的就交给`goctl`，请保持耐心，我们继续！

## 创建`jwtx.go`
在`hey-go-zero/common`创建一个文件夹`jwtx`和文件`jwtx.go`,添加如下代码

```go
package jwtx

import (
	"encoding/json"
	"net/http"

	"hey-go-zero/common/errorx"

	"github.com/tal-tech/go-zero/rest/httpx"
)

const JwtWithUserKey = "id"

func GetUserId(w http.ResponseWriter, r *http.Request) (int64, bool) {
	v := r.Context().Value(JwtWithUserKey)
	jn, ok := v.(json.Number)
	if !ok {
		httpx.Error(w, errorx.NewDescriptionError("用户信息获取失败"))
		return 0, false
	}
	vInt, err := jn.Int64()
	if err != nil {
		httpx.Error(w, errorx.NewDescriptionError(err.Error()))
		return 0, false
	}
	return vInt, true
}
```

### 填充登录逻辑
打开`service/user/api/internal/logic/noauth/loginlogic.go`文件，在`Login`中添加如下代码逻辑：

```go
if !regex.Match(req.Username, regex.Username) {
    return nil, logic.InvalidUsername
}

if !regex.Match(req.Passowrd, regex.Password) {
    return nil, logic.InvalidPassword
}

resp, err := l.svcCtx.UserModel.FindOneByUsername(req.Username)
switch err {
case nil:
    if resp.Password!=req.Passowrd{
        return nil,errorx.NewDescriptionError("密码错误")
    }
    
    jwtToken,expireAt, err := l.generateJwtToken(resp.Id,time.Now().Unix())
    if err != nil {
        return nil, err
    }
    
    return &types.UserLoginReply{
        Id:       resp.Id,
        Token:    jwtToken,
        ExpireAt: expireAt,
    }, nil
case model.ErrNotFound:
    return nil, errorx.NewDescriptionError("用户名未注册")
default:
    return nil, err
}
```
`generateJwtToken`方法：

```go
func (l *LoginLogic) generateJwtToken(id int64, iat int64) (string, int64, error) {
	claims := make(jwt.MapClaims)
	expireAt := iat + l.svcCtx.Config.Auth.AccessExpire
	claims["exp"] = expireAt
	claims["iat"] = iat
	claims[jwtx.JwtWithUserKey] = id
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	jwtToken,err:=token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
	if err != nil {
		return "", 0,err
	}
	return jwtToken,expireAt,nil
}
```

启动服务，请求一下登录协议

```shell script
$ curl -i -X POST \
    http://localhost:8888/api/user/login \
    -H 'content-type: application/json' \
    -d '{
  	"username":"songmeizi",
  	"password":"111111"
  }'
```
```text
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 04 Dec 2020 14:18:07 GMT
Content-Length: 178

{"id":1,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDcwOTUwODcsImlhdCI6MTYwNzA5MTQ4NywiaWQiOjF9.unYrI5J7o67J-FVltzbx6rH0P1LhYj13MlcYhcHcL9Y","expireAt":1607095087}
```

### 获取用户信息
和上面一样，找到对应的logic文件`service/user/api/internal/logic/auth/userinfologic.go`，找到`UserInfo`方法，发现这里没有请求参数，那么我们通过什么样式获取到当前请求户用户的
用户信息呢？
* 编辑`service/user/api/internal/logic/error.go`,添加代码
    ```go
    ErrUserNotFound = errorx.NewDescriptionError("用户不存在")
    ```
* 给`UserInfo`方法中添加请求参数`id int64`
* 找到`UserInfo`的调用方`service/user/api/internal/handler/auth/userinfohandler.go`,在方法`UserInfoHandler`中添加代码
    ```go
    id,ok:=jwtx.GetUserId(w,r)
    if !ok{
        return
    }
    ```
    完整代码
    ```go
    func UserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
    	return func(w http.ResponseWriter, r *http.Request) {
    		id,ok:=jwtx.GetUserId(w,r) // add
    		if !ok{ // add
    			return // add
    		} // add
    
    		l := logic.NewUserInfoLogic(r.Context(), ctx)
    		resp, err := l.UserInfo(id) // edit
    		if err != nil {
    			httpx.Error(w, err)
    		} else {
    			httpx.OkJson(w, resp)
    		}
    	}
    }
    ```
* 在`userinfologic`添加全局定义
    ```go
    var genderConvert = map[int64]string{
    	0: "未知",
    	1: "男",
    	2: "女",
    }
    ```
* 填充`UserInfo`方法逻辑
    
    ```go
    resp, err := l.svcCtx.UserModel.FindOne(id)
    switch err {
    case nil:
        return &types.UserInfoReply{
            Id:     resp.Id,
            Name:   resp.Name,
            Gender: genderConvert[resp.Gender],
            Role:   resp.Role,
        }, nil
    case model.ErrNotFound:
        return nil, logic.ErrUserNotFound
    default:
        return nil, err
    }
    ```
  
### 编辑用户信息
和上面一样，找到对应的logic文件`service/user/api/internal/logic/auth/userinfoeditlogic.go`，找到`UserInfoEdit`方法，这里和【获取用户信息一样】均需要在handler层获取到用户id，并传递到logic层。
最终代码如下:

* 找到`UserInfoEdit`的调用方`service/user/api/internal/handler/auth/userinfoedithandler.go`,在方法`UserInfoEditHandler`中添加代码
    ```go
    id,ok:=jwtx.GetUserId(w,r)
    if !ok{
        return
    }
    ```
    完整代码
    ```go
    func UserInfoEditHandler(ctx *svc.ServiceContext) http.HandlerFunc {
    	return func(w http.ResponseWriter, r *http.Request) {
    		var req types.UserInfoReq
    		if err := httpx.Parse(r, &req); err != nil {
    			httpx.Error(w, err)
    			return
    		}
    
    		id,ok:=jwtx.GetUserId(w,r) // add
    		if !ok{ // add
    			return  // add
    		}   // add
    		
    		l := logic.NewUserInfoEditLogic(r.Context(), ctx)
    		err := l.UserInfoEdit(id,req)   // edit
    		if err != nil {
    			httpx.Error(w, err)
    		} else {
    			httpx.Ok(w)
    		}
    	}
    }
    ```
* 填充`UserInfoEdit`方法逻辑
    
    ```go
    // 全量更新，允许字段为空
    resp, err := l.svcCtx.UserModel.FindOne(id)
    switch err {
    case nil:
        resp.Name = req.Name
        switch req.Gender {
        case "男":
            resp.Gender = 1
        case "女":
            resp.Gender = 2
        default:
            return errorx.NewInvalidParameterError("gender")
        }
        return l.svcCtx.UserModel.Update(*resp)
    case model.ErrNotFound:
        return logic.ErrUserNotFound
    default:
        return err
    }
    ```
### 添加用户校验中间件
对于获取用户信息，编辑用户信息，我们需要使用jwt鉴权，这样才能知道当前请求的用户是否合法，除此之外，我们还需要传递被修改人的用户id，
而对于这样的需求，用户的鉴权信息中的用户信息必须要和当前操作人的id是一个人，我们可以通过中间件去做一层业务拦截，由于考虑到后续也有这种
场景，这里就将用户信息校验逻辑方在`common`目录下。

在`common`下创建`middleware`文件夹，并添加`usercheckmiddleware.go`文件，填入代码：

```go
package middleware

import (
	"fmt"
	"net/http"

	"hey-go-zero/common/errorx"
	"hey-go-zero/common/jwtx"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UserCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(jwtx.JwtWithUserKey)
		xUserId := r.Header.Get("x-user-id")
		if len(xUserId) == 0 {
			httpx.Error(w, errorx.NewDescriptionError("x-user-id不能为空"))
			return
		}

		if xUserId != fmt.Sprintf("%v", v) {
			httpx.Error(w, errorx.NewDescriptionError("用户信息不一致"))
			return
		}
		next(w, r)
	}
}
```

在main函数文件`service/user/api/user.go`中使用中间件

```go
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	errHandler := errorx.Handler{}
	httpx.SetErrorHandler(errHandler.Handle())

	handler.RegisterHandlers(server, ctx)

	server.Use(middleware.UserCheck) // add: 添加用户信息校验中间件
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
```

最后请求来验证一下以上两条协议

* 修改用户信息
    ```shell script
    $ curl -i -X POST \
        http://localhost:8888/api/user/info/edit \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDcwOTUwMDksImlhdCI6MTYwNzA5MTQwOSwiaWQiOjF9.qx_t1dY3LEoQc-GtGBDASSHpyYx1iba7YrlJyGNk-nA' \
        -H 'x-user-id: 1' \
        -H 'content-type: application/json' \
        -d '{
          "name": "松妹子",
          "gender": "男"
      }'
    ```
    ```text
    HTTP/1.1 200 OK
    Date: Fri, 04 Dec 2020 15:07:59 GMT
    Content-Length: 0
    ```

* 获取用户信息

    ```shell script
    $ curl -i -X GET \
        http://localhost:8888/api/user/info/self \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDcwOTUwMDksImlhdCI6MTYwNzA5MTQwOSwiaWQiOjF9.qx_t1dY3LEoQc-GtGBDASSHpyYx1iba7YrlJyGNk-nA' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 1'
    ```
    ```text
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Fri, 04 Dec 2020 15:09:22 GMT
    Content-Length: 59
    
    {"id":1,"name":"松妹子","gender":"男","role":"student"}
    ```
 
 # 本章节贡献者
 * [songmeizi](https://github.com/songmeizi)
 
 # 技术点总结
 * [正则表达式](https://github.com/ziishaned/learn-regex)
 * [JSON Web Tokens](https://jwt.io/)
 * [Mysql](https://www.mysql.com/)
 * [Redis](https://redis.io/)
 
 # 相关推荐
 * [go-zero微服务框架](https://github.com/tal-tech/go-zero)
 * [超好用的正则在线表达式在线验证网站](https://regex101.com/)
 * [jwt中文社区](http://jwtio.online/)
 * [mysql中文文档](https://www.mysqlzh.com/)
 * [redis命令参考](http://redisdoc.com/index.html)
 
 # 结尾
 本章节完。
 
 如发现任何错误请通过Issue发起问题修复申请。
 
你可能会浏览 
* [课程模块](../../../doc/requirement/course.md)
* [选课模块](../../../doc/requirement/selection.md)
* [排课模块](../../../doc/requirement/schedule.md)
 
 
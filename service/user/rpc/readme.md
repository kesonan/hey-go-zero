# user.rpc
user rpc 用于给其他服务根据其需要提供服务能力，如查询用户信息等。

# 创建rpc目录
在`service/user`下创建rpc目录，得到目录树

```text
user
└── rpc
```

# 新建user.proto
在文件夹`service/user/rpc`上右键->`New`->`New proto file`->`选择Empty file`->`输入user`->`OK`

# 定义proto
我们修改proto文件内容如下:

```protobuf
syntax = "proto3";

package user;

message UserReq {
  // 用户id
  int64 id = 1;
}

enum gender {
  UNKNOWN = 0;
  MALE = 1;
  FEMALE = 2;
}

message UserReply {
  int64 id = 1;
  string name = 2;
  // 性别，0-未知，1-男，2-女
  gender gender = 3;
  // 角色，teacher-教师，student-学生，manager-管理员
  string role = 4;
  int64 createTime = 5;
  int64 updateTime = 6;
}

message IdsReq {
  repeated int64 ids = 1;
}

message UserListReply {
  repeated UserReply list = 1;
}

service UserService {
  // findone
  rpc findOne(UserReq) returns (UserReply);
  // findByIds
  rpc findByIds(IdsReq) returns (UserListReply);
}
```

> 温馨提示：rpc中的方法我们尽量保持逻辑简单切复用性比较高，不要做太多的逻辑判断，否则他和api业务层就没什么区别了；如：查询需要单独查询用户名称、用户角色不需要为
> 每个场景都添加一个rpc method，而是通过一个findOne就够了。

# 生成rpc服务
在`service/user/rpc/user.proto`文件上右键->`Open in Terminal`进入idea终端。

```shell script
$  goctl rpc proto -src user.proto -dir .
```
```text
protoc  -I=/Users/xxx/goland/go/hey-go-zero/service/user/rpc user.proto --go_out=plugins=grpc:/Users/xxx/goland/go/hey-go-zero/service/user/rpc/user
Done.
```

> 说明：执行goctl命令时会输出protoc真正执行的命令内容，其中`xxx`为当前计算机user名称。

我们进入`service/user/rpc`目录看一下生成后的目录树：
```shell script
$ tree
```

```text
user/rpc
├── etc
│   └── user.yaml
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── findbyidslogic.go
│   │   └── findonelogic.go
│   ├── server
│   │   └── userserviceserver.go
│   └── svc
│       └── servicecontext.go
├── readme.md
├── user
│   └── user.pb.go
├── user.go
├── user.proto
└── userservice
    └── userservice.go
```

> 说明：在你们的tree中不会有`readme.md`文件

# 添加`Mysql`和`CacheRedis`配置定义和yaml配置项
* 编打开`service/user/rpc/internal/config/config.go`，添加`Mysql`、`CacheRedis`定义

    ```go
    package config
    
    import (
    	"github.com/tal-tech/go-zero/core/stores/cache"
    	"github.com/tal-tech/go-zero/zrpc"
    )
    
    type Config struct {
    	zrpc.RpcServerConf
    	Mysql struct {
    		DataSource string
    	}
    	CacheRedis cache.CacheConf
    }
    ```
  
* 打开`service/user/rpc/etc/user-api.yaml`文件，添加`Mysql`、`CacheRedis`配置项

    ```yaml
    Name: user.rpc
    ListenOn: 127.0.0.1:8080
    Etcd:
      Hosts:
      - 127.0.0.1:2379
      Key: user.rpc
    Mysql:
      DataSource: ugozero@tcp(127.0.0.1:3306)/heygozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
    CacheRedis:
      -
        Host: 127.0.0.1:6379
        Type: node
    ```

    >说明： 我本地redis没有设置密码，因此没有配置`Password`配置项。

# ServiceContext增加`UserModel`资源
打开`service/user/rpc/internal/svc/servicecontext.go`，添加`UserModel`依赖。

```go
package svc

import (
	"hey-go-zero/service/user/model"
	"hey-go-zero/service/user/rpc/internal/config"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	c config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn:=sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		c: c,
		UserModel: model.NewUserModel(conn,c.CacheRedis),
	}
}
```

# 填充代码逻辑

## 添加一个common.go文件
由于我们不同的方法都是有不同的struct包裹，因此不同logic需要实现共同逻辑时不太灵活，我们这里用一个`静态`型的通用文件去实现这些共同逻辑，然后不同的logic中去访问。
* 在`service/user/rpc/internal/logic`目录中创建`common.go`添加`convertUserFromDbToPb`方法，填充代码：

    ```go
    func convertUserFromDbToPb(in *model.User) *user.UserReply {
    	var resp user.UserReply
    	resp.Id = in.Id
    	resp.Name = in.Name
    	resp.Gender = user.Gender(in.Gender)
    	resp.Role = in.Role
    	resp.CreateTime = in.CreateTime.UnixNano() / 1e6
    	resp.UpdateTime = in.UpdateTime.UnixNano() / 1e6
    	return &resp
    }
    ```
 
## 填充`FindOne`
* 文件位置：`service/user/rpc/internal/logic/findonelogic.go`
* 方法：`FindOne`
* 代码内容：

    ```go
    func (l *FindOneLogic) FindOne(in *user.UserReq) (*user.UserReply, error) {
    	data, err := l.svcCtx.UserModel.FindOne(in.Id)
    	switch err {
    	case nil:
    		return convertUserFromDbToPb(data), nil
    	case model.ErrNotFound:
    		return nil, status.Error(codes.NotFound, err.Error())
    	default:
    		return nil, status.Error(codes.Unknown, err.Error())
    	}
    }
    ```
  
## 填充`FindByIds`
首先我们需要在`service/user/model/usermodel.go`中添加`FindByIds`方法
* 在interface中添加`FindByIds`方法

    ```go
    FindByIds(ids []int64) ([]*User, error)
    ```
* 在default`实现`中添加`FindByIds`方法

    ```go
    func (m *defaultUserModel) FindByIds(ids []int64) ([]*User, error) {
    	query, args, err := builder.Select(userRows).From(m.table).Where(builder.Eq{"id": ids}).ToSQL()
    	if err != nil {
    		return nil, err
    	}
    
    	var resp []*User
    	err = m.CachedConn.QueryRowsNoCache(&resp, query, args...)
    
    	return resp, err
    }
    ```
* 文件位置：`service/user/rpc/internal/logic/findbyidslogic.go`
* 方法：`FindByIds`
* 代码内容：

    ```go
    func (l *FindByIdsLogic) FindByIds(in *user.IdsReq) (*user.UserListReply, error) {
    	var resp user.UserListReply
    	fx.From(func(source chan<- interface{}) {
    		for _, each := range in.Ids {
    			source <- each
    		}
    	}).Split(2000).ForEach(func(item interface{}) {
    		chunks, ok := item.([]interface{})
    		if !ok {
    			return
    		}
    
    		var ids []int64
    		for _, chunk := range chunks {
    			id, ok := chunk.(int64)
    			if !ok {
    				continue
    			}
    
    			ids = append(ids, id)
    		}
    
    		if len(ids) == 0 {
    			return
    		}
    
    		users, err := l.svcCtx.UserModel.FindByIds(ids)
    		if err != nil {
    			logx.Error(err)
    			return
    		}
    
    		for _, each := range users {
    			resp.List = append(resp.List, convertUserFromDbToPb(each))
    		}
    	})
    
    	return &resp, nil
    }
    ```
你可能会浏览
* [用户模块](../../../doc/requirement/user.md)
* [选课模块](../../../doc/requirement/selection.md)

# course.rpc
user rpc 用于给其他服务根据其需要提供服务能力，如查询课程信息等。

 创建rpc目录
在`service/course`下创建rpc目录，得到目录树

```text
course
└── rpc
```

# course.proto
在文件夹`service/couese/rpc`上右键->`New`->`New proto file`->`选择Empty file`->`输入user`->`OK`

# 定义proto
我们修改proto文件内容如下:

```protobuf
syntax = "proto3";

package course;

message IdReq {
  int64 id = 1;
}

message IdsReq {
  repeated int64 ids = 1;
}

enum GenderLimit {
  NoLimit = 0;
  Male = 1;
  Female = 2;
}

message Course {
  int64 id = 1;
  // 名称
  string name = 2;
  // 描述
  string description = 3;
  // 分类
  string classify = 4;
  // 性别限制
  GenderLimit genderLimit = 5;
  // 限制人数
  int64 memberLimit = 6;
  // 当前课程学分
  int64 credit = 7;
}

message CourseListReply {
  repeated Course list = 1;
}

service CourseService{
  // 查询课程
  rpc findOne(IdReq)returns(Course);
  // 批量获取课程
  rpc findByIds(IdsReq)returns(CourseListReply);
}
```

# 生成rpc服务
在`service/course/rpc/course.proto`文件上右键->`Open in Terminal`进入idea终端。

```shell script
$  goctl rpc proto -src course.proto -dir .
```
```text
protoc  -I=/Users/xxx/goland/go/hey-go-zero/service/course/rpc course.proto --go_out=plugins=grpc:/Users/xxx/goland/go/hey-go-zero/service/course/rpc/course
Done.
```

> 说明：执行goctl命令时会输出protoc真正执行的命令内容，其中`xxx`为当前计算机user名称。

我们进入`service/course/rpc`目录看一下生成后的目录树：
```shell script
$ tree
```
```text
course/rpc
├── course
│   └── course.pb.go
├── course.go
├── course.proto
├── courseservice
│   └── courseservice.go
├── etc
│   └── course.yaml
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── findbyidslogic.go
│   │   └── findonelogic.go
│   ├── server
│   │   └── courseserviceserver.go
│   └── svc
│       └── servicecontext.go
└── readme.md
```

> 说明：在你们的tree中不会有`readme.md`文件

# 添加`Mysql`和`CacheRedis`配置定义和yaml配置项
* 编打开`service/course/rpc/internal/config/config.go`，添加`Mysql`、`CacheRedis`定义

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
  
* 打开`service/course/rpc/etc/course.yaml`文件，添加`Mysql`、`CacheRedis`配置项

    ```yaml
    Name: course.rpc
    ListenOn: 127.0.0.1:8081
    Etcd:
      Hosts:
      - 127.0.0.1:2379
      Key: course.rpc
    Mysql:
      DataSource: ugozero@tcp(127.0.0.1:3306)/heygozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
    CacheRedis:
      -
        Host: 127.0.0.1:6379
        Type: node
    ```

    >说明： 我本地redis没有设置密码，因此没有配置`Password`配置项。

# ServiceContext增加`CourseModel`资源
打开`service/course/rpc/internal/svc/servicecontext.go`，添加`CourseModel`依赖。

```go
package svc

import (
	"hey-go-zero/service/course/model"
	"hey-go-zero/service/course/rpc/internal/config"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	c           config.Config
	CourseModel model.CourseModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		c:           c,
		CourseModel: model.NewCourseModel(conn, c.CacheRedis),
	}
}
```

# 填充逻辑

## 添加一个common.go文件
由于我们不同的方法都是有不同的struct包裹，因此不同logic需要实现共同逻辑时不太灵活，我们这里用一个`静态`型的通用文件去实现这些共同逻辑，然后不同的logic中去访问。
* 在`service/course/rpc/internal/logic`目录中创建`common.go`并添加`convertCourseFromDbToPb`方法，填充代码：

    ```go
    func convertCourseFromDbToPb(in *model.Course) *course.Course {
    	var resp course.Course
    	resp.Id = in.Id
    	resp.Name = in.Name
    	resp.Description = in.Description
    	resp.Classify = in.Classify
    	resp.GenderLimit = course.GenderLimit(in.GenderLimit)
    	resp.MemberLimit = in.MemberLimit
    	resp.Credit = in.Credit
    	return &resp
    }
    ```

## 填充`FindOne`

* 文件位置：`service/course/rpc/internal/logic/findonelogic.go`
* 方法： `FindOne`
* 代码内容：

    ```go
    func (l *FindOneLogic) FindOne(in *course.IdReq) (*course.Course, error) {
    	data, err := l.svcCtx.CourseModel.FindOne(in.Id)
    	switch err {
    	case nil:
    		return convertCourseFromDbToPb(data), nil
    	case model.ErrNotFound:
    		return nil, status.Error(codes.NotFound, err.Error())
    	default:
    		return nil, status.Error(codes.Unknown, err.Error())
    	}
    }
    ```
  
## 填充`FindByIds`

首先我们需要在`service/course/model/coursemodel.go`中添加`FindByIds`方法
* 在interface中添加`FindByIds`方法

    ```go
    FindByIds(ids []int64) ([]*Course, error)
    ```
* 在default`实现`中添加`FindByIds`方法

    ```go
    func (m *defaultCourseModel) FindByIds(ids []int64) ([]*Course, error) {
    	query, args, err := builder.Select(courseRows).From(m.table).Where(builder.Eq{"id": ids}).ToSQL()
    	if err != nil {
    		return nil, err
    	}
    
    	var resp []*Course
    	err = m.CachedConn.QueryRowsNoCache(&resp, query, args...)
    
    	return resp, err
    }
    ```
* 文件位置：`service/course/rpc/internal/logic/findbyidslogic.go`
* 方法： `FindByIds`
* 代码内容：

    ```go
    func (l *FindByIdsLogic) FindByIds(in *course.IdsReq) (*course.CourseListReply, error) {
    	var resp course.CourseListReply
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
    
    		users, err := l.svcCtx.CourseModel.FindByIds(ids)
    		if err != nil {
    			logx.Error(err)
    			return
    		}
    
    		for _, each := range users {
    			resp.List = append(resp.List, convertCourseFromDbToPb(each))
    		}
    	})
    
    	return &resp, nil
    }
    ```
你可能会浏览
* [用户模块](../../../doc/requirement/user.md)
* [选课模块](../../../doc/requirement/selection.md)
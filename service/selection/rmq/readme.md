# 选课模块rmq

# 创建目录结构
在前面的[服务目录](../../../doc/prepare/service-structure.md)章节中我们讲到了使用什么样的目录结构，这里就不赘述了，我们来创建一个`rmq`服务
* 在`service/selection`目录下创建`rmq`文件夹
* 在`rmq`目录下创建`selection.go`文件
* 然后在`rmq`目录下创建以下目录结构
    
    ``` text
    rmq
    ├── etc
    │   └── selection-rmq.yaml
    ├── internal
    │   ├── config
    │   │   └── config.go
    │   ├── logic
    │   │   └── consumer.go
    │   └── svc
    │       └── servicecontext.go
    ├── readme.md
    └── selection.go
    ```
    > 说明：在你们的结构中不会包含`readme.md`文件

# 定义配置和添加配置项
* config.go
    * 文件位置：`service/selection/rmq/internal/config/config.go`
    * 代码内容：
    
        ``` go
        package config
        
        import (
            "github.com/tal-tech/go-queue/dq"
            "github.com/tal-tech/go-zero/core/service"
        )
        
        type Config struct {
            service.ServiceConf
            Dq dq.DqConf
        }
        ```
  * etc.yaml
    * 文件位置：`service/selection/rmq/etc/selection-rmq.yaml`
    * 代码内容：
    
        ``` yaml
        Name: selection.rmq
        Log:
          Mode: console
        Dq:
          Beanstalks:
            -
              Endpoint: 127.0.0.1:11300
              Tube: course_select
            - Endpoint: 127.0.0.1:11301
              Tube: course_select
          Redis:
            Host: 127.0.0.1:6379
            Type: node
        ```

# 添加配置依赖servicecontext.go
* 文件位置：`service/selection/rmq/internal/svc/servicecontext.go`
* 代码内容：

    ``` go
    package svc
    
    import (
    	"hey-go-zero/service/selection/rmq/internal/config"
    
    	"github.com/tal-tech/go-queue/dq"
    )
    
    type ServiceContext struct {
    	Config   config.Config
    	Consumer dq.Consumer
    }
    
    func NewServiceContext(c config.Config) *ServiceContext {
    	return &ServiceContext{
    		Config:   c,
    		Consumer: dq.NewConsumer(c.Dq),
    	}
    }
    ```

# 填充消费逻辑
* 文件位置：`service/selection/rmq/internal/logic/consumer.go`
* 代码位置：

    ``` go
    package logic
    
    import (
    	"context"
    
    	"hey-go-zero/service/selection/rmq/internal/svc"
    
    	"github.com/tal-tech/go-zero/core/logx"
    )
    
    type ConsumerLogic struct {
    	logx.Logger
    	ctx    context.Context
    	svcCtx *svc.ServiceContext
    }
    
    func NewConsumerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumerLogic {
    	return &ConsumerLogic{
    		Logger: logx.WithContext(ctx),
    		ctx:    ctx,
    		svcCtx: svcCtx,
    	}
    }
    
    func (c *ConsumerLogic) Start() {
    	c.svcCtx.Consumer.Consume(c.Consume)
    }
    
    func (c *ConsumerLogic) Stop() {
    }
    
    func (c *ConsumerLogic) Consume(body []byte) {
    	msg := string(body)
    	logx.Info("consume:" + msg)
    }
    ```

# 填充main函数逻辑
* 文件位置：`service/selection/rmq/selection.go`
* 代码内容：

    ``` go
    package main
    
    import (
    	"context"
    	"flag"
    
    	"hey-go-zero/service/selection/rmq/internal/config"
    	"hey-go-zero/service/selection/rmq/internal/logic"
    	"hey-go-zero/service/selection/rmq/internal/svc"
    
    	"github.com/tal-tech/go-zero/core/conf"
    	"github.com/tal-tech/go-zero/core/logx"
    	"github.com/tal-tech/go-zero/core/service"
    )
    
    var configFile = flag.String("f", "etc/selection-rmq.yaml", "the config file")
    
    func main() {
    	flag.Parse()
    
    	var c config.Config
    	conf.MustLoad(*configFile, &c)
    	logx.MustSetup(c.Log)
    
    	ctx := svc.NewServiceContext(c)
    	l := logic.NewConsumerLogic(context.Background(), ctx)
    	sg := service.NewServiceGroup()
    	sg.Add(l)
    	sg.Start()
    }
    ```

# 测试消费

## 启动服务
前面需要启动的步骤和[选课api模块](../api/readme.md)一样,最后在启动本服务
``` shell script
$ go run selection.go
```

## 创建选课
``` shell script
$ curl -i -X POST \
    http://127.0.0.1:8890/api/selection/create \
    -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgyMTU5MjgsImlhdCI6MTYwODIxMjMyOCwiaWQiOjJ9.NKTaSB5CHyaoaWQOhwa6QM-ZhHDAojmsFK_dj9gOzCY' \
    -H 'content-type: application/json' \
    -H 'x-user-id: 2' \
    -d '{
  	"name":"测试4",
  	"maxCredit":12,
  	"startTime":1608220500,
  	"endTime":1608306192,
  	"notification":"选课开始了。"
  }'
```
> 说明：代码中消费的时间点是startTime前两小时，这里为了不等待太长时间来验证消费消息，你可以设置startTime为当前时间往后推2小时01分钟，那么消息就将会在1分钟后消费。
## 等待消费
``` shell script
{"@timestamp":"2020-12-17T21:49:10.291+08","level":"stat","content":"CPU: 0m, MEMORY: Alloc=0.6Mi, TotalAlloc=0.6Mi, Sys=70.2Mi, NumGC=0"}
{"@timestamp":"2020-12-17T21:50:10.298+08","level":"stat","content":"CPU: 0m, MEMORY: Alloc=0.6Mi, TotalAlloc=0.6Mi, Sys=70.2Mi, NumGC=0"}
{"@timestamp":"2020-12-17T21:51:10.295+08","level":"stat","content":"CPU: 0m, MEMORY: Alloc=0.6Mi, TotalAlloc=0.6Mi, Sys=70.2Mi, NumGC=0"}
{"@timestamp":"2020-12-17T21:52:10.297+08","level":"stat","content":"CPU: 0m, MEMORY: Alloc=0.6Mi, TotalAlloc=0.6Mi, Sys=70.2Mi, NumGC=0"}
{"@timestamp":"2020-12-17T21:53:10.299+08","level":"stat","content":"CPU: 0m, MEMORY: Alloc=0.6Mi, TotalAlloc=0.6Mi, Sys=70.2Mi, NumGC=0"}
{"@timestamp":"2020-12-17T21:54:10.299+08","level":"stat","content":"CPU: 0m, MEMORY: Alloc=0.6Mi, TotalAlloc=0.6Mi, Sys=70.2Mi, NumGC=0"}
选课开始了。
```

# 本章节贡献者
 * [anqiansong](https://github.com/anqiansong)
 
 # 技术点总结
 * go-zero中间件使用
    * 全局中间件
    * 指定路由组中间件
 * go-zero自定义错误
 * go-zero rpc调用
 * [go-queue](https://github.com/tal-tech/go-queue)
    * dq
 * redis(BizRedis)
    * string
    * map
 * redisLock
 
 # 相关推荐
 * [dq/kq使用说明](https://github.com/tal-tech/go-queue)
 * [beanstalkd](https://beanstalkd.github.io)
 
 # 结尾
 本章节完。
 
 如发现任何错误请通过Issue发起问题修复申请。
 
你可能会浏览 
* [用户模块](../../../doc/requirement/user.md)
* [课程模块](../../../doc/requirement/course.md)


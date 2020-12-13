# course-api创建步骤
本文档为course-api的创建流程步骤文档，通过本文档你可以了解到maintainer是如何完成对course api的实现。

# 创建api目录
在`service/source`下创建api目录得到目录树

```text
course
└── api
```

# 创建course.api文件
在`service/course/api`文件夹`右键`->`New Api File`->`输入course`->`选择Empty file`->`回车`，
然后修改course.api文件内容为

```go
info(
	title: "课程管理api"
	desc: "描述课程添加、编辑、删除、查看等协议"
	author: "松妹子"
	version: "V1.0"
)

type (
	Course {
		Name string `json:"name"`
		Description string `json:"description,optional"`
		Classify string `json:"classify,options=天文|地理|数学|物理|机械|航天|医学|信息|互联网|计算机"`
		// 性别限制，0-不限，1-男，2-女
		GenderLimit int `json:"genderLimit,options=0|1|2"`
		// 可选参数，如果不传则代表不限制人数
		MemberLimit MemberLimit `json:"memberLimit,optional"`
		StartTime int64 `json:"startTime"`
		// 学分
		Credit int `json:"credit,range=(0:6]"`
	}

	MemberLimit {
		// 男生限制人数 <=0：不限
		MaleCount int `json:"maleCount"`
		// 女生限制人数 <=0：不限
		FemaleCount int `json:"femaleCount"`
	}

	AddCourseReq {
		Course
	}
	
	EditCourseReq {
		Id int64 `path:"id"`
		Course
	}
	
	DeleteCourseReq {
		Id int64 `path:"id"`
	}
	
	CourseInfoReq {
		Id int64 `path:"id"`
	}
	
	CourseInfoReply {
		Id int64 `json:"id"`
		Course
	}
	
	CourseListReq {
		Page int `form:"page,range=(0:]"`
		Size int `form:"size,range=(0:]"`
	}
	
	CourseListReply {
		// 总条数
		Total int `json:"total"`
		// 当前返回数量，即list.length
		Size int `json:"size"`
		List []*CourseInfoReply `json:"list"`
	}
)

@server(
	jwt: Auth
	middleware: AuthMiddleware
)
service course-api {
	@handler addCourse
	post /api/course/add (AddCourseReq)
	
	@handler editCourse
	post /api/course/edit/:id (EditCourseReq)
	
	@handler deleteCourse
	post /api/course/delete/:id (DeleteCourseReq)
	
	@handler getCourseInfo
	get /api/course/:id (CourseInfoReq) returns (CourseInfoReply)
	
	@handler getCourseList
	get /api/course/list (CourseListReq) returns (CourseListReply)
}
```

# 生成代码
在Goland中生成代码有三种方式（任意一种均可）
* project面板区文件右键生成
    * 选中`course.api`文件->`右键`->`New`->`Go Zero`->`Api Code`
    * `Api Code`回车后会弹出一个文件对话框询问你需要生成服务的目标目录，默认为`course.api`所在目录，我们这里选择默认，点击`OK`确认生成。
* api文件编辑区右键生成
    * 打开`course.api`文件->`编辑区`->`右键`->`Generate..`->`Api Code`
    * `Api Code`回车后会弹出一个文件对话框询问你需要生成服务的目标目录，默认为`course.api`所在目录，我们这里选择默认，点击`OK`确认生成。
* 终端生成（推荐）
    * 选中`course.api`文件->`右键`->`Open in Terminal`
    * 执行`goctl api go -api course.api -dir .`命令即可
    
        ```shell script
        $ goctl api go -api user.api -dir .
        ```
        ```text
        Done.
        ```
接下来我们进入`service/course/api`目录，查看一下目录树结构
  
```text
course/api
├── course.api
├── course.go
├── etc
│   └── course-api.yaml
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── addcoursehandler.go
│   │   ├── deletecoursehandler.go
│   │   ├── editcoursehandler.go
│   │   ├── getcourseinfohandler.go
│   │   ├── getcourselisthandler.go
│   │   └── routes.go
│   ├── logic
│   │   ├── addcourselogic.go
│   │   ├── deletecourselogic.go
│   │   ├── editcourselogic.go
│   │   ├── getcourseinfologic.go
│   │   └── getcourselistlogic.go
│   ├── middleware
│   │   └── authmiddleware.go
│   ├── svc
│   │   └── servicecontext.go
│   └── types
│       └── types.go
└── readme.md
```

> 说明：上述目录中的注释是为了大家能够快速知道该目录结构的用途，是后期我加入的，实际生成的tree不会带注释和readme.md文件。

# 生成model层代码

### 新增model文件夹
在`service/course`目录下添加`model`文件夹

### 新建表course
在前面我们已经演示了如何[创建db和table](../../../doc/prepare/db-create.md) ，我们在`heygozero`新建一张表名为`course`,create table ddl如下：

```mysql
CREATE TABLE `course` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '书籍名称',
  `description` varchar(500) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '书籍描述',
  `classify` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '书籍分类，目前仅支持 【天文|地理|数学|物理|机械|航天|医学|信息|互联网|计算机】',
  `gender_limit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别限制 0-不限，1-男，2-女',
  `male_limit` int DEFAULT '0' COMMENT '男生限制人数 0-不限',
  `female_limit` int DEFAULT '0' COMMENT '女生限制人数 0-不限',
  `start_time` int NOT NULL DEFAULT '0' COMMENT '开课时间，时间戳，单位：毫秒',
  `credit` tinyint(1) NOT NULL DEFAULT '1' COMMENT '学分',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

### 生成代码
在文件夹`service/course/model`右键->`Open in Terminal` 进入终端

```shell script
$ goctl model mysql datasource -url="ugozero@tcp(127.0.0.1:3306)/heygozero" -table="course" -c -dir .
```
```text
Done.
```

生成后我么可以看到会产生两个文件
```text
model
├── coursemodel.go
└── vars.go
```

> 说明：本地生成的`goctl`版本为`goctl version 20201125 darwin/amd64`，早起版本生成出来的数值类型会有`int`，`int64`，而后续版本统一为`int64`了。

# 添加`Mysql`和`CacheRedis`配置定义和yaml配置项
* 编打开`service/course/api/internal/config/config.go`，添加`Mysql`、`CacheRedis`定义

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
  
* 打开`service/course/api/etc/course-api.yaml`文件，添加`Mysql`、`CacheRedis`配置项

    ```yaml
    Name: user-api
    Host: 0.0.0.0
    Port: 8889
    Auth:
      AccessSecret: 1e69481b-7405-4369-9ce3-9aaffdb56ce3
      AccessExpire: 3600
    Mysql:
      DataSource: ugozero@tcp(127.0.0.1:3306)/heygozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
    CacheRedis:
      -
        Host: 127.0.0.1:6379
        Type: node
    ```

    >说明： 我本地redis没有设置密码，因此没有配置`Password`配置项，为了防止和user api端口冲突，这里将端口修改为`8889`

# ServiceContext增加`CourseModel`资源
进入文件`service/course/api/internal/svc/servicecontext.go`，添加`CourseModel`资源依赖。

```go
package svc

import (
	"hey-go-zero/service/course/api/internal/config"
	"hey-go-zero/service/course/api/internal/middleware"
	"hey-go-zero/service/course/model"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	CourseModel    model.CourseModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware().Handle,
		CourseModel:    model.NewCourseModel(conn, c.CacheRedis),
	}
}
```

# 添加用户信息校验中间件和自定义错误码
由于在user api中已经添加过`UserCheck`用户信息校验中间件和[自定义错误码](../../../doc/gozero/http-error.md)了，这里就直接在main文件中使用即可。

```go
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	errHandler := errorx.Handler{} // add 自定义错误码
	httpx.SetErrorHandler(errHandler.Handle()) // add 自定义错误码
	
	handler.RegisterHandlers(server, ctx)
	server.Use(middleware.UserCheck) // add 用户信息校验中间件
	
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
```

# 填充课程逻辑

## 创建error.go文件
在`service/course/api/internal/logic`目录下新增`error.go`文件，填充代码

```go
var (
	errCourseNotFound = errorx.NewDescriptionError("课程不存在")
)
```

## 添加common.go
在`service/course/api/internal/logic`目录下新增`common.go`文件，填充代码

```go
import (
	"hey-go-zero/service/course/api/internal/types"
	"hey-go-zero/service/course/model"
)

func convertFromDbToLogic(data model.Course) types.Course {
	return types.Course{
		Name:        data.Name,
		Description: data.Description,
		Classify:    data.Classify,
		GenderLimit: int(data.GenderLimit),
		MemberLimit: types.MemberLimit{
			MaleCount:   int(data.MaleLimit),
			FemaleCount: int(data.FemaleLimit),
		},
		StartTime: data.StartTime,
		Credit:    int(data.Credit),
	}
}
```

## 添加课程
* 文件位置：`service/course/api/internal/logic/addcourselogic.go`
* 方法：`AddCourse`
* 代码内容

    ```go
    func (l *AddCourseLogic) AddCourse(req types.AddCourseReq) error {
        if err := l.parametersCheck(req); err != nil {
            return err
        }
    
        // 如果数量小于等于0则为不限
        if req.MemberLimit.MaleCount < 0 {
            req.MemberLimit.MaleCount = 0
        }
  
        if req.MemberLimit.FemaleCount < 0 {
            req.MemberLimit.FemaleCount = 0
        }
    
        _, err := l.svcCtx.CourseModel.FindOneByName(req.Name)
        switch err {
        case nil:
            return errorx.NewDescriptionError("课程已存在")
        case model.ErrNotFound:
            _, err = l.svcCtx.CourseModel.Insert(model.Course{
                Name:        req.Name,
                Description: req.Description,
                Classify:    req.Classify,
                GenderLimit: int64(req.GenderLimit),
                MaleLimit:   int64(req.MemberLimit.MaleCount),
                FemaleLimit: int64(req.MemberLimit.FemaleCount),
                StartTime:   req.StartTime,
                Credit:      int64(req.Credit),
            })
            return err
        default:
            return err
        }
    }
    
    func (l *AddCourseLogic) parametersCheck(req types.AddCourseReq) error {
    	wordLimitErr := func(key string, limit int) error {
    		return errorx.NewDescriptionError(fmt.Sprintf("%s不能超过%d个字符", key, limit))
    	}
    
    	if len(strings.TrimSpace(req.Name)) == 0 {
    		return errorx.NewInvalidParameterError("name")
    	}
    	
    	if utf8.RuneCountInString(req.Name) > 20 {
    		return wordLimitErr("课程名称", 20)
    	}
    
    	if utf8.RuneCountInString(req.Description) > 500 {
    		return wordLimitErr("课程描述", 500)
    	}
    
    	now := time.Now().AddDate(0, 0, 1)
    	validEarliestStartTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
    	if req.StartTime < validEarliestStartTime.Unix() {
    		return errorx.NewDescriptionError(fmt.Sprintf("开课时间不能早于%s", validEarliestStartTime.Format("2006年01月02日 03时04分05秒")))
    	}
    
    	return nil
    }
    ```

> 说明：这里主要是带着大家熟悉api的开发，就不介绍具体详细业务逻辑了，感兴趣可以自己看代码，对代码中逻辑有争议我们这里也不争论。

## 编辑课程
* 文件位置：`service/course/api/internal/logic/editcourselogic.go`
* 方法：`EditCourse`
* 代码内容

    ```go
    func (l *EditCourseLogic) EditCourse(req types.EditCourseReq) error {
        if err := l.parametersCheck(req); err != nil {
            return err
        }
        
        data, err := l.svcCtx.CourseModel.FindOne(req.Id)
        switch err {
        case nil:
            data.Name = req.Name
            data.Description = req.Description
            data.Classify = req.Classify
            data.GenderLimit = int64(req.GenderLimit)
            data.MaleLimit = int64(req.MemberLimit.MaleCount)
            data.FemaleLimit = int64(req.MemberLimit.FemaleCount)
            data.StartTime = req.StartTime
            data.Credit = int64(req.Credit)
            return l.svcCtx.CourseModel.Update(*data)
        case model.ErrNotFound:
            return errCourseNotFound
        default:
            return err
        }
    }
    
    func (l *EditCourseLogic) parametersCheck(req types.EditCourseReq) error {
    	wordLimitErr := func(key string, limit int) error {
    		return errorx.NewDescriptionError(fmt.Sprintf("%s不能超过%d个字符", key, limit))
    	}
    
    	if req.Id < 0 {
    		return errorx.NewInvalidParameterError("id")
    	}
    
    	if len(strings.TrimSpace(req.Name)) == 0 {
    		return errorx.NewInvalidParameterError("name")
    	}
    
    	if utf8.RuneCountInString(req.Name) > 20 {
    		return wordLimitErr("课程名称", 20)
    	}
    
    	if utf8.RuneCountInString(req.Description) > 500 {
    		return wordLimitErr("课程描述", 500)
    	}
    
    	now := time.Now().AddDate(0, 0, 1)
    	validEarliestStartTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
    	if req.StartTime < validEarliestStartTime.Unix() {
    		return errorx.NewDescriptionError(fmt.Sprintf("开课时间不能早于%s", validEarliestStartTime.Format("2006年01月02日 03时04分05秒")))
    	}
    
    	return nil
    }
    ```

## 删除课程
* 文件位置：`service/course/api/internal/logic/deletecourselogic.go`
* 方法：`DeleteCourse`
* 代码内容

    ```go
    func (l *DeleteCourseLogic) DeleteCourse(req types.DeleteCourseReq) error {
    	if req.Id <= 0 {
    		return errorx.NewInvalidParameterError("id")
    	}
    
    	err:=l.svcCtx.CourseModel.Delete(req.Id)
    	switch err {
    	case nil:
    		return nil
    	case model.ErrNotFound:
    		return errCourseNotFound
    	default:
    		return err
    	}
    }
    ```

## 查看课程
* 文件位置：`service/course/api/internal/logic/getcourseinfologic.go`
* 方法：`GetCourseInfo`
* 代码内容

    ```go
    func (l *GetCourseInfoLogic) GetCourseInfo(req types.CourseInfoReq) (*types.CourseInfoReply, error) {
        if req.Id <= 0 {
            return nil, errorx.NewInvalidParameterError("id")
        }
    
        data, err := l.svcCtx.CourseModel.FindOne(req.Id)
        switch err {
        case nil:
            return &types.CourseInfoReply{
                Id: data.Id,
                Course: convertFromDbToLogic(*data),
            }, nil
        case model.ErrNotFound:
            return nil, errCourseNotFound
        default:
            return nil, err
        }
    }
    ```

## 获取课程列表
* 文件位置：`service/course/model/coursemodel.go`
* 添加分页查询逻辑
    * 在interface中添加两个方法`FindAllCount`和`FindLimit`
    
        ```go
        CourseModel interface {
            Insert(data Course) (sql.Result, error)
            FindOne(id int64) (*Course, error)
            FindOneByName(name string) (*Course, error)
            Update(data Course) error
            Delete(id int64) error
            FindAllCount() (int, error) // add
            FindLimit(page, size int) ([]*Course, error) // add
        }
        ```
    * 在添加两个default实现方法
        
        ```go
        func (m *defaultCourseModel) FindAllCount() (int, error) {
        	query := fmt.Sprintf("select count(id) from %s", m.table)
        	var count int
        	err := m.CachedConn.QueryRowNoCache(&count, query)
        	return count, err
        }
        
        func (m *defaultCourseModel) FindLimit(page, size int) ([]*Course, error) {
        	query := fmt.Sprintf("select %s from %s order by id limit ?,?", courseRows, m.table)
        	var resp []*Course
        	err := m.CachedConn.QueryRowsNoCache(&resp, query, (page-1)*size, size)
        	return resp, err
        }
        ```
      
  * 文件位置：`service/course/api/internal/logic/getcourselistlogic.go`
  * 方法：`GetCourseList`
  * 代码内容：
  
    ```go
    func (l *GetCourseListLogic) GetCourseList(req types.CourseListReq) (*types.CourseListReply, error) {
    	total, err := l.svcCtx.CourseModel.FindAllCount()
    	if err != nil {
    		return nil, err
    	}
    
    	data, err := l.svcCtx.CourseModel.FindLimit(req.Page, req.Size)
    	if err != nil {
    		return nil, err
    	}
    
    	var list []*types.CourseInfoReply
    	for _, item := range data {
    		list = append(list, &types.CourseInfoReply{
    			Id:     item.Id,
    			Course: convertFromDbToLogic(*item),
    		})
    	}
    
    	return &types.CourseListReply{
    		Total: total,
    		Size:  len(list),
    		List:  list,
    	}, nil
    }
    ```
## 添加中间逻辑
由于课程在api层级的操作只允许管理员操作，因此我们需要对能操作课程人员做一下限制，在前面[用户模块](../../../doc/requirement/user.md)我们介绍了用户分为三种角色分别为T、
S、M。我们这里就限制仅M可访问该层协议，在实现中间件逻辑前，我们先手动添加一个角色用为M的用户到数据库（为了简单，这里就不单独对角色为M的用户进行管理了，直接手动插入数据库）

### 添加管理员角色用户
```mysql
$ mysql -h 127.0.0.1 -uugozero -p
  Enter password:
  Welcome to the MySQL monitor.  Commands end with ; or \g.
  Your MySQL connection id is 13
  Server version: 8.0.21 MySQL Community Server - GPL
  
  Copyright (c) 2000, 2020, Oracle and/or its affiliates. All rights reserved.
  
  Oracle is a registered trademark of Oracle Corporation and/or its
  affiliates. Other names may be trademarks of their respective
  owners.
  
  Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
  
  mysql> use heygozero
  Reading table information for completion of table and column names
  You can turn off this feature to get a quicker startup with -A
  
  Database changed
  mysql> show tables;
  +---------------------+
  | Tables_in_heygozero |
  +---------------------+
  | course              |
  | user                |
  +---------------------+
  2 rows in set (0.00 sec)
  
  mysql> insert into user (username,password,name,gender,role) value ('gozero','111111','admin',1,'manager');
  Query OK, 1 row affected (0.01 sec)
  
  mysql> select * from user where role = 'manager' limit 1;
 +----+----------+----------+-------+--------+---------+---------------------+---------------------+
 | id | username | password | name  | gender | role    | create_time         | update_time         |
 +----+----------+----------+-------+--------+---------+---------------------+---------------------+
 |  2 | gozero    | 111111   | admin |      1 | manager | 2020-12-12 22:32:37 | 2020-12-12 22:32:37 |
 +----+----------+----------+-------+--------+---------+---------------------+---------------------+
 1 row in set (0.01 sec)
 
 mysql>
```

### 创建user.rpc
在这里，我们需要访问`user`表的用户信息了，因此就需要RPC来进行微服务间的通讯，所以，在此前还要创建一个user.rpc服务来为我们传递信息，关于user.rpc逻辑请查看[user rpc创建](../../user/rpc/readme.md)

### 添加`UserRpc`配置定义和yaml配置项
* 编打开`service/course/api/internal/config/config.go`，添加`UserRpc`定义

    ```go
    package config
    
    import (
    	"github.com/tal-tech/go-zero/core/stores/cache"
    	"github.com/tal-tech/go-zero/rest"
    	"github.com/tal-tech/go-zero/zrpc"
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
    	UserRpc    zrpc.RpcClientConf
    }
    ```
  
* 打开`service/course/api/etc/course-api.yaml`文件，添加`UserRpc`配置项

    ```yaml
    Name: course-api
    Host: 0.0.0.0
    Port: 8889
    Auth:
      AccessSecret: 1e69481b-7405-4369-9ce3-9aaffdb56ce3
      AccessExpire: 3600
    Mysql:
      DataSource: ugozero@tcp(127.0.0.1:3306)/heygozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
    CacheRedis:
      -
        Host: 127.0.0.1:6379
        Type: node
    UserRpc:
      Etcd:
        Hosts:
          - 127.0.0.1:2379
        Key: user.rpc
    ```

    > 说明： 我本地redis没有设置密码，因此没有配置`Password`配置项。

### ServiceContext增加`UserRpcClient`资源
打开`service/course/api/internal/svc/servicecontext.go`，添加`UserRpcClient`依赖。

```go
package svc

import (
	"hey-go-zero/service/course/api/internal/config"
	"hey-go-zero/service/course/api/internal/middleware"
	"hey-go-zero/service/course/model"
	"hey-go-zero/service/user/rpc/userservice"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	CourseModel    model.CourseModel
	UserRpcClient  userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	userRpcClient := zrpc.MustNewClient(c.UserRpc)
	userRpcService := userservice.NewUserService(userRpcClient)
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(userRpcService).Handle,
		CourseModel:    model.NewCourseModel(conn, c.CacheRedis),
		UserRpcClient:  userRpcService,
	}
}
```

### 修改中间件代码

* 添加`userRpcClient`依赖

    ```go
    type AuthMiddleware struct {
    	userRpcClient userservice.UserService
    }
    
    func NewAuthMiddleware(userRpcClient userservice.UserService) *AuthMiddleware {
    	return &AuthMiddleware{
    		userRpcClient: userRpcClient,
    	}
    }
    ```
* 完善中间件逻辑

    ```go
    userId, ok := jwtx.GetUserId(w, r)
    if !ok {
        return
    }

    data, err := m.userRpcClient.FindOne(r.Context(), &userservice.UserReq{
        Id: userId,
    })
    if err != nil {
        st := status.Convert(err)
        if st.Code() == codes.NotFound {
            httpx.Error(w, errorx.NewDescriptionError("用户不存在"))
            return
        }

        httpx.Error(w, errorx.NewDescriptionError("用户信息获取失败"))
        return
    }

    if data.Role != "manager" {
        httpx.Error(w, errorx.NewDescriptionError("无权限访问"))
        return
    }

    next(w, r)
    ```

# 请求服务

## 启动user.api
进入文件夹`service/user/api`，右键进入Idea终端

```shell script
$ go run user.go 
```

```text
Starting server at 0.0.0.0:8888...
```

## 启动redis

```shell script
$ redis-server
```

## 启动etcd

```shell script
$ etcd
```

## 启动user.rpc
进入文件夹`service/user/rpc`，通过右键进入Idea终端

```shell script
$ go run user.go
```
```text
Starting rpc server at 127.0.0.1:8080...
```

## 启动course-api服务
进入文件夹`service/course/api`，右键进入Idea终端

```shell script
$ go run course.go
```
```text
Starting server at 0.0.0.0:8889...
```

## 访问服务
* 登录

```shell script
$ curl -i -X POST \
    http://127.0.0.1:8888/api/user/login \
    -H 'cache-control: no-cache' \
    -H 'content-type: application/json' \
    -d '{
  	"username":"gozero",
      "password":"111111"
  }'
```
```text
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 13 Dec 2020 12:48:58 GMT
Content-Length: 178

{"id":2,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4NjczMzgsImlhdCI6MTYwNzg2MzczOCwiaWQiOjJ9.NPcA430BHJ3L_V_JsJiOEt00GAPIb2PMhN0TOnbc5Xk","expireAt":1607867338}
```

* 添加课程

    我们先用一个非管理员身份去访问试一下（重新以其他身份登录获取token）
    
    ```shell script
    $ curl -i -X POST \
        http://127.0.0.1:8889/api/course/add \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4Njc1MjYsImlhdCI6MTYwNzg2MzkyNiwiaWQiOjF9.-CfZi6UQ5SFEQkyVZizPyvT6oQkZ7oAMWRVTGthWl8U' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 1' \
        -d '{
        
      }'
    ```
    ```text
    HTTP/1.1 406 Not Acceptable
    Content-Type: application/json
    Date: Sun, 13 Dec 2020 12:55:41 GMT
    Content-Length: 36
    
    {"code":-1,"desc":"无权限访问"}
    ```
    
    正如我们所期望的那样，非管理员身份是禁止访问的。所以这里我们的中间件生效了；
    下面我们以管理员身份去操作课程吧
    
    ```shell script
    $ curl -i -X POST \
        http://127.0.0.1:8889/api/course/add \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4NjczMzgsImlhdCI6MTYwNzg2MzczOCwiaWQiOjJ9.NPcA430BHJ3L_V_JsJiOEt00GAPIb2PMhN0TOnbc5Xk' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2' \
        -d '{
      	"name":"Golang语言开发",
      	"description":"Golang语言开发从入门到放弃！",
      	"classify":"计算机",
      	"genderLimit":0,
      	"startTime":1607911200,
      	"credit":1
      }'
    ```
    ```text
    HTTP/1.1 200 OK
    Date: Sun, 13 Dec 2020 13:12:22 GMT
    Content-Length: 0
    ```
    再以相同的参数去请求一次会得到
    ```text
    HTTP/1.1 406 Not Acceptable
    Content-Type: application/json
    Date: Sun, 13 Dec 2020 13:13:04 GMT
    Content-Length: 36
    
    {"code":-1,"desc":"课程已存在"}
    ```

* 编辑课程

    我们先找一个不存在的id去编辑一下，验证一下逻辑
    
    ```shell script
    $ curl -i -X POST \
        http://127.0.0.1:8889/api/course/edit/3 \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4NjczMzgsImlhdCI6MTYwNzg2MzczOCwiaWQiOjJ9.NPcA430BHJ3L_V_JsJiOEt00GAPIb2PMhN0TOnbc5Xk' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2' \
        -d '{
      	"name":"Golang语言开发",
      	"description":"Golang语言开发从入门到放弃！",
      	"classify":"计算机",
      	"genderLimit":0,
      	"startTime":1607911200,
      	"credit":1
      }'
    ```
    ```text
    HTTP/1.1 406 Not Acceptable
    Content-Type: application/json
    Date: Sun, 13 Dec 2020 13:15:30 GMT
    Content-Length: 36
    
    {"code":-1,"desc":"课程不存在"}
    ```
    正如我们预料的那样会提示`课程不存在`的错误。
    下面我们来对之前添加过的课程做一下编辑
    ```shell script
    $ curl -i -X POST \
        http://127.0.0.1:8889/api/course/edit/2 \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4NjczMzgsImlhdCI6MTYwNzg2MzczOCwiaWQiOjJ9.NPcA430BHJ3L_V_JsJiOEt00GAPIb2PMhN0TOnbc5Xk' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2' \
        -d '{
        "name":"Golang语言开发",
        "description":"Golang语言开发从入门到放弃！",
        "classify":"互联网",
        "genderLimit":0,
        "startTime":1607911200,
        "credit":1
      }'
    ```
    ```text
    HTTP/1.1 200 OK
    Date: Sun, 13 Dec 2020 13:18:31 GMT
    Content-Length: 0
    ```

* 查询课程

    ```shell script
    $ curl -i -X GET \
        http://127.0.0.1:8889/api/course/2 \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4NjczMzgsImlhdCI6MTYwNzg2MzczOCwiaWQiOjJ9.NPcA430BHJ3L_V_JsJiOEt00GAPIb2PMhN0TOnbc5Xk' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2'
    ```
    ```text
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Sun, 13 Dec 2020 13:19:54 GMT
    Content-Length: 211
    
    {"id":2,"name":"Golang语言开发","description":"Golang语言开发从入门到放弃！","classify":"互联网","genderLimit":0,"memberLimit":{"maleCount":0,"femaleCount":0},"startTime":1607911200,"credit":1}  
    ```
  
* 获取课程列表

    ```shell script
    $ curl -i -X GET \
        'http://127.0.0.1:8889/api/course/list?page=1&size=10' \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4NjczMzgsImlhdCI6MTYwNzg2MzczOCwiaWQiOjJ9.NPcA430BHJ3L_V_JsJiOEt00GAPIb2PMhN0TOnbc5Xk' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2'
    ```
    ```text
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Sun, 13 Dec 2020 13:22:51 GMT
    Content-Length: 241
    
    {"total":1,"size":1,"list":[{"id":2,"name":"Golang语言开发","description":"Golang语言开发从入门到放弃！","classify":"互联网","genderLimit":0,"memberLimit":{"maleCount":0,"femaleCount":0},"startTime":1607911200,"credit":1}]}
    ```
* 删除课程
    
    ```shell script
    curl -i -X POST \
      http://127.0.0.1:8889/api/course/delete/2 \
      -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc4NjczMzgsImlhdCI6MTYwNzg2MzczOCwiaWQiOjJ9.NPcA430BHJ3L_V_JsJiOEt00GAPIb2PMhN0TOnbc5Xk' \
      -H 'content-type: application/json' \
      -H 'x-user-id: 2'
    ```
    ```text
    HTTP/1.1 200 OK
    Date: Sun, 13 Dec 2020 13:18:31 GMT
    Content-Length: 0  
    ```

> 说明：以上请求中的id以开发人员实际数据库为准。

# 本章节贡献者
 * [songmeizi](https://github.com/songmeizi)
 
 # 技术点总结
 * go-zero中间件使用
    * 全局中间件
    * 指定路由组中间件
 * go-zero自定义错误
 * go-zero rpc调用
 
 # 相关推荐
 * [zrpc](https://github.com/tal-tech/zero-doc/blob/main/doc/zrpc.md)
 * [使用goctl创建rpc](https://github.com/tal-tech/zero-doc/blob/main/doc/goctl-rpc.md)
 * [使用goctl创建model](https://github.com/tal-tech/zero-doc/blob/main/doc/goctl-model-sql.md)
 
 # 结尾
 本章节完。
 
 如发现任何错误请通过Issue发起问题修复申请。
 
你可能会浏览 
* [用户模块](../../../doc/requirement/user.md)
* [选课模块](../../../doc/requirement/selection.md)
* [排课模块](../../../doc/requirement/schedule.md)







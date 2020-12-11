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
        // 姓名限制，0-不限，1-男，2-女
        GenderLimit int `json:"genderLimit,options=0|1|2"`
        // 可选参数，如果不传则代表不限制人数
        MemberLimit *MemberLimit `json:"memberLimit,optional"`
        StartTime int64 `json:"startTime"`
        // 学分
        Credit int `json:"credit,range=(0,6]"`
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
        Page int `json:"page,range=(0:]"`
        Size int `json:"size,range=(0:]"`
        CursorId int64 `json:"cursorId,optional"`
    }

    CourseListReply {
        CurrentPage int `json:"currentPage"`
        Size int `json:"size"`
        CursorId int64 `json:"cursorId"`
        HasMore bool `json:"hasMore"`
        List []*Course `json:"list"`
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
  
* 打开`service/course/api/etc/user-api.yaml`文件，添加`Mysql`、`CacheRedis`配置项

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
      - Host: 127.0.0.1:6379
        Type: node
    ```

    >说明： 我本地redis没有设置密码，因此没有配置`Password`配置项，为了防止和user api端口冲突，这里将端口修改为`8889`

# ServiceContext增加`CourseModel`资源
进入文件`service/course/api/internal/svc/servicecontext.go`，添加`CourseModel`资源依赖。

```go
package svc

import (
	"hey-go-zero/service/course/api/internal/config"
	"hey-go-zero/service/course/model"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	CourseModel model.CourseModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:      c,
		CourseModel: model.NewCourseModel(conn, c.CacheRedis),
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

	return l.svcCtx.CourseModel.Delete(req.Id)
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
## 添加中间逻辑





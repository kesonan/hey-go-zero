# 选课模块
本文档为selection-api的创建流程步骤文档，通过本文档你可以了解到maintainer是如何完成对selection api的实现。

> 温馨提示：在前面课程删除上其实有一些逻辑问题，我们在删除课程时并没有考虑该课程是否已经被选课占用，大家有兴趣可以去自动补充。
> 解决方案：在`service/selection`中创建一个rpc服务，提供查询已经被占用的课程，然后进行逻辑补充。

# 创建api目录
在`service/selection`下创建api目录得到目录树

```text
selection
└── api
```

# 创建selection.api文件
在`service/selection/api`文件夹`右键`->`New Api File`->`输入course`->`选择Empty file`->`回车`，
然后修改selection.api文件内容为

```go
type (
	MemberLimit {
		// 男生限制人数 <=0：不限
		MaleCount int `json:"maleCount"`
		// 女生限制人数 <=0：不限
		FemaleCount int `json:"femaleCount"`
	}
	
	Course {
		Id int64 `json:"id"`
		Name string `json:"name"`
		Description string `json:"description"`
		Classify string `json:"classify"`
		// 性别限制，0-不限，1-男，2-女
		GenderLimit int `json:"genderLimit"`
		// 可选参数，如果不传则代表不限制人数
		MemberLimit MemberLimit `json:"memberLimit"`
		StartTime int64 `json:"startTime"`
		// 学分
		Credit int `json:"credit"`
		TeacherName string `json:"teacherName"`
	}
	
	SelectionCourse {
		CourseId int64 `json:"courseId"`
		TeacherId int64 `json:"teacherId"`
	}
	
	SelectionAddCourseReq {
		SelectionId int64 `path:"selectionId"`
		List []*SelectionCourse `json:"list"`
	}
	
	CreateSelectionReq {
		Name string `json:"name"`
		MaxCredit int `json:"maxCredit,range=(0:12]"`
		StartTime int64 `json:"startTime"`
		EndTime int64 `json:"endTime"`
		Notification string `json:"notification"`
	}
	
	EditSelectionReq {
		Id int64 `path:"id"`
		CreateSelectionReq
	}
	
	SelectionIdReq {
		Id int64 `path:"id"`
	}
	
	SelectionReply {
		Name string `json:"name"`
		MaxCredit int `json:"maxCredit"`
		StartTime int64 `json:"startTime"`
		EndTime int64 `json:"endTime"`
		Notification string `json:"notification"`
		CourseList []*Course `json:"courseList"`
	}
	
	SelectCourseId {
		Id int64 `path:"id"`
	}
	
	MineCourseReply {
		List []*Course `json:"list"`
	}
)

@server(
	jwt: Auth
	middleware: ManagerCheck // 仅管理员可访问
)
service selection-api {
	@doc "创建选课"
	@handler createSelection
	post /api/selection/create (CreateSelectionReq)
	
	@doc "编辑选课"
	@handler editSelection
	post /api/selection/edit (EditSelectionReq)
	
	@doc "添加课程"
	@handler addCourse
	post /api/selection/add/course/:id (SelectionAddCourseReq)
	
	@doc "移除课程"
	@handler deleteCourse
	post /api/selection/delete/course/:id (SelectionAddCourseReq)
	
	@doc "删除选课"
	@handler deleteSelection
	post /api/selection/delete/:id (SelectionIdReq)
}

@server(
	jwt: Auth
	middleware: StudentCheck // 仅学生可访问
)
service selection-api {
	@doc "查看选课"
	@handler getSelection
	get /api/selection/info/:id returns (SelectionIdReq)
	
	@doc "选课"
	@handler select
	post /api/selection/select/:id (SelectCourseId)
	
	@doc "查看我的选课"
	@handler mineSelections
	get /api/selection/mine/list returns (MineCourseReply)
}

@server(
	jwt: Auth
	middleware: TeacherCheck // 仅教师可访问
)
service selection-api {
	@doc "查看我的任教课程"
	@handler getTeachingCourses
	get /api/selection/teaching/courses returns (MineCourseReply)
	
	@doc "查看我任教课程的学生列表"
	@handler getTeachingStudents
	get /api/selection/teaching/students/:id (SelectCourseId)
}
```

# 生成代码
在Goland中生成代码有三种方式（任意一种均可）
* project面板区文件右键生成
    * 选中`selection.api`文件->`右键`->`New`->`Go Zero`->`Api Code`
    * `Api Code`回车后会弹出一个文件对话框询问你需要生成服务的目标目录，默认为`selection.api`所在目录，我们这里选择默认，点击`OK`确认生成。
* api文件编辑区右键生成
    * 打开`course.api`文件->`编辑区`->`右键`->`Generate..`->`Api Code`
    * `Api Code`回车后会弹出一个文件对话框询问你需要生成服务的目标目录，默认为`selection.api`所在目录，我们这里选择默认，点击`OK`确认生成。
* 终端生成（推荐）
    * 选中`selection.api`文件->`右键`->`Open in Terminal`
    * 执行`goctl api go -api selection.api -dir .`命令即可
    
        ```shell script
        $ goctl api go -api selection.api -dir .
        ```
        ```text
        Done.
        ```
接下来我们进入`service/selection/api`目录，查看一下目录树结构

```shell script
$ tree
```

```text
selection/api
├── etc
│   └── selection-api.yaml
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── addcoursehandler.go
│   │   ├── createselectionhandler.go
│   │   ├── deletecoursehandler.go
│   │   ├── deleteselectionhandler.go
│   │   ├── editselectionhandler.go
│   │   ├── getselectionhandler.go
│   │   ├── getteachingcourseshandler.go
│   │   ├── getteachingstudentshandler.go
│   │   ├── mineselectionshandler.go
│   │   ├── routes.go
│   │   └── selecthandler.go
│   ├── logic
│   │   ├── addcourselogic.go
│   │   ├── createselectionlogic.go
│   │   ├── deletecourselogic.go
│   │   ├── deleteselectionlogic.go
│   │   ├── editselectionlogic.go
│   │   ├── getselectionlogic.go
│   │   ├── getteachingcourseslogic.go
│   │   ├── getteachingstudentslogic.go
│   │   ├── mineselectionslogic.go
│   │   └── selectlogic.go
│   ├── middleware
│   │   └── managercheckmiddleware.go
│   ├── svc
│   │   └── servicecontext.go
│   └── types
│       └── types.go
├── readme.md
├── selection.api
└── selection.go

```

> 说明：上述目录中的注释是为了大家能够快速知道该目录结构的用途，是后期我加入的，实际生成的tree不会带注释和readme.md文件。

# 生成model层代码

### 新建表selection

在前面我们已经演示了如何[创建db和table](../../../doc/prepare/db-create.md) ，我们在`heygozero`新建一张表名为`selection`、`selection_scourse`、`selection_student`,create table ddl如下：

```mysql
-- 选课表 --
CREATE TABLE `selection` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '选课名称',
  `max_credit` tinyint NOT NULL DEFAULT '0' COMMENT '最大可修学分',
  `start_time` bigint NOT NULL DEFAULT '0' COMMENT '选课开始时间',
  `end_time` bigint NOT NULL DEFAULT '0' COMMENT '选课结束时间',
  `notification` varchar(500) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '选课通知内容，500字以内',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 选课课程表 --
CREATE TABLE `selection_course` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `selection_id` bigint NOT NULL DEFAULT '0' COMMENT '选课任务id',
  `course_id` bigint NOT NULL DEFAULT '0' COMMENT '课程id',
  `teacher_id` bigint NOT NULL DEFAULT '0' COMMENT '任教教师',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_selection_course` (`selection_id`,`course_id`),
  UNIQUE KEY `unique_course_teacher` (`course_id`,`teacher_id`),
  KEY `idx_teacher` (`teacher_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 选课信息表 --
CREATE TABLE `selection_student` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `selection_course_id` bigint NOT NULL DEFAULT '0' COMMENT '选课任务课程id',
  `student_id` bigint NOT NULL DEFAULT '0' COMMENT '学生id',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_student_selection_course` (`student_id`,`selection_course_id`),
  KEY `idx_selection_course` (`selection_course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

# 生成代码
在文件夹`service/selection`右键->`Open in Terminal` 进入终端
```shell script
$ goctl model mysql datasource -url="ugozero@tcp(127.0.0.1:3306)/heygozero" -table="selection*" -c -dir ./model
```
```text
Done.
```

生成后我么可以看到会产生4个文件
```text
model
├── selectionmodel.go
├── selectionscoursemodel.go
├── selectionstudentmodel.go
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
  
* 打开`service/selection/api/etc/selection-api.yaml`文件，添加`Mysql`、`CacheRedis`配置项

    ```yaml
    Name: selection-api
    Host: 0.0.0.0
    Port: 8890
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

    >说明： 我本地redis没有设置密码，因此没有配置`Password`配置项，为了防止和user api端口冲突，这里将端口修改为`8889`

# 删除`service/selection/api/internal/middleware`目录下的文件，添加`usercheckmiddleware.go`

```go
package middleware

import (
	"net/http"

	"hey-go-zero/common/errorx"
	"hey-go-zero/common/jwtx"
	"hey-go-zero/service/user/rpc/userservice"

	"github.com/tal-tech/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserCheckMiddleware struct {
	userRpcClient userservice.UserService
	role          string
}

func NewManagerCheckMiddleware(role string, userRpcClient userservice.UserService) *UserCheckMiddleware {
	return &UserCheckMiddleware{
		role:          role,
		userRpcClient: userRpcClient,
	}
}

func (m *UserCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if data.Role != m.role {
			httpx.Error(w, errorx.NewDescriptionError("无权限访问"))
			return
		}

		next(w, r)
	}
}
```

> 说明：这里由于goctl中间件生成的文件创建有bug，因此需要手动做了中间文件修改，后续goctl修复后可针对每个中间件实现自己的逻辑即可。

# ServiceContext增加`SelectionModel`、`SelectionCourseModel`、`SelectionStudentModel`资源
进入文件`service/selection/api/internal/svc/servicecontext.go`，添加`SelectionModel`、`SelectionCourseModel`、`SelectionStudentModel`资源依赖。

```go
type ServiceContext struct {
	Config                config.Config
	ManagerCheck          rest.Middleware
	StudentCheck          rest.Middleware
	TeacherCheck          rest.Middleware
	SelectionModel        model.SelectionModel
	SelectionCourseModel  model.SelectionCourseModel
	SelectionStudentModel model.SelectionStudentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	userRpcClient := zrpc.MustNewClient(c.UserRpc)
	userRpcService := userservice.NewUserService(userRpcClient)
	return &ServiceContext{
		Config:                c,
		ManagerCheck:          middleware.NewManagerCheckMiddleware("manager", userRpcService).Handle,
		StudentCheck:          middleware.NewManagerCheckMiddleware("student", userRpcService).Handle,
		TeacherCheck:          middleware.NewManagerCheckMiddleware("teacher", userRpcService).Handle,
		SelectionModel:        model.NewSelectionModel(conn, c.CacheRedis),
		SelectionCourseModel:  model.NewSelectionCourseModel(conn, c.CacheRedis),
		SelectionStudentModel: model.NewSelectionStudentModel(conn, c.CacheRedis),
	}
}
```

# 添加自定义错误码
由于在user api中已经添加过[自定义错误码](../../../doc/gozero/http-error.md)了，这里就直接在main文件中使用即可。

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

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
```











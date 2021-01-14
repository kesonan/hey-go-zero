# 选课模块
本文档为selection-api的创建流程步骤文档，通过本文档你可以了解到maintainer是如何完成对selection api的实现。

> 温馨提示：在前面课程删除上其实有一些逻辑问题，我们在删除课程时并没有考虑该课程是否已经被选课占用，大家有兴趣可以去自动补充。
> 解决方案：在`service/selection`中创建一个rpc服务，提供查询已经被占用的课程，然后进行逻辑补充。

# 创建api目录
在`service/selection`下创建api目录得到目录树

``` text
selection
└── api
```

# 创建selection.api文件
在`service/selection/api`文件夹`右键`->`New Api File`->`输入course`->`选择Empty file`->`回车`，
然后修改selection.api文件内容为

``` go
type (
	Course {
		Id int64 `json:"id"`
		SelectionCourseId int64 `json:"selectionCourseId"`
		Name string `json:"name"`
		Description string `json:"description"`
		Classify string `json:"classify"`
		// 性别限制，0-不限，1-男，2-女
		GenderLimit int `json:"genderLimit"`
		// 可选参数，如果不传则代表不限制人数
		MemberLimit int `json:"memberLimit"`
		// 学分
		Credit int `json:"credit"`
		TeacherName string `json:"teacherName"`
	}
	
	SelectionCourse {
		CourseId int64 `json:"courseId"`
		TeacherId int64 `json:"teacherId"`
	}
	
	SelectionCourseReq {
		SelectionId int64 `path:"selectionId"`
		List []*SelectionCourse `json:"list"`
	}
	
	DeleteSelectionCourseReq {
		SelectionId int64 `json:"selectionId"`
		Ids []int64 `json:"ids"`
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
		Id int64 `json:"id"`
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
	post /api/selection/add/course/:selectionId (SelectionCourseReq)
	
	@doc "移除课程"
	@handler deleteCourse
	post /api/selection/delete/course (DeleteSelectionCourseReq)
	
	@doc "删除选课"
	@handler deleteSelection
	post /api/selection/delete/:id (SelectionIdReq)
}

@server(
	jwt: Auth
)
service selection-api {
	@doc "查看选课"
	@handler getSelection
	get /api/selection/info/:id (SelectionIdReq) returns (SelectionReply)
}

@server(
	jwt: Auth
	middleware: StudentCheck // 仅学生可访问
)
service selection-api {
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
    
        ``` shell script
        $ goctl api go -api selection.api -dir .
        ```
        ``` text
        Done.
        ```
接下来我们进入`service/selection/api`目录，查看一下目录树结构

``` shell script
$ tree
```

``` text
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

``` mysql
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
``` shell script
$ goctl model mysql datasource -url="ugozero@tcp(127.0.0.1:3306)/heygozero" -table="selection*" -c -dir ./model
```
``` text
Done.
```

生成后我么可以看到会产生4个文件
``` text
model
├── selectionmodel.go
├── selectionscoursemodel.go
├── selectionstudentmodel.go
└── vars.go
```

> 说明：本地生成的`goctl`版本为`goctl version 20201125 darwin/amd64`，早起版本生成出来的数值类型会有`int`，`int64`，而后续版本统一为`int64`了。

# 创建course.rpc服务
由于选课需要查询用户信息、课程信息，因此我们需要用到rpc来进行服务间的通信，在course-api我们已经创建过[user.rpc](../../user/rpc/readme.md)服务了，在这里再去创建一个course.rpc服务，
[点击这里](../../course/rpc/readme.md)进入[course.rpc](../../course/rpc/readme.md)服务创建流程。

# 添加`Mysql`和`CacheRedis`配置定义和yaml配置项
* 编打开`service/course/api/internal/config/config.go`，添加`Mysql`、`CacheRedis`、`BizRedis`、UserRpc`、`CourseRpc`、`Dq`定义

    ``` go
    package config
    
    import (
    	"github.com/tal-tech/go-queue/dq"
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
    	BizRedis   cache.NodeConf
    	UserRpc    zrpc.RpcClientConf
    	CourseRpc  zrpc.RpcClientConf
    	Dq         dq.DqConf
    }
    ```
  
* 打开`service/selection/api/etc/selection-api.yaml`文件，添加`Mysql`、`CacheRedis`配置项

    ``` yaml
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
    BizRedis:
      Host: 127.0.0.1:6379
      Type: node
      Weight: 100
    UserRpc:
      Etcd:
        Hosts:
          - 127.0.0.1:2379
        Key: user.rpc
    CourseRpc:
      Etcd:
        Hosts:
          - 127.0.0.1:2379
        Key: course.rpc
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

    >说明： 我本地redis没有设置密码，因此没有配置`Password`配置项，为了防止和user api端口冲突，这里将端口修改为`8889`



# 中间件管理

删除`service/selection/api/internal/middleware`目录下的文件，添加`usercheckmiddleware.go`

``` go
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


# ServiceContext增加`SelectionModel`、`SelectionCourseModel`等资源
进入文件`service/selection/api/internal/svc/servicecontext.go`，添加`SelectionModel`、`SelectionCourseModel`等资源依赖。

``` go
type ServiceContext struct {
	Config                config.Config
	ManagerCheck          rest.Middleware
	StudentCheck          rest.Middleware
	TeacherCheck          rest.Middleware
	SelectionModel        model.SelectionModel
	SelectionCourseModel  model.SelectionCourseModel
	SelectionStudentModel model.SelectionStudentModel
	UserService           userservice.UserService
	CourseService         courseservice.CourseService
	BizRedis              *redis.Redis
	Producer              dq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	userRpcClient := zrpc.MustNewClient(c.UserRpc)
	courseRpcClient := zrpc.MustNewClient(c.CourseRpc)
	userRpcService := userservice.NewUserService(userRpcClient)
	courseService := courseservice.NewCourseService(courseRpcClient)
	bizRedis := redis.NewRedis(c.BizRedis.Host, c.BizRedis.Type, c.BizRedis.Pass)
	return &ServiceContext{
		Config:                c,
		ManagerCheck:          middleware.NewManagerCheckMiddleware("manager", userRpcService).Handle,
		StudentCheck:          middleware.NewManagerCheckMiddleware("student", userRpcService).Handle,
		TeacherCheck:          middleware.NewManagerCheckMiddleware("teacher", userRpcService).Handle,
		SelectionModel:        model.NewSelectionModel(conn, c.CacheRedis),
		SelectionCourseModel:  model.NewSelectionCourseModel(conn, c.CacheRedis),
		SelectionStudentModel: model.NewSelectionStudentModel(conn, c.CacheRedis),
		UserService:           userRpcService,
		CourseService:         courseService,
		BizRedis:              bizRedis,
		Producer:              dq.NewProducer(c.Dq.Beanstalks),
	}
}
```

# 添加自定义错误码
由于在user api中已经添加过[自定义错误码](../../../doc/gozero/http-error.md)了，这里就直接在main文件中使用即可。

``` go
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

# 填充逻辑

### 增加error.go文件
在`service/selection/api/internal/logic`目录下创建`error.go`文件，填充代码
``` go
var (
	errCourseNotFound          = errorx.NewDescriptionError("课程不存在")
	errSelectionNotFound       = errorx.NewDescriptionError("选课任务不存在")
	errSelectionIsExists       = errorx.NewDescriptionError("选课任务已存在")
	errNotInSelectionTimeRange = errorx.NewDescriptionError("选课未开始")
	errSelectionExpired        = errorx.NewDescriptionError("选课已结束")
	errSelectionCourseNotFound = errorx.NewDescriptionError("不存在该选课课程")
	errSelectionSelected       = errorx.NewDescriptionError("你已选择过该课程，请勿重复选择")
	errSelectionNotSelected    = errorx.NewDescriptionError("你尚未选择该课程，无法取消选课")
)
```

### 添加`common.go`文件
在`service/selection/api/internal/logic`目录下添加`common.go`文件，用于存放选课创建、编辑等逻辑的公用逻辑，填充代码

``` go
func checkCourseSelection(in types.CreateSelectionReq) error {
	if len(strings.TrimSpace(in.Name)) == 0 {
		return errorx.NewInvalidParameterError("name")
	}

	if utf8.RuneCountInString(in.Name) > 20 {
		return lengthAlert("课程名称", 20)
	}

	now := time.Now()
	startTime := now.Add(2 * time.Hour)
	endTime := startTime.AddDate(0, 0, 5)
	if in.StartTime < startTime.Unix() {
		return errorx.NewDescriptionError(fmt.Sprintf("选课开始时间不能早于%s", startTime.Format("2006年01月02日 03时04分05秒")))
	}

	if in.EndTime < endTime.Unix() {
		return errorx.NewDescriptionError(fmt.Sprintf("选课结束时间不能晚于%s", startTime.Format("2006年01月02日 03时04分05秒")))
	}

	if utf8.RuneCountInString(in.Notification) > 500 {
		return lengthAlert("选课通知", 500)
	}

	return nil
}

func lengthAlert(hint string, length int) error {
	return errorx.NewDescriptionError(fmt.Sprintf("%s不能超过%d个字符", hint, length))
}
```

### 创建选课
创建选课我们用到了dq来发送选课通知，因此我们需要一个dq的consumer服务来消费该消息，具体消费服务请移步至[selection-rmq](../rmq/readme.md)

* 文件位置: `service/selection/api/internal/logic/createselectionlogic.go`
* 方法名: `CreateSelection`
* 代码内容:

    ``` go
    func (l *CreateSelectionLogic) CreateSelection(req types.CreateSelectionReq) error {
    	if err := checkCourseSelection(req); err != nil {
    		return err
    	}
    
    	_, err := l.svcCtx.SelectionModel.FindOneByName(req.Name)
    	switch err {
    	case nil:
    		return errSelectionIsExists
    	case model.ErrNotFound:
    		_, err := l.svcCtx.SelectionModel.Insert(model.Selection{
    			MaxCredit:    int64(req.MaxCredit),
    			StartTime:    req.StartTime,
    			EndTime:      req.EndTime,
    			Notification: req.Notification,
    			Name:         req.Name,
    		})
    		if err != nil {
    			return err
    		}
    
    		// dq，todo：这里建议用cron-job替代，如果用dq对于这种需要变更时间的逻辑，将导致发送了多个不同时间点的message，本案例仅用于演示dq这么使用。
    		_, err = l.svcCtx.Producer.At([]byte(req.Notification), time.Unix(req.StartTime, 0).Add(-2*time.Hour))
    
    		return err
    	default:
    		return err
    	}
    }
    ```

  
### 编辑选课任务

* 文件位置: `service/selection/api/internal/logic/editselectionlogic.go`
* 方法名: `EditSelection`
* 代码内容:

    ``` go
    func (l *EditSelectionLogic) EditSelection(req types.EditSelectionReq) error {
    	if err := checkCourseSelection(req.CreateSelectionReq); err != nil {
    		return err
    	}
    	data, err := l.svcCtx.SelectionModel.FindOne(req.Id)
    	if err != nil {
    		if err == model.ErrNotFound {
    			return errSelectionNotFound
    		}
    		return err
    	}
    
    	nameData, err := l.svcCtx.SelectionModel.FindOneByName(req.Name)
    	if err != nil {
    		if err == model.ErrNotFound {
    			return err
    		}
    	} else {
    		if nameData.Id != req.Id {
    			return errSelectionIsExists
    		}
    	}
    
    	data.Name = req.Name
    	data.MaxCredit = int64(req.MaxCredit)
    	data.StartTime = req.StartTime
    	data.EndTime = req.EndTime
    	data.Notification = req.Notification
    	err = l.svcCtx.SelectionModel.Update(*data)
    	if err != nil {
    		return err
    	}
    
    	// dq，todo：这里建议用cron-job替代，如果用dq对于这种需要变更时间的逻辑，将导致发送了多个不同时间点的message，本案例仅用于演示dq这么使用。
    	_, err = l.svcCtx.Producer.At([]byte(req.Notification), time.Unix(req.StartTime, 0).Add(-2*time.Hour))
    	return err
    }
    ```
  
### 删除选课任务

* 文件位置: `service/selection/api/internal/logic/deleteselectionlogic.go`
* 方法名: `DeleteSelection`
* 代码内容:

    ``` go
    func (l *DeleteSelectionLogic) DeleteSelection(req types.SelectionIdReq) error {
    	err := l.svcCtx.SelectionModel.Delete(req.Id)
    	switch err {
    	case nil:
    		return nil
    	case model.ErrNotFound:
    		return errSelectionNotFound
    	default:
    		return err
    	}
    }
    ```

### 修改`selectioncoursemodel.go`  
编辑`service/selection/model/selectioncoursemodel.go`添加`FindBySelectionId`、`FindByTeacherId`、`DeleteBySelectionId`方法

``` go
FindBySelectionId(selectionId int64) ([]*SelectionCourse, error)
FindByTeacherId(teacherId int64) ([]*SelectionCourse, error)
DeleteBySelectionId(selectionId int64) error
```
``` go
func (m *defaultSelectionCourseModel) FindBySelectionId(selectionId int64) ([]*SelectionCourse, error) {
	query := fmt.Sprintf("select %s from %s where selection_id = ?", selectionCourseRows, m.table)
	var resp []*SelectionCourse
	err := m.QueryRowsNoCache(&resp, query, selectionId)
	return resp, err
}

func (m *defaultSelectionCourseModel) FindByTeacherId(teacherId int64) ([]*SelectionCourse, error) {
	query := fmt.Sprintf("select %s from %s where teacher_id = ?", selectionCourseRows, m.table)
	var resp []*SelectionCourse
	err := m.QueryRowsNoCache(&resp, query, teacherId)
	return resp, err
}

func (m *defaultSelectionCourseModel) DeleteBySelectionId(selectionId int64) error {
	list, err := m.FindBySelectionId(selectionId)
	if err != nil {
		return err
	}

	keys := collection.NewSet()
	for _, item := range list {
		keys.AddStr(m.formatPrimary(item.Id))
	}
	query := fmt.Sprintf("delete from %s where selection_id = ?",m.table)
	_, err = m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		return conn.Exec(query, selectionId)
	}, keys.KeysStr()...)

	return err
}
```


### 查看选课任务

* 文件位置: `service/selection/api/internal/logic/getselectionlogic.go`
* 方法名: `GetSelection`
* 代码内容:

    ``` go
    func (l *GetSelectionLogic) GetSelection(req types.SelectionIdReq) (*types.SelectionReply, error) {
    	list, err := l.svcCtx.SelectionCourseModel.FindBySelectionId(req.Id)
    	if err != nil {
    		return nil, err
    	}
    
    	var courseList []*types.Course
    	fx.From(func(source chan<- interface{}) {
    		for _, item := range list {
    			source <- item
    		}
    	}).Walk(func(item interface{}, pipe chan<- interface{}) {
    		data := item.(*model.SelectionCourse)
    		courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{
    			Id: data.CourseId,
    		})
    		if err != nil {
    			logx.Error(err)
    			return
    		}
    
    		var teacherName string
    		userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &user.UserReq{Id: data.TeacherId})
    		if err != nil {
    			logx.Error(err)
    		} else {
    			teacherName = userInfo.Name
    		}
    
    		courseList = append(courseList, &types.Course{
    			Id:                courseInfo.Id,
    			SelectionCourseId: data.Id,
    			Name:              courseInfo.Name,
    			Description:       courseInfo.Description,
    			Classify:          courseInfo.Classify,
    			GenderLimit:       int(courseInfo.GenderLimit),
    			MemberLimit:       int(courseInfo.MemberLimit),
    			Credit:            int(courseInfo.Credit),
    			TeacherName:       teacherName,
    		})
    	})
    
    	// sort by id desc
    	sort.Slice(courseList, func(i, j int) bool {
    		return courseList[i].Id > courseList[j].Id
    	})
    
    	data, err := l.svcCtx.SelectionModel.FindOne(req.Id)
    	switch err {
    	case nil:
    		return &types.SelectionReply{
    			Id:           data.Id,
    			Name:         data.Name,
    			MaxCredit:    int(data.MaxCredit),
    			StartTime:    data.StartTime,
    			EndTime:      data.EndTime,
    			Notification: data.Notification,
    			CourseList:   nil,
    		}, nil
    	case model.ErrNotFound:
    		return nil, errSelectionNotFound
    	default:
    		return nil, err
    	}
    }
    ```

### 在selectionstudentmodel.go添加`FindBySelectionCourseId`方法

* 文件位置：`service/selection/model/selectionstudentmodel.go`
* 方法名：FindBySelectionCourseId
* interface代码内容：

    ``` go
    FindBySelectionCourseId(selectionCourseId int64) ([]*SelectionStudent, error)
    ```
* default实现代码内容：

    ``` go
    func (m *defaultSelectionStudentModel) FindBySelectionCourseId(selectionCourseId int64) ([]*SelectionStudent, error) {
    	query := fmt.Sprintf("select %s from %s where selection_course_id = ?", selectionStudentRows, m.table)
    	var resp []*SelectionStudent
    	err := m.QueryRowNoCache(&resp, query,selectionCourseId)
    	return resp, err
    }
    ```
  
### 删除选课任务
* 文件位置: `service/selection/api/internal/logic/deleteselectionlogic.go`
* 方法名: `DeleteSelection`
* 代码内容:

    ``` go
    func (l *DeleteSelectionLogic) DeleteSelection(req types.SelectionIdReq) error {
    	err := l.svcCtx.SelectionModel.Delete(req.Id)
    	switch err {
    	case nil:
    		selectCourseList, err := l.svcCtx.SelectionCourseModel.FindBySelectionId(req.Id)
    		if err != nil {
    			return err
    		}
    
    		fx.From(func(source chan<- interface{}) {
    			for _, item := range selectCourseList {
    				source <- item
    			}
    		}).Walk(func(item interface{}, pipe chan<- interface{}) {
    			data := item.(*model.SelectionCourse)
    			list, err := l.svcCtx.SelectionStudentModel.FindBySelectionCourseId(data.Id)
    			if err != nil {
    				return
    			}
    
    			for _, each := range list {
    				_ = l.svcCtx.SelectionStudentModel.Delete(each.Id)
    			}
    		}).Done()
    
    		_ = l.svcCtx.SelectionCourseModel.DeleteBySelectionId(req.Id)
    		
    		return nil
    	case model.ErrNotFound:
    		return errSelectionNotFound
    	default:
    		return err
    	}
    }
    ```

### 添加课程

* 文件位置: `service/selection/api/internal/logic/addcourselogic.go`
* 方法名: `AddCourse`
* 代码内容:

    ``` go
    func (l *AddCourseLogic) AddCourse(req types.SelectionAddCourseReq) error {
    	if req.SelectionId <= 0 {
    		return errorx.NewInvalidParameterError("selectionId")
    	}
    
    	if len(req.List) == 0 {
    		return errorx.NewInvalidParameterError("list")
    	}
    
    	for _, item := range req.List {
    		if item.TeacherId <= 0 {
    			return errorx.NewInvalidParameterError("teacherId")
    		}
    
    		if item.CourseId <= 0 {
    			return errorx.NewInvalidParameterError("courseId")
    		}
    	}
    
    	_, err := l.svcCtx.SelectionModel.FindOne(req.SelectionId)
    	if err != nil {
    		if err == model.ErrNotFound {
    			return errSelectionNotFound
    		}
    		return err
    	}
    
    	selectionCourseList, err := l.svcCtx.SelectionCourseModel.FindBySelectionId(req.SelectionId)
    	if err != nil {
    		return err
    	}
    
    	selectionCourseM := make(map[int64]struct{})
    	for _, item := range selectionCourseList {
    		selectionCourseM[item.CourseId] = struct{}{}
    	}
    
    	err = mr.MapReduceVoid(func(source chan<- interface{}) {
    		for _, item := range req.List {
    			source <- item
    		}
    	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
    		data := item.(*types.SelectionCourse)
    		if _, ok := selectionCourseM[data.CourseId]; ok {
    			cancel(errorx.NewDescriptionError("已经添加过该课程，请勿重复添加"))
    			return
    
    		}
    		_, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{
    			Id: data.CourseId,
    		})
    		if err != nil {
    			st := status.Convert(err)
    			if st.Code() == codes.NotFound {
    				cancel(errCourseNotFound)
    				return
    			}
    			cancel(err)
    			return
    		}
    
    		_, err = l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{
    			Id: data.TeacherId,
    		})
    		if err != nil {
    			st := status.Convert(err)
    			if st.Code() == codes.NotFound {
    				cancel(errCourseNotFound)
    				return
    			}
    			cancel(err)
    			return
    		}
    
    		writer.Write(struct{}{})
    	}, func(pipe <-chan interface{}, cancel func(error)) {
    		for range pipe {
    		}
    	})
    
    	if err != nil {
    		return err
    	}
    
    	fx.From(func(source chan<- interface{}) {
    		for _, item := range req.List {
    			source <- item
    		}
    	}).Walk(func(item interface{}, pipe chan<- interface{}) {
    		data := item.(*types.SelectionCourse)
    		// 不考虑事务
    		_, err := l.svcCtx.SelectionCourseModel.Insert(model.SelectionCourse{
    			SelectionId: req.SelectionId,
    			CourseId:    data.CourseId,
    			TeacherId:   data.TeacherId,
    		})
    		logx.Error(err)
    	}).Done()
    
    	return nil
    }
    ```

### 删除课程
* 文件位置: `service/selection/api/internal/logic/deletecourselogic.go`
* 方法名: `DeleteCourse`
* 代码内容:

    ``` go
    func (l *DeleteCourseLogic) DeleteCourse(req types.DeleteSelectionCourseReq) error {
    	selection, err := l.svcCtx.SelectionModel.FindOne(req.SelectionId)
    	switch err {
    	case nil:
    		if time.Now().After(time.Unix(selection.StartTime, 0)) {
    			return errorx.NewDescriptionError("该选课已发布，不能编辑课程")
    		}
    
    		err = l.svcCtx.SelectionModel.Delete(req.SelectionId)
    		if err != nil {
    			logx.Error(err)
    		}
    		for _, each := range req.Ids {
    			_ = l.svcCtx.SelectionCourseModel.Delete(each)
    		}
    		return nil
    	case model.ErrNotFound:
    		return errSelectionNotFound
    	default:
    		return err
    	}
    }
    ```

### selectionstudentmodel.go添加`FindByStudentId`方法
* 文件位置：`service/selection/model/selectionstudentmodel.go`
* 方法名称：添加`FindByStudentId`和`FindByStudentIdAndSelectionCourseId`
* interface中代码：
    ``` go
    FindByStudentId(studentId int64) ([]*SelectionStudent, error)
    FindByStudentIdAndSelectionCourseId(studentId, selectionCourseId int64) (*SelectionStudent, error)
    ```
* default实现中代码
    ``` go
    func (m *defaultSelectionStudentModel) FindByStudentId(studentId int64) ([]*SelectionStudent, error) {
        query := fmt.Sprintf("select %s from %s where student_id = ?", selectionStudentRows, m.table)
        var resp []*SelectionStudent
        err := m.QueryRowsNoCache(&resp, query,studentId)
        return resp, err
    }
    
    func (m *defaultSelectionStudentModel)FindByStudentIdAndSelectionCourseId(studentId, selectionCourseId int64) (*SelectionStudent, error){
        query := fmt.Sprintf("select %s from %s where student_id = ? and selection_course_id = ? limit 1", selectionStudentRows, m.table)
        var resp SelectionStudent
        err := m.QueryRowNoCache(&resp, query,studentId,selectionCourseId)
        return &resp, err
    }
    ```

### 修改selecthandler.go
* 文件位置： `service/selection/api/internal/handler/selecthandler.go`
* 方法名：`selectHandler`
代码内容：
    ``` go
    func selectHandler(ctx *svc.ServiceContext) http.HandlerFunc {
    	return func(w http.ResponseWriter, r *http.Request) {
    		userId, ok := jwtx.GetUserId(w, r) // add
    		if !ok {  // add
    			return  // add
    		}  // add
    
    		var req types.SelectCourseId
    		if err := httpx.Parse(r, &req); err != nil {
    			httpx.Error(w, err)
    			return
    		}
    
    		l := logic.NewSelectLogic(r.Context(), ctx)
    		err := l.Select(userId, req)  // add userId
    		if err != nil {
    			httpx.Error(w, err)
    		} else {
    			httpx.Ok(w)
    		}
    	}
    }
    ```

### 开始选课
* 文件位置: `service/selection/api/internal/logic/selectlogic.go`
* 方法名: `Select`
* 代码内容:

    ``` go
    func (l *SelectLogic) Select(userId int64, req types.SelectCourseId) error {
    	if req.Id <= 0 {
    		return errorx.NewInvalidParameterError("id")
    	}
    
    	key := fmt.Sprintf("%v", req.Id)
    	lock := redis.NewRedisLock(l.svcCtx.BizRedis, key)
    	lock.SetExpire(3)
    	ok, err := lock.Acquire()
    	if err != nil {
    		return err
    	}
    
    	// todo: 这里和中间件中有重复逻辑，都查询了一次user，开发人员可以自己优化一下。
    	userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{Id: req.Id})
    	if err != nil {
    		st := status.Convert(err)
    		if st.Code() == codes.NotFound {
    			return errorx.NewDescriptionError("用户不存在")
    		}
    
    		return errorx.NewDescriptionError(st.Message())
    	}
    
    	now := time.Now()
    	selectionCourse, err := l.svcCtx.SelectionCourseModel.FindOne(req.Id)
    	switch err {
    	case nil:
    	case model.ErrNotFound:
    		return errSelectionCourseNotFound
    	default:
    		return err
    	}
    
    	selection, err := l.svcCtx.SelectionModel.FindOne(selectionCourse.SelectionId)
    	switch err {
    	case nil:
    		if now.Before(time.Unix(selection.StartTime, 0)) {
    			return errNotInSelectionTimeRange
    		}
    
    		if now.After(time.Unix(selection.EndTime, 0)) {
    			return errSelectionExpired
    		}
    
    	case model.ErrNotFound:
    		return errSelectionNotFound
    	default:
    		return err
    	}
    
    	courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{Id: selectionCourse.CourseId})
    	if err != nil {
    		st := status.Convert(err)
    		if st.Code() == codes.NotFound {
    			return errCourseNotFound
    		}
    
    		return errorx.NewDescriptionError(st.Message())
    	}
    
    	if courseInfo.GenderLimit != course.GenderLimit_NoLimit && courseInfo.GenderLimit != course.GenderLimit(userInfo.Gender) {
    		return errorx.NewDescriptionError("性别不符合")
    	}
    
    	if !ok {
    		return errorx.NewDescriptionError("当前选课人数较多，请稍后再试")
    	}
    	defer func() {
    		_, err = lock.Release()
    		logx.Error(err)
    	}()
    
    	endTime := time.Unix(selection.EndTime, 0)
    	ok, err = l.trySelect(req.Id, userId, selection.MaxCredit, courseInfo, endTime)
    	if err != nil {
    		logx.Error(err)
    		return errorx.NewDescriptionError("选课失败，请稍后再试")
    	}
    
    	if !ok {
    		return errorx.NewDescriptionError("选课人数已满，请选择其他课程")
    	}
    	threading.GoSafe(func() {
    		_, err = l.svcCtx.SelectionStudentModel.Insert(model.SelectionStudent{
    			SelectionCourseId: req.Id,
    			StudentId:         userId,
    		})
    		if err != nil {
    			logx.Error(err)
    		}
    	})
    
    	return nil
    }
    
    func (l *SelectLogic) trySelect(selectCourseId, userId, maxCredit int64, courseInfo *courseservice.Course, expireAt time.Time) (bool, error) {
    	expire := int(expireAt.Sub(time.Now()).Seconds()) + 1
    	userKey := fmt.Sprintf("biz#user#selected#status#%v#%v", userId, selectCourseId)
    	ok, err := l.svcCtx.BizRedis.SetnxEx(userKey, "*", expire)
    	if err != nil {
    		return false, err
    	}
    
    	if !ok {
    		return false, errSelectionSelected
    	}
    
    	userCreditKey := fmt.Sprintf("biz#user#selected#credit#%v", userId)
    	credit, err := l.svcCtx.BizRedis.Incrby(userCreditKey, courseInfo.Credit)
    	if err != nil {
    		_, _ = l.svcCtx.BizRedis.Del(userKey)
    		return false, err
    	}
    
    	if credit > maxCredit {
    		_, _ = l.svcCtx.BizRedis.Del(userKey)
    		return false, errorx.NewDescriptionError(fmt.Sprintf("选择当前课程后，该学期你的学分已经超出总学分%d，请合理选择课程", maxCredit))
    	}
    	_ = l.svcCtx.BizRedis.Expire(userCreditKey, expire)
    
    	courseKey := fmt.Sprintf("biz#course#selected#data#%v", selectCourseId)
    	var field string
    	switch courseInfo.GenderLimit {
    	case course.GenderLimit_Male:
    		field = "male"
    	case course.GenderLimit_Female:
    		field = "female"
    	default:
    		field = "all"
    	}
    
    	count, err := l.svcCtx.BizRedis.Hincrby(courseKey, field, 1)
    	if err != nil {
    		_, _ = l.svcCtx.BizRedis.Del(userKey)
    		return false, err
    	}
    
    	if count > int(courseInfo.MemberLimit) {
    		_, _ = l.svcCtx.BizRedis.Hincrby(courseKey, field, -1)
    		_, _ = l.svcCtx.BizRedis.Del(userKey)
    		return false, nil
    	}
    
    	_ = l.svcCtx.BizRedis.Expire(courseKey, expire)
    	return true, nil
    }
    ```

### 我的选课列表
* 文件位置: `service/selection/api/internal/logic/mineselectionslogic.go`
* 方法：`MineSelections`
* 代码内容：

    ``` go
    func (l *MineSelectionsLogic) MineSelections(userId int64) (*types.MineCourseReply, error) {
    	studentSelectedCourseList, err := l.svcCtx.SelectionStudentModel.FindByStudentId(userId)
    	if err != nil {
    		return nil, err
    	}
    
    	var resp types.MineCourseReply
    	fx.From(func(source chan<- interface{}) {
    		for _, each := range studentSelectedCourseList {
    			source <- each
    		}
    	}).Walk(func(item interface{}, pipe chan<- interface{}) {
    		data := item.(*model.SelectionStudent)
    		var teacherName string
    		userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{Id: data.StudentId})
    		if err != nil {
    			logx.Error(errSelectionSelected)
    		} else {
    			teacherName = userInfo.Name
    		}
    
    		selectionCourse, err := l.svcCtx.SelectionCourseModel.FindOne(data.SelectionCourseId)
    		if err != nil {
    			return
    		}
    
    		courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{Id: selectionCourse.CourseId})
    		if err != nil {
    			return
    		}
    
    		resp.List = append(resp.List, &types.Course{
    			Id:                courseInfo.Id,
    			SelectionCourseId: selectionCourse.Id,
    			Name:              courseInfo.Name,
    			Description:       courseInfo.Description,
    			Classify:          courseInfo.Classify,
    			GenderLimit:       int(courseInfo.GenderLimit),
    			MemberLimit:       int(courseInfo.MemberLimit),
    			Credit:            int(courseInfo.Credit),
    			TeacherName:       teacherName,
    		})
    	}).Done()
    
    	sort.Slice(resp.List, func(i, j int) bool {
    		return resp.List[i].Id < resp.List[j].Id
    	})
    
    	return &resp, nil
    }
    ```

### 修改`getteachingcourseshandler.go`
* 文件位置：`service/selection/api/internal/handler/getteachingcourseshandler.go`
* 方法名称：`getTeachingCoursesHandler`
* 代码内容：

    ``` go
    func getTeachingCoursesHandler(ctx *svc.ServiceContext) http.HandlerFunc {
    	return func(w http.ResponseWriter, r *http.Request) {
    		userId,ok:=jwtx.GetUserId(w,r)
    		if !ok{
    			return
    		}
    
    		l := logic.NewGetTeachingCoursesLogic(r.Context(), ctx)
    		resp, err := l.GetTeachingCourses(userId)
    		if err != nil {
    			httpx.Error(w, err)
    		} else {
    			httpx.OkJson(w, resp)
    		}
    	}
    }
    ```
### 获取教师任课课程
* 文件位置：`service/selection/api/internal/logic/getteachingcourseslogic.go`
* 方法名称：`GetTeachingCourses`
* 代码内容：

    ``` go
    func (l *GetTeachingCoursesLogic) GetTeachingCourses(userId int64) (*types.MineCourseReply, error) {
    	selectCourseList, err := l.svcCtx.SelectionCourseModel.FindByTeacherId(userId)
    	if err != nil {
    		return nil, err
    	}
    
    	var teacherName string
    	userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{Id: userId})
    	if err != nil {
    		logx.Error(err)
    	} else {
    		teacherName = userInfo.Name
    	}
    
    	var resp types.MineCourseReply
    	fx.From(func(source chan<- interface{}) {
    		for _, item := range selectCourseList {
    			source <- item
    		}
    	}).Walk(func(item interface{}, pipe chan<- interface{}) {
    		data := item.(*model.SelectionCourse)
    		courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{Id: data.CourseId})
    		if err != nil {
    			logx.Error(err)
    			return
    		}
    
    		resp.List = append(resp.List, &types.Course{
    			Id:                courseInfo.Id,
    			SelectionCourseId: data.Id,
    			Name:              courseInfo.Name,
    			Description:       courseInfo.Description,
    			Classify:          courseInfo.Classify,
    			GenderLimit:       int(courseInfo.GenderLimit),
    			MemberLimit:       int(courseInfo.MemberLimit),
    			Credit:            int(courseInfo.Credit),
    			TeacherName:       teacherName,
    		})
    	}).Done()
    
    	return &resp, nil
    }
    ```
  
# 接口测试

### 准备工作
* 启动redis
    ``` shell script
    $ redis-server
    ```
* 启动etcd
    ``` shell script
    $ etcd
    ```
* 启动两个beanstalkd
    ``` shell script
    $ beanstalkd -l 127.0.0.1 -p 11300
    ```
    ``` shell script
    $ beanstalkd -l 127.0.0.1 -p 11301
    ```
* 分别一次启动`user-api`、`user-rpc`、`course-api`、`course-rpc`、`selection-api`服务

### 接口验证
* 管理员登录
    ``` shell script
    $ curl -i -X POST \
        http://127.0.0.1:8888/api/user/login \
        -H 'content-type: application/json' \
        -d '{
              "username":"gozero",
              "password":"111111"
      }'
    ```
    ``` text
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 16 Dec 2020 08:37:00 GMT
    Content-Length: 178
    
    {"id":2,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxMTE0MjAsImlhdCI6MTYwODEwNzgyMCwiaWQiOjJ9.SKMOrt22pCN2SE0qYiKkmmcJIyr3F0U7hn04pcZLmxQ","expireAt":1608111420}
    ```
  
* 注册一位教师

    ``` shell script
    $ curl -i -X POST \
        http://127.0.0.1:8888/api/user/register \
        -H 'content-type: application/json' \
        -d '{
      	"username":"teacher",
          "password":"111111",
          "role":"teacher"
      }'
    ```
    ``` text
    HTTP/1.1 200 OK
    Date: Wed, 16 Dec 2020 08:35:04 GMT
    Content-Length: 0
    ```
* 添加课程(课程模块)
    
    ``` shell script
    $ curl -i -X POST \
        http://127.0.0.1:8889/api/course/add \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxMTE0MjAsImlhdCI6MTYwODEwNzgyMCwiaWQiOjJ9.SKMOrt22pCN2SE0qYiKkmmcJIyr3F0U7hn04pcZLmxQ' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2' \
        -d '{
      	"name":"Golang",
      	"description":"Golang编程",
      	"classify":"计算机",
      	"genderLimit": 0,
      	"memberLimit":1,
      	"startTime":1608282661,
      	"credit":2
      }'
    ```
    ``` text
    HTTP/1.1 200 OK
    Date: Wed, 16 Dec 2020 09:12:00 GMT
    Content-Length: 0
    ```

    > 说明：这里可以自行多加一些课程数据，就不一一演示了。

* 创建选课任务
    
    ``` shell script
    $ curl -i -X POST \
        http://127.0.0.1:8890/api/selection/create \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxMTgzMjAsImlhdCI6MTYwODExNDcyMCwiaWQiOjJ9.aLhWGj7VfT2HHIW_dUFytQfJkEn055ANXgftArWM2ek' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2' \
        -d '{
      	"name":"2020-2021年学期选课",
      	"maxCredit":12,
      	"startTime": 1608124977,
      	"endTime":1608124977,
      	"notification":"选课开始了"
      }'
    ```
    ``` text
    HTTP/1.1 200 OK
    Date: Wed, 16 Dec 2020 10:34:30 GMT
    Content-Length: 0
    ```
  
* 添加课程（选课）

    ``` shell script
    $ curl -i -X POST \
        http://127.0.0.1:8890/api/selection/add/course/1 \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxMTgzMjAsImlhdCI6MTYwODExNDcyMCwiaWQiOjJ9.aLhWGj7VfT2HHIW_dUFytQfJkEn055ANXgftArWM2ek' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 2' \
        -d '{
      	"list":[
      		{
      			"courseId":3,
      			"teacherId":4
      		}
      	]
      }'
    ```
    ``` text
    HTTP/1.1 200 OK
    Date: Wed, 16 Dec 2020 10:39:33 GMT
    Content-Length: 0
    ```
  
* 选课

    ``` shell script
    $ curl -i -X POST \
        http://127.0.0.1:8890/api/selection/select/1 \
        -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxMTg4NTcsImlhdCI6MTYwODExNTI1NywiaWQiOjF9.jcuMpRw3S5rEu97X3JpD6xmZrdYiwDwIi_d_FScvl5k' \
        -H 'content-type: application/json' \
        -H 'x-user-id: 1' \
        -d '{
      	"list":[
      		{
      			"courseId":3,
      			"teacherId":4
      		}
      	]
      }'
    ```
    ``` text
    HTTP/1.1 400 Bad Request
    Content-Type: text/plain; charset=utf-8
    X-Content-Type-Options: nosniff
    Date: Wed, 16 Dec 2020 10:41:57 GMT
    Content-Length: 16
    
    选课未开始
    ```
    > 说明：这里需要切换为学生用户登录后选课
 
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
 * [zrpc](https://github.com/tal-tech/zero-doc/blob/main/doc/zrpc.md)
 * [使用goctl创建rpc](https://github.com/tal-tech/zero-doc/blob/main/doc/goctl-rpc.md)
 * [使用goctl创建model](https://github.com/tal-tech/zero-doc/blob/main/doc/goctl-model-sql.md)
 * [dq/kq使用说明](https://github.com/tal-tech/go-queue)
 * [beanstalkd](https://beanstalkd.github.io)
 
 # 结尾
 本章节完。
 
 如发现任何错误请通过Issue发起问题修复申请。
 
你可能会浏览 
* [用户模块](../../../doc/requirement/user.md)
* [课程模块](../../../doc/requirement/course.md)


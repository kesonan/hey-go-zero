package svc

import (
	"hey-go-zero/service/course/rpc/courseservice"

	"github.com/tal-tech/go-queue/dq"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"

	"hey-go-zero/service/selection/api/internal/config"
	"hey-go-zero/service/selection/api/internal/middleware"
	"hey-go-zero/service/selection/model"
	"hey-go-zero/service/user/rpc/userservice"
)

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

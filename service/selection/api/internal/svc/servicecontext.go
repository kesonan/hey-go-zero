package svc

import (
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

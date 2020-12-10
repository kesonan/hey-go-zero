package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"hey-go-zero/service/user/api/internal/config"
	"hey-go-zero/service/user/api/internal/middleware"
	"hey-go-zero/service/user/model"
)

type ServiceContext struct {
	Config    config.Config
	UserCheck rest.Middleware
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserCheck: middleware.NewUserCheckMiddleware().Handle,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}

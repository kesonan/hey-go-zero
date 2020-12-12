package svc

import (
	"hey-go-zero/service/user/model"
	"hey-go-zero/service/user/rpc/internal/config"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	c         config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		c:         c,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}

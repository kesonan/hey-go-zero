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

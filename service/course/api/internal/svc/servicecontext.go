//  Copyright [2020] [hey-go-zero]
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

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

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

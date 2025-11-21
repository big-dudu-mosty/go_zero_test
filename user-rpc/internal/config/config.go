package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// Postgres 数据库配置
	Postgres struct {
		// DataSource 数据库连接字符串
		// 格式: postgres://username:password@host:port/database?sslmode=disable
		DataSource string
	}
}

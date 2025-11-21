// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

// Config 应用程序配置结构
// 包含了 REST 服务的基础配置和 PostgreSQL 数据库配置
type Config struct {
	// 嵌入 go-zero 的 REST 配置，包含服务名称、监听地址、端口等
	rest.RestConf

	// Postgres 数据库配置
	Postgres struct {
		// DataSource 数据库连接字符串
		// 格式: postgres://username:password@host:port/database?sslmode=disable
		DataSource string
	}
}

package svc

import (
	"user-demo/user/internal/config"
	"user-demo/user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ServiceContext 服务上下文
// 用于存储服务运行时所需的所有依赖，如配置、数据库连接、Model 等
// 在整个请求链路中传递，便于各层访问共享资源
type ServiceContext struct {
	Config    config.Config   // 应用配置
	UserModel model.UserModel // 用户数据模型，用于操作 user 表
}

// NewServiceContext 创建服务上下文实例
// 参数:
//   - c: 应用配置对象
//
// 返回:
//   - *ServiceContext: 初始化好的服务上下文
//
// 该函数负责:
// 1. 根据配置创建 MySQL 数据库连接
// 2. 初始化各个 Model 实例
// 3. 组装并返回 ServiceContext
func NewServiceContext(c config.Config) *ServiceContext {
	// 使用配置中的连接字符串创建 MySQL 连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn), // 初始化 UserModel
	}
}

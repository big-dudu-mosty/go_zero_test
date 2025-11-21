package svc

import (
	"user-demo/user-api/internal/config"
	"user-demo/user-rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext 服务上下文
// 用于存储服务运行时所需的所有依赖，如配置、RPC 客户端等
// 在整个请求链路中传递，便于各层访问共享资源
type ServiceContext struct {
	Config  config.Config      // 应用配置
	UserRpc userclient.User    // 用户 RPC 客户端
}

// NewServiceContext 创建服务上下文实例
// 参数:
//   - c: 应用配置对象
//
// 返回:
//   - *ServiceContext: 初始化好的服务上下文
//
// 该函数负责:
// 1. 根据配置创建 RPC 客户端连接
// 2. 组装并返回 ServiceContext
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}

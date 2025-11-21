// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 应用程序配置结构
// 包含了 REST 服务的基础配置和 RPC 客户端配置
type Config struct {
	// 嵌入 go-zero 的 REST 配置，包含服务名称、监听地址、端口等
	rest.RestConf

	// UserRpc RPC 客户端配置
	UserRpc zrpc.RpcClientConf
}

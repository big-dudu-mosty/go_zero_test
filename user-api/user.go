// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"user-demo/user-api/internal/config"
	"user-demo/user-api/internal/handler"
	"user-demo/user-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// configFile 配置文件路径标志
// 默认值为 "etc/user-api.yaml"
// 可通过命令行参数 -f 指定其他配置文件，例如: ./user -f /path/to/config.yaml
var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

// main 主函数 - 应用程序入口
// 执行流程:
//  1. 解析命令行参数
//  2. 加载配置文件
//  3. 创建 REST 服务器
//  4. 初始化服务上下文（数据库连接、Model 等）
//  5. 注册所有路由处理器
//  6. 启动服务器并监听请求
func main() {
	// 解析命令行参数
	flag.Parse()

	// 加载配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c) // MustLoad 会在加载失败时 panic

	// 创建 REST 服务器实例
	// 服务器配置包含在 c.RestConf 中（Host、Port、Timeout 等）
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop() // 确保程序退出时优雅关闭服务器

	// 初始化服务上下文
	// ServiceContext 包含了所有业务逻辑需要的依赖（如数据库连接、Model 等）
	ctx := svc.NewServiceContext(c)

	// 注册所有 HTTP 路由处理器
	// 根据 .api 文件定义自动生成的路由注册
	handler.RegisterHandlers(server, ctx)

	// 打印服务启动信息
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	// 启动服务器，开始监听和处理请求
	// 这是一个阻塞调用，会一直运行直到收到退出信号
	server.Start()
}

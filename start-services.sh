#!/bin/bash

# 微服务启动脚本
# 使用方法：./start-services.sh

echo "========================================="
echo "  启动微服务"
echo "========================================="
echo ""

# 检查 PostgreSQL
echo "1. 检查 PostgreSQL 连接..."
pg_isready -h 127.0.0.1 -p 5432 -U admin -d dex_db > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✅ PostgreSQL 运行中"
else
    echo "   ⚠️  PostgreSQL 未运行或无法连接"
    echo "   请确保 PostgreSQL 已启动"
fi
echo ""

# 启动 user-rpc
echo "2. 启动 user-rpc (RPC 服务)..."
cd user-rpc
go run user.go -f etc/user.yaml > /tmp/user-rpc.log 2>&1 &
USER_RPC_PID=$!
echo "   进程 PID: $USER_RPC_PID"
echo "   监听端口: 8080"
echo "   日志文件: /tmp/user-rpc.log"
sleep 2
echo ""

# 启动 user-api
echo "3. 启动 user-api (API 网关)..."
cd ../user-api
go run user.go -f etc/user-api.yaml > /tmp/user-api.log 2>&1 &
USER_API_PID=$!
echo "   进程 PID: $USER_API_PID"
echo "   监听端口: 8888"
echo "   日志文件: /tmp/user-api.log"
sleep 2
echo ""

echo "========================================="
echo "  启动完成！"
echo "========================================="
echo ""
echo "📌 测试地址："
echo "   http://localhost:8888/user/register  (创建用户)"
echo "   http://localhost:8888/user/1         (获取用户)"
echo "   http://localhost:8888/user/list      (用户列表)"
echo ""
echo "📌 查看日志："
echo "   tail -f /tmp/user-rpc.log    (RPC 服务日志)"
echo "   tail -f /tmp/user-api.log    (API 网关日志)"
echo ""
echo "📌 停止服务："
echo "   kill $USER_RPC_PID $USER_API_PID"
echo "   或运行: pkill -f 'user.go'"
echo ""
echo "服务已在后台运行..."

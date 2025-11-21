#!/bin/bash

# 停止微服务脚本

echo "停止所有服务..."
pkill -f 'user.go'

echo "✅ 服务已停止"

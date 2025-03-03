#!/bin/bash

# 进入脚本所在目录
cd "$(dirname "$0")"

echo "停止 GoAccess 容器..."
docker-compose down

echo "停止 goaccess-web..."
pkill goaccess-web

echo "所有服务已停止！"
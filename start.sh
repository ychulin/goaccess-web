#!/bin/bash

# 进入项目目录
cd "$(dirname "$0")"

# 确保 uploads 目录存在
mkdir -p uploads

# 启动 goaccess-web
echo "启动 goaccess-web..."
./goaccess &

# 等待 2 秒确保 Go 服务启动
sleep 2

# 启动 Docker Compose（GoAccess 容器）
echo "启动 GoAccess 容器..."
docker-compose up -d

echo "所有服务已启动，请访问 Web 界面上传日志文件进行分析！"

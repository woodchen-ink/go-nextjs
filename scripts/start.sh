#!/bin/sh

# 确保时区设置正确
export TZ=Asia/Shanghai

# 创建数据目录
mkdir -p /app/data

# 设置数据目录权限为全局可写
chmod -R 777 /app/data
echo "设置数据目录权限："
ls -la /app/data

# 确保环境变量存在
PORT=${PORT:-8080}
NEXT_PORT=${NEXT_PORT:-3000}
echo "后端将使用端口: $PORT"
echo "前端将使用端口: $NEXT_PORT"

# 启动后端应用
cd /app/backend
export PORT=$PORT
./go-nextjs-backend &
BACKEND_PID=$!
echo "后端启动，PID: $BACKEND_PID"

# 启动Next.js应用
cd /app/frontend
export PORT=$NEXT_PORT
node server.js &
FRONTEND_PID=$!
echo "前端启动，PID: $FRONTEND_PID"

# 监听子进程，如果任何一个退出，则退出容器
wait $BACKEND_PID $FRONTEND_PID
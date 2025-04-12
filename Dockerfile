# 最终运行镜像
FROM node:22-alpine AS runner

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache libc6-compat tzdata

# 设置时区为中国标准时间
ENV TZ=Asia/Shanghai
  
# 复制GitHub Actions构建的Go后端
# 使用ARG指定架构，由buildx自动设置
ARG TARGETARCH
COPY go-nextjs-backend-${TARGETARCH} /app/backend/go-nextjs-backend
RUN chmod +x /app/backend/go-nextjs-backend

# 复制Next.js standalone构建结果
COPY web/.next/standalone /app/frontend/
COPY web/.next/static /app/frontend/.next/static
COPY web/public /app/frontend/public

# 添加启动脚本
COPY scripts/start.sh /app/
RUN chmod +x /app/start.sh

# 创建数据目录并设置权限
RUN mkdir -p /app/data && chmod 777 /app/data
  
# 设置工作目录环境变量
# 删除PORT环境变量，让docker-compose.yml中的设置生效
ENV NODE_ENV=production
ENV GIN_MODE=release
ENV DATA_DIR=/app/data

# 以root用户运行
USER root

# 暴露端口
EXPOSE 3000 8080

# 启动应用
CMD ["/app/start.sh"]
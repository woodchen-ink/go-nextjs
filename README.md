# go-nextjs 开发模板

一个集成了Go后端和Next.js前端的全栈开发模板，提供了开发现代Web应用的完整解决方案。

## About

本项目是一个开箱即用的全栈Web应用开发模板，专为快速构建现代化、高性能的Web应用而设计。集成了Go语言后端与Next.js前端，提供了完整的开发环境和项目结构，让开发者可以专注于业务逻辑实现，而不必花费时间在项目搭建和基础架构上。

### 为什么选择这个模板？

- **全栈解决方案**：同时提供前后端完整架构，无需额外配置
- **现代技术栈**：采用Go 1.23+和Next.js 15+等最新技术
- **开发体验优化**：集成热重载、类型检查等提升开发效率的功能
- **容器化部署**：内置Docker配置，支持一键部署
- **安全性考虑**：集成JWT认证，提供基础的安全机制
- **可扩展性强**：模块化设计，易于根据项目需求进行扩展
- **中文友好**：全中文代码注释和文档，适合中文开发者使用

无论是个人项目还是企业应用，这个模板都能帮助您快速启动开发，专注于创建有价值的功能。

## 技术栈

### 后端
- Go 1.23+
- Gin Web框架
- GORM ORM框架
- SQLite数据库
- JWT认证
- Cron定时任务

### 前端
- Next.js 15.2+
- React 19
- TypeScript
- Tailwind CSS
- Radix UI组件库

## 项目结构

```
go-nextjs
├─ config    # 配置层，处理应用配置和数据库连接
├─ cron      # 定时任务，处理周期性执行的任务
├─ data      # 数据文件，存储SQLite数据库和其他持久化数据
├─ handler   # HTTP层，处理API请求和响应
├─ main.go   # 主文件，应用入口点
├─ middleware # 中间件，处理认证、日志等横切关注点
├─ models    # 数据模型，定义数据库模型和关系
├─ pkg       # 外部包，存放通用工具和辅助函数
├─ router    # 路由，定义API路由
├─ scripts   # 脚本，包含启动和部署脚本
│  └─ start.sh # 启动脚本
├─ service   # 服务层，包含业务逻辑
└─ web       # 前端Next.js应用
   ├─ app      # Next.js应用代码
   ├─ components # React组件
   ├─ public    # 静态资源
   └─ lib       # 前端工具库
```

## 快速开始

### 本地开发

1. 克隆仓库
```bash
git clone https://github.com/woodchen-ink/go-nextjs.git
cd go-nextjs
```

2. 配置环境变量（复制.env.example到.env并修改必要配置）

3. 启动后端服务
```bash
go run main.go
```

4. 启动前端开发服务器
```bash
cd web
npm install
npm run dev
```

### Docker部署

使用Docker Compose一键部署：
```bash
docker-compose up -d
```

## 特性

- 完整的前后端分离架构
- 内置用户认证系统
- 定时任务支持
- 响应式UI设计
- 开发和生产环境配置
- Docker容器化支持

## 许可证

MIT


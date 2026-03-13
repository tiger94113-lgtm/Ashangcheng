# Asset Mall 后端部署指南

## 方案一：Render（推荐 ⭐）

### 步骤：

1. **注册账号**
   - 访问 https://render.com
   - 用 GitHub 账号登录

2. **创建仓库**
   ```bash
   # 在 backend 目录初始化 git
   cd backend
   git init
   git add .
   git commit -m "Initial commit"
   
   # 创建 GitHub 仓库并推送
   git remote add origin https://github.com/你的用户名/asset-mall-backend.git
   git push -u origin main
   ```

3. **部署到 Render**
   - 登录 https://dashboard.render.com
   - 点击 "New +" → "Web Service"
   - 选择你的 GitHub 仓库
   - 配置：
     - **Name**: asset-mall-backend
     - **Runtime**: Go
     - **Build Command**: `go build -o main .`
     - **Start Command**: `./main`
   - 点击 "Create Web Service"

4. **获取域名**
   - 部署完成后，Render 会给你一个域名
   - 例如：`https://asset-mall-backend.onrender.com`

5. **更新前端配置**
   ```javascript
   // 在前端 app.js 中修改
   const API_BASE = "https://asset-mall-backend.onrender.com/api";
   ```

---

## 方案二：Railway（推荐 ⭐⭐）

### 步骤：

1. **注册账号**
   - 访问 https://railway.app
   - 用 GitHub 登录

2. **部署**
   - 点击 "New Project"
   - 选择 "Deploy from GitHub repo"
   - 选择你的仓库
   - Railway 会自动检测 Go 项目并部署

3. **获取域名**
   - 部署完成后，Settings → Domains 查看

---

## 方案三：Fly.io（最稳定）

### 步骤：

1. **安装 Fly CLI**
   ```bash
   # Windows
   powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"
   
   # Mac/Linux
   curl -L https://fly.io/install.sh | sh
   ```

2. **登录**
   ```bash
   fly auth login
   ```

3. **部署**
   ```bash
   cd backend
   fly launch
   # 按照提示操作
   ```

---

## 方案四：Vercel（Serverless）

Vercel 更适合前端，但也可以部署 Go：

1. 创建 `vercel.json`:
```json
{
  "version": 2,
  "builds": [
    {
      "src": "main.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "main.go"
    }
  ]
}
```

2. 推送到 GitHub
3. 在 Vercel 导入项目

---

## 数据库持久化说明

### Render 免费版限制
- 磁盘不是永久的，重启后数据会丢失
- **解决方案**：使用 Render PostgreSQL（免费）或定期备份

### 使用 PostgreSQL（推荐用于生产）

1. 在 Render 创建 PostgreSQL 数据库
2. 修改 `main.go` 使用 PostgreSQL：

```go
import "github.com/lib/pq"

// 连接 PostgreSQL
db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
```

### 数据备份脚本

创建 `backup.sh`：
```bash
#!/bin/bash
# 备份 SQLite 数据库
curl -o backup-$(date +%Y%m%d).db https://your-api-url/backup
```

---

## 前端更新

部署完成后，修改前端 API 地址：

```javascript
// app.js 或 api-client.js
const API_CONFIG = {
  baseURL: "https://你的域名/api",  // ← 修改这里
};
```

---

## 监控和维护

### 健康检查
访问：`https://你的域名/api/health`

### 日志查看
- Render: Dashboard → Logs
- Railway: Dashboard → Deployments → Logs
- Fly.io: `fly logs`

### 免费额度限制

| 平台 | 限制 |
|------|------|
| Render | 15分钟无请求会休眠，启动需几秒 |
| Railway | $5/月额度，用完会暂停 |
| Fly.io | $5/月额度，按用量计费 |

---

## 推荐配置

**开发/测试**：Render（免费、简单）
**生产环境**：Railway 或 Fly.io（更稳定）

需要我帮你配置 PostgreSQL 或其他功能吗？

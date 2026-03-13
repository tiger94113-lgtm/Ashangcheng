# Asset Mall Go 后端

简单的 Go + SQLite 后端，用于存储和管理订单数据。

## 功能特性

- ✅ RESTful API
- ✅ SQLite 数据库（无需额外配置）
- ✅ CORS 支持
- ✅ 订单 CRUD 操作
- ✅ 搜索和筛选
- ✅ 统计数据

## API 接口

### 订单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/orders` | 获取订单列表 |
| GET | `/api/orders/:id` | 获取单个订单 |
| POST | `/api/orders` | 创建订单 |
| PUT | `/api/orders/:id` | 更新订单 |
| DELETE | `/api/orders/:id` | 删除订单 |

### 查询参数

**GET /api/orders**
- `wallet` - 按钱包地址筛选
- `status` - 按状态筛选 (pending/confirmed/failed)
- `search` - 搜索订单号、姓名、电话等

### 其他

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats` | 获取统计数据 |
| GET | `/api/health` | 健康检查 |

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod download
```

### 2. 运行

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 3. 测试

```bash
# 健康检查
curl http://localhost:8080/api/health

# 创建订单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "orderNo": "OD202401011200001",
    "status": "pending",
    "wallet": "0x123...",
    "name": "张三",
    "phone": "13800138000",
    "address": "北京市...",
    "usdtAmount": "1000000000000000000",
    "txHash": "0xabc..."
  }'

# 获取订单列表
curl http://localhost:8080/api/orders

# 按钱包筛选
curl "http://localhost:8080/api/orders?wallet=0x123..."

# 搜索
curl "http://localhost:8080/api/orders?search=张三"
```

## 部署

### 使用 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/orders.db .
EXPOSE 8080
CMD ["./main"]
```

### 部署到服务器

```bash
# 编译
GOOS=linux GOARCH=amd64 go build -o main .

# 上传并运行
scp main user@server:/path/to/app/
ssh user@server "cd /path/to/app && ./main"
```

### 使用 systemd

创建 `/etc/systemd/system/asset-mall.service`:

```ini
[Unit]
Description=Asset Mall Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/backend
ExecStart=/path/to/backend/main
Restart=always
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

启用服务：
```bash
sudo systemctl enable asset-mall
sudo systemctl start asset-mall
```

## 前端集成

修改前端 `app.js` 中的 API 地址：

```javascript
const API_BASE = "http://localhost:8080/api";

// 创建订单
async function saveOrderToBackend(order) {
  const response = await fetch(`${API_BASE}/orders`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(order)
  });
  return response.json();
}

// 获取订单
async function getOrdersFromBackend(wallet) {
  const response = await fetch(
    `${API_BASE}/orders?wallet=${wallet}`
  );
  return response.json();
}
```

## 数据备份

SQLite 数据库文件 `orders.db` 可以直接复制备份：

```bash
# 备份
cp orders.db orders.db.backup.$(date +%Y%m%d)

# 恢复
cp orders.db.backup.20240101 orders.db
```

## 免费部署选项

1. **Railway** (railway.app) - 免费额度充足
2. **Render** (render.com) - 免费 Web 服务
3. **Fly.io** (fly.io) - 免费额度
4. **Oracle Cloud** - 永久免费 VPS

## 注意事项

1. 生产环境建议添加身份验证
2. 定期备份数据库文件
3. 可以使用 Nginx 反向代理
4. 建议启用 HTTPS

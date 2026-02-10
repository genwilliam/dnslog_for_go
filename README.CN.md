# dnslog_for_go

中文 | [English](README.md)

## 项目简介

`dnslog_for_go` 是一个基于被动观测模型的 DNSLog 平台：**生成 → 外部触发 → 记录 → 查询/告警**。系统不主动解析 DNS，请求来源全部来自外部真实触发，核心职责是采集、存储、查询、展示与告警。

Redis 仅用于限流、异步队列、可选状态缓存与黑名单加速，不作为 `dns_records`/`dns_tokens` 主存储，MySQL 是最终事实存储。

## 功能一览

- 生成 token 与域名，并写入 `dns_tokens` 状态表
- DNS 被动捕获、记录入库（`dns_records`）
- Token 状态查询与原始记录查询
- API Key 鉴权、限流、黑名单、审计日志
- Webhook 通知（首次命中），失败退避重试
- Prometheus metrics、日志保留策略、备份/恢复
- 管理端 UI（Tokens 列表、API Keys/黑名单管理）

## 目录结构

```
dnslog_for_go/
├── cmd/dnslog/                 # 入口
├── config/                     # 配置文件与示例
├── db/migrations/              # 数据库迁移（001~006）
├── internal/
│   ├── dnslog/                 # DNS 捕获、状态表、存储、Webhook、审计
│   ├── domain/                 # HTTP 业务 handlers（生成/查询/配置）
│   ├── router/                 # 路由注册与中间件挂载
│   ├── middleware/             # 鉴权/限流/审计/黑名单/trace/metrics
│   ├── infra/                  # Redis 客户端
│   ├── metrics/                # Prometheus 指标
│   └── ...
├── pkg/                        # 公共工具（日志、响应、utils）
├── web/dnslog/                 # 前端（Vue3）
├── scripts/                    # 备份/恢复/迁移脚本
├── docs/                       # API/cron 文档
├── go.mod
├── README.CN.md
└── README.md
```

### 前端目录

- `web/dnslog/src/views`: 页面（dns-query/dnslog/tokens/security）
- `web/dnslog/src/api`: API 封装
- `web/dnslog/src/utils`: 通用工具

---

# 快速开始

## 1) 环境准备

建议版本：

- Go >= 1.21
- Node.js >= 18
- MySQL >= 8.0
- Redis >= 6.0

必备工具：

- git
- curl
- dig（macOS 可用 `brew install bind`，Ubuntu 可用 `apt install dnsutils`）

macOS 示例（Homebrew）：

```bash
brew install go node mysql redis bind git
```

Linux 示例（Ubuntu）：

```bash
sudo apt update
sudo apt install -y git curl dnsutils mysql-server redis-server
```

Windows 建议使用 WSL2（Ubuntu）。

## 2) 下载代码 + 依赖安装

```bash
git clone https://github.com/genwilliam/dnslog_for_go.git
cd dnslog_for_go

go mod tidy

cd web/dnslog
npm i
```

## 3) MySQL 初始化（本地或 Docker）

### 方式 A：本地安装 MySQL

启动 MySQL 后执行：

```bash
mysql -u root -p
```

创建数据库与用户：

```sql
CREATE DATABASE dnslog DEFAULT CHARSET utf8mb4;
CREATE USER 'dnslog'@'%' IDENTIFIED BY 'dnslog';
GRANT ALL PRIVILEGES ON dnslog.* TO 'dnslog'@'%';
FLUSH PRIVILEGES;
```

执行迁移（按顺序 001~006）：

```bash
cd /path/to/dnslog_for_go
mysql -u dnslog -p dnslog < db/migrations/001_init_dnslog.sql
mysql -u dnslog -p dnslog < db/migrations/002_tokens.sql
mysql -u dnslog -p dnslog < db/migrations/003_api_security.sql
mysql -u dnslog -p dnslog < db/migrations/004_ip_blacklist.sql
mysql -u dnslog -p dnslog < db/migrations/005_webhooks.sql
mysql -u dnslog -p dnslog < db/migrations/006_indexes.sql
```

### 方式 B：Docker 快速启动

```bash
docker run -d --name dnslog-mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=dnslog \
  -e MYSQL_USER=dnslog \
  -e MYSQL_PASSWORD=dnslog \
  -p 3306:3306 mysql:8.0
```

迁移执行方式同上（连接到 127.0.0.1:3306）。

## 4) Redis 启动与说明

Redis 用途：

- 限流
- 异步队列（审计、webhook）
- 黑名单加速
- 可选状态缓存

启动方式：

```bash
# 本地
redis-server

# 或 Docker
docker run -d --name dnslog-redis -p 6379:6379 redis:6
```

配置示例（`config/config.yaml`）：

```
redisAddr: "127.0.0.1:6379"
redisPassword: ""
redisDB: 0
```

Redis 不可用会怎样：

- 若 `rateLimitEnabled/auditEnabled/webhookEnabled` 仍为 true，后端会启动失败
- 临时解决：将这些开关设为 false 后再启动

## 5) 后端配置

复制配置：

```bash
cd /path/to/dnslog_for_go
cp config/config.example.yaml config/config.yaml
```

至少修改：

- `rootDomain`
- `mysqlDSN`
- `redisAddr`

DSN 示例：

```
mysqlDSN: "dnslog:dnslog@tcp(127.0.0.1:3306)/dnslog?parseTime=true&loc=Local&charset=utf8mb4"
```

## 6) 启动 Go 后端

```bash
cd /path/to/dnslog_for_go
go run cmd/dnslog/main.go
```

默认端口：

- HTTP：`:8080`
- DNS：`:15353`

如果要绑定 53 端口：

- macOS/Linux 需要 root 或 setcap
- 示例（Linux）：

```bash
sudo setcap 'cap_net_bind_service=+ep' /path/to/binary
```

## 7) 启动前端

```bash
cd /path/to/dnslog_for_go/web/dnslog
npm run dev
```

访问：`http://localhost:5173`

API Key 配置方式：

- 推荐：前端页面顶部输入框保存（写入 `DNSLOG_API_KEY`）
- 或环境变量：

```bash
VITE_API_KEY=your-plain-key npm run dev
```

Vite 已内置代理（`/api -> http://localhost:8080`），可在 `vite.config.ts` 查看：

```
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
  },
}
```

---

# 启动后

## 7.1 创建 API Key（bootstrap）

若 `apiKeyRequired=false` 可以跳过本步骤；默认开启时需先创建。

首次创建无需鉴权（仅当**没有 enabled key**时允许）：

```bash
curl -X POST http://localhost:8080/api/keys \
  -H 'Content-Type: application/json' \
  -d '{"name":"ops","comment":"bootstrap"}'
```

- 明文 key 只返回一次，请保存
- 如果返回 401/409：说明已有 key（或 enable=1），需要使用已有 key 或重置表

## 7.2 生成 token/domain

前端点击 Generate，或调用接口：

```bash
curl -X POST http://localhost:8080/api/tokens \
  -H 'X-API-Key: <your_key>'
```

示例返回：

```json
{
  "data": {
    "domain": "6f643ae19c.demo.com",
    "token": "6f643ae19c"
  }
}
```

## 7.3 触发 DNS（必须指定服务器与端口）

重要：不能用系统默认 DNS，必须指定 DNSLog 服务器与端口。

本地验证（监听 15353）：

```bash
dig @127.0.0.1 -p 15353 6f643ae19c.demo.com A
```

原因：DNS 服务监听 `:15353`，不指定 `@server`/`-p` 会走系统 DNS，无法命中本服务。

## 7.4 前端看到 HIT

- Token 状态从 `INIT` → `HIT`
- records 列表出现记录

如未命中：

- 检查 DNS 监听地址与端口映射
- 检查防火墙/安全组
- 确认 `rootDomain` 是否一致

---

## 配置项参考（关键项）

配置文件：`config/config.example.yaml` → `config/config.yaml`

优先级：**环境变量 > config.yaml**

| 字段                        | 默认值              | 说明                                      | 示例                                     |
| --------------------------- | ------------------- | ----------------------------------------- | ---------------------------------------- |
| rootDomain                  | demo.com            | 单个根域                                  | demo.com                                 |
| rootDomains                 | [demo.com]          | 多根域列表                                | ["demo.com","example.com"]               |
| captureAll                  | false               | 记录所有域名请求                          | true                                     |
| dnsListenAddr               | :15353              | DNS 监听地址                              | :15353                                   |
| httpListenAddr              | :8080               | HTTP 监听地址                             | :8080                                    |
| upstreamDNS                 | [8.8.8.8,223.5.5.5] | 上游 DNS 列表                             | ["8.8.8.8"]                              |
| protocol                    | udp                 | DNS 协议                                  | udp/tcp                                  |
| mysqlDSN                    | -                   | MySQL DSN                                 | user:pass@tcp(127.0.0.1:3306)/dnslog?... |
| pageSize                    | 20                  | 默认分页大小                              | 20                                       |
| maxPageSize                 | 100                 | 最大分页                                  | 100                                      |
| tokenTTLSeconds             | 3600                | token TTL（秒）                           | 3600                                     |
| apiKeyRequired              | true                | API Key 鉴权                              | true/false                               |
| rateLimitEnabled            | true                | HTTP 限流开关                             | true/false                               |
| rateLimitWindowSeconds      | 60                  | HTTP 限流窗口                             | 60                                       |
| rateLimitMaxRequests        | 60                  | HTTP 限流阈值                             | 60                                       |
| dnsRateLimitEnabled         | true                | DNS 限流开关                              | true/false                               |
| dnsRateLimitWindowSeconds   | 60                  | DNS 限流窗口                              | 60                                       |
| dnsRateLimitMaxRequests     | 1000                | DNS 限流阈值                              | 1000                                     |
| auditEnabled                | true                | 审计日志                                  | true/false                               |
| blacklistEnabled            | N/A                 | 当前无开关（黑名单通过表内 enabled 控制） | -                                        |
| webhookEnabled              | true                | Webhook 开关                              | true/false                               |
| webhookMaxRetries           | 4                   | 最大重试次数                              | 4                                        |
| webhookRetryIntervalSeconds | 30                  | 重试扫描间隔                              | 30                                       |
| webhookSecretKey            | -                   | AES-GCM 密钥（32 字节）                   | base64/hex                               |
| metricsEnabled              | true                | Metrics 开关                              | true/false                               |
| metricsPublic               | false               | Metrics 是否公开                          | true/false                               |
| redisAddr                   | 127.0.0.1:6379      | Redis 地址                                | 127.0.0.1:6379                           |
| redisPassword               |                     | Redis 密码                                | -                                        |
| redisDB                     | 0                   | Redis DB                                  | 0                                        |
| retentionEnabled            | true                | 保留策略开关                              | true/false                               |
| recordRetentionDays         | 30                  | 记录保留天数                              | 30                                       |
| retentionIntervalSeconds    | 3600                | 清理周期（秒）                            | 3600                                     |
| retentionBatchSize          | 1000                | 清理批次大小                              | 1000                                     |

环境变量（示例）：

```
MYSQL_DSN=...
REDIS_ADDR=127.0.0.1:6379
API_KEY_REQUIRED=true
RATE_LIMIT_ENABLED=true
WEBHOOK_ENABLED=true
```

---

## Redis 用途与边界

- Redis 用于限流、异步队列、可选状态缓存与黑名单加速
- Redis 不作为 `dns_records` / `dns_tokens` 主存储
- MySQL 是最终事实存储

---

## API

完整 API 文档见：`docs/api.md`

- `GET /api/tokens/{token}`：token 状态
- `GET /api/tokens/{token}/records`：raw records
- `GET /api/records`：记录检索
- `POST /api/tokens/{token}/webhook`：绑定 webhook
- `POST /api/submit`：legacy 接口

---

## 常见问题（FAQ）

1. **为什么返回 401？**
   未携带或使用了错误的 `X-API-Key`。
2. **为什么返回 403？**
   可能命中黑名单，检查 `ip_blacklist` 表。
3. **为什么返回 429？**
   触发限流，调高阈值或更换 key。
4. **前端黑屏/空白？**
   未配置 API Key 或接口路径错误，检查 `/api` 代理与 `VITE_API_KEY`。
5. **DNS 不命中？**
   dig 必须指定 `@server -p`，并确认端口/防火墙/根域一致。
6. **Redis 不可用导致启动失败？**
   关闭 `rateLimitEnabled/auditEnabled/webhookEnabled` 后再启动。
7. **token 显示 EXPIRED？**
   超过 `tokenTTLSeconds`，重新生成即可。
8. **Webhook 只触发一次？**
   当前默认 FIRST_HIT，仅首次命中触发，失败会退避重试。
9. **/metrics 访问失败？**
   `metricsPublic=false` 时需要 API Key。
10. **多实例部署注意事项？**
    应用无状态，MySQL/Redis 需共享，Webhook 接收端建议幂等。
11. **/api/keys 返回 409？**
    说明已初始化 key，需携带 `X-API-Key` 或清空/禁用表内 key。
12. **DNS 端口 53 权限不足？**
    需 root 或 `setcap` 赋权。

---

## Webhook Secret 加密与迁移

- Webhook secret 使用 AES-GCM 加密存储
- 密钥通过 `webhookSecretKey` / `WEBHOOK_SECRET_KEY` 提供（32 字节，base64 或 hex）

迁移旧明文 secret：

```bash
WEBHOOK_SECRET_KEY=your_key \
MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/dnslog?parseTime=true&loc=Local&charset=utf8mb4" \
go run scripts/migrate_webhook_secrets.go
```

## 监控

- `/metrics` 暴露 Prometheus 指标（可设置为私有）

## 备份与恢复

脚本：

```bash
scripts/backup.sh
scripts/restore.sh /path/to/backup.sql
```

cron 示例见：`docs/cron.md`

## 多实例部署

- 应用无状态，可水平扩展
- MySQL 与 Redis 需共享

## License

MIT

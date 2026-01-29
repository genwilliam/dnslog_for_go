# DNSLog API

## 基本路径
后端路由挂载在 `/` 路径下。在本地开发环境中，前端会将 `/api/*` 代理到后端（参见 Vite 代理配置）。以下示例均使用 `/api` 前缀；如果您直接调用后端而不使用代理，请移除 `/api` 前缀。

## 身份验证
所有受保护的端点都需要：

```
X-API-Key: <your_api_key>
```

## 响应格式
```json
{
"code": 200,
"message": "ok",
"data": {},
"trace_id": "...",
"timestamp": 1690000000000
}
```

## 令牌
### POST /api/tokens
生成令牌/域名（`GET /api/random-domain` 的别名）。

### GET /api/tokens
列出令牌。

查询参数：
- `status`: INIT | HIT | EXPIRED
- `keyword`: 令牌或域名关键词
- `created_start`, `created_end`: 毫秒时间戳
- `last_start`, `last_end`: 毫秒时间戳
- `orderBy`: created_at | last_seen
- `order`: asc | desc
- `page`, `pageSize`

### GET /api/tokens/{token}
获取令牌状态。

响应数据：
```json
{
"token": "abc123",
"domain": "abc123.example.com",
"status": "HIT",
"first_seen": 1690000000,
"last_seen": 1690000123,
"hit_count": 3,
"expires_at": 1690003600,
"expired": false
}
```

### GET /api/tokens/{token}/records
按令牌获取原始 DNS 记录。

查询参数：
- `page`, `pageSize`
- `order`: asc | desc

### POST /api/tokens/{token}/webhook
绑定 Webhook（仅限首次命中）。正文：
```json
{ "webhook_url": "https://example.com/hook", "secret": "abc", "mode": "FIRST_HIT" }
```

### GET /api/tokens/{token}/webhook
获取令牌 Webhook。

### DELETE /api/tokens/{token}/webhook
禁用令牌 Webhook。

（兼容性）`POST /api/tokens/{token}/webhook/disable`

## 记录
### GET /api/records
查询原始记录。

查询参数：
- `domain`、`token`
- `client_ip`、`protocol`、`qtype`
- `start`、`end`（毫秒时间戳）
- `order`（asc | desc）
- `page`、`pageSize`
- `cursor`（毫秒时间戳，用于键集分页）

## API 密钥
### POST /api/keys
创建 API 密钥（明文仅返回一次）。

正文：
```json
{ "name": "ops", "comment": "rotation-2025-01" }
```

### GET /api/keys
列出 API 密钥。

### DELETE /api/keys/{id}
禁用 API 密钥。

（兼容性）`POST /api/keys/{id}/disable`

## 黑名单
### POST /api/blacklist
将 IP 地址添加到黑名单。

正文：
```json
{ "ip": "1.2.3.4", "reason": "abuse" }
```

### GET /api/blacklist
列出黑名单条目。

### DELETE /api/blacklist/{id}
禁用黑名单条目。

（兼容性）`POST /api/blacklist/{id}/disable`

## 配置
### GET /api/config
运行时配置（如果 `publicConfig=true` 则为公开）。

## 指标
### GET /metrics
Prometheus 指标端点（除非 `metricsPublic=true`，否则受保护）。

## 旧版
### POST /api/submit
旧版查询，仅用于兼容性。

### GET /api/random-domain
旧版令牌生成（`POST /api/tokens` 的别名）。 ## Webhook 交付
Webhook 请求包含：
- `X-Event-ID`：用于保证幂等性的作业 ID
- `X-Signature`：当设置了密钥时，为 HMAC-SHA256(payload, secret)

## 错误代码
- `unauthorized`（未授权）
- `invalid_api_key`（API 密钥无效）
- `rate_limited`（请求频率受限）
- `rate_limit_error`（速率限制错误）
- `rate_limit_unavailable`（速率限制服务不可用）
- `token_not_found`（令牌未找到）
- `not_found`（未找到）
- `bad_request`（请求无效）
- `internal_error`（内部错误）
- `system_paused`（系统已暂停）
- `webhook_secret_key_required`（需要 Webhook 密钥）

## Redis 使用
- Redis 用于速率限制、异步队列、可选状态缓存和黑名单加速。
- Redis 不用作 `dns_records` 或 `dns_tokens` 的主要存储。
- MySQL 仍然是数据的权威来源。
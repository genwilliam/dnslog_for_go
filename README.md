# dnslog_for_go

[中文](README.CN.md) | English

## Overview
`dnslog_for_go` is a passive DNSLog platform: **generate → external trigger → record → query/alert**. The server never actively resolves DNS. MySQL is the source of truth; Redis is used for rate limiting, async queues, and acceleration.

## Quick Start (Beginner)

### Requirements
- Go >= 1.21
- Node.js >= 18
- MySQL >= 8.0
- Redis >= 6.0
- Tools: git, curl, dig

### 1) Clone & install
```bash
git clone https://github.com/genwilliam/dnslog_for_go.git
cd dnslog_for_go

go mod tidy

cd web/dnslog
npm i
```

### 2) MySQL
Create database and user:
```sql
CREATE DATABASE dnslog DEFAULT CHARSET utf8mb4;
CREATE USER 'dnslog'@'%' IDENTIFIED BY 'dnslog';
GRANT ALL PRIVILEGES ON dnslog.* TO 'dnslog'@'%';
FLUSH PRIVILEGES;
```
Run migrations:
```bash
mysql -u dnslog -p dnslog < db/migrations/001_init_dnslog.sql
mysql -u dnslog -p dnslog < db/migrations/002_tokens.sql
mysql -u dnslog -p dnslog < db/migrations/003_api_security.sql
mysql -u dnslog -p dnslog < db/migrations/004_ip_blacklist.sql
mysql -u dnslog -p dnslog < db/migrations/005_webhooks.sql
mysql -u dnslog -p dnslog < db/migrations/006_indexes.sql
```

### 3) Redis
```bash
redis-server
# or
# docker run -d --name dnslog-redis -p 6379:6379 redis:6
```

### 4) Configure
```bash
cp config/config.example.yaml config/config.yaml
```
Edit `rootDomain`, `mysqlDSN`, `redisAddr`.

### 5) Start backend
```bash
go run cmd/dnslog/main.go
```
Ports:
- HTTP: `:8080`
- DNS: `:15353`

### 6) Start frontend
```bash
cd web/dnslog
VITE_API_KEY=your-plain-key npm run dev
```
Visit: `http://localhost:5173`

### 7) Bootstrap API Key (first key)
If no enabled key exists, you can create the first one without auth:
```bash
curl -X POST http://localhost:8080/api/keys \
  -H 'Content-Type: application/json' \
  -d '{"name":"ops","comment":"bootstrap"}'
```

### 8) Generate token/domain
```bash
curl -X POST http://localhost:8080/api/tokens \
  -H 'X-API-Key: <your_key>'
```

### 9) Trigger DNS (must specify server/port)
```bash
dig @127.0.0.1 -p 15353 <token>.demo.com A
```

### 10) Verify
- Token status becomes HIT
- Records list shows entries

## Notes
- Redis is required when rate limit/audit/webhook are enabled.
- `/api` prefix is required for frontend calls.
- The first API Key is returned only once.
- If `apiKeyRequired=false`, the frontend will not force API Key; calls work without `X-API-Key`.

For full details and troubleshooting, see `README.CN.md` and `docs/api.md`.

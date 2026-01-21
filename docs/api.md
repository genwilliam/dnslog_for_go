# DNSLog API

## Base Path
The backend routes are mounted at `/`. In local dev, the frontend proxies `/api/*` to the backend (see Vite proxy). The examples below use `/api` prefix; if you call the backend directly without proxy, drop the `/api` prefix.

## Auth
All protected endpoints require:

```
X-API-Key: <your_api_key>
```

## Response Shape
```json
{
  "code": 200,
  "message": "ok",
  "data": {},
  "trace_id": "...",
  "timestamp": 1690000000000
}
```

## Tokens
### POST /api/tokens
Generate token/domain (alias of `GET /api/random-domain`).

### GET /api/tokens
List tokens.

Query params:
- `status`: INIT | HIT | EXPIRED
- `keyword`: token or domain keyword
- `created_start`, `created_end`: ms timestamp
- `last_start`, `last_end`: ms timestamp
- `orderBy`: created_at | last_seen
- `order`: asc | desc
- `page`, `pageSize`

### GET /api/tokens/{token}
Get token status.

Response data:
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
Raw DNS records by token.

Query params:
- `page`, `pageSize`
- `order`: asc | desc

### POST /api/tokens/{token}/webhook
Bind webhook (FIRST_HIT only).

Body:
```json
{ "webhook_url": "https://example.com/hook", "secret": "abc", "mode": "FIRST_HIT" }
```

### GET /api/tokens/{token}/webhook
Get token webhook.

### DELETE /api/tokens/{token}/webhook
Disable token webhook.

(Compatibility) `POST /api/tokens/{token}/webhook/disable`

## Records
### GET /api/records
Query raw records.

Query params:
- `domain`, `token`
- `client_ip`, `protocol`, `qtype`
- `start`, `end` (ms timestamp)
- `order` (asc | desc)
- `page`, `pageSize`
- `cursor` (ms timestamp, for keyset paging)

## API Keys
### POST /api/keys
Create API key (plaintext returned once).

Body:
```json
{ "name": "ops", "comment": "rotation-2025-01" }
```

### GET /api/keys
List API keys.

### DELETE /api/keys/{id}
Disable API key.

(Compatibility) `POST /api/keys/{id}/disable`

## Blacklist
### POST /api/blacklist
Add IP to blacklist.

Body:
```json
{ "ip": "1.2.3.4", "reason": "abuse" }
```

### GET /api/blacklist
List blacklist entries.

### DELETE /api/blacklist/{id}
Disable blacklist entry.

(Compatibility) `POST /api/blacklist/{id}/disable`

## Config
### GET /api/config
Runtime config (public if `publicConfig=true`).

## Metrics
### GET /metrics
Prometheus metrics endpoint (protected unless `metricsPublic=true`).

## Legacy
### POST /api/submit
Legacy query, compatibility only.

### GET /api/random-domain
Legacy token generation (alias of `POST /api/tokens`).

## Webhook Delivery
Webhook requests include:
- `X-Event-ID`: job id for idempotency
- `X-Signature`: HMAC-SHA256(payload, secret) when secret is set

## Error Codes
- `unauthorized`
- `invalid_api_key`
- `rate_limited`
- `rate_limit_error`
- `rate_limit_unavailable`
- `token_not_found`
- `not_found`
- `bad_request`
- `internal_error`
- `system_paused`
- `webhook_secret_key_required`

## Redis Usage
- Redis is used for rate limiting, async queues, optional status cache, and blacklist acceleration.
- Redis is not used as primary storage for `dns_records` or `dns_tokens`.
- MySQL remains the source of truth.

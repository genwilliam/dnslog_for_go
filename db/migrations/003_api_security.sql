CREATE TABLE IF NOT EXISTS api_keys (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    api_key VARCHAR(128) NOT NULL UNIQUE,
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    last_used_at BIGINT NOT NULL DEFAULT 0,
    comment VARCHAR(255) DEFAULT '',
    INDEX idx_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trace_id VARCHAR(64) NOT NULL,
    api_key_id BIGINT DEFAULT NULL,
    path VARCHAR(128) NOT NULL,
    method VARCHAR(16) NOT NULL,
    client_ip VARCHAR(64) NOT NULL,
    status_code INT NOT NULL,
    latency_ms INT NOT NULL,
    token VARCHAR(128) DEFAULT '',
    created_at BIGINT NOT NULL,
    INDEX idx_created (created_at),
    INDEX idx_path (path),
    INDEX idx_ip (client_ip)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

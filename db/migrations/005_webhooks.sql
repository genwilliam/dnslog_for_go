CREATE TABLE IF NOT EXISTS token_webhooks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(128) NOT NULL,
    webhook_url VARCHAR(512) NOT NULL,
    secret VARCHAR(128) DEFAULT '',
    mode ENUM('FIRST_HIT','EACH_HIT') NOT NULL DEFAULT 'FIRST_HIT',
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    UNIQUE KEY uk_token (token),
    INDEX idx_token (token),
    INDEX idx_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS webhook_jobs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(128) NOT NULL,
    url VARCHAR(512) NOT NULL,
    payload TEXT NOT NULL,
    secret VARCHAR(128) DEFAULT '',
    status ENUM('PENDING','SUCCESS','FAILED') NOT NULL DEFAULT 'PENDING',
    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at BIGINT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_status_retry (status, next_retry_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS dns_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token       VARCHAR(128) NOT NULL UNIQUE,
    domain      VARCHAR(255) NOT NULL,
    status      ENUM('INIT','HIT','EXPIRED') NOT NULL DEFAULT 'INIT',
    hit_count   BIGINT NOT NULL DEFAULT 0,
    first_seen  BIGINT NOT NULL DEFAULT 0,
    last_seen   BIGINT NOT NULL DEFAULT 0,
    created_at  BIGINT NOT NULL,
    updated_at  BIGINT NOT NULL,
    expires_at  BIGINT NOT NULL,
    INDEX idx_status (status),
    INDEX idx_expires (expires_at),
    INDEX idx_status_created (status, created_at),
    INDEX idx_status_last (status, last_seen)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

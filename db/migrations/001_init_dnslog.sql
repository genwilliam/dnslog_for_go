-- 创建数据库（若已存在可忽略）
CREATE DATABASE IF NOT EXISTS `dnslog`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `dnslog`;

-- DNS 请求记录表（与代码中的自动建表一致）
CREATE TABLE IF NOT EXISTS `dns_records` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `domain`     VARCHAR(255) NOT NULL,
  `client_ip`  VARCHAR(64)  DEFAULT '',
  `protocol`   VARCHAR(16)  DEFAULT '',
  `qtype`      VARCHAR(32)  DEFAULT '',
  `timestamp`  BIGINT       NOT NULL,
  `server`     VARCHAR(64)  DEFAULT '',
  `token`      VARCHAR(128) DEFAULT '',
  `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_domain` (`domain`),
  INDEX `idx_token` (`token`),
  INDEX `idx_ts` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

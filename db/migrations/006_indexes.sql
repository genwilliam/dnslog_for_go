ALTER TABLE dns_records
  ADD INDEX idx_token_ts (token, timestamp),
  ADD INDEX idx_qtype_ts (qtype, timestamp),
  ADD INDEX idx_client_ts (client_ip, timestamp),
  ADD INDEX idx_protocol_ts (protocol, timestamp);

ALTER TABLE dns_tokens
  ADD INDEX idx_status_created (status, created_at),
  ADD INDEX idx_status_last (status, last_seen);

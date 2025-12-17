use dnslog;
CREATE USER IF NOT EXISTS 'dnslog'@'%' IDENTIFIED BY 'dnslog';
GRANT ALL PRIVILEGES ON dnslog.* TO 'dnslog'@'%';
FLUSH PRIVILEGES;
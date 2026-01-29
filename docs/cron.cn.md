# Cron 示例

## 每日凌晨 02:30 备份
```
30 2 * * * /path/to/repo/scripts/backup.sh >> /var/log/dnslog_backup.log 2>&1
```

## 每周日凌晨 03:00 备份
```
0 3 * * 0 /path/to/repo/scripts/backup.sh >> /var/log/dnslog_backup.log 2>&1
```

## 保留策略说明
如果 `retentionEnabled` 设置为 true，则记录清理将在应用程序内部运行。如果禁用此设置，则需要手动安排清理任务。
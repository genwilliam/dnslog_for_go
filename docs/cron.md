# Cron Examples

## Daily backup at 02:30
```
30 2 * * * /path/to/repo/scripts/backup.sh >> /var/log/dnslog_backup.log 2>&1
```

## Weekly backup at Sunday 03:00
```
0 3 * * 0 /path/to/repo/scripts/backup.sh >> /var/log/dnslog_backup.log 2>&1
```

## Retention note
Record cleanup runs in-app when `retentionEnabled` is true. If disabled, schedule a manual cleanup job.

#!/usr/bin/env bash
set -euo pipefail

DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-}"
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_NAME="${DB_NAME:-dnslog}"
BACKUP_DIR="${BACKUP_DIR:-./backups}"

TS="$(date +%Y%m%d%H%M%S)"
FILE="${BACKUP_DIR}/dnslog_${TS}.sql"

mkdir -p "${BACKUP_DIR}"

PASS_ARG=()
if [ -n "${DB_PASS}" ]; then
  PASS_ARG=(-p"${DB_PASS}")
fi

mysqldump -h "${DB_HOST}" -P "${DB_PORT}" -u "${DB_USER}" "${PASS_ARG[@]}" "${DB_NAME}" > "${FILE}"
echo "backup saved: ${FILE}"

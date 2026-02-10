#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "usage: restore.sh /path/to/backup.sql"
  exit 1
fi

SQL_FILE="$1"
if [ ! -f "${SQL_FILE}" ]; then
  echo "file not found: ${SQL_FILE}"
  exit 1
fi

DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-}"
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_NAME="${DB_NAME:-dnslog}"

PASS_ARG=()
if [ -n "${DB_PASS}" ]; then
  PASS_ARG=(-p"${DB_PASS}")
fi

mysql -h "${DB_HOST}" -P "${DB_PORT}" -u "${DB_USER}" "${PASS_ARG[@]}" "${DB_NAME}" < "${SQL_FILE}"
echo "restore complete: ${SQL_FILE}"

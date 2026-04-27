#!/usr/bin/env bash
set -euo pipefail

APP_DIR="${APP_DIR:-/opt/apps/portfolio}"
BACKUP_DIR="${BACKUP_DIR:-/opt/backups/portfolio}"
RETENTION_DAYS="${RETENTION_DAYS:-14}"

mkdir -p "$BACKUP_DIR"

TS=$(date -u +%Y%m%dT%H%M%SZ)
DB_FILE="$BACKUP_DIR/nightly-db-$TS.sql.gz"
UP_FILE="$BACKUP_DIR/nightly-uploads-$TS.tar.gz"

cd "$APP_DIR"

docker compose exec -T postgres pg_dump \
  -U portfolio -d portfolio --clean --if-exists --no-owner --no-privileges \
  | gzip > "$DB_FILE"

DB_SIZE=$(stat -c%s "$DB_FILE")
if [ "$DB_SIZE" -lt 1024 ]; then
  echo "ERROR: db backup suspiciously small: $DB_SIZE bytes" >&2
  exit 1
fi

docker run --rm \
  -v portfolio_uploads:/data/uploads:ro \
  -v "$BACKUP_DIR":/backup \
  alpine:3.19 \
  tar -czf "/backup/$(basename "$UP_FILE")" -C /data uploads

find "$BACKUP_DIR" -name 'nightly-*' -mtime +"$RETENTION_DAYS" -delete

echo "$(date -u) ok  db=$(du -h "$DB_FILE" | cut -f1)  uploads=$(du -h "$UP_FILE" | cut -f1)"

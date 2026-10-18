#!/usr/bin/env zsh

set -Eeuo pipefail

if [[ $# -ne 1 ]]; then
  echo "使用方法: ./scripts/import_daily_prices.sh YYYY-MM-DD" >&2
  exit 1
fi

set -a
source .env
set +a

if [[ -z "${JQUANTS_API_KEY:-}" ]]; then
  echo "エラー: JQUANTS_API_KEYが設定されていません。" >&2
  exit 1
fi

if [[ -z "${MYSQL_DATABASE:-}" ]]; then
  echo "エラー: MYSQL_DATABASEが設定されていません。" >&2
  exit 1
fi

if [[ -z "${MYSQL_USER:-}" ]]; then
  echo "エラー: MYSQL_USERが設定されていません。" >&2
  exit 1
fi

if [[ -z "${MYSQL_PASSWORD:-}" ]]; then
  echo "エラー: MYSQL_PASSWORDが設定されていません。" >&2
  exit 1
fi

cd apps/api
DB_HOST=127.0.0.1 \
  DB_PORT=4306 \
  DB_NAME="${MYSQL_DATABASE}" \
  DB_USER="${MYSQL_USER}" \
  DB_PASSWORD="${MYSQL_PASSWORD}" \
  go run ./cmd/importdailyprices "$1"

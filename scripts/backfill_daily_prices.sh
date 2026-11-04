#!/usr/bin/env bash

set -Eeuo pipefail

if [[ $# -ne 2 ]]; then
  echo "使用方法: ./scripts/backfill_daily_prices.sh FROM_DATE TO_DATE" >&2
  echo "例: ./scripts/backfill_daily_prices.sh 2024-10-19 2026-10-19" >&2
  exit 1
fi

docker compose run --rm backfill_daily_prices "$1" "$2"

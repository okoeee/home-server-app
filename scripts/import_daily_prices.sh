#!/usr/bin/env bash

set -Eeuo pipefail

if [[ $# -ne 1 ]]; then
  echo "使用方法: ./scripts/import_daily_prices.sh YYYY-MM-DD" >&2
  exit 1
fi

docker compose run --rm -T import_daily_prices "$1"

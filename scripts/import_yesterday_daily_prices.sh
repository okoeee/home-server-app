#!/usr/bin/env bash

set -Eeuo pipefail

target_date=$(TZ=Asia/Tokyo date --date=yesterday +%F)

./scripts/import_daily_prices.sh "$target_date"

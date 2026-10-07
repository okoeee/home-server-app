#!/usr/bin/env bash

set -Eeuo pipefail

if [[ "$(git branch --show-current)" != "main" ]]; then
  echo "エラー: mainブランチではありません。" >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "エラー: コミットされていない変更があります。" >&2
  exit 1
fi

echo "GitHubから最新版を取得します..."
git pull --ff-only origin main

echo "Dockerイメージをビルドして起動します..."
docker compose up --build -d

echo "APIの起動を確認します..."
curl \
  --fail \
  --silent \
  --show-error \
  --retry 10 \
  --retry-delay 1 \
  --retry-connrefused \
  http://localhost:8080/ping

echo
echo "デプロイが完了しました。"
docker compose ps

# home-server-app

Raspberry Pi 5で動かす自宅サーバー用アプリケーションです。

現在は、GoとGinで実装したAPIをDocker Composeで起動できます。

## 必要な環境

- Git
- Docker
- Docker Compose

## ディレクトリ構成

```text
.
├── apps/
│   └── api/
│       ├── Dockerfile
│       ├── go.mod
│       ├── go.sum
│       └── main.go
├── scripts/
│   └── deploy.sh
├── compose.yaml
└── README.md
```

## デプロイ

mainブランチへマージした変更は、ラズパイのプロジェクトルートで次のスクリプトを実行して反映します。

```bash
./scripts/deploy.sh
```

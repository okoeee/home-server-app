# home-server-app

Raspberry Pi 5で動かす自宅サーバー用アプリケーションです。

現在は、GoとGinで実装したAPIをDocker Composeで起動できます。

## 必要な環境

- Git
- Go 1.27.1
- Docker
- Docker Compose

## ディレクトリ構成

```text
.
├── .github/
│   └── workflows/
│       └── api-ci.yaml
├── apps/
│   └── api/
│       ├── Dockerfile
│       ├── go.mod
│       ├── go.sum
│       ├── handler/
│       │   └── ping.go
│       ├── main.go
│       └── router/
│           ├── router.go
│           └── router_test.go
├── scripts/
│   └── deploy.sh
├── compose.yaml
└── README.md
```

## テスト

API全体のテストは次のコマンドで実行します。

```bash
cd apps/api
go test ./...
```

## デプロイ

mainブランチへマージした変更は、プロジェクトルートで次のスクリプトを実行して反映します。

```bash
./scripts/deploy.sh
```

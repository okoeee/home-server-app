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
│   ├── api/
│   │   ├── Dockerfile
│   │   ├── db/
│   │   │   └── sqlcによる生成ファイル
│   │   ├── domain/
│   │   │   ├── model/
│   │   │   └── repository/
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── handler/
│   │   ├── query/
│   │   │   ├── daily_prices.sql
│   │   │   └── stocks.sql
│   │   ├── main.go
│   │   ├── router/
│   │   ├── sqlc.yaml
│   │   └── usecase/
│   │       └── stock/
│   └── database/
│       └── mysql/
│           └── init/
│               ├── table_schema.sql
│               └── table_seed.sql
├── documents/
│   └── stock-price-design.md
├── scripts/
│   ├── deploy.sh
│   └── gen_sqlc.sh
├── compose.yaml
└── README.md
```

株価データの設計は、[documents/stock-price-design.md](documents/stock-price-design.md)を参照してください。

## sqlc

`apps/database/mysql/init/table_schema.sql`または`apps/api/query/`内のSQLを変更した場合は、プロジェクトルートで次のコマンドを実行します。

```bash
./scripts/gen_sqlc.sh
```

生成された`apps/api/db/`内のGoファイルはGitで管理し、直接編集しません。

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

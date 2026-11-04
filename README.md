# home-server-app

Raspberry Pi 5で動かす自宅サーバー用アプリケーションです。

現在は、GoとGinで実装したAPIをDocker Composeで起動できます。

フロントエンドはSvelte・TypeScript・Viteで実装しています。

## 必要な環境

- Git
- Make
- Go 1.27.1
- Node.js（フロント開発用。Viteの要件は20.19以上の20系、または22.12以上）
- npm（OpenAPIの生成には10以上）
- Docker
- Docker Compose

## ディレクトリ構成

```text
.
├── .github/
│   └── workflows/
│       └── api-test.yaml
├── apps/
│   ├── api/
│   │   ├── Dockerfile
│   │   ├── cmd/
│   │   │   ├── backfilldailyprices/
│   │   │   ├── importdailyprices/
│   │   │   └── importstocks/
│   │   ├── db/
│   │   │   └── sqlcによる生成ファイル
│   │   ├── domain/
│   │   │   ├── model/
│   │   │   └── repository/
│   │   ├── external/
│   │   │   ├── jquants/
│   │   │   └── openmeteo/
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── handler/
│   │   ├── infrastructure/
│   │   │   └── mysql/
│   │   ├── query/
│   │   │   ├── daily_prices.sql
│   │   │   └── stocks.sql
│   │   ├── main.go
│   │   ├── oapi-codegen.yaml
│   │   ├── openapi/
│   │   │   └── models.gen.go
│   │   ├── router/
│   │   ├── sqlc.yaml
│   │   └── usecase/
│   │       ├── dailyprice/
│   │       └── stock/
│   ├── database/
│   │   └── mysql/
│   │       └── init/
│   │           ├── table_schema.sql
│   │           └── table_seed.sql
│   └── web/
│       ├── .dockerignore
│       ├── Caddyfile
│       ├── Dockerfile
│       ├── package.json
│       ├── package-lock.json
│       ├── public/
│       ├── src/
│       │   └── api/generated/schema.d.ts
│       └── vite.config.ts
├── documents/
│   ├── architecture-decisions.md
│   ├── home-dashboard.md
│   ├── naming-conventions.md
│   ├── server-metrics-adr.md
│   ├── stock-price-design.md
│   ├── testing-policy.md
│   ├── ui-library-comparison.md
│   └── weather-design.md
├── openapi/
│   ├── openapi.yaml
│   ├── package.json
│   ├── package-lock.json
│   ├── paths/
│   ├── responses/
│   └── schemas/
├── scripts/
│   ├── backfill_daily_prices.sh
│   ├── deploy.sh
│   ├── import_daily_prices.sh
│   └── import_stocks.sh
├── compose.yaml
├── Makefile
└── README.md
```

## 開発環境の起動

Dockerを起動し、プロジェクトルートにDB設定を含む`.env`を用意して実行します。

```bash
make
```

sqlc生成、OpenAPI用依存のインストールとGo・TypeScriptモデル生成、Dockerイメージのビルドと`web`・`api`・`db`の起動を順に実行します。途中で失敗した場合は後続の処理を実行しません。Go・Node.js・npmはホスト側にも必要で、初回の依存取得にはネットワーク接続を使用します。

起動後は`http://localhost:3000`で画面を確認できます。既存のDBデータは保持しますが、株価データのインポートは実行しません。停止は`docker compose stop`です。`make dev`も同じ処理を実行します。

APIのコードを変更した場合は、`docker compose up -d --build api`で再ビルド・起動して変更を反映します。

## フロントエンドの開発

プロジェクトルートから次のコマンドを実行します。

```bash
cd apps/web
npm ci
npm run dev
```

ターミナルに表示されたURL（通常は`http://localhost:5173`）をブラウザで開きます。停止は`Ctrl+C`です。APIを利用する場合は、先にプロジェクトルートで`docker compose up -d api`を実行します。Viteは`/api/*`から`/api`を除き、`http://localhost:8080`へ転送します。

`apps/web`で型チェックと本番用ビルドを実行できます。

```bash
npm run check
npm run build
```

フロントのSvelte・TypeScript・CSSとテストの整形は、`apps/web`で実行します。

```bash
npm run format
npm run format:check
```

`format`は整形して保存し、`format:check`は整形済みかを確認します。設定は`.prettierrc.json`で共有し、OpenAPIの生成コードは`.prettierignore`で対象外にしています。エディターのPrettier連携で保存時の整形を有効にすれば、同じ設定を使って自動整形できます。

API通信のテストは`apps/web`で`npm test`を実行します。テストにはNode.js 22.12以上が必要です。

ビルド結果は`apps/web/dist/`に出力されます。現在はnpmを使用し、`package-lock.json`をGitで管理します。`node_modules/`と`dist/`はignore済みです。

## フロントエンドのDocker起動

プロジェクトルートにDB設定を含む`.env`を用意し、ビルド・起動します。

```bash
docker compose up -d --build web
```

依存するAPIとMySQLも起動します。`http://localhost:3000`で画面を確認できます。API転送は次のコマンドで確認します。

```bash
curl http://localhost:3000/api/ping
```

`{"message":"pong"}`が返れば成功です。フロントだけを停止する場合は`docker compose stop web`を実行します。

DockerfileはNode.jsで本番用ファイルをビルドし、生成された`dist/`と`Caddyfile`をCaddyイメージへコピーします。HTTPの80番ポートから配信し、`/api/*`は`/api`を除いて`api:8080`へ転送します。起動コマンドはCaddy公式イメージのものを使用します。APIはMySQLのヘルスチェック成功後に起動します。`web`はAPIの起動後に開始しますが、APIの準備完了までは待ちません。

## ラズパイ本体のメトリクス取得

ラズパイにSSHで接続し、プロジェクトルートで実行します。

```bash
docker compose --profile monitoring up -d node_exporter
curl --fail --silent --show-error http://127.0.0.1:9100/metrics
```

Node Exporterのみを起動します。`monitoring`プロファイルに分けているため、通常の`make`や`docker compose up -d`では起動しません。待ち受けはラズパイ自身の`127.0.0.1:9100`に限定しています。MacのDocker Desktopは本体ではなくLinux VMの情報になるため、取得値はラズパイ実機で確認します。

ログの確認と停止は次のコマンドで行います。

```bash
docker compose logs --tail 50 node_exporter
docker compose --profile monitoring stop node_exporter
```

取得項目と実機での確認方法は[メトリクス取得のADR](documents/server-metrics-adr.md#実機での確認手順)を参照してください。

## OpenAPIとGo・TypeScriptモデルの生成

API仕様は[openapi/openapi.yaml](openapi/openapi.yaml)、画面・設計の詳細は`documents/`で管理します。

`openapi/`内の仕様または`apps/api/oapi-codegen.yaml`を変更した場合は、プロジェクトルートで実行します。

```bash
make generate-openapi
```

OpenAPI用の依存を`npm --prefix openapi ci`でインストールしてから生成します。フロントの依存関係とは分けて管理します。

Redocly CLI 2.71.2で分割した仕様を`openapi/dist/openapi.yaml`へ結合し、次のツールでモデルを生成します。

| 対象 | ツール | 生成先 |
| --- | --- | --- |
| Go | `oapi-codegen v2.8.1`（`go run`） | `apps/api/openapi/models.gen.go` |
| TypeScript | `openapi-typescript 7.13.0` | `apps/web/src/api/generated/schema.d.ts` |

ツールの取得・バージョン解決にネットワーク接続を使用します。生成ファイルはGit管理し、直接編集しません。結合結果の`openapi/dist/`と依存関係の`openapi/node_modules/`はGit管理しません。仕様・設定・ロックファイル・生成コード・利用側の変更は同じPRに含めます。

## sqlc

`apps/database/mysql/init/table_schema.sql`または`apps/api/query/`内のSQLを変更した場合は、プロジェクトルートで次のコマンドを実行します。

```bash
make generate-sqlc
```

生成された`apps/api/db/`内のGoファイルはGitで管理し、直接編集しません。

## 銘柄マスターの取得

MySQLコンテナを起動した状態で、プロジェクトルートから次のコマンドを実行します。

```bash
./scripts/import_stocks.sh
```

スクリプトは`.env`からJ-QuantsのAPIキーとDB設定を読み込み、ホスト側からMySQLの`4306`番ポートへ接続して銘柄マスターを保存します。

## 株価日足の取得と保存

MySQLコンテナを起動し、`.env`にJ-QuantsのAPIキーとDB設定を用意した状態で、プロジェクトルートから取得日を指定して実行します。

```bash
./scripts/import_daily_prices.sh 2026-10-15
```

指定日は`YYYY-MM-DD`形式で入力します。J-Quantsから取得した株価日足を`daily_prices`へUPSERTし、保存件数をログへ出力します。保存処理は1つのトランザクションで実行します。

## 過去の株価日足の取得と保存

MySQLコンテナを起動し、`.env`にJ-QuantsのAPIキーとDB設定を用意した状態で、プロジェクトルートから開始日と終了日を指定して実行します。

```bash
./scripts/backfill_daily_prices.sh 2024-10-19 2026-10-19
```

開始日と終了日は`YYYY-MM-DD`形式で入力します。指定期間内の営業日ごとにJ-Quantsから株価日足を取得し、`daily_prices`へUPSERTします。保存処理は1営業日単位のトランザクションで実行し、日付と保存件数をログへ出力します。同じ期間の再実行も可能です。

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

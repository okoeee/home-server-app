# home-server-app

Raspberry Pi 5で動かす自宅サーバー用アプリケーションです。

現在は、GoとGinで実装したAPIをDocker Composeで起動できます。

フロントエンドはSvelte・TypeScript・Viteで実装し、ローカルの開発サーバーまたはComposeの`web`サービスで騰落率を表示できます。

## 必要な環境

- Git
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
│   │   │   └── jquants/
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
│   ├── naming-conventions.md
│   ├── stock-price-design.md
│   └── testing-policy.md
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
│   ├── gen_openapi.sh
│   ├── gen_sqlc.sh
│   ├── import_daily_prices.sh
│   └── import_stocks.sh
├── compose.yaml
└── README.md
```

株価データの設計は、[documents/stock-price-design.md](documents/stock-price-design.md)を参照してください。

構成と技術選定の理由は、[documents/architecture-decisions.md](documents/architecture-decisions.md)を参照してください。

命名の指針は、[documents/naming-conventions.md](documents/naming-conventions.md)を参照してください。

テストの範囲・命名・書き方は、[documents/testing-policy.md](documents/testing-policy.md)を参照してください。

## フロントエンドの開発

プロジェクトルートから次のコマンドを実行します。

```bash
cd apps/web
npm ci
npm run dev
```

ターミナルに表示されたURL（通常は`http://localhost:5173`）をブラウザで開きます。停止は`Ctrl+C`です。APIを利用する場合は、先にプロジェクトルートで`docker compose up -d api`を実行します。Viteは`/api/*`から`/api`を除き、`http://localhost:8080`へ転送します。

ページを開くと、DBに保存された最新取引日の騰落率を1回のAPI通信で取得します。フロント側で数値として降順に並べ、先頭10件を表示します（10件未満なら全件）。同率の場合は取得順を保ちます。日付入力フォームは設けず、APIが返した対象日と表示件数を併記します。読み込み中・データなし・取得エラーも表示します。比較元は銘柄ごとの直近の保存済み日足です。APIは引き続き全件を返します。

`apps/web`で型チェックと本番用ビルドを実行できます。

```bash
npm run check
npm run build
```

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

DockerfileはNode.jsで本番用ファイルをビルドし、生成された`dist/`と`Caddyfile`をCaddyイメージへコピーします。HTTPの80番ポートから配信し、`/api/*`は`/api`を除いて`api:8080`へ転送します。起動コマンドはCaddy公式イメージのものを使用します。`depends_on`は起動順を指定するもので、APIとDBの準備完了は待ちません。

## OpenAPIとGo・TypeScriptモデルの生成

API仕様の入口は[openapi/openapi.yaml](openapi/openapi.yaml)です。`paths/`にエンドポイント、`schemas/`にモデル、`responses/`に共通レスポンスを分け、`$ref`で参照します。Go APIへ直接接続する場合のパスを定義し、Caddy・Vite経由では先頭に`/api`を付けます。

初回、またはOpenAPI用の依存関係を更新した場合は、プロジェクトルートでインストールします。フロントの依存関係とは分けて管理します。

```bash
npm --prefix openapi ci
```

`openapi/`内の仕様または`apps/api/oapi-codegen.yaml`を変更した場合は、プロジェクトルートで実行します。

```bash
./scripts/gen_openapi.sh
```

Redocly CLI 2.71.2で分割した仕様を`openapi/dist/openapi.yaml`へ結合し、次のツールでモデルを生成します。

| 対象 | ツール | 生成先 |
| --- | --- | --- |
| Go | `oapi-codegen v2.8.1`（`go run`） | `apps/api/openapi/models.gen.go` |
| TypeScript | `openapi-typescript 7.13.0` | `apps/web/src/api/generated/schema.d.ts` |

ツールの取得・バージョン解決にネットワーク接続を使用します。生成ファイルはGit管理し、直接編集しません。結合結果の`openapi/dist/`と依存関係の`openapi/node_modules/`はGit管理しません。仕様・設定・ロックファイル・生成コード・利用側の変更は同じPRに含めます。

TypeScriptは`--root-types --root-types-no-schema-prefix`を指定し、`DailyPriceChangeResponse`などをトップレベルの型エイリアスとして生成します。利用側は`src/api/generated/schema.d.ts`から直接importし、手書きの型エイリアスは設けません。表示用の上位10件の選択処理は引き続き`src/models/daily-price-change.ts`で管理します。

モデル名は仕様全体で一意にし、共通モデルは参照で再利用します。同名で内容が異なるスキーマを結合するときは、`--component-renaming-conflicts-severity=error`により自動改名せず失敗させます。同じYAMLマップ内のキーの重複もパースエラーになります。

生成するのはAPI用モデルだけです。Handler・ルート・通信処理・入力検証は手書きのままにし、Domainやsqlcのモデルとは分けます。日付・数値は既存のJSON表現に合わせて文字列で扱います。最新取引日の`null`と一覧の空配列も維持します。モデル生成だけで実行時の入力・レスポンス検証が追加されるわけではありません。

CIでの再生成差分チェックは、次の段階で追加します。

## sqlc

`apps/database/mysql/init/table_schema.sql`または`apps/api/query/`内のSQLを変更した場合は、プロジェクトルートで次のコマンドを実行します。

```bash
./scripts/gen_sqlc.sh
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

## 騰落率の取得API

データ取得APIのGo側のパスは`/v1`で始まります。Caddy・Vite経由では`/api/v1/...`を使用し、転送時には`/api`だけを除去します。バージョンなしの旧データ取得パスは提供しません。稼働確認用の`/ping`はバージョン外に残します。

### 最新の保存済み取引日

```bash
curl 'http://localhost:3000/api/v1/daily-price-changes/latest'
```

レスポンス例:

```json
{
  "trade_date": "2026-10-19",
  "changes": [
    {
      "symbol": "12340",
      "trade_date": "2026-10-19",
      "close_price": "110.5",
      "previous_trade_date": "2026-10-16",
      "previous_close_price": "100.0",
      "change_rate": "10.5"
    }
  ]
}
```

「最新」は`daily_prices`内の最も新しい取引日であり、今日や市場の最新営業日を保証するものではありません。API内部で対象日を決定し、その日の騰落率を返します。取り込みの遅れを確認できるよう、画面にも対象日を表示します。

DBが空なら`200`と`{"trade_date":null,"changes":[]}`を返します。最新日に比較可能な銘柄がなくても、過去日には戻さず、その日付と空配列を返します。取得・計算失敗は`500`です。

数値は文字列で返します。`change_rate`は、調整済み終値の差を比較元の調整済み終値で割り、100を掛けた百分率です（例: `"10.5"`は10.5%上昇）。比較元は銘柄ごとの対象日より前の直近の保存済み日足であり、必ずしも前営業日ではありません。

APIは対象日と比較元の両方のデータがある全銘柄を返し、並べ替えや上位10銘柄の抽出はフロント側で行います。騰落率の取得は最新取得のみを公開し、日付指定の`/v1/daily-price-changes`は提供しません。

## 出来高倍率の取得API

```bash
curl 'http://localhost:8080/v1/daily-volume-ratios?date=2026-10-19'
```

レスポンス例:

```json
[
  {
    "symbol": "12340",
    "trade_date": "2026-10-19",
    "volume": "200.0",
    "ratio": "2"
  }
]
```

`volume`は対象日の調整済み出来高、`ratio`は対象日より前の直近20件の調整済み出来高の平均に対する倍率です（`"2"`は2倍）。数値は文字列で返します。履歴は銘柄ごとの保存済み日足を使うため、20営業日より長い期間になることがあります。

履歴が20件に満たない銘柄、平均出来高がゼロの銘柄は除外します。対象がなければ`200`と空配列`[]`、日付の未指定・不正は`400`、取得失敗は`500`を返します。並べ替えや件数制限は行いません。

## テスト

API全体のテストは次のコマンドで実行します。

```bash
cd apps/api
go test ./...
```

`apps/api/router`のHTTPテストはUsecaseをstubに置き換え、パス・入力検証・JSON・ステータスコードを確認します。実DBや外部APIは使用せず、計算や処理の組み合わせはDomain・Usecaseのテストで確認します。

## デプロイ

mainブランチへマージした変更は、プロジェクトルートで次のスクリプトを実行して反映します。

```bash
./scripts/deploy.sh
```

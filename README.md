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
│   │   ├── router/
│   │   ├── sqlc.yaml
│   │   └── usecase/
│   │       ├── dailyprice/
│   │       └── stock/
│   └── database/
│       └── mysql/
│           └── init/
│               ├── table_schema.sql
│               └── table_seed.sql
├── documents/
│   ├── naming-conventions.md
│   └── stock-price-design.md
├── scripts/
│   ├── backfill_daily_prices.sh
│   ├── deploy.sh
│   ├── gen_sqlc.sh
│   ├── import_daily_prices.sh
│   └── import_stocks.sh
├── compose.yaml
└── README.md
```

株価データの設計は、[documents/stock-price-design.md](documents/stock-price-design.md)を参照してください。

命名の指針は、[documents/naming-conventions.md](documents/naming-conventions.md)を参照してください。

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

APIの起動後、対象日を`YYYY-MM-DD`形式で指定して取得します。

```bash
curl 'http://localhost:8080/daily-price-changes?date=2026-10-19'
```

レスポンス例:

```json
[
  {
    "symbol": "12340",
    "trade_date": "2026-10-19",
    "close_price": "110.5",
    "previous_trade_date": "2026-10-16",
    "previous_close_price": "100.0",
    "change_rate": "10.5"
  }
]
```

数値は文字列で返します。`change_rate`は、調整済み終値の差を比較元の調整済み終値で割り、100を掛けた百分率です（例: `"10.5"`は10.5%上昇）。比較元は銘柄ごとの対象日より前の直近の保存済み日足であり、必ずしも前営業日ではありません。

対象日と比較元の両方のデータがある全銘柄を返します。並べ替えや上位10銘柄の抽出はまだ行いません。対象がなければ`200`と空配列`[]`、日付の未指定・不正は`400`、取得・計算の失敗は`500`を返します。

## 出来高倍率の取得API

```bash
curl 'http://localhost:8080/daily-volume-ratios?date=2026-10-19'
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

## デプロイ

mainブランチへマージした変更は、プロジェクトルートで次のスクリプトを実行して反映します。

```bash
./scripts/deploy.sh
```

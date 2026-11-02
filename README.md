# home-server-app

Raspberry Pi 5で動かす自宅サーバー用アプリケーションです。

現在は、GoとGinで実装したAPIをDocker Composeで起動できます。

フロントエンドはSvelte・TypeScript・Viteで実装し、ローカルの開発サーバーまたはComposeの`web`サービスで天気予報と騰落率を表示できます。

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

ページを開くと、DBに保存された最新取引日の騰落率を1回のAPI通信で取得します。フロント側で数値として降順に並べ、先頭10件を表示します（10件未満なら全件）。同率の場合は取得順を保ちます。日付入力フォームは設けず、APIが返した対象日と表示件数を併記します。読み込み中・データなし・取得エラーも表示します。比較元は銘柄ごとの直近の保存済み日足です。APIは引き続き全件を返します。

`apps/web`で型チェックと本番用ビルドを実行できます。

```bash
npm run check
npm run build
```

ホームでは天気APIも1回呼び出し、東京都の5日分の天気・最高／最低気温・平均湿度と取得日時を表示します。天気名はenumから日本語へ変換し、欠損値は「—」と表示します。日付はAPIの暦日を保って`10/31（土）`のように曜日付きで表示し、取得日時は日本時間で表示します。天気と騰落率の読み込み・エラー状態は独立しています。

Tailwind CSS 4・shadcn-svelteのCardとLucideを利用し、ダークな天気・株価カードを表示します。天気アイコンは16種類のenumに対応し、日本語名と並べます。共通UIは`src/lib/components/ui/`、機能別コンポーネントは`src/components/`に置きます。`components.json`がshadcn-svelteの設定です。色・角丸・フォントは`src/app.css`で調整し、ダーク表示は`index.html`の`dark`クラスで固定しています。

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

## OpenAPIとGo・TypeScriptモデルの生成

API仕様の入口は[openapi/openapi.yaml](openapi/openapi.yaml)です。`paths/`にエンドポイント、`schemas/`にモデル、`responses/`に共通レスポンスを分け、`$ref`で参照します。Go APIへ直接接続する場合のパスを定義し、Caddy・Vite経由では先頭に`/api`を付けます。

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

TypeScriptは`--root-types --root-types-no-schema-prefix`を指定し、`DailyPriceChangeResponse`などをトップレベルの型エイリアスとして生成します。利用側は`src/api/generated/schema.d.ts`から直接importし、手書きの型エイリアスは設けません。表示用の上位10件の選択処理は引き続き`src/models/daily-price-change.ts`で管理します。

モデル名は仕様全体で一意にし、共通モデルは参照で再利用します。同名で内容が異なるスキーマを結合するときは、`--component-renaming-conflicts-severity=error`により自動改名せず失敗させます。同じYAMLマップ内のキーの重複もパースエラーになります。

生成するのはAPI用モデルだけです。Handler・ルート・通信処理・入力検証は手書きのままにし、Domainやsqlcのモデルとは分けます。日付・数値は既存のJSON表現に合わせて文字列で扱います。最新取引日の`null`と一覧の空配列も維持します。モデル生成だけで実行時の入力・レスポンス検証が追加されるわけではありません。

CIでの再生成差分チェックは、次の段階で追加します。

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

`apps/api/router`のHTTPテストはUsecaseをstubに置き換え、パス・入力検証・JSON・ステータスコードを確認します。実DBや外部APIは使用せず、計算や処理の組み合わせはDomain・Usecaseのテストで確認します。

## デプロイ

mainブランチへマージした変更は、プロジェクトルートで次のスクリプトを実行して反映します。

```bash
./scripts/deploy.sh
```

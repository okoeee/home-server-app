# ADR：ラズパイ本体のメトリクス取得

- 日付：2026-11-03
- 更新：2026-11-17（CPU・メモリに限定した共通構成とmakeでの監視起動）
- 状態：採用。Go APIとホームへの表示を実装済み。Mac・ラズパイの両方でメトリクスAPIの200応答を確認済み。ラズパイではCPU数・メモリ総量もOSの値と一致。復旧後の実機での画面表示は未確認

## 背景

ホームにラズパイ本体のCPU使用率・メモリ使用率を表示したい。最初は情報を取得できることを実機で確認し、その後にAPIと画面を追加する。

監視対象はラズパイ本体全体とする。コンテナごとの使用量、サービスの応答確認、ディスク使用率は今回の対象に含めない。CPU温度も対象から外し、センサーの有無・種類に応じた分岐や欠損値の扱いを省く。

## 決定

Node ExporterをDocker Composeで起動し、本体のメトリクスをHTTPの`/metrics`で公開する。Prometheusが定期収集し、時系列として保存する。

| 構成要素      | 役割                                           | 導入時期   |
|---------------|------------------------------------------------|------------|
| Node Exporter | 本体のCPU・メモリなどを取得して公開する       | 設定済み   |
| Prometheus    | 定期収集・時系列保存・使用率の計算を行う       | 設定済み   |
| Go API        | Prometheusへ問い合わせ、画面用のJSONへ変換する | 実装済み   |
| Svelte        | 自作のサーバーカードに表示する                 | 実装済み   |

問い合わせの流れは以下とする。矢印はリクエストの方向を示す。

```text
定期収集：Prometheus → Node Exporter → LinuxのCPU・メモリ情報
画面表示：Svelte → Go API → Prometheus
```

自作の画面を利用するため、Grafanaは導入しない。メトリクスはPrometheus側で保存し、アプリのMySQLには保存しない。

## 選定理由と比較

| 選択肢                                  | 利点・考慮点                                                                    | 判断                        |
|-----------------------------------------|---------------------------------------------------------------------------------|-----------------------------|
| Goで`/proc`・`/sys`を直接読む           | 常駐ソフトを減らせるが、読み取りやCPU使用率の計算を自作する必要がある           | 採用しない                  |
| Node ExporterをOSへ直接インストールする | ホストの情報を参照しやすいが、OS側に別のインストール・起動管理が必要            | Composeによる管理を優先する |
| Node Exporterをコンテナで起動する       | 既存のComposeで管理でき、収集とアプリの責務を分けられる。CPU・メモリはLinuxの`/proc`から取得する | 採用する                    |

Node ExporterはLinuxホストの情報を公開するためのソフトで、CPU・メモリに対応するcollectorを備える。アプリ側でOS固有の計測処理を持たず、後からPrometheusによる収集につなげられるため選ぶ。[公式README](https://github.com/prometheus/node_exporter#collectors)

Node Exporter自体は履歴を保存しない。CPUの稼働時間は累積値として取得するため、CPU使用率の表示には複数時点の差分が必要になる。Prometheusで収集・保存し、一定期間の使用率の計算を任せる。[Node Exporterによるホスト監視の公式ガイド](https://prometheus.io/docs/guides/node-exporter/)、[Prometheusの概要](https://prometheus.io/docs/introduction/overview/)

## コンテナから本体を参照する方法

CPU・メモリ情報は、コンテナから見えるLinuxの`/proc/stat`・`/proc/meminfo`を参照する。これらはコンテナのCPU・メモリ制限の使用量を表すものではない。ラズパイでは本体、MacのDocker DesktopではLinux VMの値を取得する。

`compose.yaml`の`node_exporter`サービスで、[collectorの指定](https://github.com/prometheus/node_exporter#collectors)を使い、対象をCPU・メモリに限定する。

- `--collector.disable-defaults`で既定のcollectorを無効にし、`--collector.cpu`と`--collector.meminfo`だけを有効にする。
- 既定の`/proc`からCPU・メモリを読み取るため、ホストのbind mountや`pid: host`は追加しない。
- ネットワークはComposeのbridgeを利用する。コンテナ単位やネットワーク・ディスクなどの収集は今回の対象に含めない。
- Mac・ラズパイで同じCompose設定を使う。OS別の上書きファイルは作成しない。
- イメージはバージョンを固定する。`privileged`やDockerソケットのマウントは追加しない。

当初のホスト全体を参照する構成は、CPU・メモリだけの取得には不要だった。Macで利用できない`rslave`などを省き、現在の対象に必要な設定へ絞る。

イメージは`quay.io/prometheus/node-exporter:v1.12.1`を使用する。[採用バージョンのリリース](https://github.com/prometheus/node_exporter/releases/tag/v1.12.1)

Node ExporterとPrometheusをAPIと同じComposeネットワークへ接続する。Prometheusは`node_exporter:9100`から収集し、APIは`PROMETHEUS_URL=http://prometheus:9090`へ問い合わせる。APIをホストネットワークへ移さず、既存の`db:3306`への接続も維持する。APIにはPrometheusへの`depends_on`を追加せず、監視が停止していても他のAPIを起動できるようにする。

各コンテナは`:9090`・`:9100`で待ち受け、ホストへのポート公開だけを`127.0.0.1:9090:9090`・`127.0.0.1:9100:9100`に限定する。同じComposeネットワーク内からはサービス名で接続できるが、LAN向けには公開しない。管理画面は従来どおりSSHトンネルで確認し、Caddyへ直接転送しない。[Dockerのポート公開](https://docs.docker.com/engine/network/port-publishing/)、[Node Exporterの参照先オプション](https://github.com/prometheus/node_exporter/blob/v1.12.1/collector/paths.go)

当初は両方をホストネットワークで起動していたが、ホストのループバックだけに待ち受けるPrometheusへAPIのコンテナから接続できなかったため変更した。収集先の変更に伴い、`instance`ラベルも`127.0.0.1:9100`から`node_exporter:9100`へ変わる。既存の履歴は削除せず、新しい系列のCPUサンプルがたまるまで1分程度待つ。

監視サービスはComposeの`monitoring`プロファイルで管理し、`make`の開発起動で有効にする。画面がメトリクスAPIを呼ぶため、通常の開発起動にもPrometheus・Node Exporterを含める。起動コマンドは`docker compose --profile monitoring up -d --build web prometheus`とし、依存するAPI・DB・Node Exporterも起動する。直接`docker compose up -d`を実行した場合は監視サービスを起動しない。[Composeのプロファイル](https://docs.docker.com/compose/how-tos/profiles/)

`scripts/deploy.sh`でも同じプロファイルを有効にし、リリース時に監視サービスの設定変更も反映する。

## 最初の確認範囲

ラズパイ上でNode Exporterを起動し、次を確認する。

- `/metrics`へアクセスできる。
- CPUの累積稼働時間と、メモリの総量・利用可能量を取得できる。
- 本体の情報であることを、OS側の値と比較して確認する。

MacのDocker DesktopはLinux VM上でコンテナを動かす。そのため、そこで参照できる情報をMac本体やラズパイの情報として扱わず、受け入れ確認はラズパイ実機で行う。[Docker Desktopの実行環境](https://docs.docker.com/desktop/features/vmm/)

アラートは今回の対象に含めない。画面向けAPIの仕様は以下で定義する。

## 画面向けAPIの設計

`GET /v1/server/metrics`でラズパイ1台の現在値を返す。パラメーターは設けず、ブラウザはGo APIだけに接続する。仕様の正本は[OpenAPIのパス定義](../openapi/paths/server-metrics.yaml)と[レスポンス定義](../openapi/schemas/server-metrics.yaml)とする。

| フィールド                | 内容                                                                     |
|---------------------------|--------------------------------------------------------------------------|
| `fetched_at`              | Prometheusへ問い合わせた基準日時（RFC 3339）。元データの収集時刻ではない |
| `cpu_usage_percent`       | 直近1分の全コア平均使用率（0〜100%）                                     |
| `memory_usage_percent`    | 本体のメモリ使用率（0〜100%）。利用可能量を基に計算する                  |

リクエストごとにPrometheusの`GET /api/v1/query`へ問い合わせる。全クエリへ同じ`time`を渡し、その評価時刻を`fetched_at`で返す。Go側のキャッシュやMySQLへの保存は追加しない。15秒周期の収集と履歴保存はPrometheusが担当しているため、まずはその最新値を利用する。[PrometheusのInstant queries](https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries)

### 算出と取得値の検証

対象は現在の収集設定に合わせ、`job="node",instance="node_exporter:9100"`で限定する。複数ホストや変更前の系列を混ぜた平均や合計にしない。

| 対象         | PromQL                                                                                                                                                |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| 収集の成否   | `up{job="node",instance="node_exporter:9100"}` |
| CPU使用率    | `100 * (1 - avg(rate(node_cpu_seconds_total{job="node",instance="node_exporter:9100",mode="idle"}[1m])))` |
| メモリ使用率 | `100 * (1 - node_memory_MemAvailable_bytes{job="node",instance="node_exporter:9100"} / node_memory_MemTotal_bytes{job="node",instance="node_exporter:9100"})` |

CPUは「各コアの直近1分のアイドル時間の割合を平均し、100%から引く」。メモリは「利用可能量を総量で割った割合を100%から引く」。Goでは総量・利用可能量を個別に問い合わせ、範囲の検証後にこのメモリ使用率を計算する。株価の金額と異なり、表示用の近似値なのでJSONの数値とGoの`float64`を使う。APIでは表示桁数への丸めを行わない。[rateの公式仕様](https://prometheus.io/docs/prometheus/latest/querying/functions/#rate)

`up=1`だけでは必須データが有効とは限らない。CPUは直近1分のサンプルから計算できること、メモリは総量が正で利用可能量が総量以下であることを確認する。

ホームで使用状況を確認する用途では厳密な収集時刻の上限は設けない。最新サンプルの選択と系列の失効はPrometheusに任せ、Goでは収集の成否・必要なデータの有無・数値の範囲を確認する。追加の問い合わせと条件分岐を減らすため、`timestamp()`による1分以内のチェックは行わない。[Prometheusの失効処理](https://prometheus.io/docs/prometheus/latest/querying/basics/#staleness)

### 取得できない場合

| HTTPステータス | 条件                                                                                                     |
|----------------|----------------------------------------------------------------------------------------------------------|
| `200`          | CPU・メモリが有効 |
| `502`          | Prometheusへの接続失敗・タイムアウト・クエリエラー、不正なJSON・応答形式・数値。NaNや無限大も返さない    |
| `503`          | Node Exporterが未収集・停止・収集失敗、CPUのサンプル不足、CPU・メモリの系列不足 |
| `500`          | 予期しない内部エラー                                                                                     |

エラーは既存の`ErrorResponse`で返し、詳細はGo APIのログへ記録する。欠損を0%で補完したり、`up=0`の収集失敗時に保存済みの値で成功を返したりしない。`monitoring`を起動していない開発環境でも、株価・天気のAPIへ影響させず、このAPIだけで取得失敗を扱う。

### 実装構成

- `external/prometheus/client.go`：HTTP通信とPrometheusの応答形式の検証。
- `external/prometheus/metrics_reader.go`：同じ評価時刻での取得、収集の成否・データの有無・範囲の検証、内部モデルへの変換。
- `domain/model/server_metrics.go`：画面向けの内部モデル。Prometheusのラベルや応答形式は持たない。
- `handler/server_metrics.go`：Readerの呼び出しとOpenAPIのレスポンス・HTTPステータスへの変換。

キャッシュなどのアプリ固有の処理はないため、呼び出しを中継するだけのUsecaseは追加しない。標準の`net/http`で問い合わせ、新しいライブラリは追加しない。`PROMETHEUS_URL`未指定でGoを直接起動した場合は`http://localhost:9090`へ接続する。複数クエリ全体のタイムアウトは10秒とする。

Go側のテストでは正常系に加え、収集失敗・起動直後のデータ不足・不正な応答を確認する。実機での確認結果は後述する。

### フロントの表示

`api/server-metrics.ts`から`/api/v1/server/metrics`へ問い合わせ、OpenAPIで生成した`ServerMetricsResponse`をそのまま利用する。`App.svelte`が画面表示時に1回取得し、`ServerMetricsCard.svelte`へ結果・読み込み状態・エラーを渡す。自動更新は行わない。

カードは既存の`DashboardCard`を使い、CPU・メモリ使用率を小数1桁の数値とバーで表示する。取得基準日時は日本時間で表示し、元データの収集時刻とは区別する。PCでは天気カードの隣、狭い画面では縦に配置する。読み込み中・取得失敗・データなしを表示し、サーバー情報の取得失敗で天気・株価の表示を妨げない。

## 実機での確認手順

ラズパイへSSHで接続し、プロジェクトルートで[READMEの手順](../README.md#サーバーのメトリクス取得)に従って起動する。Prometheusの依存先としてNode Exporterも起動する。以降もラズパイ上で実行する。

```bash
curl --fail --silent --show-error http://127.0.0.1:9100/metrics \
  | grep -E '^(node_cpu_seconds_total|node_memory_(MemTotal|MemAvailable)_bytes|node_scrape_collector_success)'
```

| メトリクス | 確認する内容 |
| --- | --- |
| `node_cpu_seconds_total` | CPU・動作モードごとの累積秒数。CPU使用率そのものではない |
| `node_memory_MemTotal_bytes` | メモリ総量（バイト） |
| `node_memory_MemAvailable_bytes` | 利用可能なメモリ量（バイト） |
| `node_scrape_collector_success` | 対象collectorの収集成否。成功は1、失敗は0 |

OS側の値は次のように確認する。

```bash
head -n 5 /proc/stat
grep -E '^(MemTotal|MemAvailable):' /proc/meminfo
```

CPUは`/proc/stat`のクロック刻みを秒へ変換した値として、メモリは`/proc/meminfo`のkB表記を1024倍したバイトとして比較する。取得時刻の差で値が動くことも考慮する。

## Prometheusの収集・保存設定

`compose.yaml`の`prometheus`サービスと、`apps/monitoring/prometheus/prometheus.yaml`で管理する。

| 設定 | 値・理由 |
| --- | --- |
| イメージ | `quay.io/prometheus/prometheus:v3.15.1`。バージョンを固定する |
| 収集対象 | `node_exporter:9100`。Node Exporterの`/metrics`を既定のHTTPで取得する |
| ジョブ名 | `node`。収集した系列に`job="node"`が付く |
| 収集間隔 | 15秒。短い期間のCPU使用率を計算でき、秒単位の収集より負荷を抑えられる初期値とする |
| 保存期間 | 7日。まず直近の推移を確認する範囲に留める。実際の使用量を見て調整する |
| 保存先 | 専用のnamed volume `prometheus-data`を`/prometheus`へマウントする |

Prometheusは時系列用DB（TSDB）を内蔵しており、MySQLの接続設定は不要。再起動やコンテナの再作成後も、専用ボリュームの履歴を利用する。7日は削除の基準であり、期限を過ぎたデータはブロック単位で順次削除される。ディスク使用量の厳密な上限ではない。

`depends_on`でNode Exporterの起動を先に行うが、収集可能になるまで待つ設定ではない。収集に失敗しても、Prometheusが次の周期で再試行する。

設定変更はファイルを保存して`docker compose --profile monitoring restart prometheus`で反映する。イメージや起動引数を変更した場合は`docker compose --profile monitoring up -d prometheus`でコンテナを再作成する。

参考：[採用バージョン](https://github.com/prometheus/prometheus/releases/tag/v3.15.1)、[Dockerでの設定と保存先](https://prometheus.io/docs/prometheus/latest/installation/)、[収集設定](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)、[TSDBと保持期間](https://prometheus.io/docs/prometheus/latest/storage/)

## Prometheusでの確認手順

ラズパイ上で起動後、まず準備完了を確認する。準備が完了していなければ、少し待って再実行する。

```bash
curl --fail --silent --show-error http://127.0.0.1:9090/-/ready
```

次に、15秒程度待ってから収集状態を確認する。

```bash
curl --fail --silent --show-error --get \
  --data-urlencode 'query=up{job="node",instance="node_exporter:9100"}' \
  http://127.0.0.1:9090/api/v1/query
```

`data.result`内の`value`が`[取得時刻, "1"]`なら収集成功、`"0"`なら失敗。空配列の場合は、まだ収集されていない可能性がある。失敗時は`http://127.0.0.1:9090/api/v1/targets`の`lastError`とNode Exporterのログを確認する。

`up=1`は収集リクエストの成功を示す。CPU・メモリの個別collectorの成功は、既存の`node_scrape_collector_success`も併せて確認する。

PCのブラウザで確認する場合は、PC側のターミナルでSSHトンネルを開く。接続先は普段SSHに使うユーザー名・ホスト名またはIPアドレスに置き換える。

```bash
ssh -N -L 9090:127.0.0.1:9090 kazuma@home-pi
```

接続を維持した状態でPCの`http://localhost:9090`を開き、以下のPromQLを1つずつ実行する。

| 確認対象 | PromQL |
| --- | --- |
| 収集の成否 | `up{job="node",instance="node_exporter:9100"}` |
| CPU使用率（直近1分、全コア平均、%） | `100 * (1 - avg(rate(node_cpu_seconds_total{job="node",instance="node_exporter:9100",mode="idle"}[1m])))` |
| メモリ使用率（%） | `100 * (1 - node_memory_MemAvailable_bytes{job="node",instance="node_exporter:9100"} / node_memory_MemTotal_bytes{job="node",instance="node_exporter:9100"})` |

CPU使用率は各コアのアイドル時間の割合を求め、その平均を100%から引く。起動直後はサンプル不足で結果が出ないため、1分程度収集してから確認する。メモリ使用率は利用可能量（`MemAvailable`）を引いて求め、単純な空き容量（`MemFree`）は使わない。

## 実装時の検証

Go全体の`go test -race ./...`、フロントの型チェック、OpenAPI生成、Composeの構文検証、公式イメージの`promtool check config`は成功した。

専用の一時Compose環境で、Node Exporter・Prometheus・GoのReaderとRouterを接続した。サービス名による収集と問い合わせで`200`、Node Exporter停止時に`503`、Prometheus停止時に`502`となることを確認した。既存のコンテナ・データは変更せず、一時環境は検証後に削除した。

Macの一時コンテナで、ホストのマウント・PID共有なしにCPU・メモリの収集が成功することを確認した。コンテナを1 CPU・64 MiBに制限しても、取得値はLinux VMの10 CPU・約8.2 GBのメモリ総量だった。

共通設定を反映したMac上で`make`を実行し、sqlc・OpenAPI生成、ビルド、5サービスの起動まで成功した。Prometheusの準備完了と、Caddy経由の`/api/v1/server/metrics`が200でCPU・メモリ使用率を返すことも確認した。

2026-11-17にラズパイの監視コンテナへ現在の設定を反映し、Caddy経由の200応答を確認した。Node Exporterが返す4 CPU・メモリ総量はOS側の値と一致し、既存ボリュームの過去の履歴も保持されていた。リリース時に古い監視コンテナが残った原因と対応は[トラブルシューティング](troubleshooting/api-502-after-deploy.md)に記録する。

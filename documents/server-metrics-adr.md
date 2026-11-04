# ADR：ラズパイ本体のメトリクス取得

- 日付：2026-11-03
- 更新：2026-11-04（Prometheusの追加）
- 状態：採用。Node Exporter・Prometheusの設定を追加済み。Prometheusによるラズパイ実機の収集確認は未実施

## 背景

ホームにラズパイ本体のCPU使用率・メモリ使用率・温度を表示したい。最初は情報を取得できることを実機で確認し、その後にAPIと画面を追加する。

監視対象はラズパイ本体全体とする。コンテナごとの使用量、サービスの応答確認、ディスク使用率は今回の対象に含めない。

## 決定

Node ExporterをDocker Composeで起動し、本体のメトリクスをHTTPの`/metrics`で公開する。Prometheusが定期収集し、時系列として保存する。

| 構成要素      | 役割                                           | 導入時期   |
|---------------|------------------------------------------------|------------|
| Node Exporter | 本体のCPU・メモリ・温度などを取得して公開する  | 設定済み   |
| Prometheus    | 定期収集・時系列保存・使用率の計算を行う       | 設定済み   |
| Go API        | Prometheusへ問い合わせ、画面用のJSONへ変換する | 後続の実装 |
| Svelte        | 自作のサーバーカードに表示する                 | 後続の実装 |

将来の問い合わせの流れは以下とする。矢印はリクエストの方向を示す。

```text
定期収集：Prometheus → Node Exporter → ラズパイ本体の情報
画面表示：Svelte → Go API → Prometheus
```

自作の画面を利用するため、Grafanaは導入しない。メトリクスはPrometheus側で保存し、アプリのMySQLには保存しない。Go APIとSvelteへの接続は後続で実装する。

## 選定理由と比較

| 選択肢                                  | 利点・考慮点                                                                    | 判断                        |
|-----------------------------------------|---------------------------------------------------------------------------------|-----------------------------|
| Goで`/proc`・`/sys`を直接読む           | 常駐ソフトを減らせるが、読み取りやCPU使用率の計算を自作する必要がある           | 採用しない                  |
| Node ExporterをOSへ直接インストールする | ホストの情報を参照しやすいが、OS側に別のインストール・起動管理が必要            | Composeによる管理を優先する |
| Node Exporterをコンテナで起動する       | 既存のComposeで管理でき、収集とアプリの責務を分けられる。ホスト参照の設定が必要 | 採用する                    |

Node ExporterはLinuxホストの情報を公開するためのソフトで、CPU・メモリ・温度に対応するcollectorを備える。アプリ側でOS固有の計測処理を持たず、後からPrometheusによる収集につなげられるため選ぶ。[公式README](https://github.com/prometheus/node_exporter#collectors)

Node Exporter自体は履歴を保存しない。CPUの稼働時間は累積値として取得するため、CPU使用率の表示には複数時点の差分が必要になる。Prometheusで収集・保存し、一定期間の使用率の計算を任せる。[Node Exporterによるホスト監視の公式ガイド](https://prometheus.io/docs/guides/node-exporter/)、[Prometheusの概要](https://prometheus.io/docs/introduction/overview/)

## コンテナから本体を参照する方法

Linux上のNode Exporterは、主に`/proc`のCPU・メモリ情報と、`/sys`の温度センサー情報を利用する。コンテナ化する場合は、コンテナ自身と本体の情報を混同しないようにする。

`compose.yaml`の`node_exporter`サービスで、[公式のDocker構成](https://github.com/prometheus/node_exporter#docker)を基に次を設定する。

- ホストのルートを`/host`へ読み取り専用でbind mountし、`--path.rootfs=/host`を指定する。
- `network_mode: host`と`pid: host`を利用し、ホストの名前空間を参照する。
- ホスト参照の設定はNode Exporterに限定し、Go APIのコンテナへ追加しない。
- ホストの情報を広く参照するため、イメージはバージョンを固定する。`privileged`やDockerソケットのマウントは追加しない。

イメージは`quay.io/prometheus/node-exporter:v1.12.1`を使用する。[採用バージョンのリリース](https://github.com/prometheus/node_exporter/releases/tag/v1.12.1)

Node Exporterは`--web.listen-address=127.0.0.1:9100`でラズパイ自身からのみ接続可能にする。Prometheusも`network_mode: host`を指定し、同じホストの`127.0.0.1:9100`へ接続する。通常のComposeネットワークのサービス名では接続しない。

Prometheusの待ち受けは`127.0.0.1:9090`とし、管理画面はSSHトンネル経由で確認する。Caddyへの転送は追加しない。Go APIは通常のコンテナネットワークにあるため、現状ではPrometheusへ接続できない。API連携時に接続経路と待ち受け先を設計する。

通常の開発起動へ影響させないため、Composeの`monitoring`プロファイルで明示的に起動する設計とする。既存の`make`による開発環境起動には含めない。[Composeのプロファイル](https://docs.docker.com/compose/how-tos/profiles/)

## 最初の確認範囲

ラズパイ上でNode Exporterを起動し、次を確認する。

- `/metrics`へアクセスできる。
- CPUの累積稼働時間と、メモリの総量・利用可能量を取得できる。
- 温度センサーのメトリクスを取得でき、どのセンサーがCPU温度か識別できる。
- 本体の情報であることを、OS側の値と比較して確認する。

温度はOS・カーネル・センサーの公開状況によって取得できる項目が異なる。実機でcollectorとラベルを確認してから採用する値を決める。欠損を0℃として扱わない。

MacのDocker DesktopはLinux VM上でコンテナを動かす。そのため、そこで参照できる情報をMac本体やラズパイの情報として扱わず、温度を含む受け入れ確認はラズパイ実機で行う。[Docker Desktopの実行環境](https://docs.docker.com/desktop/features/vmm/)

この段階ではGo API・OpenAPI・自作の画面・アラートを追加しない。

## 実機での確認手順

ラズパイへSSHで接続し、プロジェクトルートで[READMEの手順](../README.md#ラズパイ本体のメトリクス取得)に従って起動する。Prometheusの依存先としてNode Exporterも起動する。以降もラズパイ上で実行する。

```bash
curl --fail --silent --show-error http://127.0.0.1:9100/metrics \
  | grep -E '^(node_cpu_seconds_total|node_memory_(MemTotal|MemAvailable)_bytes|node_hwmon_(temp_celsius|chip_names|sensor_label)|node_thermal_zone_temp|node_scrape_collector_success)'
```

| メトリクス | 確認する内容 |
| --- | --- |
| `node_cpu_seconds_total` | CPU・動作モードごとの累積秒数。CPU使用率そのものではない |
| `node_memory_MemTotal_bytes` | メモリ総量（バイト） |
| `node_memory_MemAvailable_bytes` | 利用可能なメモリ量（バイト） |
| `node_hwmon_temp_celsius` / `node_thermal_zone_temp` | 温度（℃）。公開される系列とセンサーのラベルは実機で確認する |
| `node_scrape_collector_success` | 対象collectorの収集成否。成功は1、失敗は0 |

OS側の値は次のように確認する。

```bash
head -n 5 /proc/stat
grep -E '^(MemTotal|MemAvailable):' /proc/meminfo
grep -H . /sys/class/thermal/thermal_zone*/type /sys/class/thermal/thermal_zone*/temp
```

CPUは`/proc/stat`のクロック刻みを秒へ変換した値として、メモリは`/proc/meminfo`のkB表記を1024倍したバイトとして比較する。`/sys/class/thermal`の温度はミリ℃表記なので、1000で割って比較する。取得時刻の差で値が動くことも考慮する。

温度系列がない場合は、Node ExporterのログとOS側のセンサー公開状況を確認する。複数センサーがある場合は種類・ラベルを調べ、CPU温度に対応する系列を特定する。

## Prometheusの収集・保存設定

`compose.yaml`の`prometheus`サービスと、`apps/monitoring/prometheus/prometheus.yaml`で管理する。

| 設定 | 値・理由 |
| --- | --- |
| イメージ | `quay.io/prometheus/prometheus:v3.15.1`。バージョンを固定する |
| 収集対象 | `127.0.0.1:9100`。Node Exporterの`/metrics`を既定のHTTPで取得する |
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
  --data-urlencode 'query=up{job="node"}' \
  http://127.0.0.1:9090/api/v1/query
```

`data.result`内の`value`が`[取得時刻, "1"]`なら収集成功、`"0"`なら失敗。空配列の場合は、まだ収集されていない可能性がある。失敗時は`http://127.0.0.1:9090/api/v1/targets`の`lastError`とNode Exporterのログを確認する。

`up=1`は収集リクエストの成功を示す。温度など個別collectorの成功は、既存の`node_scrape_collector_success`も併せて確認する。

PCのブラウザで確認する場合は、PC側のターミナルでSSHトンネルを開く。接続先は普段SSHに使うユーザー名・ホスト名またはIPアドレスに置き換える。

```bash
ssh -N -L 9090:127.0.0.1:9090 kazuma@home-pi
```

接続を維持した状態でPCの`http://localhost:9090`を開き、以下のPromQLを1つずつ実行する。

| 確認対象 | PromQL |
| --- | --- |
| 収集の成否 | `up{job="node"}` |
| CPU使用率（直近1分、全コア平均、%） | `100 * (1 - avg(rate(node_cpu_seconds_total{job="node",mode="idle"}[1m])))` |
| メモリ使用率（%） | `100 * (1 - node_memory_MemAvailable_bytes{job="node"} / node_memory_MemTotal_bytes{job="node"})` |
| 温度（hwmon、℃） | `node_hwmon_temp_celsius{job="node"}` |
| 温度（thermal zone、℃） | `node_thermal_zone_temp{job="node"}` |

CPU使用率は各コアのアイドル時間の割合を求め、その平均を100%から引く。起動直後はサンプル不足で結果が出ないため、1分程度収集してから確認する。メモリ使用率は利用可能量（`MemAvailable`）を引いて求め、単純な空き容量（`MemFree`）は使わない。

温度の2つのクエリは、OSが公開する系列を調べるためのもの。両方の取得を必須とはせず、実機でCPUに対応するセンサーを特定する。

## 実装時の検証

Composeの構文・プロファイルと、公式イメージに同梱された`promtool check config`による設定検証は成功した。隔離した一時コンテナで、本番と同じ収集設定・Prometheusの起動引数を使い、Node Exporterとのループバック経由の収集、CPU・メモリ使用率の計算、再起動後の履歴保持を確認した。

隔離環境では2つのコンテナが同じネットワーク名前空間を共有する構成で検証した。ラズパイのホストネットワークと温度センサーについては、実機での確認を別途行う。

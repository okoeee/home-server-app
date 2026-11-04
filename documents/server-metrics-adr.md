# ADR：ラズパイ本体のメトリクス取得

- 日付：2026-11-03
- 状態：採用。Compose設定を追加済み。ラズパイ実機での取得確認は未実施

## 背景

ホームにラズパイ本体のCPU使用率・メモリ使用率・温度を表示したい。最初は情報を取得できることを実機で確認し、その後にAPIと画面を追加する。

監視対象はラズパイ本体全体とする。コンテナごとの使用量、サービスの応答確認、ディスク使用率は今回の対象に含めない。

## 決定

Node ExporterをDocker Composeで起動し、本体のメトリクスをHTTPの`/metrics`で取得する。最初はNode Exporterのみを追加する。

| 構成要素      | 役割                                           | 導入時期   |
|---------------|------------------------------------------------|------------|
| Node Exporter | 本体のCPU・メモリ・温度などを取得して公開する  | 最初に導入 |
| Prometheus    | 定期収集・時系列保存・使用率の計算を行う       | 次の段階   |
| Go API        | Prometheusへ問い合わせ、画面用のJSONへ変換する | 後続の実装 |
| Svelte        | 自作のサーバーカードに表示する                 | 後続の実装 |

将来の問い合わせの流れは以下とする。矢印はリクエストの方向を示す。

```text
定期収集：Prometheus → Node Exporter → ラズパイ本体の情報
画面表示：Svelte → Go API → Prometheus
```

自作の画面を利用するため、Grafanaは導入しない。メトリクスはPrometheus側で保存し、アプリのMySQLには保存しない方針とする。保持期間や収集間隔はPrometheus導入時に決める。

## 選定理由と比較

| 選択肢                                  | 利点・考慮点                                                                    | 判断                        |
|-----------------------------------------|---------------------------------------------------------------------------------|-----------------------------|
| Goで`/proc`・`/sys`を直接読む           | 常駐ソフトを減らせるが、読み取りやCPU使用率の計算を自作する必要がある           | 採用しない                  |
| Node ExporterをOSへ直接インストールする | ホストの情報を参照しやすいが、OS側に別のインストール・起動管理が必要            | Composeによる管理を優先する |
| Node Exporterをコンテナで起動する       | 既存のComposeで管理でき、収集とアプリの責務を分けられる。ホスト参照の設定が必要 | 採用する                    |

Node ExporterはLinuxホストの情報を公開するためのソフトで、CPU・メモリ・温度に対応するcollectorを備える。アプリ側でOS固有の計測処理を持たず、後からPrometheusによる収集につなげられるため選ぶ。[公式README](https://github.com/prometheus/node_exporter#collectors)

Node Exporter自体は履歴を保存しない。CPUの稼働時間は累積値として取得するため、CPU使用率の表示には複数時点の差分が必要になる。後続でPrometheusを追加し、一定期間の使用率の計算を任せる。[Node Exporterによるホスト監視の公式ガイド](https://prometheus.io/docs/guides/node-exporter/)、[Prometheusの概要](https://prometheus.io/docs/introduction/overview/)

## コンテナから本体を参照する方法

Linux上のNode Exporterは、主に`/proc`のCPU・メモリ情報と、`/sys`の温度センサー情報を利用する。コンテナ化する場合は、コンテナ自身と本体の情報を混同しないようにする。

`compose.yaml`の`node_exporter`サービスで、[公式のDocker構成](https://github.com/prometheus/node_exporter#docker)を基に次を設定する。

- ホストのルートを`/host`へ読み取り専用でbind mountし、`--path.rootfs=/host`を指定する。
- `network_mode: host`と`pid: host`を利用し、ホストの名前空間を参照する。
- ホスト参照の設定はNode Exporterに限定し、Go APIのコンテナへ追加しない。
- ホストの情報を広く参照するため、イメージはバージョンを固定する。`privileged`やDockerソケットのマウントは追加しない。

イメージは`quay.io/prometheus/node-exporter:v1.12.1`を使用する。[採用バージョンのリリース](https://github.com/prometheus/node_exporter/releases/tag/v1.12.1)

最初は`--web.listen-address=127.0.0.1:9100`でラズパイ自身からのみ接続可能にする。SSHでラズパイへ接続して`http://127.0.0.1:9100/metrics`を確認し、Caddyやブラウザへの公開は行わない。Prometheus導入時にコンテナ間の接続方法と待ち受け先を決める。

通常の開発起動へ影響させないため、Composeの`monitoring`プロファイルで明示的に起動する設計とする。既存の`make`による開発環境起動には含めない。[Composeのプロファイル](https://docs.docker.com/compose/how-tos/profiles/)

## 最初の確認範囲

ラズパイ上でNode Exporterを起動し、次を確認する。

- `/metrics`へアクセスできる。
- CPUの累積稼働時間と、メモリの総量・利用可能量を取得できる。
- 温度センサーのメトリクスを取得でき、どのセンサーがCPU温度か識別できる。
- 本体の情報であることを、OS側の値と比較して確認する。

温度はOS・カーネル・センサーの公開状況によって取得できる項目が異なる。実機でcollectorとラベルを確認してから採用する値を決める。欠損を0℃として扱わない。

MacのDocker DesktopはLinux VM上でコンテナを動かす。そのため、そこで参照できる情報をMac本体やラズパイの情報として扱わず、温度を含む受け入れ確認はラズパイ実機で行う。[Docker Desktopの実行環境](https://docs.docker.com/desktop/features/vmm/)

この段階ではPrometheus・Go API・OpenAPI・画面・アラートを追加しない。

## 実機での確認手順

ラズパイへSSHで接続し、プロジェクトルートで[READMEの手順](../README.md#ラズパイ本体のメトリクス取得)に従って起動する。以降もラズパイ上で実行する。

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

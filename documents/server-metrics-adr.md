# ADR：ラズパイ本体のメトリクス取得

- 日付：2026-11-03
- 状態：採用方針を決定。未実装

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

実装時は[公式のDocker構成](https://github.com/prometheus/node_exporter#docker)を基に、次を設定する。

- ホストのルートを`/host`へ読み取り専用でbind mountし、`--path.rootfs=/host`を指定する。
- `network_mode: host`と`pid: host`を利用し、ホストの名前空間を参照する。
- ホスト参照の設定はNode Exporterに限定し、Go APIのコンテナへ追加しない。
- ホストの情報を広く参照するため、イメージはバージョンを固定する。`privileged`やDockerソケットのマウントは追加しない。

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

この段階ではPrometheus・Go API・OpenAPI・画面・アラートを追加しない。起動手順は実装時に整備する。

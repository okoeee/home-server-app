# リリース後に天気・メトリクスAPIが502になる

確認日：2026-11-17。対象はRaspberry Pi 5上のDocker Compose環境。

## 問題1：メトリクスAPIが502になる

APIログには次のエラーがあった。

```text
get server metrics: prometheus query failed: ... lookup prometheus on 127.0.0.11:53: no such host
```

### 原因

**リリース時に監視サービスの設定変更が反映されていなかった。**

APIは現在のComposeネットワークで動いていたが、PrometheusとNode Exporterは13日前の`network_mode: host`のままだった。Prometheusが起動していても、APIからComposeのサービス名`prometheus`では接続できない。

`compose.yaml`は更新済みだったが、`scripts/deploy.sh`の`docker compose up --build -d`は、`monitoring`プロファイルのサービスを対象に含めていなかった。そのため、監視コンテナだけが古い設定で動き続けた。

### 解決策

ラズパイのプロジェクトルートで、監視サービスへ現在の設定を反映する。

```bash
docker compose --profile monitoring up -d prometheus
```

- `--profile monitoring`：監視サービスを起動対象に含める。
- `up -d prometheus`：Prometheusと依存先のNode Exporterをバックグラウンドで起動する。設定が変わったコンテナは再作成する。

DB・APIの再作成や、保存用ボリュームの削除は不要。Prometheusの履歴は既存の`prometheus-data`を利用する。

起動後に準備完了とAPIの応答を確認する。

```bash
curl --fail --silent --show-error http://127.0.0.1:9090/-/ready
curl -i --max-time 15 http://127.0.0.1:3000/api/v1/server/metrics
```

CPUのサンプルがたまるまでは503になる場合があるため、1分程度待って再確認する。

実機では監視2サービスの再作成後に200応答を確認した。CPU数・メモリ総量はOSの値と一致し、再作成前の系列の履歴も保持されていた。

### 再発防止

`scripts/deploy.sh`でも`--profile monitoring`を指定し、アプリと監視サービスを一緒に更新する。`make`の開発起動も同じプロファイルを有効にする。

## 問題2：天気APIも502になった

### 原因

APIログにOpen-Meteoへの通信の`context deadline exceeded`と`context canceled`が記録されていた。現在の実装では、通信が10秒でタイムアウトした場合やリクエストがキャンセルされた場合も502になる。

**タイムアウトの根本原因は未特定。** メトリクスのサービス名解決エラーとは別の問題として扱う。

### 現在の状態と調査方法

調査時は、天気APIに対する修正や再起動を行わずに200応答と5日分の予報取得を確認できた。監視コンテナの再作成が天気を復旧させたとは判断しない。

再発時は、再起動する前に失敗時点のログを確認する。

```bash
docker compose logs --since 10m --tail 50 api
```

`get weather forecast:`のエラーから、名前解決・TLS・タイムアウト・外部APIのHTTPエラー・応答データの検証失敗を切り分ける。

### 通信区間のログを読む

通信計測を含むAPIをリリースした後は、Open-Meteoの取得に失敗したログに`network_trace`が付く。反映にはAPIイメージの再ビルド・コンテナの更新が必要。以前の失敗をさかのぼって計測することはできない。

例えば、次のログならTCP接続までは完了し、TLSの開始後に待ち時間が発生している。

```text
network_trace={total=10000ms events=[get_conn(api.open-meteo.com:443)@0ms, dns_start(api.open-meteo.com)@0ms, dns_done(ok)@7ms, connect_start(tcp 192.0.2.1:443)@7ms, connect_done(tcp 192.0.2.1:443)(ok)@137ms, tls_start@137ms]}
```

これは読み方を示す例であり、実機で観測したログではない。`@...ms`は取得開始からの経過時間。上の例はDNSに約7ms、TCP接続に約130ms、その後のTLSで失敗時点まで約9.9秒待ったと読める。

| 記録されたイベント | 切り分けの目安 |
| --- | --- |
| `dns_start`の後に`dns_done`がない | 名前解決の完了を待っていた可能性がある |
| `connect_start`の後に対応する`connect_done`がない | その接続先へのTCP接続が完了していない |
| `tls_start`の後に`tls_done`がない | TLSハンドシェイクが完了していない |
| `request_written(ok)`の後に`first_response_byte`がない | 送信後、応答ヘッダーの最初の1バイトを待っている |
| `first_response_byte`はあるが`headers_received`がない | 応答ヘッダー全体の受信が完了していない |
| `headers_received`・`body_read_start`があり、`decode forecast response`で失敗 | 本文の読み取りやJSON解析で失敗。タイムアウトか不正なJSONかは元のエラーも確認する |
| `headers_received(status=429)`など | 外部APIからHTTPエラーが返った |

`got_conn(reused=true)`なら既存の接続を使っているため、DNS・TCP・TLSのイベントがなくても異常ではない。複数の接続先へ並行して接続することもあるため、開始・完了は接続先アドレスも見て対応付ける。

ログは失敗時点の記録なので、その後に完了するコールバックは含まれない。この計測で遅延した区間は絞れるが、自宅回線・経路・Open-Meteoのどこに根本原因があるかは、それだけでは断定できない。現時点では根本原因は未特定で、タイムアウト延長や自動リトライは追加していない。

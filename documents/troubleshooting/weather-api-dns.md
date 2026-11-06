# GitHubに接続できない・天気APIが502になる

発生・復旧確認：2026-11-05〜06。対象はRaspberry Pi 5上のDocker Compose環境。

**GitHubへ接続できない原因は、ホストの名前解決の失敗だった。天気APIの502は復旧したが、直接の原因は特定できていない。** この2つを分けて記録する。

以下のコマンドは、ラズパイへSSH接続し、プロジェクトルートで実行する。

## 問題1：git pullが名前解決エラーで失敗する

デプロイ時にGitHubへ接続できず、次のエラーになった。

```text
ssh: Could not resolve hostname github.com: Temporary failure in name resolution
```

### 原因

ラズパイが`github.com`を接続先のIPアドレスに変換できていなかった。これはDNSによる名前解決の失敗で、SSHの認証より前に発生している。

判断の根拠は、外部IPへの`ping -c 3 1.1.1.1`は成功し、`getent hosts github.com`と`getent hosts api.open-meteo.com`はどちらも結果が表示されなかったこと。IP宛ての通信はできていたが、名前から接続先を調べる処理が失敗していた。

利用していたDNSは、`/etc/resolv.conf`の`nameserver 192.168.11.1`で指定された自宅ルーターだった。ただし、**ルーターのDNS機能と、その先のDNSサーバーのどちらが不調だったか、なぜ突然発生したかは分かっていない。**

### 解決策：DNSの問い合わせ先を一時的に変更する

ルーターをDNSの問い合わせ先として使う経路を避け、CloudflareのDNSへ変更したところ、ホストの名前解決が復旧し、`git pull`も成功した。

```bash
sudo nmcli device modify wlan0 \
  ipv4.ignore-auto-dns yes \
  ipv4.dns "1.1.1.1 1.0.0.1"
```

| 指定                         | 意味                                             |
|------------------------------|--------------------------------------------------|
| `device modify wlan0`        | Wi-Fi用デバイス`wlan0`の現在の設定を一時変更する |
| `ipv4.ignore-auto-dns yes`   | DHCPなどで自動配布されたIPv4のDNSを使わない      |
| `ipv4.dns "1.1.1.1 1.0.0.1"` | CloudflareのDNSを問い合わせ先として指定する      |

`wlan0`はWi-Fiの接続口、DNSは名前を調べる問い合わせ先。Wi-Fi接続やインターネットへの通信経路を切り替えたわけではなく、通信は引き続き自宅ルーターを通る。

変更後は、次のコマンドでIPアドレスが表示されることを確認する。

```bash
getent hosts github.com
getent hosts api.open-meteo.com
```

`getent hosts`はOSの名前解決を使って接続先のIPアドレスを調べる。この変更は保存済み接続プロファイルには残らず、DNSを恒久的に変更するものではない。

## 問題2：ホストの名前解決が直っても、天気APIが502を返す

Go APIへ直接アクセスしても、`502`と`{"message":"天気予報を取得できませんでした"}`が返った。Caddyを通さずに同じ応答になったため、Go API側の予報取得・検証で失敗していることが分かった。

### 原因

**502の直接原因は不明。失敗した時点では、元のエラーをAPIログへ出していなかったため特定できなかった。**

ホストでDNSを変更した後もAPIが失敗していたため、起動済みコンテナが以前のDNS設定を使い続けていた可能性はある。ただし、コンテナ内の名前解決エラーを記録できておらず、この仮説を裏付ける証拠はない。タイムアウトや外部APIの応答・データ検証の失敗も、この実装では502になる。

ホストで名前解決が成功したことだけでは、コンテナ内の外部通信が正常かは判断できない。

### 解決策：APIコンテナを再作成する

ホストの名前解決を復旧したうえで、最新コードをビルドし、APIコンテナを再作成したところ、天気APIが200を返し、5日分の予報を取得できた。

```bash
git pull --ff-only
docker compose up -d --no-deps --build --force-recreate api
curl -i --max-time 15 http://127.0.0.1:8080/v1/weather/forecast
```

| 指定                    | 意味                                                                                   |
|-------------------------|----------------------------------------------------------------------------------------|
| `git pull --ff-only`    | 現在のブランチの追跡先から更新する。マージコミットが必要なら停止する                   |
| `up -d`                 | コンテナを作成し、バックグラウンドで起動する                                           |
| `--no-deps`             | 依存するDBなどを起動対象に含めず、APIだけを対象にする。MySQLは起動済みであることが前提 |
| `--build`               | ソースの変更をイメージへ反映する。今回は診断用ログの追加を反映した                     |
| `--force-recreate`      | 変更の検出結果にかかわらずコンテナを作り直し、実行環境を取り直す                       |
| `curl -i --max-time 15` | HTTPステータス・ヘッダー・本文を表示し、最大15秒で打ち切る                             |

**これは今回復旧した対処であり、502全般に対する解決策ではない。** DNS変更とコンテナ再作成を行っているため、何がAPI復旧に決定的だったかは断定できない。Goの天気取得ロジック自体は変更していない。

### 再発した場合：再作成する前にエラーの詳細を確認する

今回、取得・検証のエラーを`get weather forecast:`付きでAPIログへ出すよう変更した。次回は、再作成などの対処を行う前に詳細を確認し、その原因に応じて対応する。

```bash
docker compose logs --tail 50 api
```

`logs`はコンテナのログを表示し、`--tail 50`は末尾50行に絞る。`get weather forecast:`から始まる行で、DNS・タイムアウト・TLS・外部HTTPステータス・データ検証など、どこで失敗したかを確認する。画面には内部エラーの詳細を返さない。

## 一時変更したDNSを元に戻す

保存済み接続プロファイルは変更していないため、元の設定を再適用する。

```bash
sudo nmcli device reapply wlan0
cat /etc/resolv.conf
getent hosts github.com
getent hosts api.open-meteo.com
```

`device reapply`は接続を切断せずに保存済みプロファイルを適用する。`cat /etc/resolv.conf`で問い合わせ先が元の`192.168.11.1`へ戻ったことを確認し、`getent hosts`で名前解決も成功することを確認する。

名前解決が成功したら、APIコンテナも再作成して確認する。コード変更はないため、`--build`は不要。

```bash
docker compose up -d --no-deps --force-recreate api
curl -i --max-time 15 http://127.0.0.1:8080/v1/weather/forecast
```

元のDNSに戻した後の動作は、この記録時点では未確認。再び名前解決に失敗する場合は、ルーターや上流DNSの調査が必要になる。

## 参考資料

- [NetworkManager：一時変更と再適用](https://networkmanager.dev/docs/api/latest/nmcli.html)
- [NetworkManager：Reapplyと保存済みプロファイル](https://networkmanager.dev/docs/api/latest/gdbus-org.freedesktop.NetworkManager.Device.html)
- [Cloudflare：DNSサーバーのIPアドレス](https://developers.cloudflare.com/1.1.1.1/ip-addresses/)
- [Docker：コンテナのDNS](https://docs.docker.com/engine/network/#dns-services)
- [Docker Compose：upのオプション](https://docs.docker.com/reference/cli/docker/compose/up/)

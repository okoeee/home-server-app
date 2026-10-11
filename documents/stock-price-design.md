# 株価データ設計

## 目的

国内株式の日足を保存し、新高値などの分析に利用する。

## データ取得元

[J-Quants API](https://jpx-jquants.com/) v2の株価四本値を使用する。

J-Quants APIから、証券コード、取引日、始値、高値、安値、終値、出来高を取得できる。株式分割などの影響を避けるため、価格と出来高は調整済みの値を保存する。

## データの流れ

```text
J-Quants API → Go → MySQL → Go API → クライアント
```

## テーブル

`daily_prices`に1銘柄・1取引日につき1件の日足を保存する。

| カラム | 内容 |
| --- | --- |
| `symbol` | 証券コード |
| `trade_date` | 取引日 |
| `open_price` | 調整済み始値 |
| `high_price` | 調整済み高値 |
| `low_price` | 調整済み安値 |
| `close_price` | 調整済み終値 |
| `volume` | 調整済み出来高 |

`symbol`と`trade_date`を複合主キーとし、同じ銘柄・取引日の重複を防ぐ。

テーブル定義は`apps/database/mysql/init/table_schema.sql`で管理する。

## 最初に作成するAPI

```text
GET /stocks/:symbol/prices
```

指定した証券コードの日足を取引日の昇順で返す。


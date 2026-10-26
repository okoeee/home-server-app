# 命名ガイド

このプロジェクトで名前を決めるときの指針。役割や意味を優先し、大文字・小文字、略語、ファイル名の形式は各言語の慣習と既存の構成に合わせる。

## 基本方針

- 名前から対象と役割が分かるようにする。同じ意味には同じ単語を使う。
- 利用側から見た役割で名前を付ける。内部で使うDBやライブラリの名前は、実装を区別するときに使う。
- 文脈から分かる情報は繰り返さず、意味が曖昧になるほど省略しない。
- `Manager`、`Helper`、`Util`などは、具体的な役割を表せないか先に検討する。
- 既存の名前の変更は、関連する機能の変更や合意したリファクタリングの範囲で行う。

## 型・インターフェースの役割

次の使い分けはプロジェクト内の指針であり、言語共通の厳密な規則ではない。名前に合わせて不要な型や層を追加しない。

| 名前                  | 役割                                                         | 例                              |
|-----------------------|--------------------------------------------------------------|---------------------------------|
| `Repository`          | 集約を単位とする取得・保存                                   | `DailyPriceRepository`          |
| `Reader`              | 利用側が必要とする読み取り操作                               | `DailyVolumeReader`             |
| `Queries` / `Querier` | 条件を指定する問い合わせ操作の集合、またはそれを実行するもの | sqlcの`Queries`、`StockQuerier` |
| `Calculator`          | 入力から計算結果を導出する                                   | `DailyVolumeRatioCalculator`    |
| `Importer`            | 外部データを取り込み、内部へ保存する処理をまとめる           | `StockImporter`                 |
| `Client`              | 外部サービスとの通信を扱う                                   | `JQuantsClient`                 |
| `Handler`             | リクエストを受け、処理を呼び出してレスポンスを返す           | `DailyPriceChangeHandler`       |
| `Response`            | APIなどの返却形式を表す                                      | `DailyPriceChangeResponse`      |

`Reader`は読み取り能力、`Queries`は問い合わせ操作の集合を強調する。どちらもDB専用ではない。両方が当てはまる場合は、利用側の目的と既存の命名に合わせる。

Repositoryの実装が、利用側の小さなReaderインターフェースを満たしてもよい。実装の名前と、利用側が要求する役割は区別する。

## 変数・フィールド・引数

利用側での役割を表す。型が`DailyVolumeReader`ならフィールドと引数は`reader`、`DailyVolumeRatioCalculator`なら`calculator`を基本とする。

単一の値と複数の値を区別する（`volume` / `volumes`）。対応表はキーと値の関係が分かる名前にする（`volumesBySymbol`など）。

## 関数・メソッド

操作は動詞で表し、対象や条件が必要なら補う。周囲の型やモジュールから分かる情報は省略できる。

| 名前 | 意味 |
| --- | --- |
| `ListVolumesByTradeDate` | 指定日の出来高を複数取得する |
| `Calculate` | 計算する |
| `Import` | 取り込む |
| `Upsert` | 未登録なら追加し、登録済みなら更新する |
| `VolumeRatio` | モデルが持つ値から出来高倍率を導出する |

導出値や判定は、言語の慣習に応じて名詞や`Is`・`Has`などで表してもよい。名前が同じなら、同じ意味や条件で使う。

## データ・定数

- 計算用の入力と計算結果を区別する（`DailyVolumeComparison` / `DailyVolumeRatio`）。
- 数値の単位や基準が曖昧になる場合は、名前または型で表す。倍率と百分率を混同しない。
- 判定は肯定形を基本とする（`isListed`など）。
- 意味のある固定値には、用途を表す名前を付ける（`volumeRatioHistoryCount`など）。
- 境界を含むかが重要な引数は、`beforeDate`などで意味を示し、必要に応じてコメントを添える。

## 参考資料

- [Effective Go：Interface names](https://go.dev/doc/effective_go#interface-names)：振る舞いを表すインターフェース名の参考。
- [Go Code Review Comments：Interfaces](https://go.dev/wiki/CodeReviewComments#interfaces)：利用側でインターフェースを定義する設計の参考。

上記はGo向けの資料。言語に固有の規則は、他の言語へ一律に適用しない。

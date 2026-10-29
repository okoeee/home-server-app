# テスト方針

## 目的

振る舞いの正しさを確認しつつ、実装変更のたびに多数のテストを直す状態を避ける。同じ条件をすべての層で網羅せず、各テストが保証する責務を決める。

以下は本プロジェクトでの方針であり、言語や設計手法に共通する必須規則ではない。既存テストは、関連する変更やテスト整理のPRで段階的に合わせる。

## 層ごとの範囲

| 層                              | 保証すること           | 確認する例                                           | 他の層へ任せること                               |
|---------------------------------|------------------------|------------------------------------------------------|--------------------------------------------------|
| Domain                          | 業務ルール・計算       | 正常系、分岐、境界値、ゼロ、小数                     | HTTP、DB、処理の接続                             |
| Usecase                         | 処理の組み合わせと判断 | 銘柄の対応付け、計算不可の除外、取得失敗の伝播       | 計算式の全パターン、JSONの形式                   |
| Handler・Routerを通すHTTPテスト | HTTPとしての契約と接続 | パス、代表的な正常JSON、入力不正の400、処理失敗の500 | 計算式や除外条件の網羅、Repositoryの内部呼び出し |
| Infrastructure                  | 外部仕様との接続や変換 | SQLの取得・保存、外部レスポンスの変換                | Domainの計算式                                   |

Domainは、意味のあるルール・分岐・境界条件を網羅する。あらゆる入力の組み合わせやカバレッジ100%を一律の目標にしない。

Usecaseは、依存先を呼ぶだけならテストを無理に増やさない。条件分岐、集約、マッピング、除外など、自身が担う振る舞いに応じて追加する。

HTTPテストはRouterからHandlerを通す形を基本とし、同じ内容のRouter単体・Handler単体テストを重ねない。各エンドポイントの代表的な正常系でルート登録とレスポンスを確認し、入力検証やエラー変換に必要なケースを追加する。

HTTPテストでは、計算の入力を大量に作らず、Usecaseの結果・失敗を小さなstubで用意する方法を優先する。ただし、stub化のためだけに過剰な抽象化を追加しない。実Usecaseを通すテストも、確認目的が明確なら許容する。

## 依存先とテストデータ

- 外部通信や実DBが必要なテストは、単体テストと区別して必要な環境を明示する。
- テスト用のstub・fakeは、利用側が要求する最小のインターフェースを満たす。
- 検証には、原則として戻り値・状態・HTTPレスポンスなど、外から確認できる結果を使う。
- 呼び出し回数や順序は、重複保存・二重通知・トランザクションなど、正しさに関わる場合だけ検証する。
- 引数の渡し間違いも、結果の検証で検出できるなら、同じ内容の呼び出し検証を重ねない。
- シナリオに必要なデータだけを用意し、期待値を本番と同じ計算ロジックで導出しない。

## テスト名

Goのテスト関数名は`Test{関数名}_{シナリオ}`とする。シナリオは条件と期待する振る舞いが分かる名前にする。

```text
TestVolumeRatio_InsufficientHistory
TestCalculate_ExcludesUncalculableSymbols
TestGetDailyVolumeRatios_InvalidDateReturnsBadRequest
```

同じパッケージに同名メソッドのテストがある場合は、対象を区別するため型名も含める（例：`TestDailyVolumeRatioCalculatorCalculate_ReaderError`）。

同じ構造のケースはテーブル駆動テストにまとめてよい。その場合は親を`TestVolumeRatio_Scenarios`などとし、`t.Run`の名前で各条件と期待する結果を表す。サブテスト名は日本語でもよい。

Go以外でも対象とシナリオを表す方針は共通とし、実際の命名形式は言語・テストランナーに合わせる。

## Arrange・Act・Assert

各テストはArrange–Act–Assert（AAA）で構成し、各フェーズの先頭に`// Arrange`、`// Act`、`// Assert`を記載する。テーブル駆動テストでは、各サブテスト内に記載する。

- Arrange：入力、依存先、期待値を準備する。
- Act：検証対象の操作を実行する。
- Assert：結果を検証する。

```go
func TestVolumeRatio_InsufficientHistory(t *testing.T) {
    // Arrange
    comparison := DailyVolumeComparison{
        Volume: decimal.NewFromInt(100),
        PreviousVolumes: []decimal.Decimal{},
    }

    // Act
    ratio, ok := comparison.VolumeRatio()

    // Assert
    if ok {
        t.Error("VolumeRatio() ok = true, want false")
    }
    if !ratio.IsZero() {
        t.Errorf("VolumeRatio() = %s, want 0", ratio)
    }
}
```

Goでは標準の`testing`を基本とする。後続の検証を続けられない場合は`t.Fatal`・`t.Fatalf`、続けられる場合は`t.Error`・`t.Errorf`を使う。`decimal.Decimal`は内部表現ではなく、`Equal`や`IsZero`で数値として比較する。

## テストを追加・整理するときの確認

- このテストでしか検出できない不具合、または保証したい振る舞いがあるか。
- 下位のテストと同じ条件を網羅し直していないか。
- 正しいリファクタリングでも壊れる内部構造を固定していないか。
- 不要な共通化や、準備のためだけの大量のデータがないか。

## 参考資料

- [The Practical Test Pyramid：Avoid Test Duplication](https://martinfowler.com/articles/practical-test-pyramid.html#AvoidTestDuplication)：下位と上位のテストの役割を分け、重複を避ける考え方。
- [Google Testing Blog：Testing State vs. Testing Interactions](https://testing.googleblog.com/2013/03/testing-on-toilet-testing-state-vs.html)：結果の検証を優先し、呼び出しの検証が必要な条件を見極める考え方。
- [Mocks Aren’t Stubs](https://martinfowler.com/articles/mocksArentStubs.html)：テストダブルと、実装への結合についての比較。

層ごとの範囲、命名形式、AAAのコメントは、これらの原則と今回の合意をもとに定めたプロジェクト内の規則。

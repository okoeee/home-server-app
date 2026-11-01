# UIキット・アイコンの比較

調査日：2026-11-01。公式ドキュメント・公式リポジトリを参照した選定案で、まだ導入・採用はしていない。

## 前提と判断軸

- 現在はSvelte 5・TypeScript・Vite 8。SvelteKit・Tailwind CSS・UIキットは未導入。
- まず天気カードを整え、その後サーバー監視・GitHub・筋トレなどのカードにも広げたい。
- 個人で保守できること、必要な部品から追加できること、実装を理解しやすいことを優先する。
- UIキットはカード・ボタン・入力・ダイアログなどの見た目と操作、アイコンは情報の図示を担当する。別々に選定できる。
- 導入負担・今回との相性は、このリポジトリを前提にした判断。バンドルサイズ・描画性能の実測や導入検証は未実施。

## UIキットの比較

| 候補                                                                                  | 方式・導入条件                                                                             | 強み                                                                                | 今回の負担・注意点                                                                                       | 判断                                                 |
|---------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| [shadcn-svelte](https://www.shadcn-svelte.com/docs)                                   | 部品のコードをプロジェクトへ追加する方式。Svelte 5・Tailwind 4。Viteの導入ガイドあり       | 必要な部品を追加し、手元で見た目や構造を編集できる。操作系部品はBits UIを基盤にする | Tailwind・CSS変数・import aliasなどの初期設定が必要。コピーした部品の変更・更新は自分で管理する          | UIキットを導入する場合の第一候補                     |
| [Skeleton](https://www.skeleton.dev/docs/svelte/get-started/installation/vite-svelte) | TailwindのデザインシステムとSvelte部品。公式条件はSvelte 5以上・Vite 6以上・Tailwind 4以上 | テーマ・色・余白などをまとめてそろえやすく、ダッシュボード全体の統一に向く          | TailwindとSkeleton独自のテーマ・ユーティリティを覚える。今回の小さなカードだけなら設定範囲が広い         | テーマやダークモードを早く整えたい場合の候補         |
| [Flowbite Svelte](https://flowbite-svelte.com/docs/pages/quickstart)                  | TailwindベースのSvelte部品をパッケージから利用。Viteでの案内あり                           | カード・表・フォームなどの完成した部品を利用できる                                  | Tailwind・Flowbiteプラグイン・CSSの探索対象を設定する。独自の見た目への調整方法を覚える                  | 既製の部品を組み合わせたい場合の候補                 |
| [daisyUI](https://daisyui.com/docs/install/vite/)                                     | Tailwindプラグイン。`card`・`btn`などのCSSクラスをHTMLに付ける                             | Svelte固有のAPIが少なく、テーマ付きの見た目を比較的少ない記述で作れる               | Svelteの操作ロジックを提供するライブラリではない。操作系UIの状態やフォーカスなどはHTML・自分の実装で扱う | CSSクラス中心で簡単に整えたい場合の候補              |
| [Bits UI](https://www.bits-ui.com/docs/introduction)                                  | 見た目をほぼ持たないSvelte用の操作部品。Tailwindは必須ではない                             | キーボード操作・フォーカス管理などを利用しながら、既存CSSを継続できる               | 配色・カード・レイアウトのデザインは自分で作る。表示だけの天気カードには出番が少ない                     | 今回は単独導入せず、必要な操作部品ができたときに検討 |
| 現在のCSSを継続                                                                       | 新しいUI依存なし                                                                           | 最小の変更で済み、既存構成をそのまま理解できる                                      | カードが増えると色・余白・状態表現を自分で共通化する必要がある                                           | アイコンだけ先に入れる場合の有効な選択肢             |

[shadcn-svelteの現行要件](https://www.shadcn-svelte.com/docs/migration)はSvelte 5・Tailwind 4。[Flowbite Svelteの対応](https://github.com/themesberg/flowbite-svelte/discussions/1529)もSvelte 5・Tailwind 4が案内されている。SvelteKitへの移行を前提にはしない。

### アクセシビリティと更新方法

- [Bits UI](https://www.bits-ui.com/docs/introduction)はキーボード操作・フォーカス管理・スクリーンリーダー対応を提供する。shadcn-svelteではその基盤を利用できる。ただしラベル・色のコントラスト・部品の組み合わせはアプリ側でも確認する。
- Skeleton・Flowbite Svelteなどのパッケージ型は、依存バージョンの更新と移行ガイドに沿って保守する。
- [shadcn-svelte](https://www.shadcn-svelte.com/docs)は、Bits UIなどの依存更新と、手元に追加したUIコードの変更を分けて管理する。npmの更新だけでコピー済みコードの修正がすべて反映されるわけではない。
- daisyUIや独自CSSの採用だけで、複雑な操作部品のアクセシビリティが完成するわけではない。実際に使う部品単位で確認する。
- 調査時点の公式ガイドに対応環境・移行手順があることを確認した。継続保守の保証や人気順の評価はしていない。

## アイコンの比較

| 候補                                                          | 利用方法・見た目                                                                            | 天気との相性                                                                                   | 今回の判断                                                               |
|---------------------------------------------------------------|---------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|
| [Lucide](https://lucide.dev/guide/svelte/getting-started)     | `@lucide/svelte`のSVGコンポーネント。線画で色・サイズを調整でき、使ったアイコンを取り込む   | 晴れ・雲・雨・雪・雷などの基本的な図示に向く。細かい気象区分は同じ絵柄を共有し、日本語名で補う | 第一候補。サーバー・活動・設定などにも同じ見た目を使える                 |
| [Tabler Icons](https://docs.tabler.io/icons/libraries/svelte) | `@tabler/icons-svelte`のSVGコンポーネント。線幅・色・サイズを変更でき、個別importに対応     | 天気と一般的なUIの両方に利用できる                                                             | Lucideと同等の有力候補。好みの絵柄で選べる                               |
| [Meteocons](https://github.com/basmilius/meteocons)           | 天気専用のSVG。静止画・アニメーション・Lottieの配布がある。静止SVGは`@meteocons/svg-static` | 天気の細かな違いを表しやすく、色付きの予報カードに向く                                         | 天気の見やすさ・表現を優先する場合の第一候補。一般UI用アイコンは別に必要 |
| [Weather Icons](https://erikflowers.github.io/weather-icons/) | 天気専用のアイコンフォントとCSS                                                             | 気象向けの種類が多く、WMOコード向けの対応表もある                                              | 利用可能だが、今回はSVGを必要分だけ扱える上記候補を優先                  |

Lucide・Tablerは公式に個別利用とtree shakingを案内している。ただし実際の配信量はimport方法とビルド結果で確認し、未計測のサイズを比較材料にしない。

[Lucideの現行ガイド](https://lucide.dev/guide/svelte/getting-started)のパッケージ名は`@lucide/svelte`。古い記事の`lucide-svelte`と混同しない。また、[v1ではGitHubなどのブランドアイコンが削除](https://lucide.dev/guide/svelte/migration)されているため、将来GitHubのロゴ自体が必要なら別途用意する。

## ライセンス

比較したOSS本体は無料で利用できる。テンプレート・有料素材などの別商品は、この比較の対象外。配布物に必要なライセンス表記は保持する。

| 対象            | ライセンス・出典                                                                                    |
|-----------------|-----------------------------------------------------------------------------------------------------|
| shadcn-svelte   | [MIT](https://github.com/huntabyte/shadcn-svelte/blob/main/LICENSE.md)                              |
| Skeleton        | [MIT](https://github.com/skeletonlabs/skeleton/blob/main/LICENSE)                                   |
| Flowbite Svelte | [MIT](https://github.com/themesberg/flowbite-svelte#license)                                        |
| daisyUI         | [MIT](https://github.com/saadeghi/daisyui/blob/master/LICENSE)                                      |
| Bits UI         | [MIT](https://github.com/huntabyte/bits-ui#license)                                                 |
| Lucide          | [ISC](https://lucide.dev/license)                                                                   |
| Tabler Icons    | [MIT](https://github.com/tabler/tabler-icons)                                                       |
| Meteocons       | [MIT](https://github.com/basmilius/meteocons#license)                                               |
| Weather Icons   | [フォント：SIL OFL 1.1、コード：MIT、資料：CC BY 3.0](https://erikflowers.github.io/weather-icons/) |

## 今回の選定案

UIキットを導入するなら **shadcn-svelte＋Lucide** を第一候補とする。これは今回の学習・保守方針に基づく提案で、採用済みの決定ではない。

- 既存のSvelte・Viteを継続できる。
- 必要なCard・Badgeなどから追加でき、部品のコードを読んで調整できる。
- 将来、設定フォームやダイアログが必要になったときも同じ基盤を使える。
- Lucideで天気と一般UIの見た目をそろえられる。
- 初期設定とTailwindの学習、追加したUIコードの管理が増えることは受け入れる必要がある。

「設定を増やさず天気アイコンを表示する」だけなら、現在のCSS＋Lucideが最小。色付きの天気表現が好みならMeteocons、テーマ全体を素早く整えたいならSkeletonを優先する。

## 導入する場合の範囲

1. アイコンを先に追加し、既存の天気カードで確認する。UIキットとは分けてコミットできる。
2. `WeatherCondition`からアイコンへの対応はフロントに置く。16種類すべてを網羅し、`unknown`には不明を表すアイコンを使う。霧雨・着氷性の雨などで絵柄を共有しても、日本語名は残す。
3. 日別予報は現在時刻ではなく1日を表すため、最初は昼夜で切り替えず固定の絵柄にする。天気APIは変更しない。
4. アイコンはローカルのビルド・配信に含め、表示時の外部アイコンAPIやCDN取得を追加しない。最初は静止アイコンとする。
5. 日本語の天気名と並べる装飾アイコンには`aria-hidden="true"`を付け、アイコンだけに意味を持たせない。
6. UIキット採用時はTailwindなどの基盤と、天気カードへの適用を分ける。既存の騰落率表の意味や表示件数は変えない。
7. 導入前後で`npm test`・`npm run check`・`npm run build`を実行し、狭い画面・文字のコントラスト・読み込み／エラー状態を確認する。インストール時に実際のSvelte・Vite・TypeScriptとの依存条件を確認し、ロックファイルを管理する。

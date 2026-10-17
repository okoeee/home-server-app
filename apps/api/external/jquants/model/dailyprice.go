package model

import "github.com/shopspring/decimal"

// DailyPrice J-Quantsから取得した株価日足
// [https://jpx-jquants.com/ja/spec/eq-bars-daily]
type DailyPrice struct {
	// 取引日
	Date string `json:"Date"`

	// 銘柄コード
	Code string `json:"Code"`

	// 始値
	Open *decimal.Decimal `json:"O"`

	// 高値
	High *decimal.Decimal `json:"H"`

	// 安値
	Low *decimal.Decimal `json:"L"`

	// 終値
	Close *decimal.Decimal `json:"C"`

	// 出来高
	Volume *decimal.Decimal `json:"Vo"`

	// 売買代金
	TurnoverValue *decimal.Decimal `json:"Va"`

	// 株式分割や併合などを反映する調整係数
	AdjustmentFactor decimal.Decimal `json:"AdjFactor"`

	// 調整済み始値
	AdjustedOpen *decimal.Decimal `json:"AdjO"`

	// 調整済み高値
	AdjustedHigh *decimal.Decimal `json:"AdjH"`

	// 調整済み安値
	AdjustedLow *decimal.Decimal `json:"AdjL"`

	// 調整済み終値
	AdjustedClose *decimal.Decimal `json:"AdjC"`

	// 調整済み出来高
	AdjustedVolume *decimal.Decimal `json:"AdjVo"`
}

// ListedDailyPrices J-Quantsの日足APIレスポンス全体
type ListedDailyPrices struct {
	// 取得した日足の一覧
	Data []DailyPrice `json:"data"`

	// 次のページを取得するためのキー
	PaginationKey string `json:"pagination_key"`
}

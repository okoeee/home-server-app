package model

// DailyPrice J-Quantsから取得した株価日足
// [https://jpx-jquants.com/ja/spec/eq-bars-daily]
type DailyPrice struct {
	// 取引日
	Date string `json:"Date"`

	// 銘柄コード
	Code string `json:"Code"`

	// 始値
	Open *float64 `json:"O"`

	// 高値
	High *float64 `json:"H"`

	// 安値
	Low *float64 `json:"L"`

	// 終値
	Close *float64 `json:"C"`

	// 出来高
	Volume *float64 `json:"Vo"`

	// 売買代金
	TurnoverValue *float64 `json:"Va"`

	// 株式分割や併合などを反映する調整係数
	AdjustmentFactor float64 `json:"AdjFactor"`

	// 調整済み始値
	AdjustedOpen *float64 `json:"AdjO"`

	// 調整済み高値
	AdjustedHigh *float64 `json:"AdjH"`

	// 調整済み安値
	AdjustedLow *float64 `json:"AdjL"`

	// 調整済み終値
	AdjustedClose *float64 `json:"AdjC"`

	// 調整済み出来高
	AdjustedVolume *float64 `json:"AdjVo"`
}

// ListedDailyPrices J-Quantsの日足APIレスポンス全体
type ListedDailyPrices struct {
	// 取得した日足の一覧
	Data []DailyPrice `json:"data"`

	// 次のページを取得するためのキー
	PaginationKey string `json:"pagination_key"`
}

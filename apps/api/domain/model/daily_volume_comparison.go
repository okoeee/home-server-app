package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const volumeRatioHistoryCount = 20

type DailyVolumeComparison struct {
	Symbol          string
	TradeDate       time.Time
	Volume          decimal.Decimal
	PreviousVolumes []decimal.Decimal
}

// VolumeRatio VolumeRatioは過去20件の平均に対する倍率を返す
// 履歴が20件でない場合や平均がゼロの場合、計算不可としてfalseを返す
func (comparison DailyVolumeComparison) VolumeRatio() (decimal.Decimal, bool) {
	if len(comparison.PreviousVolumes) != volumeRatioHistoryCount {
		return decimal.Zero, false
	}

	total := decimal.Zero
	for _, volume := range comparison.PreviousVolumes {
		total = total.Add(volume)
	}

	average := total.Div(decimal.NewFromInt(volumeRatioHistoryCount))
	if average.IsZero() {
		return decimal.Zero, false
	}

	return comparison.Volume.Div(average), true
}

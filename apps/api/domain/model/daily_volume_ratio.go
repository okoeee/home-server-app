package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyVolumeRatio struct {
	Symbol    string
	TradeDate time.Time
	Volume    decimal.Decimal
	Ratio     decimal.Decimal
}

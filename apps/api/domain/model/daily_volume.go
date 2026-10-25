package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyVolume struct {
	Symbol    string
	TradeDate time.Time
	Volume    decimal.Decimal
}

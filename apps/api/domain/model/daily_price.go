package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyPrice struct {
	Symbol     string
	TradeDate  time.Time
	OpenPrice  decimal.Decimal
	HighPrice  decimal.Decimal
	LowPrice   decimal.Decimal
	ClosePrice decimal.Decimal
	Volume     decimal.Decimal
}

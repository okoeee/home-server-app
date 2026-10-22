package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyPriceComparison struct {
	Symbol             string
	TradeDate          time.Time
	ClosePrice         decimal.Decimal
	PreviousTradeDate  time.Time
	PreviousClosePrice decimal.Decimal
}

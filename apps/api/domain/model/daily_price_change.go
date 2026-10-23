package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyPriceChange struct {
	Symbol             string
	TradeDate          time.Time
	ClosePrice         decimal.Decimal
	PreviousTradeDate  time.Time
	PreviousClosePrice decimal.Decimal
	ChangeRate         decimal.Decimal
}

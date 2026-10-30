package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyPriceChange struct {
	Symbol             string
	CompanyName        *string
	TradeDate          time.Time
	ClosePrice         decimal.Decimal
	PreviousTradeDate  time.Time
	PreviousClosePrice decimal.Decimal
	ChangeRate         decimal.Decimal
}

package model

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type DailyPriceComparison struct {
	Symbol             string
	CompanyName        *string
	TradeDate          time.Time
	ClosePrice         decimal.Decimal
	PreviousTradeDate  time.Time
	PreviousClosePrice decimal.Decimal
}

func (comparison DailyPriceComparison) ChangeRate() (
	decimal.Decimal,
	error,
) {

	if comparison.PreviousClosePrice.IsZero() {
		return decimal.Zero, fmt.Errorf("previous close price is zero: symbol=%s", comparison.Symbol)
	}

	changeRate := comparison.ClosePrice.Sub(comparison.PreviousClosePrice).
		Div(comparison.PreviousClosePrice).
		Mul(decimal.NewFromInt(100))

	return changeRate, nil

}

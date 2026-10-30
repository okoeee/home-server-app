package dailyprice

import (
	"api/domain/model"
	"context"
	"time"
)

type DailyPriceComparisonReader interface {
	GetLatestTradeDate(ctx context.Context) (time.Time, bool, error)
	ListByTradeDate(ctx context.Context, tradeDate time.Time) ([]model.DailyPriceComparison, error)
}

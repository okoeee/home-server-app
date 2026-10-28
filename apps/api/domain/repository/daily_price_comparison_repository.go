package repository

import (
	"api/domain/model"
	"context"
	"time"
)

type DailyPriceComparisonRepository interface {
	GetLatestTradeDate(context.Context) (time.Time, bool, error)
	ListByTradeDate(
		context.Context,
		time.Time,
	) ([]model.DailyPriceComparison, error)
}

package repository

import (
	"api/domain/model"
	"context"
	"time"
)

type DailyPriceComparisonRepository interface {
	ListByTradeDate(
		context.Context,
		time.Time,
	) ([]model.DailyPriceComparison, error)
}

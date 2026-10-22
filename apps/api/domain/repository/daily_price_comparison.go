package repository

import (
	"api/domain/model"
	"context"
	"time"
)

type DailyPriceComparison interface {
	ListByTradeDate(
		context.Context,
		time.Time,
	) ([]model.DailyPriceComparison, error)
}

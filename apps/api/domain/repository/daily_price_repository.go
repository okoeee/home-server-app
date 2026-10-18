package repository

import (
	"api/domain/model"
	"context"
)

type DailyPriceRepository interface {
	Upsert(context.Context, model.DailyPrice) error
}

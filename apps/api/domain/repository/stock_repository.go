package repository

import (
	"api/domain/model"
	"context"
)

type StockRepository interface {
	Upsert(context.Context, model.Stock) error
}

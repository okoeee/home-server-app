package mysql

import (
	"api/db"
	"api/domain/model"
	"context"
)

type DailyPriceRepository struct {
	queries *db.Queries
}

func NewDailyPriceRepository(queries *db.Queries) *DailyPriceRepository {
	return &DailyPriceRepository{queries}
}

func (repository *DailyPriceRepository) Upsert(
	ctx context.Context,
	dailyPrice model.DailyPrice,
) error {

	return repository.queries.UpsertDailyPrice(
		ctx,
		db.UpsertDailyPriceParams{
			Symbol:     dailyPrice.Symbol,
			TradeDate:  dailyPrice.TradeDate,
			OpenPrice:  dailyPrice.OpenPrice,
			HighPrice:  dailyPrice.HighPrice,
			LowPrice:   dailyPrice.LowPrice,
			ClosePrice: dailyPrice.ClosePrice,
			Volume:     dailyPrice.Volume,
		},
	)

}

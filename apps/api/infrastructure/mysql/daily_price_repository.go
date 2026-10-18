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
			OpenPrice:  dailyPrice.OpenPrice.StringFixed(1),
			HighPrice:  dailyPrice.HighPrice.StringFixed(1),
			LowPrice:   dailyPrice.LowPrice.StringFixed(1),
			ClosePrice: dailyPrice.ClosePrice.StringFixed(1),
			Volume:     dailyPrice.Volume.StringFixed(1),
		},
	)

}

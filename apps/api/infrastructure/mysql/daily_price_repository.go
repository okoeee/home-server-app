package mysql

import (
	"api/db"
	"api/domain/model"
	"context"
	"fmt"
	"time"
)

type DailyPriceRepository struct {
	queries *db.Queries
}

func NewDailyPriceRepository(queries *db.Queries) *DailyPriceRepository {
	return &DailyPriceRepository{queries}
}

func (repository *DailyPriceRepository) ListByTradeDate(
	ctx context.Context,
	date time.Time,
) ([]model.DailyPriceComparison, error) {

	rows, err := repository.queries.ListDailyPricesWithPreviousClose(
		ctx,
		date,
	)
	if err != nil {
		return nil, fmt.Errorf("daily price repository list by trade date %w", err)
	}

	dailyPriceComparisons := make([]model.DailyPriceComparison, 0, len(rows))

	for _, row := range rows {
		dailyPriceComparisons = append(dailyPriceComparisons, model.DailyPriceComparison{
			Symbol:             row.Symbol,
			TradeDate:          row.TradeDate,
			ClosePrice:         row.ClosePrice,
			PreviousTradeDate:  row.PreviousTradeDate,
			PreviousClosePrice: row.PreviousClosePrice,
		})
	}

	return dailyPriceComparisons, nil
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

package mysql

import (
	"api/db"
	"api/domain/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DailyPriceRepository struct {
	queries *db.Queries
}

func NewDailyPriceRepository(queries *db.Queries) *DailyPriceRepository {
	return &DailyPriceRepository{queries}
}

func (repository *DailyPriceRepository) GetLatestTradeDate(ctx context.Context) (time.Time, bool, error) {
	date, err := repository.queries.GetLatestDailyPriceTradeDate(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, fmt.Errorf("get latest daily price trade date: %w", err)
	}
	return date, true, nil
}

func (repository *DailyPriceRepository) ListVolumesByTradeDate(
	ctx context.Context,
	tradeDate time.Time,
) ([]model.DailyVolume, error) {
	rows, err := repository.queries.ListDailyVolumesByTradeDate(ctx, tradeDate)
	if err != nil {
		return nil, fmt.Errorf("list daily volumes for %s: %w", tradeDate.Format(time.DateOnly), err)
	}

	volumes := make([]model.DailyVolume, 0, len(rows))
	for _, row := range rows {
		volumes = append(volumes, model.DailyVolume{
			Symbol:    row.Symbol,
			TradeDate: row.TradeDate,
			Volume:    row.Volume,
		})
	}
	return volumes, nil
}

func (repository *DailyPriceRepository) ListVolumeHistory(
	ctx context.Context,
	beforeDate time.Time,
) ([]model.DailyVolume, error) {
	rows, err := repository.queries.ListDailyVolumeHistory(ctx, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("list daily volume history before %s: %w", beforeDate.Format(time.DateOnly), err)
	}

	volumes := make([]model.DailyVolume, 0, len(rows))
	for _, row := range rows {
		volumes = append(volumes, model.DailyVolume{
			Symbol:    row.Symbol,
			TradeDate: row.TradeDate,
			Volume:    row.Volume,
		})
	}
	return volumes, nil
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

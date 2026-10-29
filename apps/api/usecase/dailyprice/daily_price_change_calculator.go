package dailyprice

import (
	"api/domain/model"
	"api/domain/repository"
	"context"
	"fmt"
	"time"
)

type DailyPriceChangeCalculator struct {
	repository repository.DailyPriceComparisonRepository
}

func NewDailyPriceChangeCalculator(
	repository repository.DailyPriceComparisonRepository,
) *DailyPriceChangeCalculator {
	return &DailyPriceChangeCalculator{
		repository: repository,
	}
}

func (calculator *DailyPriceChangeCalculator) CalculateLatest(ctx context.Context) (model.LatestDailyPriceChanges, error) {
	tradeDate, found, err := calculator.repository.GetLatestTradeDate(ctx)
	if err != nil {
		return model.LatestDailyPriceChanges{}, fmt.Errorf("get latest trade date: %w", err)
	}
	if !found {
		return model.LatestDailyPriceChanges{Changes: []model.DailyPriceChange{}}, nil
	}
	changes, err := calculator.calculate(ctx, tradeDate)
	if err != nil {
		return model.LatestDailyPriceChanges{}, err
	}
	return model.LatestDailyPriceChanges{TradeDate: &tradeDate, Changes: changes}, nil
}

func (calculator *DailyPriceChangeCalculator) calculate(
	ctx context.Context,
	tradeDate time.Time,
) ([]model.DailyPriceChange, error) {

	comparisons, err := calculator.repository.ListByTradeDate(ctx, tradeDate)
	if err != nil {
		return nil, fmt.Errorf("daily price change calculator: %w", err)
	}

	results := make([]model.DailyPriceChange, 0, len(comparisons))

	for _, comparison := range comparisons {

		dailyPriceChange, err := comparison.ChangeRate()

		if err != nil {
			return nil, fmt.Errorf("daily price change calculator: %w", err)
		}

		results = append(results, model.DailyPriceChange{
			Symbol:             comparison.Symbol,
			TradeDate:          comparison.TradeDate,
			ClosePrice:         comparison.ClosePrice,
			PreviousTradeDate:  comparison.PreviousTradeDate,
			PreviousClosePrice: comparison.PreviousClosePrice,
			ChangeRate:         dailyPriceChange,
		})
	}

	return results, nil

}

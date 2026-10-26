package dailyprice

import (
	"api/domain/model"
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type DailyVolumeReader interface {
	ListVolumesByTradeDate(ctx context.Context, tradeDate time.Time) ([]model.DailyVolume, error)
	ListVolumeHistory(ctx context.Context, beforeDate time.Time) ([]model.DailyVolume, error)
}

type DailyVolumeRatioCalculator struct {
	repository DailyVolumeReader
}

func NewDailyVolumeRatioCalculator(repository DailyVolumeReader) *DailyVolumeRatioCalculator {
	return &DailyVolumeRatioCalculator{
		repository: repository,
	}
}

func (calculate *DailyVolumeRatioCalculator) Calculate(
	ctx context.Context,
	tradeDate time.Time,
) ([]model.DailyVolumeRatio, error) {

	targetDateVolumes, err := calculate.repository.ListVolumesByTradeDate(ctx, tradeDate)
	if err != nil {
		return nil, fmt.Errorf("list volumes by tradeDate: %w", err)
	}
	volumeHistories, err := calculate.repository.ListVolumeHistory(ctx, tradeDate)
	if err != nil {
		return nil, fmt.Errorf("list volume history by tradeDate: %w", err)
	}

	symbolToVolumeHistories := make(map[string][]decimal.Decimal)
	for _, history := range volumeHistories {
		symbolToVolumeHistories[history.Symbol] = append(symbolToVolumeHistories[history.Symbol], history.Volume)
	}

	dailyVolumeComparisons := make([]model.DailyVolumeComparison, 0, len(targetDateVolumes))
	for _, targetDateVolume := range targetDateVolumes {
		volumeHistories := symbolToVolumeHistories[targetDateVolume.Symbol]

		dailyVolumeComparisons = append(dailyVolumeComparisons, model.DailyVolumeComparison{
			Symbol:          targetDateVolume.Symbol,
			TradeDate:       targetDateVolume.TradeDate,
			Volume:          targetDateVolume.Volume,
			PreviousVolumes: volumeHistories,
		})
	}

	results := make([]model.DailyVolumeRatio, 0, len(dailyVolumeComparisons))
	for _, comparison := range dailyVolumeComparisons {
		volumeRatio, ok := comparison.VolumeRatio()

		if ok {
			results = append(results, model.DailyVolumeRatio{
				Symbol:    comparison.Symbol,
				TradeDate: comparison.TradeDate,
				Volume:    comparison.Volume,
				Ratio:     volumeRatio,
			})
		}
	}

	return results, nil
}

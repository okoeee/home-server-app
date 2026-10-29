package dailyprice

import (
	"api/domain/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

type stubRepository struct {
	volumes    []model.DailyVolume
	history    []model.DailyVolume
	volumesErr error
	historyErr error
}

func (s *stubRepository) ListVolumesByTradeDate(ctx context.Context, tradeDate time.Time) ([]model.DailyVolume, error) {
	return s.volumes, s.volumesErr
}

func (s *stubRepository) ListVolumeHistory(ctx context.Context, beforeDate time.Time) ([]model.DailyVolume, error) {
	return s.history, s.historyErr
}

func TestDailyVolumeRatioCalculatorCalculate_ExcludesUncalculableSymbols(t *testing.T) {
	// Arrange
	tradeDate := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	repository := &stubRepository{
		volumes: []model.DailyVolume{
			{Symbol: "valid", TradeDate: tradeDate, Volume: decimal.NewFromInt(200)},
			{Symbol: "short", TradeDate: tradeDate, Volume: decimal.NewFromInt(200)},
			{Symbol: "zero", TradeDate: tradeDate, Volume: decimal.NewFromInt(200)},
			{Symbol: "missing", TradeDate: tradeDate, Volume: decimal.NewFromInt(200)},
		},
	}
	for index := range 20 {
		date := tradeDate.AddDate(0, 0, -index-1)
		repository.history = append(repository.history,
			model.DailyVolume{Symbol: "valid", TradeDate: date, Volume: decimal.NewFromInt(100)},
			model.DailyVolume{Symbol: "zero", TradeDate: date, Volume: decimal.Zero},
			model.DailyVolume{Symbol: "history_only", TradeDate: date, Volume: decimal.NewFromInt(100)},
		)
		if index < 19 {
			repository.history = append(repository.history,
				model.DailyVolume{Symbol: "short", TradeDate: date, Volume: decimal.NewFromInt(100)},
			)
		}
	}
	calculator := NewDailyVolumeRatioCalculator(repository)
	wantRatio := decimal.NewFromInt(2)

	// Act
	got, err := calculator.Calculate(context.Background(), tradeDate)

	// Assert
	if err != nil {
		t.Fatalf("Calculate() unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("Calculate() returned %d results, want 1: %+v", len(got), got)
	}
	if got[0].Symbol != "valid" || !got[0].Ratio.Equal(wantRatio) {
		t.Errorf("Calculate() = %+v, want only symbol valid with ratio 2", got)
	}
}

func TestDailyVolumeRatioCalculatorCalculate_EmptyResults(t *testing.T) {
	tradeDate := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		volumes []model.DailyVolume
	}{
		{name: "対象日の出来高がない"},
		{
			name: "全銘柄が計算不可",
			volumes: []model.DailyVolume{
				{Symbol: "12340", TradeDate: tradeDate, Volume: decimal.NewFromInt(100)},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			repository := &stubRepository{volumes: test.volumes}
			calculator := NewDailyVolumeRatioCalculator(repository)

			// Act
			got, err := calculator.Calculate(context.Background(), tradeDate)

			// Assert
			if err != nil {
				t.Fatalf("Calculate() unexpected error: %v", err)
			}
			if len(got) != 0 {
				t.Errorf("Calculate() = %+v, want no results", got)
			}
		})
	}
}

func TestDailyVolumeRatioCalculatorCalculate_PropagatesReaderErrors(t *testing.T) {
	wantErr := errors.New("read failed")
	tests := []struct {
		name       string
		volumesErr error
		historyErr error
	}{
		{name: "対象日の取得に失敗", volumesErr: wantErr},
		{name: "履歴の取得に失敗", historyErr: wantErr},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			repository := &stubRepository{volumesErr: test.volumesErr, historyErr: test.historyErr}
			tradeDate := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
			calculator := NewDailyVolumeRatioCalculator(repository)

			// Act
			got, err := calculator.Calculate(context.Background(), tradeDate)

			// Assert
			if !errors.Is(err, wantErr) {
				t.Errorf("Calculate() error = %v, want wrapped %v", err, wantErr)
			}
			if got != nil {
				t.Errorf("Calculate() = %+v, want nil on error", got)
			}
		})
	}
}

func TestDailyVolumeRatioCalculatorCalculate_MapsHistoryBySymbol(t *testing.T) {
	// Arrange
	tradeDate := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	repository := &stubRepository{
		volumes: []model.DailyVolume{
			{Symbol: "12340", TradeDate: tradeDate, Volume: decimal.NewFromInt(200)},
			{Symbol: "56780", TradeDate: tradeDate, Volume: decimal.NewFromInt(150)},
		},
	}
	for index := range 20 {
		date := tradeDate.AddDate(0, 0, -index-1)
		repository.history = append(repository.history,
			model.DailyVolume{Symbol: "12340", TradeDate: date, Volume: decimal.NewFromInt(100)},
			model.DailyVolume{Symbol: "56780", TradeDate: date, Volume: decimal.NewFromInt(300)},
		)
	}

	calculator := NewDailyVolumeRatioCalculator(repository)
	wantRatios := []decimal.Decimal{decimal.NewFromInt(2), decimal.RequireFromString("0.5")}

	// Act
	got, err := calculator.Calculate(context.Background(), tradeDate)

	// Assert
	if err != nil {
		t.Fatalf("Calculate() unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("Calculate() returned %d results, want 2", len(got))
	}

	for index, result := range got {
		wantVolume := repository.volumes[index]
		if result.Symbol != wantVolume.Symbol || !result.TradeDate.Equal(tradeDate) || !result.Volume.Equal(wantVolume.Volume) {
			t.Errorf("result[%d] = %+v, want volume data %+v", index, result, wantVolume)
		}
		if !result.Ratio.Equal(wantRatios[index]) {
			t.Errorf("symbol %s: ratio = %s, want %s", result.Symbol, result.Ratio, wantRatios[index])
		}
	}
}

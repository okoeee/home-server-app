package dailyprice

import (
	"api/domain/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

type priceComparisonRepositoryStub struct {
	latestDate  time.Time
	latestFound bool
	latestErr   error
	comparisons map[time.Time][]model.DailyPriceComparison
	listErr     error
}

func (stub priceComparisonRepositoryStub) GetLatestTradeDate(context.Context) (time.Time, bool, error) {
	return stub.latestDate, stub.latestFound, stub.latestErr
}

func (stub priceComparisonRepositoryStub) ListByTradeDate(_ context.Context, date time.Time) ([]model.DailyPriceComparison, error) {
	return stub.comparisons[date], stub.listErr
}

func TestDailyPriceChangeCalculatorCalculate_MapsComparisonsToResults(t *testing.T) {
	// Arrange
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	previousDate := date.AddDate(0, 0, -3)
	comparisons := []model.DailyPriceComparison{
		{
			Symbol: "12340", TradeDate: date, ClosePrice: decimal.NewFromInt(110),
			PreviousTradeDate: previousDate, PreviousClosePrice: decimal.NewFromInt(100),
		},
		{
			Symbol: "56780", TradeDate: date, ClosePrice: decimal.NewFromInt(190),
			PreviousTradeDate: previousDate.AddDate(0, 0, -1), PreviousClosePrice: decimal.NewFromInt(200),
		},
	}
	calculator := NewDailyPriceChangeCalculator(priceComparisonRepositoryStub{
		comparisons: map[time.Time][]model.DailyPriceComparison{date: comparisons},
	})
	wantRates := []decimal.Decimal{decimal.NewFromInt(10), decimal.NewFromInt(-5)}

	// Act
	got, err := calculator.Calculate(context.Background(), date)

	// Assert
	if err != nil {
		t.Fatalf("Calculate() unexpected error: %v", err)
	}
	if len(got) != len(comparisons) {
		t.Fatalf("len(results) = %d, want %d", len(got), len(comparisons))
	}
	for index, result := range got {
		input := comparisons[index]
		if result.Symbol != input.Symbol || !result.TradeDate.Equal(input.TradeDate) ||
			!result.ClosePrice.Equal(input.ClosePrice) || !result.PreviousTradeDate.Equal(input.PreviousTradeDate) ||
			!result.PreviousClosePrice.Equal(input.PreviousClosePrice) || !result.ChangeRate.Equal(wantRates[index]) {
			t.Errorf("result[%d] = %+v, want comparison %+v with rate %s", index, result, input, wantRates[index])
		}
	}
}

func TestDailyPriceChangeCalculatorCalculate_NoComparisonsReturnsEmptyResults(t *testing.T) {
	// Arrange
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	calculator := NewDailyPriceChangeCalculator(priceComparisonRepositoryStub{})

	// Act
	got, err := calculator.Calculate(context.Background(), date)

	// Assert
	if err != nil {
		t.Fatalf("Calculate() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Calculate() = %+v, want no results", got)
	}
}

func TestDailyPriceChangeCalculatorCalculate_ErrorsReturnNoPartialResults(t *testing.T) {
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	readErr := errors.New("read failed")
	tests := []struct {
		name       string
		repository priceComparisonRepositoryStub
		wantCause  error
	}{
		{
			name:       "取得エラーを伝播する",
			repository: priceComparisonRepositoryStub{listErr: readErr},
			wantCause:  readErr,
		},
		{
			name: "途中で計算に失敗しても一部の結果を返さない",
			repository: priceComparisonRepositoryStub{
				comparisons: map[time.Time][]model.DailyPriceComparison{
					date: {
						{Symbol: "12340", TradeDate: date, ClosePrice: decimal.NewFromInt(110), PreviousClosePrice: decimal.NewFromInt(100)},
						{Symbol: "56780", TradeDate: date, PreviousClosePrice: decimal.Zero},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			calculator := NewDailyPriceChangeCalculator(test.repository)

			// Act
			got, err := calculator.Calculate(context.Background(), date)

			// Assert
			if err == nil {
				t.Fatal("Calculate() error = nil, want error")
			}
			if test.wantCause != nil && !errors.Is(err, test.wantCause) {
				t.Errorf("Calculate() error = %v, want wrapped %v", err, test.wantCause)
			}
			if got != nil {
				t.Errorf("Calculate() = %+v, want nil on error", got)
			}
		})
	}
}

func TestCalculateLatest_UsesLatestStoredDate(t *testing.T) {
	// Arrange
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	previousDate := date.AddDate(0, 0, -3)
	calculator := NewDailyPriceChangeCalculator(priceComparisonRepositoryStub{
		latestDate: date, latestFound: true,
		comparisons: map[time.Time][]model.DailyPriceComparison{
			date:         {{Symbol: "12340", TradeDate: date, ClosePrice: decimal.NewFromInt(110), PreviousTradeDate: previousDate, PreviousClosePrice: decimal.NewFromInt(100)}},
			previousDate: {{Symbol: "56780", TradeDate: previousDate, ClosePrice: decimal.NewFromInt(100), PreviousClosePrice: decimal.NewFromInt(100)}},
		},
	})

	// Act
	got, err := calculator.CalculateLatest(context.Background())

	// Assert
	if err != nil {
		t.Fatalf("CalculateLatest() unexpected error: %v", err)
	}
	if got.TradeDate == nil || !got.TradeDate.Equal(date) {
		t.Fatalf("trade date = %v, want %v", got.TradeDate, date)
	}
	if len(got.Changes) != 1 {
		t.Fatalf("len(changes) = %d, want 1", len(got.Changes))
	}
	if got.Changes[0].Symbol != "12340" || !got.Changes[0].TradeDate.Equal(date) || !got.Changes[0].ChangeRate.Equal(decimal.NewFromInt(10)) {
		t.Errorf("changes = %+v, want latest date's symbol 12340 with rate 10", got.Changes)
	}
}

func TestCalculateLatest_EmptyResults(t *testing.T) {
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		found bool
	}{
		{name: "DBが空なら対象日なしで返す"},
		{name: "比較可能な銘柄がなくても最新日を維持する", found: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			calculator := NewDailyPriceChangeCalculator(priceComparisonRepositoryStub{latestDate: date, latestFound: test.found})

			// Act
			got, err := calculator.CalculateLatest(context.Background())

			// Assert
			if err != nil {
				t.Fatalf("CalculateLatest() unexpected error: %v", err)
			}
			if test.found {
				if got.TradeDate == nil || !got.TradeDate.Equal(date) {
					t.Errorf("trade date = %v, want %v", got.TradeDate, date)
				}
			} else if got.TradeDate != nil {
				t.Errorf("trade date = %v, want nil", got.TradeDate)
			}
			if len(got.Changes) != 0 {
				t.Errorf("changes = %+v, want no results", got.Changes)
			}
		})
	}
}

func TestCalculateLatest_PropagatesRepositoryErrors(t *testing.T) {
	readErr := errors.New("read failed")
	tests := []struct {
		name       string
		repository priceComparisonRepositoryStub
	}{
		{name: "最新日の取得エラー", repository: priceComparisonRepositoryStub{latestErr: readErr}},
		{name: "日足の取得エラー", repository: priceComparisonRepositoryStub{latestFound: true, listErr: readErr}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			calculator := NewDailyPriceChangeCalculator(test.repository)

			// Act
			_, err := calculator.CalculateLatest(context.Background())

			// Assert
			if !errors.Is(err, readErr) {
				t.Errorf("CalculateLatest() error = %v, want wrapped %v", err, readErr)
			}
		})
	}
}

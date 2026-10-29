package model

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestChangeRate_Scenarios(t *testing.T) {
	tests := []struct {
		name               string
		closePrice         string
		previousClosePrice string
		want               string
		wantErr            bool
	}{
		{
			name:               "終値が10パーセント上昇した場合",
			closePrice:         "110",
			previousClosePrice: "100",
			want:               "10",
		},
		{
			name:               "終値が10パーセント下落した場合",
			closePrice:         "90",
			previousClosePrice: "100",
			want:               "-10",
		},
		{
			name:               "終値が変化しなかった場合",
			closePrice:         "100",
			previousClosePrice: "100",
			want:               "0",
		},
		{
			name:               "小数を含む場合",
			closePrice:         "110.5",
			previousClosePrice: "100",
			want:               "10.5",
		},
		{
			name:               "前回終値がゼロの場合",
			closePrice:         "100",
			previousClosePrice: "0",
			wantErr:            true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			comparison := DailyPriceComparison{
				Symbol:             "12340",
				ClosePrice:         decimal.RequireFromString(test.closePrice),
				PreviousClosePrice: decimal.RequireFromString(test.previousClosePrice),
			}
			want := decimal.Zero
			if !test.wantErr {
				want = decimal.RequireFromString(test.want)
			}

			// Act
			got, err := comparison.ChangeRate()

			// Assert
			if test.wantErr {
				if err == nil {
					t.Fatal("ChangeRate() error = nil, want error")
				}

				return
			}
			if err != nil {
				t.Fatalf("ChangeRate() unexpected error: %v", err)
			}

			if !got.Equal(want) {
				t.Errorf(
					"ChangeRate() = %s, want %s",
					got.String(),
					want.String(),
				)
			}
		})
	}
}

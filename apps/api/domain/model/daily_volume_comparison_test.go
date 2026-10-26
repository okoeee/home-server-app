package model

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestDailyVolumeComparisonVolumeRatio(t *testing.T) {
	tests := []struct {
		name            string
		volume          string
		previousVolumes []string
		want            string
		wantOK          bool
	}{
		{
			name:            "過去20件の平均の2倍",
			volume:          "21",
			previousVolumes: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"},
			want:            "2",
			wantOK:          true,
		},
		{name: "平均と同じなら1倍", volume: "100", previousVolumes: repeatedVolumes("100", 20), want: "1", wantOK: true},
		{name: "平均の半分なら0.5倍", volume: "50", previousVolumes: repeatedVolumes("100", 20), want: "0.5", wantOK: true},
		{name: "小数の出来高を計算できる", volume: "150.3", previousVolumes: repeatedVolumes("100.2", 20), want: "1.5", wantOK: true},
		{name: "対象日の出来高がゼロなら0倍", volume: "0", previousVolumes: repeatedVolumes("100", 20), want: "0", wantOK: true},
		{
			name:            "過去の出来高ゼロも平均の件数に含める",
			volume:          "100",
			previousVolumes: append(repeatedVolumes("0", 10), repeatedVolumes("100", 10)...),
			want:            "2",
			wantOK:          true,
		},
		{name: "履歴がなければ計算不可", volume: "100", want: "0"},
		{name: "履歴が19件なら計算不可", volume: "100", previousVolumes: repeatedVolumes("100", 19), want: "0"},
		{name: "履歴が21件なら計算不可", volume: "100", previousVolumes: repeatedVolumes("100", 21), want: "0"},
		{name: "平均出来高がゼロなら計算不可", volume: "100", previousVolumes: repeatedVolumes("0", 20), want: "0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			history := make([]decimal.Decimal, 0, len(test.previousVolumes))
			for _, volume := range test.previousVolumes {
				history = append(history, decimal.RequireFromString(volume))
			}
			comparison := DailyVolumeComparison{
				Symbol:          "12340",
				Volume:          decimal.RequireFromString(test.volume),
				PreviousVolumes: history,
			}
			got, ok := comparison.VolumeRatio()
			if ok != test.wantOK {
				t.Errorf("VolumeRatio() ok = %v, want %v", ok, test.wantOK)
			}
			want := decimal.RequireFromString(test.want)
			if !got.Equal(want) {
				t.Errorf("VolumeRatio() = %s, want %s", got, want)
			}
		})
	}
}

func repeatedVolumes(volume string, count int) []string {
	volumes := make([]string, count)
	for index := range volumes {
		volumes[index] = volume
	}
	return volumes
}

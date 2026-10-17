package mapper

import (
	domainmodel "api/domain/model"
	jquantsmodel "api/external/jquants/model"
	"fmt"
	"time"
)

func ToDailyPrices(dailyPrices []jquantsmodel.DailyPrice) ([]domainmodel.DailyPrice, error) {
	result := make([]domainmodel.DailyPrice, 0, len(dailyPrices))

	for _, dailyPrice := range dailyPrices {
		// 保存・分析には調整済み値だけを使用するため、必要な調整済み項目が欠ける日足は変換しない。
		if dailyPrice.AdjustedOpen == nil ||
			dailyPrice.AdjustedHigh == nil ||
			dailyPrice.AdjustedLow == nil ||
			dailyPrice.AdjustedClose == nil ||
			dailyPrice.AdjustedVolume == nil {
			continue
		}

		tradeDate, err := time.Parse(time.DateOnly, dailyPrice.Date)
		if err != nil {
			return nil, fmt.Errorf(
				"parse trade date for symbol %s: %w",
				dailyPrice.Code,
				err,
			)
		}

		result = append(result, domainmodel.DailyPrice{
			Symbol:     dailyPrice.Code,
			TradeDate:  tradeDate,
			OpenPrice:  *dailyPrice.AdjustedOpen,
			HighPrice:  *dailyPrice.AdjustedHigh,
			LowPrice:   *dailyPrice.AdjustedLow,
			ClosePrice: *dailyPrice.AdjustedClose,
			Volume:     *dailyPrice.AdjustedVolume,
		})
	}

	return result, nil
}

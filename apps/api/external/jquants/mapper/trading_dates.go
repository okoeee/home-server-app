package mapper

import (
	domainmodel "api/domain/model"
	jquantsmodel "api/external/jquants/model"
	"fmt"
	"time"
)

func ToTradingDates(
	listedTradingCalendar jquantsmodel.ListedTradingCalendar,
) (domainmodel.TradingDates, error) {

	result := make([]time.Time, 0, len(listedTradingCalendar.Data))

	for _, tradingCalendar := range listedTradingCalendar.Data {

		date, err := time.Parse(time.DateOnly, tradingCalendar.Date)
		if err != nil {
			return domainmodel.TradingDates{}, fmt.Errorf(
				"failed to parse date %q: %w",
				tradingCalendar.Date,
				err,
			)
		}

		result = append(result, date)

	}

	return domainmodel.TradingDates{
		Date: result,
	}, nil

}

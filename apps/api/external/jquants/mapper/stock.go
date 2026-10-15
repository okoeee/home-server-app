package mapper

import (
	domainmodel "api/domain/model"
	jquantsmodel "api/external/jquants/model"
	"fmt"
	"time"
)

func ToStocks(equityMasters []jquantsmodel.EquityMaster) ([]domainmodel.Stock, error) {
	stocks := make([]domainmodel.Stock, 0, len(equityMasters))

	for _, equityMaster := range equityMasters {
		effectiveDate, err := time.Parse(time.DateOnly, equityMaster.Date)
		if err != nil {
			return nil, fmt.Errorf(
				"parse effective date for symbol %s: %w",
				equityMaster.Code,
				err,
			)
		}

		stocks = append(stocks, domainmodel.Stock{
			Symbol:              equityMaster.Code,
			EffectiveDate:       effectiveDate,
			CompanyName:         equityMaster.CompanyName,
			Sector33Code:        equityMaster.Sector33Code,
			Sector33Name:        equityMaster.Sector33Name,
			MarketCode:          equityMaster.MarketCode,
			MarketName:          equityMaster.MarketName,
			ProductCategoryCode: equityMaster.ProductCategoryCode,
		})
	}

	return stocks, nil
}

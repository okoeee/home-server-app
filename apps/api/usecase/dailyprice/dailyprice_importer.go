package dailyprice

import (
	"api/domain/model"
	"api/domain/repository"
	"context"
	"fmt"
)

type DailyPriceImporter struct {
	repository repository.DailyPriceRepository
}

func NewDailyPriceImporter(repository repository.DailyPriceRepository) *DailyPriceImporter {
	return &DailyPriceImporter{repository}
}

func (importer *DailyPriceImporter) Import(
	ctx context.Context,
	dailyPrices []model.DailyPrice,
) error {

	for _, dailyPrice := range dailyPrices {

		err := importer.repository.Upsert(ctx, dailyPrice)
		if err != nil {
			return fmt.Errorf("upsert daily price %s: %w", dailyPrice.Symbol, err)
		}

	}

	return nil

}

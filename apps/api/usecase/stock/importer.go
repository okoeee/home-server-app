package stock

import (
	"api/domain/model"
	"api/domain/repository"
	"context"
	"fmt"
	"sort"
)

type StockImporter struct {
	repository repository.StockRepository
}

func NewStockImporter(repository repository.StockRepository) *StockImporter {
	return &StockImporter{repository: repository}
}

func (stockImporter *StockImporter) Import(ctx context.Context, stocks []model.Stock) error {
	orderedStocks := append([]model.Stock(nil), stocks...)

	sort.Slice(orderedStocks, func(left, right int) bool {
		return orderedStocks[left].Symbol < orderedStocks[right].Symbol
	})

	for _, stock := range orderedStocks {
		if err := stockImporter.repository.Upsert(ctx, stock); err != nil {
			return fmt.Errorf("upsert stock %s: %w", stock.Symbol, err)
		}
	}

	return nil
}

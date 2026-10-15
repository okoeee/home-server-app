package mysql

import (
	"api/db"
	"api/domain/model"
	"context"
)

type StockRepository struct {
	queries *db.Queries
}

func NewStockRepository(queries *db.Queries) *StockRepository {
	return &StockRepository{queries: queries}
}

func (repository *StockRepository) Upsert(
	ctx context.Context,
	stock model.Stock,
) error {

	return repository.queries.UpsertStock(
		ctx,
		db.UpsertStockParams{
			Symbol:              stock.Symbol,
			EffectiveDate:       stock.EffectiveDate,
			CompanyName:         stock.CompanyName,
			Sector33Code:        stock.Sector33Code,
			Sector33Name:        stock.Sector33Name,
			MarketCode:          stock.MarketCode,
			MarketName:          stock.MarketName,
			ProductCategoryCode: stock.ProductCategoryCode,
		},
	)

}

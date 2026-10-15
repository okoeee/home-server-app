package main

import (
	"api/db"
	"api/domain/model"
	"api/external/jquants/fetch"
	"api/external/jquants/mapper"
	jquantsmodel "api/external/jquants/model"
	"api/infrastructure/mysql"
	"api/usecase/stock"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {

	ctx := context.Background()

	jQuantsAPIKey := os.Getenv("JQUANTS_API_KEY")

	equityMasters, err := fetchEquityMasters(ctx, jQuantsAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	stocks, err := mapper.ToStocks(equityMasters)
	if err != nil {
		log.Fatal(err)
	}

	err = save(ctx, stocks)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Stocks imported: count=%d", len(stocks))

}

func fetchEquityMasters(
	ctx context.Context,
	jQuantsAPIKey string,
) ([]jquantsmodel.EquityMaster, error) {

	client, err := fetch.NewClient(jQuantsAPIKey)
	if err != nil {
		return nil, fmt.Errorf("error creating J-Quants client: %w", err)
	}

	equityMasters, err := client.FetchEquityMasterList(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching equity masters: %w", err)
	}

	return equityMasters, nil

}

func save(
	ctx context.Context,
	stocks []model.Stock,
) error {

	databaseConnection, err := mysql.Open()
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	defer databaseConnection.Close()

	transaction, err := databaseConnection.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating database transaction: %w", err)
	}

	defer transaction.Rollback()

	queries := db.New(transaction)
	stockRepository := mysql.NewStockRepository(queries)
	stockImporter := stock.NewStockImporter(stockRepository)

	err = stockImporter.Import(ctx, stocks)
	if err != nil {
		return fmt.Errorf("error importing stocks: %w", err)
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("error committing database transaction: %w", err)
	}

	return nil

}

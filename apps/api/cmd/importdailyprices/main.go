package main

import (
	"api/db"
	"api/domain/model"
	"api/external/jquants/fetch"
	"api/external/jquants/mapper"
	jquantsmodel "api/external/jquants/model"
	"api/infrastructure/mysql"
	"api/usecase/dailyprice"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("usage: go run ./cmd/importdailyprices YYYY-MM-DD")
	}

	date, err := time.Parse(time.DateOnly, os.Args[1])
	if err != nil {
		log.Fatalf("date format error: %v", err)
	}

	apiKey := os.Getenv("JQUANTS_API_KEY")

	ctx := context.Background()

	jQuantsDailyPrices, err := fetchDailyPrices(ctx, apiKey, date)
	if err != nil {
		log.Fatalf("fetch daily prices error: %v", err)
	}

	dailyPrices, err := mapper.ToDailyPrices(jQuantsDailyPrices)
	if err != nil {
		log.Fatalf("mapper to daily prices error: %v", err)
	}

	err = save(ctx, dailyPrices)
	if err != nil {
		log.Fatalf("save daily prices error: %v", err)
	}

	log.Printf("daily prices imported: date=%s, count=%d", date.Format(time.DateOnly), len(dailyPrices))

}

func fetchDailyPrices(
	ctx context.Context,
	jQuantsAPIKey string,
	date time.Time,
) ([]jquantsmodel.DailyPrice, error) {

	client, err := fetch.NewClient(jQuantsAPIKey)
	if err != nil {
		return []jquantsmodel.DailyPrice{}, fmt.Errorf("fetch daily prices client error: %w", err)
	}

	jQuantsDailyPrices, err := client.FetchDailyPrices(ctx, date)
	if err != nil {
		return []jquantsmodel.DailyPrice{}, fmt.Errorf("fetch daily prices error: %w", err)
	}

	return jQuantsDailyPrices, nil

}

func save(
	ctx context.Context,
	dailyPrices []model.DailyPrice,
) error {

	databaseConnection, err := mysql.Open()
	if err != nil {
		return fmt.Errorf("open database connection error: %w", err)
	}
	defer func() {
		err := databaseConnection.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	transaction, err := databaseConnection.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin database transaction error: %w", err)
	}
	defer func() {
		_ = transaction.Rollback()
	}()

	queries := db.New(transaction)
	dailyPriceRepository := mysql.NewDailyPriceRepository(queries)

	importer := dailyprice.NewDailyPriceImporter(dailyPriceRepository)

	err = importer.Import(ctx, dailyPrices)
	if err != nil {
		return fmt.Errorf("import daily prices error: %w", err)
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("commit database transaction error: %w", err)
	}

	return nil

}
